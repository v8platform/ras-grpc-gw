package main

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/satori/go.uuid"
)

func Encode(u uuid.UUID) string {
	return base58.Encode(u.Bytes())
}

func Decode(s string) (uuid.UUID, error) {
	return uuid.FromBytes(base58.Decode(s))
}
