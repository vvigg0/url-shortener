package generator

import (
	"crypto/rand"
	"math/big"
)

const alphabet string = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func GenerateShortCode(n int) (string, error) {
	b := make([]byte, n)

	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		b[i] = alphabet[num.Int64()]
	}
	return string(b), nil
}
