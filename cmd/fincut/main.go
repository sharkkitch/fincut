package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/fincut/internal/filter"
)

func main() {
	var patternFlag string
	flag.StringVar(&patternFlag, "filter", "", "Comma-separated filter patterns (prefix with '!' to invert)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: fincut [options] [file...]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nReads from stdin if no files are provided.\n")
	}
	flag.Parse()

	patterns := []string{}
	if patternFlag != "" {
		for _, p := range strings.Split(patternFlag, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				patterns = append(patterns, p)
			}
		}
	}

	pipeline, err := filter.NewPipeline(patterns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fincut: %v\n", err)
		os.Exit(1)
	}

	files := flag.Args()
	if len(files) == 0 {
		if err := processReader(os.Stdin, pipeline); err != nil {
			fmt.Fprintf(os.Stderr, "fincut: %v\n", err)
			os.Exit(1)
		}
		return
	}

	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fincut: %v\n", err)
			os.Exit(1)
		}
		if err := processReader(f, pipeline); err != nil {
			f.Close()
			fmt.Fprintf(os.Stderr, "fincut: %v\n", err)
			os.Exit(1)
		}
		f.Close()
	}
}

func processReader(r interface{ Read([]byte) (int, error) }, p *filter.Pipeline) error {
	scanner := bufio.NewScanner(r.(interface {
		Read([]byte) (int, error)
	}))
	for scanner.Scan() {
		line := scanner.Text()
		if p.Match(line) {
			fmt.Println(line)
		}
	}
	return scanner.Err()
}
