package signer

type Signer interface {
	Sign(key string) (string, error)
}
