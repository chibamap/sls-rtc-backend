package webpush

import (
	"fmt"
	"io"
	"log"

	wp "github.com/SherClockHolmes/webpush-go"
)

type WebpushSendParameters struct {
	// message
	Message *string
	// VAPID keys (Required)
	PrivateKey *string
	PublicKey  *string
	//
	Endpoint *string
	Auth     *string
	P256dh   *string
}

func Send(params *WebpushSendParameters) error {

	subscription := &wp.Subscription{
		Endpoint: *params.Endpoint,
		Keys: wp.Keys{
			P256dh: *params.P256dh,
			Auth:   *params.Auth,
		},
	}

	message := []byte(*params.Message)
	// message := []byte("test mesasge")

	var vapidPrivateKey string
	gconfig := getGlobalConfig()
	if params.PrivateKey != nil {
		vapidPrivateKey = *params.PrivateKey
	} else {
		vapidPrivateKey = gconfig.PrivateKey
	}
	var vapidPublicKey string
	if params.PublicKey != nil {
		vapidPublicKey = *params.PublicKey
	} else {
		vapidPublicKey = gconfig.PublicKey
	}

	res, err := wp.SendNotification(message, subscription, &wp.Options{
		Subscriber:      "example@example.com",
		VAPIDPrivateKey: vapidPrivateKey,
		VAPIDPublicKey:  vapidPublicKey,
		TTL:             30,
	})
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("status was not normal %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	log.Printf("status[%d] %s", res.StatusCode, body)
	return nil
}
