package webpush

import (
	"encoding/base64"
	"log"
	"math/big"
	"os"
	"strings"
)

type WebpushConfig struct {
	PublicKey  string
	PrivateKey string
}

var webpushConfig *WebpushConfig

func init() {
	webpushConfig = &WebpushConfig{}

	privateKey := os.Getenv("VAPID_PRIVATE_KEY")
	privateKeyInt, success := new(big.Int).SetString(privateKey, 10)
	if !success {
		log.Fatal("failed to parse big int")
	}

	webpushConfig.PrivateKey = base64.RawURLEncoding.EncodeToString(privateKeyInt.Bytes())

	publicKey := os.Getenv("VAPID_PUBLIC_KEY")
	webpushConfig.PublicKey = strings.Trim(publicKey, "=")
}

func getGlobalConfig() *WebpushConfig {
	return webpushConfig
}
