package gocsv_test

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/rickbassham/gocsv"
)

func TestEncoder(t *testing.T) {
	val := simpleTest{
		StringVal: "this is a string",
		IntVal:    12345,
	}

	expected := "this is a string,12345\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(&val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_NonPointer(t *testing.T) {
	val := simpleTest{
		StringVal: "this is a string",
		IntVal:    12345,
	}

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != gocsv.ErrInvalidType {
		t.Error("expected ErrInvalidType")
	}
}

func TestEncoder_WithHeader(t *testing.T) {
	val := simpleTest{
		StringVal: "this is a string",
		IntVal:    12345,
	}

	expected := "12345,this is a string\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw).WithHeader([]string{"n", "str"})

	err := enc.Encode(&val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_WithHeaderSubset(t *testing.T) {
	val := simpleTest{
		StringVal: "this is a string",
		IntVal:    12345,
	}

	expected := "12345\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw).WithHeader([]string{"n"}).WithAllowMissingColumns()

	err := enc.Encode(&val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_WithHeaderUnexpectedSubset(t *testing.T) {
	val := simpleTest{
		StringVal: "this is a string",
		IntVal:    12345,
	}

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw).WithHeader([]string{"n"})

	err := enc.Encode(&val)
	if err != gocsv.ErrMissingColumn {
		t.Error("expected error ErrMissingColumn")
		t.FailNow()
	}
}

func TestEncoder_PointerValue(t *testing.T) {
	testVal := "I'm a string"
	val := simpleTestPointer{
		StringVal: &testVal,
		IntVal:    12345,
	}

	expected := "I'm a string,12345\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(&val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_NilValue(t *testing.T) {
	val := simpleTestPointer{
		StringVal: nil,
		IntVal:    12345,
	}

	expected := ",12345\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(&val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_WithNilValue(t *testing.T) {
	val := simpleTestPointer{
		StringVal: nil,
		IntVal:    12345,
	}

	expected := "\" --\",12345\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw).WithNilValue(" --")

	err := enc.Encode(&val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_NoTag(t *testing.T) {
	val := simpleTestNoTag{
		StringVal: "this is a string value",
		IntVal:    12345,
	}

	expected := "this is a string value,12345\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(&val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_ValueMarshaller(t *testing.T) {
	val := valueMarshalerTest{
		StringVal: "this is a string value",
		Test:      valueMarshaller("i'm also a string, but with a prefix"),
	}

	expected := "this is a string value,\"prefix: i'm also a string, but with a prefix\"\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(&val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_ValueMarshallerPointer(t *testing.T) {
	testVal := valueMarshaller("i'm also a string, but with a prefix")
	val := valueMarshalerPointerTest{
		StringVal: "this is a string value",
		Test:      &testVal,
	}

	expected := "this is a string value,\"prefix: i'm also a string, but with a prefix\"\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(&val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_ValueMarshallerPointerNil(t *testing.T) {
	val := valueMarshalerPointerTest{
		StringVal: "this is a string value",
		Test:      nil,
	}

	expected := "this is a string value,\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(&val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_Map(t *testing.T) {
	val := map[string]string{
		"str": "I'm a string!",
		"n":   "12345",
	}

	expected := "I'm a string!,12345\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw).WithHeader([]string{"str", "n"})

	err := enc.Encode(val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_MapPointer(t *testing.T) {
	val := &map[string]string{
		"str": "I'm a string!",
		"n":   "12345",
	}

	expected := "I'm a string!,12345\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw).WithHeader([]string{"str", "n"})

	err := enc.Encode(val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_MapNoHeader(t *testing.T) {
	val := map[string]string{
		"str": "I'm a string!",
		"n":   "12345",
	}

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(&val)
	if err != gocsv.ErrMissingHeader {
		t.Error("expected ErrMissingHeader")
		t.FailNow()
	}
}

func TestEncoder_Marshaler(t *testing.T) {
	val := &marshalerTest{
		StringVal: "I'm a string!",
		OtherVal:  12345,
	}

	expected := "I'm a string!,12345\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_MarshalerError(t *testing.T) {
	val := &marshalerTestError{
		StringVal: "I'm a string!",
		OtherVal:  12345,
	}

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err == nil || err.Error() != "i'm an error" {
		t.Error("expected i'm an error")
		t.FailNow()
	}
}

func TestEncoder_Bool(t *testing.T) {
	val := &boolTest{
		True:  true,
		False: false,
	}

	expected := "true,false\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_Int(t *testing.T) {
	val := &intTest{
		A: -15,
		B: -15,
		C: 255,
		D: 255,
	}

	expected := "-15,-f,255,377\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_IntBadBase(t *testing.T) {
	val := &intTestBadBase{
		A: -1,
	}

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != gocsv.ErrInvalidIntBase {
		t.Error("expected ErrInvalidIntBase")
		t.FailNow()
	}
}

func TestEncoder_UintBadBase(t *testing.T) {
	val := &uintTestBadBase{
		A: 1,
	}

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != gocsv.ErrInvalidIntBase {
		t.Error("expected ErrInvalidIntBase")
		t.FailNow()
	}
}

func TestEncoder_Float(t *testing.T) {
	val := &floatTest{
		A: 0.0000005,
		B: 1234567890.1234567,
		C: 1234567890.1234567,
	}

	expected := "0.0000005,1234567890.12,1234567890.12346\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_FloatBadPrecision(t *testing.T) {
	val := &floatTestBadPrecision{
		A: 0.0000005,
	}

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != gocsv.ErrInvalidFloatPrecision {
		t.Error("expected ErrInvalidFloatPrecision")
		t.FailNow()
	}
}

func TestEncoder_Time(t *testing.T) {
	val := &timeTest{
		A: time.Date(2019, 03, 9, 0, 0, 0, 0, time.UTC),
		B: time.Date(2019, 03, 9, 0, 0, 0, 0, time.UTC),
	}

	expected := "2019-03-09T00:00:00Z,2019-03-09\n"

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	csvw.Flush()

	actual := b.String()

	if actual != expected {
		t.Errorf("expected: %s got: %s", expected, actual)
	}
}

func TestEncoder_InvalidSlice(t *testing.T) {
	val := &invalidSliceTest{
		A: []string{"1", "2"},
	}

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != gocsv.ErrInvalidDestType {
		t.Error("expected ErrInvalidDestType")
		t.FailNow()
	}
}

func TestEncoder_InvalidType(t *testing.T) {
	val := &invalidTypeTest{
		A: strings.NewReader(""),
	}

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != gocsv.ErrInvalidDestType {
		t.Error("expected ErrInvalidDestType")
		t.FailNow()
	}
}

func TestEncoder_InvalidStruct(t *testing.T) {
	val := &invalidStructTest{}

	b := strings.Builder{}
	csvw := csv.NewWriter(&b)
	enc := gocsv.NewEncoder(csvw)

	err := enc.Encode(val)
	if err != gocsv.ErrInvalidDestType {
		t.Error("expected ErrInvalidDestType")
		t.FailNow()
	}
}
