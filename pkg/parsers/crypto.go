package parsers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"errors"
)

func aesDecrypt(cipherText string, key string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	if len(data) < 16 {
		return nil, errors.New("ciphertext too short")
	}
	if bytes.HasPrefix(data, []byte("Salted__")) {
		salt := data[8:16]
		data = data[16:]
		keyBytes, iv := evpBytesToKey([]byte(key), salt, 32, 16)
		block, err := aes.NewCipher(keyBytes)
		if err != nil {
			return nil, err
		}
		mode := cipher.NewCBCDecrypter(block, iv)
		plain := make([]byte, len(data))
		mode.CryptBlocks(plain, data)
		return pkcs7Strip(plain)
	}
	keyBytes := []byte(key)
	switch {
	case len(keyBytes) < 32:
		keyBytes = append(keyBytes, make([]byte, 32-len(keyBytes))...)
	case len(keyBytes) > 32:
		keyBytes = keyBytes[:32]
	}
	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext smaller than block size")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plain := make([]byte, len(data))
	mode.CryptBlocks(plain, data)
	return pkcs7Strip(plain)
}

func pkcs7Strip(buf []byte) ([]byte, error) {
	if len(buf) == 0 {
		return nil, errors.New("pkcs7: empty buffer")
	}
	pad := int(buf[len(buf)-1])
	if pad == 0 || pad > aes.BlockSize || pad > len(buf) {
		return nil, errors.New("pkcs7: invalid padding")
	}
	return buf[:len(buf)-pad], nil
}

func evpBytesToKey(password, salt []byte, keyLen, ivLen int) ([]byte, []byte) {
	var result []byte
	var last []byte
	for len(result) < keyLen+ivLen {
		h := md5.New()
		h.Write(last)
		h.Write(password)
		h.Write(salt)
		last = h.Sum(nil)
		result = append(result, last...)
	}
	return result[:keyLen], result[keyLen : keyLen+ivLen]
}
