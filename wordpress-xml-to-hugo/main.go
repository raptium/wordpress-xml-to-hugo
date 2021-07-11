// wordpress-xml-to-hugo parses an XML export from WordPress and generates Markdown files for Hugo

package main

import (
	"flag"
	"github.com/raptium/wordpress-xml-to-hugo/pkg/converter"
	"github.com/raptium/wordpress-xml-to-hugo/pkg/parser"
	"log"
)

// commandline processing only. Everything else is in pkg/converter
func main() {
	options := converter.ParseOptions()
	if !options.IsValid() {
		flag.Usage()
		return
	}

	parsed, err := parser.Parse(options.InputFile)
	if err != nil {
		log.Panicf("failed to parse input file: %v", err)
	}

	c := converter.NewConverter(options)

	c.Convert(parsed.Channel.Items, options.OutputDirectory)
	log.Printf("parsed and converted a file with %d items\n", len(parsed.Channel.Items))
}
