package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"

	"github.com/google/uuid"
)

// RandomBytes helper
func RandomBytes() []byte {
	result := make([]byte, RandomInt())

	_, err := rand.Read(result)
	if err != nil {
		panic(err)
	}

	return result
}

// RandomString helper
func RandomString(n int) string {
	result := make([]byte, n)

	_, err := rand.Read(result)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(result)
}

// RandomSeed generates random integer in ragne
func RandomSeed(min, max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)-int64(min)+1))
	if err != nil {
		return max
	}
	return min + int(n.Int64())
}

// RandomInt helper
func RandomInt() int64 {
	const (
		minValue = 1024
		maxValue = 1048576
	)

	nBig, err := rand.Int(rand.Reader, big.NewInt(maxValue))
	if err != nil {
		// This function is used in tests only, so panic here is OK
		panic(err)
	}
	if nBig.Int64() < minValue {
		return minValue
	}

	return nBig.Int64()
}

// NewObjectID new object id
func NewObjectID() string {
	return uuid.New().String()
}
