package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"crypto/md5"
)

var (
	secretKey = []byte("vnnaEPK8CJbXGuSk2qa9Zh2VetP")
)

func RandomInt(min, max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(max-min) + min
}

func Sha256(message []byte) string {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hex.EncodeToString(expectedMAC)
}

func MakeTrackSecret(message string) string {
	msg := []byte(message)
	hash := Sha256(msg)
	return hash
}

func Base64Decode(str string) string {
	sDec, _ := base64.StdEncoding.DecodeString(str)
	return string(sDec)
}

func Base64Encode(data string) string {
	enc := base64.StdEncoding.EncodeToString([]byte(data))
	return enc
}

func MD5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
