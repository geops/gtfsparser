// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfsparser

import (
	"encoding/csv"
	"io"
	"os"
)

type CsvParser struct {
	header  []string
	reader  *csv.Reader
	filestr string
	Curline int
}

func NewCsvParser(filestr string) (CsvParser, error) {
	file, err := os.Open(filestr)
	if err != nil {
		return CsvParser{}, err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	p := CsvParser{reader: reader, filestr: filestr}
	p.ParseHeader()

	return p, nil
}

func (p *CsvParser) ParseRecord() (map[string]string, error) {
	l, e := p.ParseCsvLine()

	if e != nil {
		return nil, ParseError{p.filestr, p.Curline, e.Error()}
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
