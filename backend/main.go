package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ozaki-physics/raison-me/capital"
	"github.com/ozaki-physics/raison-me/delight"
	"github.com/ozaki-physics/raison-me/growth"
	"github.com/ozaki-physics/raison-me/info"
	"github.com/ozaki-physics/raison-me/regung"
	"github.com/ozaki-physics/raison-me/seed"
	"github.com/ozaki-physics/raison-me/zeit"
)

func main() {
	// fmt.Println("hello world!")
	// helloworld.Main()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// 直接 / だけでアクセスされたときは 意図的に まだ 404 にしておく
	r.Handle("/", http.HandlerFunc(http.NotFound))

	// TODO: 将来的に Router という interface を作ってもいいかも
	r.Mount("/capital", capital.Router())
	r.Mount("/delight", delight.Router())
	r.Mount("/growth", growth.Router())
	r.Mount("/info", info.Router())
	r.Mount("/regung", regung.Router())
	r.Mount("/seed", seed.Router())
	r.Mount("/zeit", zeit.Router())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	// サーバ起動
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
