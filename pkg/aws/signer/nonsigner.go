package signer

import "fmt"

type nonSigner struct {
	CloudfrontEndpoint string
}

func NewNonSigner(endpoint string) (Signer, error) {
	return &nonSigner{
		CloudfrontEndpoint: endpoint,
	}, nil
}

func (s *nonSigner) Sign(key string) (string, error) {
	return fmt.Sprintf("%s/%s", s.CloudfrontEndpoint, key), nil
}
