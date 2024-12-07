package i18n

import (
	"github.com/leonelquinteros/gotext"
	"golang.org/x/text/language"
)

type catalog struct {
	cache map[string]*gotext.Po
}

func (c *catalog) set(l language.Tag, pos *gotext.Po) {
	b, _ := l.Base()
	c.cache[b.String()] = pos
}
func (c *catalog) get(l language.Tag) (*gotext.Po, bool) {
	b, _ := l.Base()
	v, ok := c.cache[b.String()]
	return v, ok
}
