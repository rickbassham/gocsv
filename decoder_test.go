package gocsv_test

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/rickbassham/gocsv"
)

func TestDecoder_MustPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic i'm an error")
		}
	}()

	gocsv.Must(nil, errors.New("i'm an error"))
}

func TestDecoder(t *testing.T) {
	data := strings.NewReader("this is a string,12345\n")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"})

	testVal := simpleTest{}

	err := dec.Decode(&testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.StringVal != "this is a string" {
		t.Errorf("testVal.StringVal expected: %s got: %s", "this is a string", testVal.StringVal)
	}

	if testVal.IntVal != 12345 {
		t.Errorf("testVal.IntVal expected: %d got: %d", 12345, testVal.IntVal)
	}
}

func TestDecoder_ReadHeader(t *testing.T) {
	data := strings.NewReader("str,n\nthis is a string,12345\n")
	r := csv.NewReader(data)
	dec := gocsv.Must(gocsv.NewDecoder(r).ReadHeader())

	hdr := dec.Header()
	if len(hdr) != 2 {
		t.Errorf("hdr len expected: 2; got %d", len(hdr))
		return
	}

	if hdr[0] != "str" {
		t.Errorf("hdr[0] expected \"str\"; got \"%s\"", hdr[0])
		return
	}

	if hdr[1] != "n" {
		t.Errorf("hdr[1] expected \"str\"; got \"%s\"", hdr[1])
		return
	}

	testVal := simpleTest{}

	err := dec.Decode(&testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.StringVal != "this is a string" {
		t.Errorf("testVal.StringVal expected: %s got: %s", "this is a string", testVal.StringVal)
	}

	if testVal.IntVal != 12345 {
		t.Errorf("testVal.IntVal expected: %d got: %d", 12345, testVal.IntVal)
	}
}

func TestDecoder_ReadHeaderError(t *testing.T) {
	data := &bytes.Buffer{}
	r := csv.NewReader(data)
	_, err := gocsv.NewDecoder(r).ReadHeader()

	if err != io.EOF {
		t.Error("expected io.EOF")
	}
}

func TestDecoder_MissingHeader(t *testing.T) {
	data := strings.NewReader("this is a string,12345\n")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r)

	testVal := simpleTest{}

	err := dec.Decode(&testVal)
	if err != gocsv.ErrMissingHeader {
		t.Error("expected gocsv.ErrMissingHeader")
		return
	}
}

func TestDecoder_ReadError(t *testing.T) {
	data := &bytes.Buffer{}
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"})

	testVal := simpleTest{}

	err := dec.Decode(&testVal)
	if err != io.EOF {
		t.Error("expected io.EOF")
	}
}

func TestDecoder_NonPointer(t *testing.T) {
	data := strings.NewReader("string,1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"})

	testVal := simpleTest{}

	err := dec.Decode(testVal)
	if err != gocsv.ErrInvalidType {
		t.Error("expected gocsv.ErrInvalidType")
	}
}

