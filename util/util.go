package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
	"time"

	"github.com/mr-tron/base58"
	"github.com/shopspring/decimal"
)

func EncryptKey(text string, passphrase string) (string, error) {
	data := []byte(text)
	key := []byte(createMD5Hash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// fmt.Println("block #1: ", block)

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	noneSize := gcm.NonceSize()
	// fmt.Println("noneSize: ", noneSize)
	nonce := make([]byte, noneSize)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	// s := []byte(seed)
	// if len(s) < len(nonce) {
	// 	for i := 0; i < len(s); i++ {
	// 		nonce[i] = s[i]
	// 	}
	// } else {
	// 	nonce = s
	// }
	// fmt.Println("nonce[:noneSize]: ", nonce[:noneSize])
	// seal := gcm.Seal(nonce[:noneSize], nonce[:noneSize], data, nil)
	seal := gcm.Seal(nonce, nonce, data, nil)

	return base58.Encode(seal), nil
}
func DecryptKey(text string, passphrase string) (string, error) {
	pt, err := base58.Decode(text)
	if err != nil {
		return "", err
	}

	data := []byte(pt)
	key := []byte(createMD5Hash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// fmt.Println("block #2: ", block)

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	// fmt.Println("gcm #2: ", gcm)

	nonceSize := gcm.NonceSize()
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func CreateSHA1Hash(s string) string {
	str := sha1.Sum([]byte(s))
	return hex.EncodeToString(str[:])
}
func Timestamp() int64 {
	return time.Now().UTC().UnixNano() / 1000000
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

func createMD5Hash(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
