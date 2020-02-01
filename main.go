package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

// https://github.com/PuerkitoBio/goquery

type config struct {
	input     *os.File
	output    *os.File
	verbosity bool
	query     string
}

var logger *log.Logger

type nopWriter struct{}

func (*nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func parseArgs() (*config, error) {
	config := &config{
		input:     os.Stdin,
		output:    os.Stdout,
		verbosity: false,
	}

	var (
		input     = flag.String("i", "", "input file, default = stdin, skippable")
		output    = flag.String("o", "", "output file, default = stdout, skippable")
		verbosity = flag.Bool("v", false, "verbosity, default = false. if true, say logs")
		query     = flag.String("q", "", "query")
		err       error
	)
	flag.Parse()

	if *input != "" {
		config.input, err = os.OpenFile(*input, os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
	}

	if *output != "" {
		config.output, err = os.OpenFile(*output, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
	}

	if *verbosity != false {
		config.verbosity = *verbosity
	}
	setupLogger(config)

	if *query == "" {
		return nil, fmt.Errorf("query must satisfied")
	}
	config.query = *query

	return config, nil
}

func setupLogger(c *config) {
	writer := io.Writer(&nopWriter{})
	if c.verbosity {
		writer = os.Stderr
	}
	logger = log.New(writer, "", log.LstdFlags)
}

func main() {
	config, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(config.input)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(config.query).Each(func(i int, s *goquery.Selection) {
		h, err := s.Html()
		if err != nil {
			logger.Printf("failed to get html: %s", err)
		}
		config.output.WriteString(h + "\n")
	})
}
