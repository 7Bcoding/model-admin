package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

const (
	DefaultNonce = "JUSTIN"
	DefaultKey   = "QUlETi1HU0xC"
)

// AuthHeader 包含认证所需的所有头部信息
type AuthHeader struct {
	Sign      string
	Nonce     string
	Timestamp string
}

// GenerateAuthHeaders 生成认证头信息
func GenerateAuthHeaders() AuthHeader {
	timeInt := time.Now().UnixMilli()
	timeStr := strconv.FormatInt(timeInt, 10)
	sign := GetSign(DefaultKey, DefaultNonce, timeStr)

	return AuthHeader{
		Sign:      sign,
		Nonce:     DefaultNonce,
		Timestamp: timeStr,
	}
}

// GetSign 计算签名
func GetSign(key, nonce, timestamp string) string {
	keyBytes := []byte(key)
	message := []byte(nonce + timestamp)

	h := hmac.New(sha256.New, keyBytes)
	h.Write(message)

	return hex.EncodeToString(h.Sum(nil))
}
