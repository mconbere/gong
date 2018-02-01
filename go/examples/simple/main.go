package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mconbere/gong/go/gong"
)

func main() {
	g := &gong.Gong{}
	err := g.SetTemplateBaseDir("templates/")
	if err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.GetHead)
	r.Use(gong.Inject(g))

	gong.URL(r, "/", indexView())
	gong.URL(r, "/other", otherView())

	http.Handle("/", r)
	http.ListenAndServe(":8081", nil)
}