func TestDecoder_WithHeaderSubset(t *testing.T) {
	data := strings.NewReader("1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"n"}).WithAllowMissingColumns()

	testVal := simpleTest{}

	err := dec.Decode(&testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.IntVal != 1234 {
		t.Errorf("testVal.IntVal expected %d but got %d", 1234, testVal.IntVal)
	}
}

func TestDecoder_WithHeaderUnexpectedSubset(t *testing.T) {
	data := strings.NewReader("1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"n"})

	testVal := simpleTest{}

	err := dec.Decode(&testVal)
	if err != gocsv.ErrMissingColumn {
		t.Error("expected gocsv.ErrMissingColumn")
	}
}

func TestDecoder_PointerValue(t *testing.T) {
	data := strings.NewReader("string val,1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"})

	testVal := simpleTestPointer{}

	err := dec.Decode(&testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.StringVal == nil {
		t.Error("testVal.StringVal nil")
		return
	}

	if *testVal.StringVal != "string val" {
		t.Errorf("testVal.StringVal expected %s but got %s", "string val", *testVal.StringVal)
	}

	if testVal.IntVal != 1234 {
		t.Errorf("testVal.IntVal expected %d but got %d", 1234, testVal.IntVal)
	}
}

func TestDecoder_NilValue(t *testing.T) {
	data := strings.NewReader(",1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"})

	testVal := simpleTestPointer{}

	err := dec.Decode(&testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.StringVal != nil {
		t.Error("testVal.StringVal should be nil")
	}

	if testVal.IntVal != 1234 {
		t.Errorf("testVal.IntVal expected %d but got %d", 1234, testVal.IntVal)
	}
}

func TestDecoder_WithNilValue(t *testing.T) {
	data := strings.NewReader(" --,1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"}).WithNilValue(" --")

	testVal := simpleTestPointer{}

	err := dec.Decode(&testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.StringVal != nil {
		t.Error("testVal.StringVal should be nil")
	}

	if testVal.IntVal != 1234 {
		t.Errorf("testVal.IntVal expected %d but got %d", 1234, testVal.IntVal)
	}
}

func TestDecoder_NoTag(t *testing.T) {
	data := strings.NewReader("string val,1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"})

	testVal := simpleTestNoTag{
		IntVal2: 2,
		IntVal3: 3,
	}

	err := dec.Decode(&testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.StringVal != "string val" {
		t.Errorf("testVal.StringVal expected string val but got %s", testVal.StringVal)
	}

	if testVal.IntVal != 1234 {
		t.Errorf("testVal.IntVal expected %d but got %d", 1234, testVal.IntVal)
	}

	if testVal.IntVal2 != 2 {
		t.Errorf("testVal.IntVal2 expected %d but got %d", 2, testVal.IntVal2)
	}

	if testVal.IntVal3 != 3 {
		t.Errorf("testVal.IntVal3 expected %d but got %d", 3, testVal.IntVal3)
	}
}

func TestDecoder_ValueUnmarshaler(t *testing.T) {
	data := strings.NewReader("string val,other")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "prefixed"})

	testVal := valueMarshalerTest{}

	err := dec.Decode(&testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.StringVal != "string val" {
		t.Errorf("testVal.StringVal expected string val but got %s", testVal.StringVal)
	}

	if testVal.Test != "prefix other" {
		t.Errorf("testVal.Test expected %s but got %s", "prefix other", testVal.Test)
	}
}

func TestDecoder_ValueUnmarshalerPointer(t *testing.T) {
	data := strings.NewReader("string val,other")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "prefixed"})

	testVal := valueMarshalerPointerTest{}

	err := dec.Decode(&testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.StringVal != "string val" {
		t.Errorf("testVal.StringVal expected string val but got %s", testVal.StringVal)
	}

	if testVal.Test == nil {
		t.Error("testVal.Test should not be nil")
		return
	}

	if *testVal.Test != "prefix other" {
		t.Errorf("testVal.Test expected %s but got %s", "prefix other", *testVal.Test)
	}
}

func TestDecoder_Map(t *testing.T) {
	data := strings.NewReader("string val,1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"})

	testVal := map[string]string{}

	err := dec.Decode(testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal["str"] != "string val" {
		t.Errorf("testVal[str] expected string val but got %s", testVal["str"])
	}

	if testVal["n"] != "1234" {
		t.Errorf("testVal[n] expected %s but got %s", "1234", testVal["n"])
	}
}

func TestDecoder_MapPointer(t *testing.T) {
	data := strings.NewReader("string val,1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"})

	testVal := &map[string]string{}

	err := dec.Decode(testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if (*testVal)["str"] != "string val" {
		t.Errorf("testVal[str] expected string val but got %s", (*testVal)["str"])
	}

	if (*testVal)["n"] != "1234" {
		t.Errorf("testVal[n] expected %s but got %s", "1234", (*testVal)["n"])
	}
}

func TestDecoder_MapNoHeader(t *testing.T) {
	data := strings.NewReader("string val,1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r)

	testVal := map[string]string{}

	err := dec.Decode(testVal)
	if err != gocsv.ErrMissingHeader {
		t.Error("expected ErrMissingHeader")
	}
}

func TestDecoder_Unmarshaler(t *testing.T) {
	data := strings.NewReader("string val,1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"})

	testVal := &marshalerTest{}

	err := dec.Decode(testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.StringVal != "string val" {
		t.Errorf("testVal.StringVal expected string val but got %s", testVal.StringVal)
	}

	if testVal.OtherVal != 1234 {
		t.Errorf("testVal.OtherVal expected %s but got %d", "1234", testVal.OtherVal)
	}
}

func TestDecoder_MapUnmarshaler(t *testing.T) {
	data := strings.NewReader("string val,1234")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"str", "n"})

	testVal := &mapMarshalerTest{}

	err := dec.Decode(testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.StringVal != "string val" {
		t.Errorf("testVal.StringVal expected string val but got %s", testVal.StringVal)
	}

	if testVal.OtherVal != 1234 {
		t.Errorf("testVal.OtherVal expected %s but got %d", "1234", testVal.OtherVal)
	}
}

func TestDecoder_Bool(t *testing.T) {
	data := strings.NewReader("true,false")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"t", "f"})

	testVal := &boolTest{}

	err := dec.Decode(testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if !testVal.True {
		t.Errorf("testVal.True expected true but got %v", testVal.True)
	}

	if testVal.False {
		t.Errorf("testVal.False expected false but got %v", testVal.False)
	}
}

