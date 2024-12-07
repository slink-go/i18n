package i18n

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"os"
	"testing"
)

func prepare() {
	_ = os.Setenv("LOCALES_PATH", "testdata/locales")
	Initialize()
}

func TestInitialization(t *testing.T) {
	prepare()
	tag := language.MustParse("ru")
	po, ok := cat.get(tag)
	assert.True(t, ok)
	assert.NotNil(t, po)
	assert.True(t, po.IsTranslated("existing localized message: %v"))
	assert.True(t, po.IsTranslated("other localized message: %v"))
	assert.False(t, po.IsTranslated("unknown localized message: %v"))
}
func TestUnknownLang(t *testing.T) {
	prepare()
	tag := language.MustParse("zh")
	po, ok := cat.get(tag)
	assert.False(t, ok)
	assert.Nil(t, po)
}
func TestModifiedLang(t *testing.T) {
	prepare()
	tag := language.MustParse("ru-RU")
	po, ok := cat.get(tag)
	assert.True(t, ok)
	assert.NotNil(t, po)
}

func TestExistingTranslation(t *testing.T) {
	prepare()
	res := T("ru-RU", "existing localized message: %v", 1)
	assert.Equal(t, res, "существующее локализованное сообщение: 1")
	res = T("ru", "existing localized message: %v", 1)
	assert.Equal(t, res, "существующее локализованное сообщение: 1")
	res = T("RU", "existing localized message: %v", 1)
	assert.Equal(t, res, "существующее локализованное сообщение: 1")
}
func TestOtherExistingTranslation(t *testing.T) {
	prepare()
	res := T("ru-RU", "other localized message: %v", 1)
	assert.Equal(t, res, "другое локализованное сообщение: 1")
	res = T("ru", "other localized message: %v", 1)
	assert.Equal(t, res, "другое локализованное сообщение: 1")
	res = T("RU", "other localized message: %v", 1)
	assert.Equal(t, res, "другое локализованное сообщение: 1")
}
