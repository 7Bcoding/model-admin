package config

import (
	"encoding/base64"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	Config.AESKey = "a7c234ae72f64405"
	rawStrs := []string{"Xz8.kJ-Qo!2v_Rt@uY.5sP*gN3aC9x"}
	for _, rawStr := range rawStrs {
		encryptedPassword, err := AesEncrypt([]byte(rawStr), []byte(Config.AESKey))
		if err != nil {
			t.Fatalf("Failed to encrypt password: %v", err)
		}
		pass64 := base64.StdEncoding.EncodeToString(encryptedPassword)
		t.Logf("Encrypted password: %s", pass64)
	}
	// rawPassword := "sddaicll"
	// encryptedPassword, err := AesEncrypt([]byte(rawPassword), []byte(Config.AESKey))
	// if err != nil {
	// 	t.Fatalf("Failed to encrypt password: %v", err)
	// }
	// pass64 := base64.StdEncoding.EncodeToString(encryptedPassword)
	// t.Logf("Encrypted password: %s", pass64)

	// decryptedPassword, err := AesDecrypt(encryptedPassword, []byte(Config.AESKey))
	// if err != nil {
	// 	t.Fatalf("Failed to decrypt password: %v", err)
	// }
	// t.Logf("Decrypted password: %s", decryptedPassword)
}
