package main

import (
	"log"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/mconbere/gong/go/gong"
)

func indexView() *gong.View {
	return gong.NewView().FromTemplateFile("index.html").
		GetFunc(func(w http.ResponseWriter, r *http.Request) {
			gong.NotFound(w, r)
			return

			err := gong.RenderView(r.Context(), pongo2.Context{"name": "value"}, w)
			if err != nil {
				log.Printf("Server error, render: %v\n", err)
			}
		})
}

func otherView() *gong.View {
	return gong.NewView().FromTemplateFile("other.html").
		GetFunc(func(w http.ResponseWriter, r *http.Request) {
			err := gong.RenderView(r.Context(), pongo2.Context{"foo": "bar"}, w)
			if err != nil {
				gong.ServerError(w, r, http.StatusInternalServerError, "Huge mistake")
				return
			}
		})
}
