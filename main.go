package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ozaki-physics/raison-me/capital"
)

func main() {
	// fmt.Println("hello world!")
	// helloworld.Main()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/", staticFileServer())
	// 404 のときの処理
	// r.NotFound(func(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(404)
	// w.Write([]byte("404 page not found"))
	// })

	r.Route("/capital", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("これは capital だよ\n"))
		})
		r.Mount("/", capital.CryptoAsset())

	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	// サーバ起動
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

// 静的ファイルの配信をするルータ
// TODO: テスト中
func staticFileServer() chi.Router {
	r := chi.NewRouter()
	staticFileServer := http.FileServer(http.Dir("web"))
	// 動く(http://localhost:8080/robots.txt)
	// r.Mount("/", staticFileServer)

	// 動かない
	// r.Handle("/web", http.StripPrefix("/web", staticFileServer))

	// 動く(http://localhost:8080/web/robots.txt)
	r.Mount("/web", http.StripPrefix("/web", staticFileServer))

	// 動く(http://localhost:8080/web/web/robots.txt)
	// r.Mount("/web", http.StripPrefix("/web/web", staticFileServer))

	return r
}
