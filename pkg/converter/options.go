package converter

import "flag"

type Options struct {
	InputFile       string
	OutputDirectory string
	SiteUrl         string
	PostDirectory   string
}

func ParseOptions() *Options {
	options := &Options{}
	flag.StringVar(&options.InputFile, "f", "", "Input file name")
	flag.StringVar(&options.OutputDirectory, "o", "", "Output directory")
	flag.StringVar(&options.SiteUrl, "s", "", "Site url, eg. https://www.example.net")
	flag.StringVar(&options.PostDirectory, "post-dir", "post", "directory for posts")
	flag.Parse()
	return options
}

func (o Options) IsValid() bool {
	return o.InputFile != "" && o.OutputDirectory != "" && o.SiteUrl != ""
}
