//go:generate go-bindata -o server/assets.go -pkg server -ignore ".DS_Store" -ignore ".scss" assets/...
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vjeantet/bitfan/webui/server"
)

func main() {

	httpServerMux := http.NewServeMux()

	httpServerMux.Handle(
		"/ui/",
		server.Handler(
			"assets",
			"ui",
			"127.0.0.1:5123",
		),
	)

	addr := "127.0.0.1:8080"
	if port := os.Getenv("PORT"); len(port) > 0 {
		fmt.Printf("Environment variable PORT=\"%s\"", port)
		addr = "127.0.0.1:" + port
	}
	http.ListenAndServe(addr, httpServerMux)
}