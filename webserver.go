package main

import (
	"fmt"
	"log"
	"ms-sg/config"
	"ms-sg/server/web"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	host := config.File.MustValue("web_server", "host", "127.0.0.1")
	port := config.File.MustValue("web_server", "port", "9992")
	g := gin.Default()
	web.Init(g)
	s := &http.Server{
		Addr: fmt.Sprintf("%s:%s", host, port),
		Handler: g,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1<<20,
	}
	err := s.ListenAndServe()
	log.Println(err)
}