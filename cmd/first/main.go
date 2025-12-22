package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/genians/minsuyu/cmd/first/app"
)

func main() {
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello from VM!")
		})

		_ = http.ListenAndServe(":3245", nil)
	}()

	app.Run()

	// select {}
}
