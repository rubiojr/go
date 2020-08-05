package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	b64 "encoding/base64"

	"golang.org/x/crypto/scrypt"
)

// Encrypt with AES+GCM using scrypt for key derivation with a 32 byte
// random salt, encoding the resulting ciphertext with Base64
func Base64AESGCMEncrypt(key string, data []byte) (string, error) {
	newkey, salt, err := deriveKey([]byte(key), nil)
	if err != nil {
		return "", err
	}

	blockCipher, err := aes.NewCipher(newkey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	ciphertext = append(ciphertext, salt...)

	return b64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt a base64 encoded string encrypted with Base64Encrypt.
func Base64AESGCMDecrypt(key string, data []byte) ([]byte, error) {
	ddata, _ := b64.StdEncoding.DecodeString(string(data))
	salt, data := ddata[len(ddata)-32:], ddata[:len(ddata)-32]
	newkey, _, err := deriveKey([]byte(key), salt)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(newkey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func deriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
