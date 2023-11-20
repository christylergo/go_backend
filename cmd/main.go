package main

import (
	"net/http"

	"example.com/go_backend/internal/authen"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	authen.AddAuthenticationRouter(r)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	srv.ListenAndServe()
	// authen.RunMain()
}
