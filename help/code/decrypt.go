package code

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"github.com/webchen/gotools/help/logs"
)

// 这样无效，转为字符串之后，IV的长度变成了32
func getIVstr() string {
	ivArr := [16]byte{
		0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7,
		0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}
	var buffer bytes.Buffer
	for _, v := range ivArr {
		buffer.WriteString(string(v))
	}
	return buffer.String()
	// ðñòóôõö÷øùúûüýþÿ
}

func getIV() []byte {
	/*
		return []byte{
			0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7,
			0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}
	*/

	//return []byte{'ð', 'ñ', 'ò', 'ó', 'ô', 'õ', 'ö', '÷', 'ø', 'ù', 'ú', 'û', 'ü', 'ý', 'þ', 'ÿ'}

	return []byte{'\xf0', '\xf1', '\xf2', '\xf3', '\xf4', '\xf5', '\xf6', '\xf7', '\xf8', '\xf9', '\xfa', '\xfb', '\xfc', '\xfd', '\xfe', '\xff'}
}

// AesCtrEncrypt AES的CTR加密模式
func AesCtrEncrypt(ct string, key string) string {
	plainText, _ := hex.DecodeString(ct)
	keyHex, _ := hex.DecodeString(key)
	block, err := aes.NewCipher(keyHex)
	if err != nil {
		logs.Error("AesCtrEncrypt Error", err)
		return ""
	}
	iv := getIV()
	stream := cipher.NewCTR(block, iv)

	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, []byte(plainText))

	return hex.EncodeToString(dst)
}

// AesCtrDecrypt 解密
func AesCtrDecrypt(encryptData string, key string) string {
	return AesCtrEncrypt(encryptData, key)
}
