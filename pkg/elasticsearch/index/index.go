package index

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/pafkiuq/backend/pkg/net/http/httputil"
	"go.uber.org/zap"
)

type Index interface {
	// Create adds the document
	//
	// Errors if already exists
	Create(ctx context.Context, docID string, doc interface{}) error
	// Upset creates a new document if docID doesn't exist;
	// simply updates if it does.
	Upsert(ctx context.Context, docID string, doc interface{}) error
	// Update allows for more advanced update calls
	Update(ctx context.Context, docID, params, script string, upsert interface{}) error
	// Updates documents if returned by query
	//
	// Updates are sent in the form of a script written in painless
	UpdateByQuery(ctx context.Context, query, params, script string) error
	// Delete deletes the given docID
	//
	// Nil error if it doesn't exist
	Delete(ctx context.Context, docID string) error
	// DeleteByQuery deletes documents that match the given query
	DeleteByQuery(ctx context.Context, query string) error

	// Creates a new PIT and returns its id
	//
	// keepAlive rounded down to the second
	OpenPointInTime(ctx context.Context, keepAlive string) (string, error)
	// Deletes the PIT with the given id
	ClosePointInTime(ctx context.Context, id string) error

	Search(ctx context.Context, body string) (*SearchResponse, error)
	// SearchWithPIT executes a search request with the given query, aggregations,
	// pit (point-in-time), and size.
	SearchWithPIT(ctx context.Context, body string) (*SearchResponse, error)
	Aggregate(ctx context.Context, body string) (*AggregateResponse, error)

	UNSAFE_RESET_INDEX_REST(*zap.Logger) http.HandlerFunc

	// Deletes the old index and creates a new one with the same name
	UNSAFE_RESET(context.Context) error
}

type index struct {
	client  *elasticsearch.Client
	name    string
	mapping string
	log     *zap.Logger
}

func NewIndex(ctx context.Context, client *elasticsearch.Client, name string, mapping string, log *zap.Logger) (Index, error) {
	resp, err := client.Indices.Exists([]string{name})
	if err != nil {
		log.Error("[NewIndex] error in checking index",
			zap.String("index name", name),
			zap.Error(err))
		return nil, err
	}

	// Create if doesn't exist
	if resp.IsError() {
		if resp.StatusCode != http.StatusNotFound {
			var e map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
				log.Error("[NewIndex] error in reading body",
					zap.String("index name", name),
					zap.Error(err),
				)
				return nil, err
			}

			err = fmt.Errorf("[%s] %s: %s", resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)

			log.Error("[NewIndex] error in checking index",
				zap.String("index name", name),
				zap.Error(err),
			)
			return nil, err
		}

		err := createIndex(ctx, client, name, mapping)
		if err != nil {
			log.Error("[NewIndex] error in creating index",
				zap.String("index name", name),
				zap.Error(err),
			)
			return nil, err
		}

		log.Info("[NewIndex] index created",
			zap.String("index name", name),
			zap.String("index mapping", mapping),
		)
	}

	log.Info("[NewIndex] connected",
		zap.String("index name", name),
		zap.String("index mapping", mapping),
	)

	return &index{
		client:  client,
		name:    name,
		mapping: mapping,
		log:     log,
	}, nil
}

func createIndex(ctx context.Context, client *elasticsearch.Client, name, mapping string) error {
	resp, err := client.Indices.Create(name,
		client.Indices.Create.WithContext(ctx),
		client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return err
		}

		return fmt.Errorf("[%s] %s: %s", resp.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}
	return nil
}

// Create adds the document
//
// Errors if already exists
func (i *index) Create(ctx context.Context, docID string, doc interface{}) error {
	request, err := json.Marshal(doc)
	if err != nil {
		i.log.Error("[Create] error in marshaling doc",
			zap.String("index", i.name),
			zap.String("docID", docID),
			zap.Any("doc", doc),
			zap.Error(err),
		)
		return err
	}

	res, err := esapi.CreateRequest{
		Index:      i.name,
		DocumentID: docID,
		Body:       bytes.NewReader(request),
	}.Do(ctx, i.client)
	if err != nil {
		i.log.Error("[Create] error in do",
			zap.String("index", i.name),
			zap.String("docID", docID),
			zap.Any("request", request),
			zap.Error(err),
		)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			i.log.Error("[Create] error in decoding error",
				zap.String("index", i.name),
				zap.String("docID", docID),
				zap.Any("request", request),
				zap.Error(err),
			)
			return err
		}

		err = fmt.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)

		i.log.Error("[Create] error in elasticsearch",
			zap.String("index", i.name),
			zap.String("docID", docID),
			zap.Any("request", request),
			zap.Error(err),
		)
		return err
	}

	i.log.Info("[Create] success",
		zap.String("index", i.name),
		zap.String("docID", docID),
		zap.Any("doc", doc),
	)
	return nil
}

const upsertRequest = `{
	"doc": %s,
	"doc_as_upsert": true
}`

