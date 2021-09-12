package main

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
)

func Encode(u uuid.UUID) string {
	return base58.Encode(u[:])
}

func Decode(s string) (uuid.UUID, error) {
	return uuid.FromBytes(base58.Decode(s))
}
