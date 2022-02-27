package handlers

import (
	"github.com/xy3/photos"
	"net/http"
)

type route struct {
	path        string
	handlerFunc func(w http.ResponseWriter, r *http.Request)
	middleware  []photos.Middleware
}

var (
	routes = []route{
		{
			path:        "/user/signup",
			handlerFunc: Signup,
			middleware:  nil,
		}, {
			path:        "/user/signin",
			handlerFunc: Signin,
			middleware:  nil,
		}, {
			path:        "/photo",
			handlerFunc: Photo,
			middleware:  []photos.Middleware{photos.BasicAuthMiddleware},
		}, {
			path:        "/photo/info",
			handlerFunc: PhotoInfo,
			middleware:  []photos.Middleware{photos.BasicAuthMiddleware},
		}, {
			path:        "/photo/list",
			handlerFunc: PhotoList,
			middleware:  []photos.Middleware{photos.BasicAuthMiddleware},
		}, {
			path:        "/user",
			handlerFunc: User,
			middleware:  []photos.Middleware{photos.BasicAuthMiddleware},
		},
	}
)

func SetupRoutes(mux *http.ServeMux) {
	for _, r := range routes {
		handlerFunc := r.handlerFunc
		// here we can apply multiple middleware to a handler func
		if r.middleware != nil {
			for _, middlewareFunc := range r.middleware {
				handlerFunc = middlewareFunc(handlerFunc)
			}
		}
		mux.HandleFunc(photos.Config.BasePath+r.path, handlerFunc)
	}
}