// Inserts if no document exists, updates if it does
func (i *index) Upsert(ctx context.Context, docID string, doc interface{}) error {
	body, err := json.Marshal(doc)
	if err != nil {
		i.log.Error("[Upsert] error in marshaling doc",
			zap.String("index", i.name),
			zap.String("docID", docID),
			zap.Any("doc", doc),
			zap.Error(err),
		)
		return err
	}

	request := fmt.Sprintf(upsertRequest, body)
	res, err := i.client.Update(i.name,
		docID,
		strings.NewReader(request),
		i.client.Update.WithContext(ctx),
	)
	if err != nil {
		i.log.Error("[Upsert] error in do",
			zap.String("index", i.name),
			zap.String("docID", docID),
			zap.String("request", request),
			zap.Error(err),
		)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			i.log.Error("[Upsert] error in decoding error",
				zap.String("index", i.name),
				zap.String("docID", docID),
				zap.String("request", request),
				zap.Error(err),
			)
			return err
		}

		if _, ok := e["error"].(string); ok {
			i.log.Error("[Upsert] error in elasticsearch",
				zap.String("index", i.name),
				zap.String("docID", docID),
				zap.Any("request", request),
				zap.Error(err),
			)
			return errors.New(e["error"].(string))
		}

		err = fmt.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)

		if res.StatusCode == http.StatusConflict {
			i.log.Warn("[Upsert] conflict",
				zap.Any("index", i.name),
				zap.String("docID", docID),
				zap.Any("request", request),
				zap.Error(err),
			)
			return NewVersionConflictError(err)
		}

		i.log.Error("[Upsert] error in elasticsearch",
			zap.String("index", i.name),
			zap.String("docID", docID),
			zap.Any("request", request),
			zap.Error(err),
		)
		return err
	}

	i.log.Info("[Upsert] success",
		zap.String("index", i.name),
		zap.String("docID", docID),
		zap.Any("doc", doc),
		zap.Error(err),
	)
	return nil
}

const updateRequest = `{
	"script": {
		"params": %s,
		"source": %q,
		"lang": "painless"
	},
	"upsert": %s
}`

func (i *index) Update(ctx context.Context, docID, params, script string, upsert interface{}) error {
	body, err := json.Marshal(upsert)
	if err != nil {
		return err
	}

	request := fmt.Sprintf(updateRequest, params, script, body)
	i.log.Debug("[Update]",
		zap.String("request", request),
	)

	res, err := i.client.Update(i.name,
		docID,
		strings.NewReader(request),
		i.client.Update.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			i.log.Error("[Update] error in decoding error",
				zap.String("index", i.name),
				zap.String("docID", docID),
				zap.String("request", request),
				zap.Error(err),
			)
			return err
		}

		if _, ok := e["error"].(string); ok {
			i.log.Error("[Update] error in elasticsearch",
				zap.String("index", i.name),
				zap.String("docID", docID),
				zap.Any("request", request),
				zap.Error(err),
			)
			return errors.New(e["error"].(string))
		}

		err = fmt.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)

		if res.StatusCode == http.StatusConflict {
			i.log.Warn("[Update] conflict",
				zap.Any("index", i.name),
				zap.String("docID", docID),
				zap.Any("request", request),
				zap.Error(err),
			)
			return NewVersionConflictError(err)
		}

		i.log.Error("[Update] error in elasticsearch",
			zap.String("index", i.name),
			zap.String("docID", docID),
			zap.Any("request", request),
			zap.Error(err),
		)
		return err
	}

	i.log.Info("[Update] success",
		zap.String("name", i.name),
		zap.String("request", request),
	)

	return nil
}

const updateByQueryRequest = `{
	"script": {
		"params": %s,
		"source": %q,
		"lang": "painless"
	},
	"query": %s
}`

