package auth

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Hash(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

func HashSum(data ...string) string {
	hasher := md5.New()
	for _, d := range data {
		hasher.Write([]byte(d))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func Compare(first, second string) bool {
	if c := strings.Compare(first, second); c == 0 {
		return true
	}
	return false
}
