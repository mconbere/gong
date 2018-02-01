// Package gong is a simple web framework using chi (http router) and pongo2 (django-syntax template engine). Gong relies heavily on context.Context.
//
// Example:
//
// init() {
//   g := &gong.Gong{}
//   err := g.SetTemplateBaseDir("templates/")
//   if err != nil {
//     panic(err)
//   }
//   r := chi.NewRouter()
//   r.Use(gong.Inject(g))
//
//   gong.URL(r, "/", gong.NewView().FromTemplateFile("index.html").
//     GetFunc(func(w http.ResponseWriter, r *http.Request) {
//       RenderView(r.Context(), pongo2.Context{"name": "value"}, w)
//     }))
//   http.Handle("/", r)
// }
package gong

import (
	"context"
	"net/http"

	"github.com/flosch/pongo2"
)

type ContextKey struct {
	name string
}

func (c *ContextKey) String() string {
	return c.name
}

var contextKey = &ContextKey{name: "gong key"}

func FromContext(ctx context.Context) *Gong {
	v, ok := ctx.Value(contextKey).(*Gong)
	if !ok {
		return nil
	}
	return v
}

func ContextWith(ctx context.Context, g *Gong) context.Context {
	return context.WithValue(ctx, contextKey, g)
}

type Gong struct {
	templateSet *pongo2.TemplateSet
}

func (g *Gong) SetTemplateBaseDir(dir string) error {
	loader, err := pongo2.NewLocalFileSystemLoader(dir)
	if err != nil {
		return err
	}
	g.templateSet = pongo2.NewSet("gong", loader)
	return nil
}

func Inject(g *Gong) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ContextWith(r.Context(), g)))
		})
	}
}

func (g *Gong) TemplateFromFile(filename string) (*pongo2.Template, error) {
	return g.templateSet.FromFile(filename)
}
