package main

import (
	"fmt"

	"github.com/hogehoge-banana/sls-rtc-backend/pkg/service/webpush"
)

func main() {

	publicKey := "BFjC7jrOkbhPusXPqBi1ltpKS7tvPn5tjlDIC63rLw1MK7_HhQN4KycX6_RiuzRpxn3eJtBe1aFCcMD72KtGMcc="
	privateKey := "43772159654668669220459069873425773141257647526502194455077172281818679546385"

	endpoint := "https://fcm.googleapis.com/fcm/send/cBgCXaphzNk:APA91bHGaANx18lniKq0AkeYnGl4euMav0bAMKNztSUwWFNHaAKk2Jf6vaDcpOjvIdkFrXgCv8LT6R9JFbSfCJESEP1MWzgUexYEcXE0KR_9H4RNM8mV86ZMDJZn16AMZFmgs4XDXCep"

	message := "Hello, this is from go web push"
	params := &webpush.WebpushSendParameters{
		Message:    &message,
		Endpoint:   &endpoint,
		PrivateKey: &privateKey,
		PublicKey:  &publicKey,
	}
	err := webpush.Send(params)
	fmt.Println(err)
}
