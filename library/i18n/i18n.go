package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type Detector interface {
	Detect(tags []language.Tag, key message.Reference, args ...any) string
}

func NewDetector(b *catalog.Builder) Detector {
	if b == nil {
		b = catalog.NewBuilder(catalog.Fallback(language.English))
	}
	tags := b.Languages()
	matcher := language.NewMatcher(tags, language.PreferSameScript(true))
	option := message.Catalog(b)

	return &localizeDetector{
		builder: b,
		matcher: matcher,
		option:  option,
	}
}

type localizeDetector struct {
	builder *catalog.Builder
	matcher language.Matcher
	option  message.Option
}

func (ld *localizeDetector) Detect(tags []language.Tag, key message.Reference, args ...any) string {
	tag, _, _ := ld.matcher.Match(tags...)
	printer := message.NewPrinter(tag, ld.option)
	return printer.Sprintf(key, args...)
}
