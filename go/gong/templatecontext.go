package gong

import (
	"context"
	"io"

	"github.com/flosch/pongo2"
)

var templateContextContextKey = &ContextKey{name: "gong template context key"}

func TemplateContextFromContext(ctx context.Context) pongo2.Context {
	v, ok := ctx.Value(templateContextContextKey).(pongo2.Context)
	if !ok {
		return pongo2.Context{}
	}
	nv := pongo2.Context{}
	nv.Update(v)
	return nv
}

func ContextWithTemplateContext(ctx context.Context, c pongo2.Context) context.Context {
	nc := pongo2.Context{}
	nc.Update(c)
	return context.WithValue(ctx, templateContextContextKey, nc)
}

func Render(ctx context.Context, c pongo2.Context, t *pongo2.Template, w io.Writer) error {
	bc := TemplateContextFromContext(ctx)
	bc = bc.Update(c)
	return t.ExecuteWriter(bc, w)
}
