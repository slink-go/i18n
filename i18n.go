package i18n

import (
	"errors"
	"fmt"
	"github.com/leonelquinteros/gotext"
	"go.slink.ws/logging"
	"golang.org/x/text/language"
	"os"
	"strings"
	"sync"
)

var cat = &catalog{
	cache: make(map[string]*gotext.Po),
}

func getLocalesPath() (string, error) {

	localesPath := os.Getenv("LOCALES_PATH")
	if localesPath == "" {
		return "", errors.New("env LOCALES_PATH not set")
	}

	dirInfo, err := os.Stat(localesPath)
	if os.IsNotExist(err) {
		return "", errors.New("locales path does not exist")
	}

	if !dirInfo.IsDir() {
		return "", errors.New("locales path is not a directory")
	}

	return localesPath, nil
}
func getSupportedLanguages(dir string) ([]language.Tag, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []language.Tag
	for _, entry := range entries {
		if entry.IsDir() {
			tag, err := language.Parse(entry.Name())
			if err == nil {
				result = append(result, tag)
			}
		}
	}
	return result, nil
}
func readFile(parts ...string) ([]byte, error) {
	p := strings.Join(parts, fmt.Sprintf("%c", os.PathSeparator))
	return os.ReadFile(p)
}
func load(dir string, tag language.Tag) (*gotext.Po, error) {
	entries, err := os.ReadDir(fmt.Sprintf("%s%c%s", dir, os.PathSeparator, tag.String()))
	if err != nil {
		return nil, err
	}
	sb := strings.Builder{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".po") {
			continue
		}
		d, err := readFile(dir, tag.String(), entry.Name())
		if err != nil {
			logging.GetLogger("i18n").Warning("error reading file: %s", err)
			continue
		}
		sb.WriteString(string(d))
		sb.WriteString("\n")
	}
	if sb.Len() > 0 {
		result := gotext.NewPo()
		result.Parse([]byte(sb.String()))
		return result, nil
	}
	return nil, errors.New("no PO-files found")
}

func Initialize(lpath string) (err error) {
	once := sync.Once{}
	once.Do(func() {

		//lpath, err := getLocalesPath()
		//if err != nil {
		//	panic(err)
		//}
		var locales []language.Tag
		locales, err = getSupportedLanguages(lpath)
		if err != nil {
			return
		}

		for _, l := range locales {
			var po *gotext.Po
			po, err = load(lpath, l)
			if err != nil {
				logging.GetLogger("i18n").Warning("could not load PO-file: %s", err)
				continue
			}
			cat.set(l, po)
		}
		_ = lpath
	})
	return nil
}
func T(lang, key string, args ...interface{}) string {
	tag, err := language.Parse(lang)
	if err != nil {
		return fmt.Sprintf(key, args...)
	}
	po, ok := cat.get(tag)
	if !ok {
		return fmt.Sprintf(key, args...)
	}
	return po.Get(key, args...)
}
