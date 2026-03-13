package utils

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// ReadRequestBody 读取请求体并返回字节数组，同时恢复请求体以便后续使用
func ReadRequestBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// 恢复请求体
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

// LogError 记录错误信息
func LogError(message string, err error) {
	log.Printf("ERROR: %s: %v", message, err)
}
