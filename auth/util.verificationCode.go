package auth

import (
	"math/rand"
	"time"

	"github.com/spf13/viper"
)

const charset = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var linkExpiration = viper.GetInt("VERIFICATION_CODE.EXPIRATION")
var size = viper.GetInt("VERIFICATION_CODE.SIZE")

func generateCode() string {
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
