package helpers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	locales = make(map[string]map[string]string)
	names   = make(map[string]string)
)

func load(dir string, filename string) error {
	ext := filepath.Ext(filename)
	locale := strings.TrimSuffix(filename, ext)
	m := make(map[string]interface{})
	buf, err := ioutil.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		return err
	}
	mloc := make(map[string]string)
	for k, v := range m {
		mloc[k] = v.(string)
	}
	locales[locale] = mloc
	names[locale] = Translate("_name", locale)
	return nil
}

func LoadLocales(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
			if e := load(dir, info.Name()); e != nil {
				return e
			}
		}
		return nil
	})
	return err
}

func Translate(s string, lang string) string {
	if ls, ok := locales[lang]; ok {
		if t, ok := ls[s]; ok {
			return t
		}
	}
	return s
}

func Locales() []string {
	ls := make([]string, 0)
	for k := range locales {
		ls = append(ls, k)
	}
	sort.Strings(ls)
	return ls
}

func NameLocale(l string) string {
	s, ok := names[l]
	if !ok {
		return l
	}
	return s
}
