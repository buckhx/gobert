package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

// TextHeader is the required column name in the CSV
// It is case-sensitive
const TextHeader = "text"

// readCSV will read a CSV file and parse the records into maps keyed by column headers, requires a column with the header "text"
func readCSV(path string, d rune) ([]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows := csv.NewReader(f)
	rows.Comma = d
	headers, err := rows.Read()
	if err != nil {
		return nil, err
	}
	tdx := -1
	for i, h := range headers {
		if h == TextHeader {
			tdx = i
			break
		}
	}
	if tdx < 0 {
		return nil, fmt.Errorf("File Missing TextHeader %q", TextHeader)
	}
	recs := []map[string]string{}
	for {
		row, err := rows.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Warning:", err)
		}
		rec := make(map[string]string, len(row))
		for i, v := range row {
			rec[headers[i]] = v
		}
		recs = append(recs, rec)
	}
	return recs, nil
}
