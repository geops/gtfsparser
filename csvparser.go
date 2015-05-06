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
	header  []string
	reader  *csv.Reader
	Curline int
}

func NewCsvParser(file io.Reader) CsvParser {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	p := CsvParser{reader: reader}
	p.ParseHeader()

	return p
}

func (p *CsvParser) ParseRecord() map[string]string {
	l := p.ParseCsvLine()

	if l == nil {
		return nil
	}

	record := make(map[string]string)

	for i, e := range p.header {
		if i >= len(l) {
			record[e] = ""
		} else {
			record[e] = l[i]
		}
	}

	return record
}

func (p *CsvParser) ParseCsvLine() []string {
	record, err := p.reader.Read()

	p.Curline++

	if err == io.EOF {
		return nil
	} else if err != nil {
		panic(err.Error())
	}
	return record
}

func (p *CsvParser) ParseHeader() {
	p.header = p.ParseCsvLine()
}
