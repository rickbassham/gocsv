package gocsv

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Writer defines the functions needed to write CSVs. This is satisfied by
// csv.Writer.
type Writer interface {
	Write([]string) error
}

// Encoder is used to encode structs to a csv.
type Encoder struct {
	w                   Writer
	hdr                 map[string]int
	nilVal              string
	allowMissingColumns bool
}

// ValueMarshaller is any type that can marshal it's own csv value.
type ValueMarshaller interface {
	MarshalCSVValue() string
}

// Marshaler can be implemented by your structs for custom marshaling logic.
type Marshaler interface {
	MarshalCSV() ([]string, error)
}

// NewEncoder creates a new Encoder.
func NewEncoder(w Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

// WithHeader specifies the header to use when writing the file. This is useful
// if you want to write the csv columns in an order other than the default
// order of the struct.
func (enc *Encoder) WithHeader(h []string) *Encoder {
	hdr := map[string]int{}

	for i, v := range h {
		hdr[v] = i
	}

	enc.hdr = hdr

	return enc
}

// WithNilValue sets the string to use for a value if it is nil.
func (enc *Encoder) WithNilValue(val string) *Encoder {
	enc.nilVal = val
	return enc
}

// WithAllowMissingColumns prevents the Encoder from returning an error if a column
// defined in the struct is missing from the header.
func (enc *Encoder) WithAllowMissingColumns() *Encoder {
	enc.allowMissingColumns = true
	return enc
}

// Encode converts v to a csv, calling MarshalCSV if v implements Marshaler or
// using reflection and any csv struct tags. If a field does not have a csv tag,
// it will be skipped.
func (enc *Encoder) Encode(v interface{}) error {
	if m, ok := v.(Marshaler); ok {
		return enc.encodeMarshaler(m)
	}

	if m, ok := v.(map[string]string); ok {
		return enc.encodeMap(m)
	}

	if m, ok := v.(*map[string]string); ok {
		return enc.encodeMap(*m)
	}

	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		return ErrInvalidType
	}

	t = t.Elem()
	val := reflect.Indirect(reflect.ValueOf(v))

	if len(enc.hdr) == 0 {
		enc.buildHeader(t)
	}

	line := make([]string, len(enc.hdr))

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag, _ := parseTag(f.Tag.Get("csv"))
		if tag == "" || tag == "-" {
			continue
		}

		index, ok := enc.hdr[tag]
		if !ok {
			if enc.allowMissingColumns {
				continue
			} else {
				return ErrMissingColumn
			}
		}

		valf := val.FieldByName(f.Name)
		kind := f.Type.Kind()

		if kind == reflect.Ptr {
			if valf.IsNil() {
				line[index] = enc.nilVal
				continue
			}
		}

		if m, ok := valf.Addr().Interface().(ValueMarshaller); ok {
			line[index] = m.MarshalCSVValue()
			continue
		}

		if m, ok := valf.Interface().(ValueMarshaller); ok {
			line[index] = m.MarshalCSVValue()
			continue
		}

		if kind == reflect.Ptr {
			kind = f.Type.Elem().Kind()
			valf = reflect.Indirect(valf)
		}

		switch kind {
		case reflect.String:
			line[index] = valf.String()
		case reflect.Bool:
			line[index] = strconv.FormatBool(valf.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if baseStr, ok := f.Tag.Lookup("base"); ok {
				base, err := strconv.ParseInt(baseStr, 10, 32)
				if err != nil {
					return ErrInvalidIntBase
				}
				line[index] = strconv.FormatInt(valf.Int(), int(base))
			} else {
				line[index] = strconv.FormatInt(valf.Int(), 10)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if baseStr, ok := f.Tag.Lookup("base"); ok {
				base, err := strconv.ParseInt(baseStr, 10, 32)
				if err != nil {
					return ErrInvalidIntBase
				}
				line[index] = strconv.FormatUint(valf.Uint(), int(base))
			} else {
				line[index] = strconv.FormatUint(valf.Uint(), 10)
			}
		case reflect.Float32, reflect.Float64:
			format := f.Tag.Get("format")
			if format == "" {
				precisionStr := f.Tag.Get("precision")
				if precisionStr == "" {
					precisionStr = "-1"
				}

				precision, err := strconv.ParseInt(precisionStr, 10, 32)
				if err != nil {
					return ErrInvalidFloatPrecision
				}

				line[index] = strconv.FormatFloat(valf.Float(), 'f', int(precision), 64)
			} else {
				line[index] = fmt.Sprintf(format, valf.Float())
			}
		case reflect.Struct:
			if valf.Type() == reflect.TypeOf(time.Time{}) {
				format := f.Tag.Get("format")
				if format == "" {
					format = time.RFC3339
				}

				line[index] = valf.Interface().(time.Time).Format(format)
			} else {
				return ErrInvalidDestType
			}
		default:
			return ErrInvalidDestType
		}
	}

	return enc.w.Write(line)
}

func (enc *Encoder) buildHeader(t reflect.Type) {
	hdr := map[string]int{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag, _ := parseTag(f.Tag.Get("csv"))
		if tag == "" || tag == "-" {
			continue
		}

		hdr[tag] = i
	}

	enc.hdr = hdr
}

func (enc *Encoder) encodeMarshaler(m Marshaler) error {
	line, err := m.MarshalCSV()
	if err != nil {
		return err
	}

	return enc.w.Write(line)
}

func (enc *Encoder) encodeMap(m map[string]string) error {
	if len(enc.hdr) == 0 {
		return ErrMissingHeader
	}

	line := make([]string, len(enc.hdr))
	for k, i := range enc.hdr {
		line[i] = m[k]
	}
	return enc.w.Write(line)
}
