package gocsv_test

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

type simpleTest struct {
	StringVal string `csv:"str"`
	IntVal    int    `csv:"n"`
}

type simpleTestPointer struct {
	StringVal *string `csv:"str,omitempty"`
	IntVal    int     `csv:"n"`
}

type simpleTestNoTag struct {
	StringVal string `csv:"str"`
	IntVal    int    `csv:"n"`
	IntVal2   int    `csv:"-"`
	IntVal3   int
}

type valueMarshaller string

func (m *valueMarshaller) MarshalCSVValue() string {
	return fmt.Sprintf("prefix: %s", string(*m))
}

func (m *valueMarshaller) UnmarshalCSVValue(val string) error {
	*m = valueMarshaller("prefix " + val)
	return nil
}

type valueMarshalerTest struct {
	StringVal string          `csv:"str"`
	Test      valueMarshaller `csv:"prefixed"`
}

type valueMarshalerPointerTest struct {
	StringVal string           `csv:"str"`
	Test      *valueMarshaller `csv:"prefixed"`
}

type marshalerTest struct {
	StringVal string `csv:"str"`
	OtherVal  int    `csv:"n"`
}

func (m *marshalerTest) MarshalCSV() ([]string, error) {
	return []string{
		m.StringVal,
		strconv.Itoa(m.OtherVal),
	}, nil
}

func (m *marshalerTest) UnmarshalCSV(line []string) error {
	var err error

	m.StringVal = line[0]
	m.OtherVal, err = strconv.Atoi(line[1])

	return err
}

type mapMarshalerTest struct {
	StringVal string `csv:"str"`
	OtherVal  int    `csv:"n"`
}

func (m *mapMarshalerTest) UnmarshalCSVMap(line map[string]string) error {
	var err error

	m.StringVal = line["str"]
	m.OtherVal, err = strconv.Atoi(line["n"])

	return err
}

type marshalerTestError struct {
	StringVal string `csv:"str"`
	OtherVal  int    `csv:"n"`
}

func (m *marshalerTestError) MarshalCSV() ([]string, error) {
	return nil, errors.New("i'm an error")
}

type boolTest struct {
	True  bool `csv:"t"`
	False bool `csv:"f"`
}

type intTest struct {
	A int  `csv:"a"`
	B int  `csv:"b" base:"16"`
	C uint `csv:"c"`
	D uint `csv:"d" base:"8"`
}

type intTestBadBase struct {
	A int `csv:"a" base:"z"`
}

type uintTestBadBase struct {
	A uint `csv:"a" base:"z"`
}

type floatTest struct {
	A float64 `csv:"a"`
	B float64 `csv:"b" format:"%9.2f"`
	C float64 `csv:"c" precision:"5"`
}

type floatTestBadPrecision struct {
	A float64 `csv:"c" precision:"Z"`
}

type timeTest struct {
	A time.Time `csv:"a"`
	B time.Time `csv:"b" format:"2006-01-02"`
}

type invalidSliceTest struct {
	A []string `csv:"a"`
}

type invalidTypeTest struct {
	A io.Reader `csv:"a"`
}

type invalidStructTest struct {
	A struct {
		C string
	} `csv:"a"`
}
