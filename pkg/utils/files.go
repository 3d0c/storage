package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
)

// MakeSHA256Hash calculates SHA256 hash for input value
func MakeSHA256Hash(value []byte) string {
	const (
		salt = "1234"
	)

	h := hmac.New(sha256.New, []byte(salt))

	h.Write(value)

	return hex.EncodeToString(h.Sum(nil))
}

// BuildFilePath builds filepath like that:
// - calculates hash from objectID, e.g.: 7ed0097d7e9ee73c
// - returns /{storage-dir}/73/d0/objectID
func BuildFilePath(rootPath, objectID string) string {
	hash := MakeSHA256Hash([]byte(objectID))
	return filepath.Join(rootPath, hash[0:2], hash[2:4], objectID)
}
