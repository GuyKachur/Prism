package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/happierall/l"
)

var Logger http.Handler
var OutboundClient *http.Client

func handleError(w http.ResponseWriter, err error) {
	//LOG???
	l.Error("Something went wrong!")
	l.Debug(err)
	http.Error(w, err.Error(), 500)
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message":"pong"}`))
}

func NewServer() {
	OutboundClient = &http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	r := chi.NewRouter()
	//gochi suggested middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", healthCheck)
	r.Mount("/refract", artHandler())
	http.ListenAndServe(":9090", r)
}
