# gocsv
--
    import "github.com/rickbassham/gocsv"

Package gocsv provides a flexible encoder and decoder for csv files, allowing
you to marshal and unmarshal them from and to structs.

[![Documentation](https://godoc.org/github.com/rickbassham/gocsv?status.svg)](http://godoc.org/github.com/rickbassham/gocsv)

## Usage

```go
const (
	// ErrInvalidFloatPrecision is returned during encoding if the float precision
	// in the struct tag cannot be converted to an int.
	ErrInvalidFloatPrecision = Error("gocsv: invalid float precision in struct tag")

	// ErrInvalidIntBase is returned during encoding if the int base in the struct
	// tag cannot be converted to an int.
	ErrInvalidIntBase = Error("gocsv: invalid int base in struct tag")

	// ErrInvalidType is returned if you try to Encode or Decode a non-pointer value.
	ErrInvalidType = Error("gocsv: invalid type; must be a pointer")

	// ErrMissingColumn is returned if your CSV doesn't contain a field specified in the struct.
	ErrMissingColumn = Error("gocsv: missing column in csv")

	// ErrInvalidDestType is returned if you try to Encode or Decode a column that is not a simple type.
	// Valid types are string, all varieties of int, float, bool, and time.Time
	ErrInvalidDestType = Error("gocsv: invalid destination type; must be a simple type or time.Time")

	// ErrMissingHeader is returned when you try to Decode to a struct, but the Decoder doesn't have a valid
	// header yet. Use WithHeader or ReadHeader to let the Decoder know what the file looks like.
	ErrMissingHeader = Error("gocsv: missing header; use WithHeader or ReadHeader functions first")

	// ErrNonPointerReceiver is returned when you have implemented ValueUnmarshaler with a non-pointer
	// receiver.
	ErrNonPointerReceiver = Error("gocsv: reciever for ValueUnmarshaler must be a pointer")
)
```

#### type Decoder

```go
type Decoder struct {
}
```

Decoder will allow you to decode lines of CSV to structs.

#### func  Must

```go
func Must(dec *Decoder, err error) *Decoder
```
Must will panic if err is not nil.

#### func  NewDecoder

```go
func NewDecoder(r Reader) *Decoder
```
NewDecoder returns a new Decoder.

#### func (*Decoder) Decode

```go
func (dec *Decoder) Decode(v interface{}) error
```
Decode will read a line from the Reader and populate the fields in the struct
passed in.

#### func (*Decoder) ReadHeader

```go
func (dec *Decoder) ReadHeader() (*Decoder, error)
```
ReadHeader will read one line from the Reader and use the line as a header.

#### func (*Decoder) WithAllowMissingColumns

```go
func (dec *Decoder) WithAllowMissingColumns() *Decoder
```
WithAllowMissingColumns prevents the Decoder from returning an error if a column
defined in the struct is missing from the csv.

#### func (*Decoder) WithHeader

```go
func (dec *Decoder) WithHeader(h []string) *Decoder
```
WithHeader allows you to specify a header to use. Useful if your CSV doesn't
have a header line, or if you just don't like the header that is there.

#### func (*Decoder) WithNilValue

```go
func (dec *Decoder) WithNilValue(val string) *Decoder
```
WithNilValue will set the empty value for the Decoder.

#### type Encoder

```go
type Encoder struct {
}
```

Encoder is used to encode structs to a csv.

#### func  NewEncoder

```go
func NewEncoder(w Writer) *Encoder
```
NewEncoder creates a new Encoder.

#### func (*Encoder) Encode

```go
func (enc *Encoder) Encode(v interface{}) error
```
Encode converts v to a csv, calling MarshalCSV if v implements Marshaler or
using reflection and any csv struct tags. If a field does not have a csv tag, it
will be skipped.

#### func (*Encoder) WithAllowMissingColumns

```go
func (enc *Encoder) WithAllowMissingColumns() *Encoder
```
WithAllowMissingColumns prevents the Encoder from returning an error if a column
defined in the struct is missing from the header.

#### func (*Encoder) WithHeader

```go
func (enc *Encoder) WithHeader(h []string) *Encoder
```
WithHeader specifies the header to use when writing the file. This is useful if
you want to write the csv columns in an order other than the default order of
the struct.

#### func (*Encoder) WithNilValue

```go
func (enc *Encoder) WithNilValue(val string) *Encoder
```
WithNilValue sets the string to use for a value if it is nil.

#### type Error

```go
type Error string
```

Error represents any error that can be returned by the gocsv package.

#### func (Error) Error

```go
func (err Error) Error() string
```

#### type MapUnmarshaler

```go
type MapUnmarshaler interface {
	UnmarshalCSVMap(map[string]string) error
}
```

MapUnmarshaler is an interface you can implement in your struct for fine control
over the decoding of the CSV. Instead of receiving a []string, you will get a
map[string]string.

#### type Marshaler

```go
type Marshaler interface {
	MarshalCSV() ([]string, error)
}
```

Marshaler can be implemented by your structs for custom marshaling logic.

#### type Reader

```go
type Reader interface {
	Read() (record []string, err error)
}
```

Reader is a simple reader interface. csv.Reader satisfies this.

#### type Unmarshaler

```go
type Unmarshaler interface {
	UnmarshalCSV([]string) error
}
```

Unmarshaler is an interface you can implement in your struct for fine control
over the decoding of the CSV.

#### type ValueMarshaller

```go
type ValueMarshaller interface {
	MarshalCSVValue() string
}
```

ValueMarshaller is any type that can marshal it's own csv value.

#### type ValueUnmarshaler

```go
type ValueUnmarshaler interface {
	UnmarshalCSVValue(string) error
}
```

ValueUnmarshaler is any type that can unmarshal it's own csv value.

#### type Writer

```go
type Writer interface {
	Write([]string) error
}
```

Writer defines the functions needed to write CSVs. This is satisfied by
csv.Writer.
