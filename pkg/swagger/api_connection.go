/*
 * ビデオミーティングスペース管理API
 *
 * ビデオミーティングスペース管理API https://github.com/hogehoge-banana/sls-rtc-backend
 *
 * API version: 1.0.0
 * Contact: apiteam@swagger.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hogehoge-banana/sls-rtc-backend/pkg/service/connection"
	"github.com/hogehoge-banana/sls-rtc-backend/pkg/service/webpush"
)

func Connect(w http.ResponseWriter, r *http.Request) {
	var pushSub PushSubscription
	if err := json.NewDecoder(r.Body).Decode(&pushSub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p := connection.NewParticipant()
	p.P256dh = &pushSub.P256dh
	p.Auth = &pushSub.Auth
	p.Endpoint = &pushSub.Endpoint
	p.Save()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Println(err)
	}
}

func ConnectParticipantIdPingGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	participantId, ok := vars["participantId"]
	if !ok {
		respondError(w, http.StatusBadRequest, &ApiResponse{
			Message: "participantId is missing",
		})
	}
	p, err := connection.FindParticipant(participantId)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusBadRequest, &ApiResponse{
			Message: err.Error(),
		})
	}
	if p == nil {
		notfound(w)
	}
	message := "pong"
	webpush.Send(&webpush.WebpushSendParameters{
		Message:  &message,
		Endpoint: p.Endpoint,
		P256dh:   p.P256dh,
		Auth:     p.Auth,
	})

	w.WriteHeader(http.StatusOK)
}

func respondError(w http.ResponseWriter, status int, res *ApiResponse) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
	}
}

func notfound(w http.ResponseWriter) {
	respondError(w, http.StatusNotFound, &ApiResponse{
		Message: "Not found",
	})
}