func TestDecoder_Int(t *testing.T) {
	data := strings.NewReader("-15,-f,255,377")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"a", "b", "c", "d"})

	testVal := &intTest{}

	err := dec.Decode(testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.A != -15 {
		t.Errorf("testVal.A expected %d but got %d", -15, testVal.A)
	}

	if testVal.B != -15 {
		t.Errorf("testVal.B expected %d but got %d", -15, testVal.B)
	}

	if testVal.C != 255 {
		t.Errorf("testVal.C expected %d but got %d", -15, testVal.C)
	}

	if testVal.D != 255 {
		t.Errorf("testVal.D expected %d but got %d", -15, testVal.D)
	}
}

func TestDecoder_Float(t *testing.T) {
	data := strings.NewReader("0.0000005,1234567890.12,1234567890.12346")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"a", "b", "c"})

	testVal := &floatTest{}

	err := dec.Decode(testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if testVal.A != 0.0000005 {
		t.Errorf("testVal.A expected %f but got %f", 0.0000005, testVal.A)
	}

	if testVal.B != 1234567890.12 {
		t.Errorf("testVal.B expected %f but got %f", 1234567890.12, testVal.B)
	}

	if testVal.C != 1234567890.12346 {
		t.Errorf("testVal.C expected %f but got %f", 1234567890.12346, testVal.C)
	}
}

func TestDecoder_Time(t *testing.T) {
	data := strings.NewReader("2019-03-09T00:00:00Z,2019-03-09")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"a", "b"})

	testVal := &timeTest{}

	err := dec.Decode(testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if !testVal.A.Equal(time.Date(2019, 03, 9, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("testVal.A expected %s but got %s", time.Date(2019, 03, 9, 0, 0, 0, 0, time.UTC).String(), testVal.A.String())
	}

	if !testVal.B.Equal(time.Date(2019, 03, 9, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("testVal.B expected %s but got %s", time.Date(2019, 03, 9, 0, 0, 0, 0, time.UTC).String(), testVal.B.String())
	}
}

func TestDecoder_TimeWithTimezone(t *testing.T) {
	data := strings.NewReader("2019-03-09T00:00:00Z,2019-03-09")
	r := csv.NewReader(data)
	dec := gocsv.NewDecoder(r).WithHeader([]string{"a", "b"})

	testVal := &timeZoneTest{}

	err := dec.Decode(testVal)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if !testVal.A.Equal(time.Date(2019, 03, 9, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("testVal.A expected %s but got %s", time.Date(2019, 03, 9, 0, 0, 0, 0, time.UTC).String(), testVal.A.String())
	}

	if !testVal.B.Equal(time.Date(2019, 03, 9, 6, 0, 0, 0, time.UTC)) {
		t.Errorf("testVal.B expected %s but got %s", time.Date(2019, 03, 9, 6, 0, 0, 0, time.UTC).String(), testVal.B.String())
	}
}
