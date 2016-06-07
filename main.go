package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
)

var inPath = flag.String("in", "", "required: input file containing one URL per line")
var outPath = flag.String("out", "", "required: output file path")
var numWorkers = flag.Int("workers", 5, "number of concurrent fetchers")
var limit = flag.Int("limit", 0, "limit number of urls to process")

func main() {
	flag.Parse()

	if *inPath == "" || *outPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	urls, err := readURLs(*inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	if *limit > 0 && *limit < len(urls) {
		urls = urls[:*limit]
	}

	outfile, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	defer outfile.Close()

	outwriter := bufio.NewWriter(outfile)

	resolver := NewResolver(*numWorkers)

	var wg sync.WaitGroup
	go func(writer *bufio.Writer) {
		wg.Add(1)
		defer wg.Done()
		for result := range resolver.Resolved {
			writeProgress(writer, result)
		}
	}(outwriter)

	resolver.Start(urls)
	resolver.Wait()
	resolver.Stop()
	wg.Wait()
	fmt.Fprintf(os.Stderr, "\n")
}

func writeProgress(writer *bufio.Writer, result *Result) {
	fmt.Fprintf(os.Stderr, ".")

	out, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	writer.Write(out)
	writer.WriteRune('\n')
	writer.Flush()
}

func readURLs(path string) ([]string, error) {
	urls := []string{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	return urls, nil
}
