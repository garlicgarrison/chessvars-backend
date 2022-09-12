package format

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Returns a random array of alpha-num of specified size
//
func Random(size uint) string {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("firestore: crypto/rand.Read error: %v", err))
	}
	for i, byt := range b {
		b[i] = alphanum[int(byt)%len(alphanum)]
	}
	return string(b)
}

// UniqueID implementation from firestore
func UniqueID() string {
	return Random(20)
}

func SHA256Base64(content interface{}) (string, error) {
	b, err := json.Marshal(content)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(b)
	return base64.RawURLEncoding.EncodeToString(hash[:]), nil
}

// Hash hashes the argument
//
// It will panic if the hashing function
// errors
func Hash(toHash interface{}) string {
	identifier, err := SHA256Base64(toHash)
	if err != nil {
		panic(err)
	}

	return identifier
}
