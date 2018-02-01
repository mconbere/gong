package gong

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/go-chi/chi"
)

type View struct {
	tmplBuilder func(ctx context.Context) (*pongo2.Template, error)
	get         http.Handler
	post        http.Handler
}

func (v *View) Get(get http.Handler) *View {
	v = &(*v)
	v.get = get
	return v
}

func (v *View) GetFunc(getf http.HandlerFunc) *View {
	return v.Get(getf)
}

func (v *View) Post(post http.Handler) *View {
	v = &(*v)
	v.post = post
	return v
}

func (v *View) PostFunc(postf http.HandlerFunc) *View {
	return v.Post(postf)
}

func URL(r chi.Router, pattern string, v *View) {
	if v.get != nil {
		r.Get(pattern, func(w http.ResponseWriter, r *http.Request) {
			v.get.ServeHTTP(w, r.WithContext(ContextWithView(r.Context(), v)))
		})
	}
	if v.post != nil {
		r.Post(pattern, func(w http.ResponseWriter, r *http.Request) {
			v.post.ServeHTTP(w, r.WithContext(ContextWithView(r.Context(), v)))
		})
	}
}

func (v *View) Render(ctx context.Context, c pongo2.Context, w io.Writer) error {
	tmpl, err := v.tmplBuilder(ctx)
	if err != nil {
		return err
	}
	return Render(ctx, c, tmpl, w)
}

var viewContextKey = &ContextKey{name: "view context key"}

func ContextWithView(ctx context.Context, v *View) context.Context {
	return context.WithValue(ctx, viewContextKey, v)
}

func ViewFromContext(ctx context.Context) *View {
	v, ok := ctx.Value(viewContextKey).(*View)
	if !ok {
		return nil
	}
	return v
}

func RenderView(ctx context.Context, c pongo2.Context, w io.Writer) error {
	v := ViewFromContext(ctx)
	if v == nil {
		return errors.New("no view in context")
	}
	return v.Render(ctx, c, w)
}

func NewView() *View {
	return &View{}
}

func (v *View) FromTemplateFile(filename string) *View {
	v = &(*v)
	v.tmplBuilder = func(ctx context.Context) (*pongo2.Template, error) {
		g := FromContext(ctx)
		if g == nil {
			return nil, errors.New("no gong in context")
		}
		return g.TemplateFromFile(filename)
	}
	return v
}

type errorHandler struct {
	v *View
}

func (e *errorHandler) ServeError(w http.ResponseWriter, r *http.Request, message string, status int) error {
	c := pongo2.Context{
		"status":  status,
		"message": message,
	}
	return e.v.Render(r.Context(), c, w)
}

func ViewErrorHandler(filename string) func(http.Handler) http.Handler {
	v := NewView().FromTemplateFile(filename)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = WithServerErrorHandler(ctx, &errorHandler{v: v})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
