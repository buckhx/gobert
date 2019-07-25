package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

const TextHeader = "text"

// readCSV will read a CSV file and parse the records into maps keyed by column headers, requires a column with the header "text"
func readCSV(path string) ([]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows := csv.NewReader(f)
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
			return nil, err
		}
		rec := make(map[string]string, len(row))
		for i, v := range row {
			rec[headers[i]] = v
		}
		recs = append(recs, rec)
	}
	return recs, nil
}
