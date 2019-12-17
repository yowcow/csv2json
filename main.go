package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

var inFile *string
var outFile *string

var r io.Reader
var w io.Writer

func init() {
	inFile = flag.String("i", "", "input file to read (default: STDIN)")
	outFile = flag.String("o", "", "output file to write (default: STDOUT)")
	flag.Parse()
}

func main() {
	if *inFile != "" {
		f, err := os.Open(*inFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed opening file to read:", err)
			os.Exit(2)
		}
		defer f.Close()
		r = f
	} else {
		r = os.Stdin
	}

	if *outFile != "" {
		f, err := os.OpenFile(*outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed opening file to write:", err)
			os.Exit(3)
		}
		defer f.Close()
		w = f
	} else {
		w = os.Stdout
	}

	err := doParse(r, w)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed parsing input:", err)
		os.Exit(4)
	}
}

func doParse(r io.Reader, w io.Writer) error {
	rdr := csv.NewReader(r)
	enc := json.NewEncoder(w)

	cols, err := rdr.Read()
	if err != nil {
		return err
	}

	for {
		row, err := rdr.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		enc.Encode(buildMap(cols, row))
	}
}

func buildMap(cols []string, row []string) map[string]string {
	m := make(map[string]string)
	for i, col := range cols {
		m[col] = row[i]
	}
	return m
}
