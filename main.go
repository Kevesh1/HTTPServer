package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func Start() {
	router := newRouter()
	router.GET("/for/:id/demonstration/:otherid", func(w *utils.ResponseWriter, r *http.Request) {
		fmt.Println(r.Context().Value(ContextKey("id")))      // logs the first path parameter
		fmt.Println(r.Context().Value(ContextKey("otherid"))) // logs the second path parameter
		fmt.Println(r.FormValue("name"))                      // logs a query parameter
	})

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
	}
	fmt.Println("🐱‍💻 BeanGo server started on", l.Addr().String())
	if err := http.Serve(l, router); err != nil {
		fmt.Printf("server closed: %s\n", err)
	}
	os.Exit(1)
}
