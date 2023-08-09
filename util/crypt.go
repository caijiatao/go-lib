package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
)

const AesKey = "SICS-AIREC-AUTHE"

// Encrypt 使用 AES 加密数据
func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(AesKey))
	if err != nil {
		return "", err
	}

	// 填充明文数据，确保长度为16的倍数
	plainData := []byte(plainText)
	padding := aes.BlockSize - len(plainData)%aes.BlockSize
	padText := append(plainData, bytes.Repeat([]byte{byte(padding)}, padding)...)

	// 创建ECB模式的加密器
	ecb := NewECBEncrypter(block)

	// 加密
	cipherText := make([]byte, len(padText))
	err = safeCryptBlocks(ecb, cipherText, padText)
	if err != nil {
		return "", err
	}

	// 返回加密后的Base64字符串
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt 使用 AES 解密数据
func Decrypt(cipherText string) (string, error) {
	cipherData, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(AesKey))
	if err != nil {
		return "", err
	}

	// 创建ECB模式的解密器
	ecb := NewECBDecrypter(block)

	// 解密
	plainText := make([]byte, len(cipherData))
	err = safeCryptBlocks(ecb, plainText, cipherData)
	if err != nil {
		return "", err
	}

	// 去除填充数据
	padding := int(plainText[len(plainText)-1])
	plainText = plainText[:len(plainText)-padding]

	return string(plainText), nil
}

func safeCryptBlocks(mode cipher.BlockMode, dst, src []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// 将 panic 转换为错误并返回
			panicErr := fmt.Errorf("panic: %v", r)
			err = errors.WithStack(panicErr)
		}
	}()
	mode.CryptBlocks(dst, src)
	return nil
}
