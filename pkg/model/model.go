package model

import (
	"gopkg.in/yaml.v2"
	"io"
	"time"
)

type FrontMatterContent struct {
	FrontMatter FrontMatter
	Content     string
}

func (fmc FrontMatterContent) WriteTo(writer io.Writer) error {
	frontMatter, err := yaml.Marshal(&fmc.FrontMatter)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte("---\n"))
	if err != nil {
		return err
	}
	_, err = writer.Write(frontMatter)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte("---\n"))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(fmc.Content))
	return err
}

type FrontMatter struct {
	Aliases       []string `yaml:"aliases,omitempty"`
	Description   string   `yaml:"description,omitempty"`
	Draft         bool     `yaml:"draft"`
	IsCJKLanguage bool     `yaml:"isCJKLanguage"`
	Date          Time     `yaml:"date"`
	LastMod       Time     `yaml:"lastmod"`
	PublishDate   Time     `yaml:"publishDate"`
	Summary       string   `yaml:"summary,omitempty"`
	Title         string   `yaml:"title"`
	Type          string   `yaml:"type"`
	Url           string   `yaml:"url"`
	Tags          []string `yaml:"tags"`
	Categories    []string `yaml:"categories"`
}

type Time struct {
	time.Time
}

func (t Time) MarshalYAML() (interface{}, error) {
	return t.Local().Format("2006-01-02T15:04:05-07:00"), nil
}
