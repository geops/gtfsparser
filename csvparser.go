// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfsparser

import (
	"encoding/csv"
	"io"
)

type CsvParser struct {
	header   []string
	reader   *csv.Reader
	filename string
	Curline  int
}

func NewCsvParser(file io.Reader) CsvParser {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	p := CsvParser{reader: reader, filename: "a"}
	p.ParseHeader()

	return p
}

func (p *CsvParser) ParseRecord() (map[string]string, error) {
	l, e := p.ParseCsvLine()

	if e != nil {
		return nil, ParseError{p.filename, p.Curline, e.Error()}
	}

	if l == nil {
		return nil, nil
	}

	record := make(map[string]string)

	for i, e := range p.header {
		if i >= len(l) {
			record[e] = ""
		} else {
			record[e] = l[i]
		}
	}

	return record, nil
}

func (p *CsvParser) ParseCsvLine() ([]string, error) {
	record, err := p.reader.Read()

	p.Curline++

	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return record, nil
}

func (p *CsvParser) ParseHeader() error {
	var e error
	p.header, e = p.ParseCsvLine()
	if e != nil {
		return e
	}
	return nil
}
