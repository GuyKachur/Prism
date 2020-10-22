package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/happierall/l"
	"github.com/pkg/errors"
)

// var Logger http.Handler
var OutboundClient *http.Client

// var dc *database.Datastore

func HandleError(w http.ResponseWriter, err error) {
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
	r.Mount("/", artHandler())
	http.ListenAndServe(":9090", PanicHandler(r))
}

func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				er := errors.Errorf("Panic! in the go code... %v", err)
				http.Error(w, er.Error(), 500)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
