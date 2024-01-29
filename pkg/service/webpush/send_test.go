package webpush

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_key(t *testing.T) {
	// vapid
	privateKeyString := os.Getenv("VAPID_PRIVATE_KEY")

	privateKeyInt, success := new(big.Int).SetString(privateKeyString, 10)
	if !success {
		log.Fatal("failed to parse big int")
	}

	privateKeyBytes := privateKeyInt.Bytes()
	privateKeybase64 := base64.RawURLEncoding.EncodeToString(privateKeyBytes)
	privateKeyStdE := base64.StdEncoding.EncodeToString(privateKeyBytes)

	fmt.Println("private key RawURLEncoding ", privateKeybase64)
	fmt.Println("private key    StdEncoding ", privateKeyStdE)
	time.Sleep(3 * time.Second)
	fmt.Println("slept well.")
}

func Test_Send(t *testing.T) {
	// vapid
	publicKey := os.Getenv("VAPID_PUBLIC_KEY")
	privateKey := os.Getenv("VAPID_PRIVATE_KEY")

	privateKeyInt, success := new(big.Int).SetString(privateKey, 10)
	if !success {
		log.Fatal("failed to parse big int")
	}

	privateKeyParam := base64.RawURLEncoding.EncodeToString(privateKeyInt.Bytes())

	publicKeyForParam := strings.Trim(publicKey, "=")

	endpoint := "https://fcm.googleapis.com/fcm/send/ekUPAOMcLR8:APA91bGCl7sZbXaSt4kDOFeXuS_mdObn-LJkH-HcsAwlVpjxCZMYFE8xKkMAhGsvUo3oArr6goKMNXGctHuvXaI_EFNoXi6ZOm2Rd2bAQyJu3OgP5jsmuomKW31gnfvBZgvoaOD2flcG"
	// client parameters
	p256dh := "BI0fhAEvr3jNcG6JVjTpFrCEkltcEYiRspdcAlAVqgkiIe_meUmVOh-FwhCeEHiy308sCncMOO6i-bVVceFjsro"
	auth := "NRako11Ne9oOLF9SDgGF8Q"

	time.Sleep(1 * time.Second)
	fmt.Println("Run push.....")

	for i := 0; i < 5; i++ {
		message := fmt.Sprintf("Hello[%d], this is from go web push", i)
		fmt.Println("> push: ", message)
		params := &WebpushSendParameters{
			Message:    &message,
			Endpoint:   &endpoint,
			PrivateKey: &privateKeyParam,
			PublicKey:  &publicKeyForParam,
			P256dh:     &p256dh,
			Auth:       &auth,
		}
		err := Send(params)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("pushed.")
}
