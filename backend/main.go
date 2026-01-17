package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ozaki-physics/raison-me/capital"
	"github.com/ozaki-physics/raison-me/delight"
	"github.com/ozaki-physics/raison-me/growth"
	"github.com/ozaki-physics/raison-me/info"
	"github.com/ozaki-physics/raison-me/regung"
	"github.com/ozaki-physics/raison-me/seed"
	globalConfig "github.com/ozaki-physics/raison-me/share/config"
	"github.com/ozaki-physics/raison-me/trybigquery"
	"github.com/ozaki-physics/raison-me/trysupabase"
	"github.com/ozaki-physics/raison-me/zeit"
)

func main() {
	// fmt.Println("hello world!")
	// helloworld.Main()

	// trybigquery.Try()
	// trylocaldb.Try()
	// trysupabase.Try()
	Run()
}

func Run() {
	globalConfig := globalConfig.NewConfig()
	// log.Printf("globalConfig: %v", globalConfig)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// 静的ファイル の 配信
	r.Mount("/", staticFileRouter())
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

	r.HandleFunc("/test-supabase", func(w http.ResponseWriter, req *http.Request) {
		records := trysupabase.Try()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("trysupabase executed" + "\n" + records))
	})
	r.HandleFunc("/test-bigquery", func(w http.ResponseWriter, req *http.Request) {
		records := trybigquery.Try()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("trybigquery executed" + "\n" + records))
	})

	port := globalConfig.GetPort()
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

// 静的ファイル の 配信をする Router
func staticFileRouter() chi.Router {
	r := chi.NewRouter()
	r.Mount("/", http.FileServer(http.Dir("web")))
	return r
}
