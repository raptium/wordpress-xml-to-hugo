// opinionated conversion from WordPress to Hugo
package converter

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/raptium/wordpress-xml-to-hugo/pkg/model"
	"net/url"
	"time"

	// a fork of github.com/grokify/wordpress-xml-go with parsing of comments added
	wp "github.com/raptium/wordpress-xml-to-hugo/pkg/parser"
	"os"
	"path/filepath"
	"strings"
)

type ContentConverter = func(string) string

type WpConverter struct {
	options          *Options
	contentConverter ContentConverter
}

func NewConverter(options *Options) *WpConverter {
	mdConverter := md.NewConverter("", true, nil)
	convertContent := func(c string) string {
		converted, err := mdConverter.ConvertString(c)
		if err != nil {
			return c
		}
		return converted
	}

	converter := &WpConverter{
		options:          options,
		contentConverter: convertContent,
	}
	return converter
}

// Convert all items
func (wc *WpConverter) Convert(items []wp.Item, targetBaseDir string) error {
	var fmc *model.FrontMatterContent
	var err error
	for _, item := range items {
		fmc = nil
		if isPost(item) {
			fmc, err = wc.buildContent(item)
			if err != nil {
				return err
			}
		}
		if fmc == nil {
			continue
		}
		if err := wc.writeItem(fmc); err != nil {
			return err
		}
	}
	return nil
}

// convert an item according to a template
func (wc *WpConverter) writeItem(fmc *model.FrontMatterContent) error {
	baseDir := filepath.Join(wc.options.OutputDirectory, fmc.FrontMatter.Type)
	if fmc.FrontMatter.Type == "post" && wc.options.PostDirectory != "" {
		baseDir = filepath.Join(wc.options.OutputDirectory, wc.options.PostDirectory)
	}
	filePath := strings.TrimSuffix(filepath.Join(baseDir, fmc.FrontMatter.Url), "/") + ".md"
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	err = fmc.WriteTo(f)
	if err != nil {
		return err
	}

	return nil
}

func (wc *WpConverter) buildContent(item wp.Item) (*model.FrontMatterContent, error) {

	tags := make([]string, 0)
	categories := make([]string, 0)
	for _, category := range item.Categories {
		if category.Domain == "category" {
			categories = append(categories, category.DisplayName)
		}
		if category.Domain == "post_tag" {
			tags = append(tags, category.DisplayName)
		}
	}

	date, err := time.Parse("2006-01-02 15:04:05", item.PostDateGmt)
	if err != nil {
		return nil, err
	}
	lastMod, err := time.Parse("2006-01-02 15:04:05", item.PostModifiedGmt)
	if err != nil {
		return nil, err
	}
	publishDate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
	if err != nil {
		return nil, err
	}

	summary := item.Excerpt
	if summary == "" {
		moreIndex := strings.Index(item.Content, "<!--more-->")
		if moreIndex != -1 {
			summary = item.Content[0:moreIndex]
		}
	}

	item.Content = wc.contentConverter(item.Content)
	summary = wc.contentConverter(summary)

	path, err := url.PathUnescape(item.Link)
	if err != nil {
		return nil, err
	}
	path = "/" + strings.TrimPrefix(strings.TrimPrefix(path, wc.options.SiteUrl), "/")

	frontMatter := model.FrontMatter{
		Type:          "post",
		Title:         item.Title,
		Url:           path,
		Tags:          tags,
		Categories:    categories,
		Date:          model.Time{Time: date},
		LastMod:       model.Time{Time: lastMod},
		PublishDate:   model.Time{Time: publishDate},
		Draft:         item.Status != "publish",
		IsCJKLanguage: true,
		Summary:       summary,
	}

	return &model.FrontMatterContent{
		FrontMatter: frontMatter,
		Content:     item.Content,
	}, nil
}

// takes a func as handler to make it testable
func HandleComments(commentDir string, item wp.Item, handler func(wp.Comment, string, int) error) error {
	// capture replyTo relationships
	repliesTo := make(map[int]int)
	for _, c := range item.Comments {
		repliesTo[c.Id] = c.Parent
	}
	// determine names and write files
	for _, c := range item.Comments {
		commentFileName, indentLevel := GetCommentFileNameAndIndentLevel(repliesTo, c, commentDir)
		err := handler(c, commentFileName, indentLevel)
		if err != nil {
			return err
		}

	}
	return nil
}

// construct comment filename reflecting replyTo relationship, determine indent level
func GetCommentFileNameAndIndentLevel(repliesTo map[int]int, c wp.Comment, commentDir string) (string, int) {
	id := c.Id
	name := fmt.Sprintf("_%d", id)
	loop := true
	depth := -1
	for loop {
		parentId := repliesTo[id]
		name = fmt.Sprintf("_%d%s", parentId, name)
		if parentId == 0 {
			loop = false
		} else {
			id = parentId
		}
		depth++
	}
	commentFileName := commentDir +
		string(filepath.Separator) +
		fmt.Sprintf("comment%s.json", name)
	return commentFileName, depth
}

// write the comment
//func convertComment(comment wp.Comment, commentFileName string, indentLevel int) error {
//	// set indentation
//	comment.IndentLevel = indentLevel
//
//	// own comments may need replacements
//	comment = FixCommentAuthor(comment)
//	comment.AuthorUrl = UrlReplacer1.Replace(comment.AuthorUrl)
//	comment.AuthorUrl = UrlReplacer2.Replace(comment.AuthorUrl)
//	comment.Content = QuotesReplacer.Replace(comment.Content)
//	comment.Content = EmojiReplacer.Replace(comment.Content)
//
//	// open comment file
//	f, err := os.OpenFile(commentFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	// write file
//	err = CommentTemplate.Execute(f, comment)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
