package templates

import "context"

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
