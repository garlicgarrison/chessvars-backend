package signer

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
)

const (
	RESOURCE_TIME time.Duration = time.Hour
)

type SignerConfig struct {
	CloudfrontPrivateKeyPath string
	CloudfrontKeyID          string

	CloudfrontEndpoint string
	ResourceDuration   time.Duration
}

type signer struct {
	endpoint  string
	duration  time.Duration
	urlSigner *sign.URLSigner
}

func NewSigner(cfg *SignerConfig) (Signer, error) {
	// Get private key stored at FWAYGOPATH/.keys/cloudfront_private_key.pem
	privateKeyFile, err := os.Open(cfg.CloudfrontPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	defer privateKeyFile.Close()

	privateKeyBytes, err := ioutil.ReadAll(privateKeyFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("Failed to decode PEM block")
	} else if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("Failed to decode: PEM block is not a private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	urlSigner := sign.NewURLSigner(cfg.CloudfrontKeyID, privateKey)

	return &signer{
		endpoint:  cfg.CloudfrontEndpoint,
		duration:  cfg.ResourceDuration,
		urlSigner: urlSigner,
	}, nil
}

func (s *signer) Sign(key string) (string, error) {
	signed, err := s.urlSigner.Sign(s.endpoint+"/"+key, time.Now().Add(s.duration))
	if err != nil {
		return "", err
	}
	return signed, nil
}
