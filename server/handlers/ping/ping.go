package ping

import "net/http"

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong 😁"))
	}
}
