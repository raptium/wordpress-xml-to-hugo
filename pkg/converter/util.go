// UTILITY types and functions
package converter

import (
	wp "github.com/raptium/wordpress-xml-to-hugo/pkg/parser"
	"os"
	"path/filepath"
)

// true if item is a blog post
func isPost(item wp.Item) bool {
	return item.PostType == "post"
}

// creates a sub-path under a base path and returns its path
func CreateSubPath(basePath string, subPath string) (string, error) {
	resultPath, err := filepath.Abs(filepath.Join(basePath, subPath))
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(resultPath, 0755)
	if err != nil {
		return "", err
	}
	return resultPath, nil
}
