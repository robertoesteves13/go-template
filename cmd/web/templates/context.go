package templates

import "context"

// To setup some page data such as an title, it uses the context key/value
// data structure so it avoids prop drilling and coupling. You can add more
// keys if you need it.

type templateKey int

const (
	TemplateTitle templateKey = iota
	TemplateDescription
)


func title(ctx context.Context) string {
	title, ok := ctx.Value(TemplateTitle).(string)
	if !ok {
		return ""
	}

	return title
}

func description(ctx context.Context) string {
	title, ok := ctx.Value(TemplateDescription).(string)
	if !ok {
		return ""
	}

	return title
}
