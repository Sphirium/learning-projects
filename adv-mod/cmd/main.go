package main

import (
	"adv-mod/configs"
	"adv-mod/internal/auth"
	"adv-mod/pkg/db"
	"fmt"
	"net/http"
)

// func hello(w http.ResponseWriter, req *http.Request) {
// 	fmt.Println("Hello")
// 	// fmt.Fprintf(w, "Hello, World!")
// }

func main() {
	conf := configs.LoadConfig()
	_ = db.NewDb(conf)
	router := http.NewServeMux()
	auth.NewHelloHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	// router.HandleFunc("/hello", hello)

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()

}
