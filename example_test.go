package gocsv_test

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/rickbassham/gocsv"
)

func ExampleDecoder() {
	rdr := strings.NewReader("a,b,c\n1,2,3\n4,5,6")
	csvrdr := csv.NewReader(rdr)

	dec := gocsv.Must(gocsv.NewDecoder(csvrdr).ReadHeader())

	var err error

	for {
		val := struct {
			A string  `csv:"a"`
			B int     `csv:"b"`
			C float64 `csv:"c"`
		}{}

		err = dec.Decode(&val)
		if err != nil {
			break
		}

		println(fmt.Sprintf("%#v", val))
	}

	if err != io.EOF {
		println(err.Error())
	}
}
