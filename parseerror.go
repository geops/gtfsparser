// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfsparser

import (
	"fmt"
)

type ParseError struct {
    file string
    line int
    msg string
}

func (e ParseError) Error() string {
    return fmt.Sprintf("%s:%d - %s", e.file, e.line, e.msg)
}