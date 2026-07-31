package router

import "net/http"

func SetupRoutes(
	mux *http.ServeMux,
	presignHandler http.Handler,
) {
	mux.Handle("/presign", presignHandler)
}
