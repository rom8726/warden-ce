package crypt

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

func EncryptAESGCM(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

func DecryptAESGCM(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

// Sign signs the data with the private key.
func Sign(data []byte, privateKey crypto.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(data)

	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
		if err != nil {
			return nil, fmt.Errorf("sign with RSA: %w", err)
		}

		return signature, nil
	case *ecdsa.PrivateKey:
		signature, err := ecdsa.SignASN1(rand.Reader, key, hashed[:])
		if err != nil {
			return nil, fmt.Errorf("sign with ECDSA: %w", err)
		}

		return signature, nil
	default:
		return nil, fmt.Errorf("unsupported private key type: %T", privateKey)
	}
}

// Verify verifies the signature of the data with the public key.
func Verify(data, signature []byte, publicKey crypto.PublicKey) error {
	hashed := sha256.Sum256(data)

	switch key := publicKey.(type) {
	case *rsa.PublicKey:
		err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature)
		if err != nil {
			return fmt.Errorf("verify with RSA: %w", err)
		}

		return nil
	case *ecdsa.PublicKey:
		valid := ecdsa.VerifyASN1(key, hashed[:], signature)
		if !valid {
			return errors.New("ECDSA signature verification failed")
		}

		return nil
	default:
		return fmt.Errorf("unsupported public key type: %T", publicKey)
	}
}
