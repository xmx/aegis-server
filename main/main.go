package main

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/handler/initapi"
	"github.com/xmx/aegis-server/handler/shipx"
)

func main() {
	sh := ship.Default()
	sh.HandleError = shipx.HandleError
	base := sh.Group("/api")
	route := shipx.NewRouter(base, base)

	testingAPI := initapi.Testing()
	testingAPI.Register(route)

	http.ListenAndServe(":9900", sh)
}
