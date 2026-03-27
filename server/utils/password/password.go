package password

import (
	"crypto/rand"
	"math/big"
)

const (
	// 密码字符集（不包含容易混淆的字符如 0, O, l, 1）
	passwordChars = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
)

// GenerateRandomPassword 生成指定长度的随机密码
func GenerateRandomPassword(length int) (string, error) {
	if length <= 0 {
		length = 8
	}

	password := make([]byte, length)
	charsLen := big.NewInt(int64(len(passwordChars)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsLen)
		if err != nil {
			return "", err
		}
		password[i] = passwordChars[num.Int64()]
	}

	return string(password), nil
}
