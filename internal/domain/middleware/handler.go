package middleware

import "github.com/valyala/fasthttp"

type Handler interface {
	AccessControl(next fasthttp.RequestHandler, roles ...string) fasthttp.RequestHandler
	AvailableServices(next fasthttp.RequestHandler, services ...string) fasthttp.RequestHandler
}

const (
	DB    = "db"
	MINIO = "minio"
)
