package query

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 调用 aes 算法加密并编码
func Encrypt(encodingAesKey, msg string) (string, error) {
	// aes 加密
	encryptedMsg, err := aesCBCEncrypt(encodingAesKey, []byte(msg))
	if err != nil {
		return "", err
	}
	// base64 编码
	base64Msg := make([]byte, base64.StdEncoding.EncodedLen(len(encryptedMsg)))
	base64.StdEncoding.Encode(base64Msg, encryptedMsg)
	return string(base64Msg), nil
}

// 解码后调用 aes 算法解密
func Decrypt(encodingAesKey, msg string) (string, error) {
	// base64 解码
	cipherText, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}
	// aes 解密
	plainText, err := aesCBCDecrypt(encodingAesKey, []byte(cipherText))
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func aesCBCEncrypt(encodingKey string, plaintextMsg []byte) ([]byte, error) {
	block, aesKey, err := decodeAesKey(encodingKey)
	if err != nil {
		return []byte{}, err
	}
	plaintextMsg = pkcs5Padding(plaintextMsg, block.BlockSize())
	cipherText := make([]byte, len(plaintextMsg))
	iv := aesKey[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plaintextMsg)
	return cipherText, nil
}

func aesCBCDecrypt(encodingKey string, encryptedMsg []byte) ([]byte, error) {
	block, aesKey, err := decodeAesKey(encodingKey)
	if err != nil {
		return []byte{}, err
	}
	iv := aesKey[:aes.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)
	decryptedMsg := make([]byte, len(encryptedMsg))
	mode.CryptBlocks(decryptedMsg, encryptedMsg)
	decryptedMsg = pkcs5UnPadding(decryptedMsg)
	return decryptedMsg, nil
}

func decodeAesKey(encodingAesKey string) (cipher.Block, []byte, error) {
	aesKey, err := base64.StdEncoding.DecodeString(encodingAesKey + "=")
	if err != nil {
		return nil, aesKey, err
	}
	block, e := aes.NewCipher(aesKey)
	if e != nil {
		return nil, aesKey, err
	}
	return block, aesKey, nil
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
