package handlers

import (
	"encoding/pem"
	"testing"
)

func TestCert(t *testing.T) {
	block, _ := pem.Decode([]byte(cert))
	if block == nil {
		t.Fatalf("failed to decode PEM block containing private key,block is nil")
	}
	if block.Type != "RSA PRIVATE KEY" {
		t.Fatalf("failed to decode PEM block containing private key, block.Type is %s", block.Type)
	}
	t.Log(block.Bytes)
}
