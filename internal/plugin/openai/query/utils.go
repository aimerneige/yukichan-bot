package query

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"

	"github.com/google/uuid"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterNums = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenRandomLetterString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func GenRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterNums[rand.Int63()%int64(len(letterNums))]
	}
	return string(b)
}

func CalcSign(token, timestamp, nonce string, body []byte) string {
	if len(token) == 0 {
		return ""
	}
	bodyMd5 := MD5(body)
	return MD5([]byte(token + timestamp + nonce + bodyMd5))
}

func MD5(data []byte) string {
	return hex.EncodeToString(MD5Byte(data))
}

func MD5Byte(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func GetUUid() string {
	random, _ := uuid.NewRandom()
	return random.String()
}
