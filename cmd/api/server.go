package main

import (
	"log"

	"github.com/valyala/fasthttp"
)

func initServer(name string, maxRequestBodySize int) (server *fasthttp.Server) {
	log.Println("Init server")
	server = &fasthttp.Server{
		Name:               name,
		MaxRequestBodySize: maxRequestBodySize,
	}
	log.Println("Init server done")
	return
}
