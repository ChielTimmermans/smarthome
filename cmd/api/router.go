package main

import (
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	cors "github.com/AdhityaRamadhanus/fasthttpcors"
)

func initRouter(s *fasthttp.Server, h *Handler, config *ConfigRouter, configCORS *ConfigCORS) *router.Router {
	log.Println("Init router")

	r := router.New()

	r.RedirectTrailingSlash = config.RedirectTrailingSlash
	r.RedirectFixedPath = config.RedirectFixedPath
	r.HandleMethodNotAllowed = config.HandleMethodNotAllowed
	r.HandleOPTIONS = config.HandleOPTIONS

	CORS := initCors(configCORS)

	s.Handler = CORS.CorsMiddleware(s.Handler)

	initRoutes(r, h)
	log.Println("Init router done")

	return r
}

func initRoutes(r *router.Router, h *Handler) {
	log.Println(h.user)
	r.GET("/login", h.user.Login)
}

func initCors(config *ConfigCORS) (corsHandler *cors.CorsHandler) {
	log.Println("Init cors")
	// WithCors build cors
	corsHandler = cors.NewCorsHandler(cors.Options{
		// if you leave allowedOrigins empty then fasthttpcors will treat it as "*"
		AllowedOrigins: config.AllowedOrigins, // Only allow example.com to access the resource

		AllowCredentials: config.AllowCredentials,

		// if you leave allowedHeaders empty then fasthttpcors will accept any non-simple headers
		AllowedHeaders: config.AllowedHeaders,

		// if you leave this empty, only simple method will be accepted
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}, // only allow get or post to resource
		AllowMaxAge:    config.AllowMaxAge,                                  // cache the preflight result
		Debug:          config.Debug,                                        // turn on when strange cors behavior
	})
	log.Println("Init cors done")
	return
}
