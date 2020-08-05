package crypto

import (
	"testing"
)

func TestBase64AESGCMEncrypt(t *testing.T) {
	encrypted, err := Base64AESGCMEncrypt("secret", []byte("test data"))
	if err != nil {
		t.Fatal(err)
	}

	if encrypted == "test data" {
		t.Error("error: data wasn't encrypted")
	}

	decrypted, err := Base64AESGCMDecrypt("secret", []byte(encrypted))
	if err != nil {
		t.Fatal(err)
	}

	if string(decrypted) != "test data" {
		t.Error("error: decrypted data doesn't match")
	}
}
