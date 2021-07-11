// UTILITY types and functions
package converter

import (
	wp "github.com/raptium/wordpress-xml-to-hugo/pkg/parser"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// true if item is a blog post
func isPost(item wp.Item) bool {
	return item.PostType == "post"
}

// a type for making explicit what replaces what
type Replacement struct {
	From, To string
}

// strings.Replacer wants a flat list of strings
func MakeReplacer(rep ...Replacement) *strings.Replacer {
	var flatRep []string
	for _, rpl := range rep {
		flatRep = append(flatRep, rpl.From, rpl.To)
	}
	return strings.NewReplacer(flatRep...)
}

// create a parsed template, panics on failure
func MakeParsedTemplate(name string, src string) *template.Template {
	t := template.New(name)
	tp, err := t.Parse(src)
	if err != nil {
		panic(err)
	}
	return tp
}

// creates a sub-path under a base path and returns its path
func CreateSubPath(basePath string, subPath string) string {
	resultPath, err := filepath.Abs(filepath.Join(basePath, subPath))
	if err != nil {
		log.Panicln(err)
	}
	err = os.MkdirAll(resultPath, 0755)
	if err != nil {
		log.Panicln(err)
	}
	return resultPath
}
