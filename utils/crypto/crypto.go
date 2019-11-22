package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

// Returns a md5 for a given string
func GetMd5(s string) string {
	hash := md5.New()
	defer hash.Reset()

	_, err := hash.Write([]byte(s))
	if err != nil {
		panic(err.Error())
	}
	return hex.EncodeToString(hash.Sum(nil))
}
