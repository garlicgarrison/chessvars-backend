package videos

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HandleUploadEmbedded(w http.ResponseWriter, r *http.Request) {
	reqBuffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var request UploadEmbeddedVideoRequest
	err = json.Unmarshal(reqBuffer, &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	
}