func (i *index) UpdateByQuery(ctx context.Context, query, params, script string) error {
	request := fmt.Sprintf(updateByQueryRequest, params, script, query)
	i.log.Debug("[UpdateByQuery]",
		zap.String("request", request),
	)

	resp, err := i.client.UpdateByQuery([]string{i.name},
		i.client.UpdateByQuery.WithBody(strings.NewReader(request)),
		i.client.UpdateByQuery.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return err
		}

		if e["error"] == nil {
			err = fmt.Errorf("code: %s", resp.Status())
			i.log.Error("[UpdateByQuery] error in update by query",
				zap.Error(err),
			)
			return err
		}

		if _, ok := e["error"].(string); ok {
			return errors.New(e["error"].(string))
		}

		return fmt.Errorf("[%s] %s: %s", resp.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	i.log.Info("[UpdateByQuery] success",
		zap.String("name", i.name),
		zap.String("request", request),
	)

	return nil
}

const deleteByQueryRequest = `
{
	"query": %s
}`

func (i *index) DeleteByQuery(ctx context.Context, query string) error {
	request := fmt.Sprintf(deleteByQueryRequest, query)
	i.log.Debug("[DeleteByQuery]",
		zap.String("request", request),
	)

	resp, err := i.client.DeleteByQuery([]string{i.name}, strings.NewReader(request))
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		if resp.StatusCode != http.StatusNotFound {
			var e map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
				return err
			}

			return fmt.Errorf("[%s] %s: %s", resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

// Delete deletes the document
func (i *index) Delete(ctx context.Context, docID string) error {
	resp, err := esapi.DeleteRequest{
		Index:      i.name,
		DocumentID: docID,
	}.Do(ctx, i.client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		if resp.StatusCode != http.StatusNotFound {
			var e map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
				return err
			}

			return fmt.Errorf("[%s] %s: %s", resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

func (i *index) OpenPointInTime(ctx context.Context, keepAlive string) (string, error) {
	i.log.Info("[OpenPointInTime]",
		zap.String("index", i.name),
	)

	resp, err := i.client.OpenPointInTime(
		[]string{i.name},
		keepAlive,
		i.client.OpenPointInTime.WithContext(ctx),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return "", err
		}

		return "", fmt.Errorf("[%s] %s: %s", resp.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	i.log.Info("[OpenPointInTime] create PIT Response", zap.String("body", string(body)))

	var response struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.ID, nil
}

const closePointInTimeRequest = `
{
	"id": %q
}`

func (i *index) ClosePointInTime(ctx context.Context, id string) error {
	request := fmt.Sprintf(closePointInTimeRequest, id)

	resp, err := i.client.ClosePointInTime(
		i.client.ClosePointInTime.WithContext(ctx),
		i.client.ClosePointInTime.WithBody(strings.NewReader(request)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		if resp.StatusCode != http.StatusNotFound {
			var e map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
				return err
			}

			return fmt.Errorf("[%s] %s: %s", resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

//lint:ignore U1000 unused
const searchRequest = `
{
	"query": %s,
	"aggs": %s,
	"pit": {
		"id": %q,
		"keep_alive": %q
	},
	"size": %d
}`

type AggregateResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float32 `json:"max_score"`
		Hits     []struct {
			Index      string          `json:"_index"`
			ID         string          `json:"_id"`
			Score      float32         `json:"_score"`
			Source     json.RawMessage `json:"_source"`
			Highlights json.RawMessage `json:"highlight"`
			Sort       []interface{}   `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
	Aggregations struct {
		Search struct {
			DocCountErrUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount      int `json:"sum_other_doc_count"`
			Buckets               []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
				SubAgg   struct {
					DocCountErrUpperBound int `json:"doc_count_error_upper_bound"`
					SumOtherDocCount      int `json:"sum_other_doc_count"`
					Buckets               []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"sub_agg"`
			} `json:"buckets"`
		} `json:"search"`
	} `json:"aggregations"`
}

type SearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float32 `json:"max_score"`
		Hits     []struct {
			Index      string          `json:"_index"`
			ID         string          `json:"_id"`
			Score      float32         `json:"_score"`
			Source     json.RawMessage `json:"_source"`
			Highlights json.RawMessage `json:"highlight"`
			Sort       []interface{}   `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}

func (i *index) Search(ctx context.Context, body string) (*SearchResponse, error) {
	res, err := i.client.Search(
		i.client.Search.WithContext(ctx),
		i.client.Search.WithIndex(i.name),
		i.client.Search.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response SearchResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	i.log.Info("Search Response", zap.String("resp", string(b)))

	return &response, nil
}

// Search returns results from calling Search on the index
// with the given body
//
func (i *index) SearchWithPIT(ctx context.Context, body string) (*SearchResponse, error) {
	res, err := i.client.Search(
		i.client.Search.WithContext(ctx),
		// DO NOT SPECIFY INDEX WITH POINT IN TIME
		i.client.Search.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	i.log.Info("Search Response", zap.String("resp", string(b)))

	var response SearchResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (i *index) Aggregate(ctx context.Context, body string) (*AggregateResponse, error) {
	res, err := i.client.Search(
		i.client.Search.WithContext(ctx),
		i.client.Search.WithBody(strings.NewReader(body)),
		i.client.Search.WithIndex(i.name),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	i.log.Info("Aggregate Response", zap.String("resp", string(b)))

	var response AggregateResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (i *index) UNSAFE_RESET_INDEX_REST(l *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		err := i.UNSAFE_RESET(req.Context())
		if err != nil {
			l.Error("[UNSAFE_RESET_INDEX_REST] error in unsafe_reset()",
				zap.Error(err),
				zap.Duration("dur", time.Since(start)),
			)
			httputil.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		l.Info("[UNSAFE_RESET_INDEX_REST] success",
			zap.Duration("dur", time.Since(start)),
		)
		httputil.JSONSuccess(w, http.StatusOK, nil)
	}
}

func (i *index) UNSAFE_RESET(ctx context.Context) error {
	res, err := i.client.Indices.Delete(
		[]string{i.name},
		i.client.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		return err
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	i.log.Info("[unsafe_reset] index deleted",
		zap.String("index", i.name),
	)

	return createIndex(ctx, i.client, i.name, i.mapping)
}
