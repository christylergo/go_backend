package main

import (
	"net/http"
	"time"

	"example.com/go_backend/internal/authentication"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	authentication.AddAuthenticationRouter(r)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}
	srv.ListenAndServe()
	// authen.RunMain()
}
