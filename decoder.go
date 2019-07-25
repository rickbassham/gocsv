package gocsv

import (
	"reflect"
	"strconv"
	"time"
)

// Reader is a simple reader interface. csv.Reader satisfies this.
type Reader interface {
	Read() (record []string, err error)
}

// Decoder will allow you to decode lines of CSV to structs.
type Decoder struct {
	r                   Reader
	hdr                 map[string]int
	nilVal              string
	allowMissingColumns bool
}

// ValueUnmarshaler is any type that can unmarshal it's own csv value.
type ValueUnmarshaler interface {
	UnmarshalCSVValue(string) error
}

// Unmarshaler is an interface you can implement in your struct for fine control over the decoding of the CSV.
type Unmarshaler interface {
	UnmarshalCSV([]string) error
}

// MapUnmarshaler is an interface you can implement in your struct for fine control over the decoding of the CSV.
// Instead of receiving a []string, you will get a map[string]string.
type MapUnmarshaler interface {
	UnmarshalCSVMap(map[string]string) error
}

// Must will panic if err is not nil.
func Must(dec *Decoder, err error) *Decoder {
	if err != nil {
		panic(err.Error())
	}

	return dec
}

// NewDecoder returns a new Decoder.
func NewDecoder(r Reader) *Decoder {
	return &Decoder{
		hdr: nil,
		r:   r,
	}
}

// ReadHeader will read one line from the Reader and use the line as a header.
func (dec *Decoder) ReadHeader() (*Decoder, error) {
	line, err := dec.r.Read()
	if err != nil {
		return nil, err
	}

	return dec.WithHeader(line), nil
}

// WithHeader allows you to specify a header to use. Useful if your CSV doesn't have a header line,
// or if you just don't like the header that is there.
func (dec *Decoder) WithHeader(h []string) *Decoder {
	hdr := map[string]int{}
	for i := 0; i < len(h); i++ {
		hdr[h[i]] = i
	}

	dec.hdr = hdr

	return dec
}

// Header will return the fields used as the header for the decoder.
func (dec *Decoder) Header() []string {
	hdr := make([]string, len(dec.hdr))

	for key, i := range dec.hdr {
		hdr[i] = key
	}

	return hdr
}

// WithNilValue will set the empty value for the Decoder.
func (dec *Decoder) WithNilValue(val string) *Decoder {
	dec.nilVal = val
	return dec
}

// WithAllowMissingColumns prevents the Decoder from returning an error if a column
// defined in the struct is missing from the csv.
func (dec *Decoder) WithAllowMissingColumns() *Decoder {
	dec.allowMissingColumns = true
	return dec
}

// Decode will read a line from the Reader and populate the fields in the struct passed in.
func (dec *Decoder) Decode(v interface{}) error {
	if dec.hdr == nil {
		return ErrMissingHeader
	}

	line, err := dec.r.Read()
	if err != nil {
		return err
	}

	if u, ok := v.(Unmarshaler); ok {
		return u.UnmarshalCSV(line)
	}

	if u, ok := v.(MapUnmarshaler); ok {
		return dec.decodeMapUnmarshaler(line, u)
	}

	if u, ok := v.(map[string]string); ok {
		return dec.decodeMap(line, u)
	}

	if u, ok := v.(*map[string]string); ok {
		return dec.decodeMap(line, *u)
	}

	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		return ErrInvalidType
	}

	t = t.Elem()
	val := reflect.Indirect(reflect.ValueOf(v))

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag, tagOptions := parseTag(f.Tag.Get("csv"))
		if tag == "" || tag == "-" {
			continue
		}

		index, ok := dec.hdr[tag]
		if !ok {
			if dec.allowMissingColumns {
				continue
			} else {
				return ErrMissingColumn
			}
		}

		if line[index] == dec.nilVal && omitEmpty(tagOptions) {
			continue
		}

		valf := val.FieldByName(f.Name)
		kind := f.Type.Kind()

		if kind == reflect.Ptr {
			kind = f.Type.Elem().Kind()
			valf.Set(reflect.New(f.Type.Elem()))
			valf = reflect.Indirect(valf)
		}

		i := valf.Interface()

		if u, ok := i.(ValueUnmarshaler); ok {
			if kind != reflect.Ptr {
				return ErrNonPointerReceiver
			}

			err := u.UnmarshalCSVValue(line[index])
			if err != nil {
				return err
			}
			continue
		}

		if u, ok := valf.Addr().Interface().(ValueUnmarshaler); ok {
			err := u.UnmarshalCSVValue(line[index])
			if err != nil {
				return err
			}
			continue
		}

		switch kind {
		case reflect.String:
			valf.SetString(line[index])
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(line[index])
			if err != nil {
				return err
			}
			valf.SetBool(boolVal)
		case reflect.Int:
			err := decodeInt(valf, f.Tag, 0, line[index])
			if err != nil {
				return err
			}
		case reflect.Int8:
			err := decodeInt(valf, f.Tag, 8, line[index])
			if err != nil {
				return err
			}
		case reflect.Int16:
			err := decodeInt(valf, f.Tag, 16, line[index])
			if err != nil {
				return err
			}
		case reflect.Int32:
			err := decodeInt(valf, f.Tag, 32, line[index])
			if err != nil {
				return err
			}
		case reflect.Int64:
			err := decodeInt(valf, f.Tag, 64, line[index])
			if err != nil {
				return err
			}
		case reflect.Uint:
			err := decodeUint(valf, f.Tag, 0, line[index])
			if err != nil {
				return err
			}
		case reflect.Uint8:
			err := decodeUint(valf, f.Tag, 8, line[index])
			if err != nil {
				return err
			}
		case reflect.Uint16:
			err := decodeUint(valf, f.Tag, 16, line[index])
			if err != nil {
				return err
			}
		case reflect.Uint32:
			err := decodeUint(valf, f.Tag, 32, line[index])
			if err != nil {
				return err
			}
		case reflect.Uint64:
			err := decodeUint(valf, f.Tag, 64, line[index])
			if err != nil {
				return err
			}
		case reflect.Float32:
			floatVal, err := strconv.ParseFloat(line[index], 32)
			if err != nil {
				return err
			}
			valf.SetFloat(floatVal)
		case reflect.Float64:
			floatVal, err := strconv.ParseFloat(line[index], 64)
			if err != nil {
				return err
			}
			valf.SetFloat(floatVal)
		case reflect.Struct:
			if valf.Type() == reflect.TypeOf(time.Time{}) {
				format := f.Tag.Get("format")
				if format == "" {
					format = time.RFC3339
				}

				timeVal, err := time.Parse(format, line[index])
				if err != nil {
					return err
				}
				valf.Set(reflect.ValueOf(timeVal))
			} else {
				return ErrInvalidDestType
			}
		default:
			return ErrInvalidDestType
		}
	}

	return nil
}

func (dec *Decoder) decodeMapUnmarshaler(line []string, u MapUnmarshaler) error {
	m := map[string]string{}
	for k, v := range dec.hdr {
		m[k] = line[v]
	}

	return u.UnmarshalCSVMap(m)
}

func (dec *Decoder) decodeMap(line []string, u map[string]string) error {
	for k, v := range dec.hdr {
		u[k] = line[v]
	}
	return nil
}

func decodeInt(valf reflect.Value, tag reflect.StructTag, bitSize int, value string) error {
	b, err := base(tag)
	if err != nil {
		return err
	}

	intVal, err := strconv.ParseInt(value, b, bitSize)
	if err != nil {
		return err
	}
	valf.SetInt(intVal)

	return nil
}

func decodeUint(valf reflect.Value, tag reflect.StructTag, bitSize int, value string) error {
	b, err := base(tag)
	if err != nil {
		return err
	}

	intVal, err := strconv.ParseUint(value, b, bitSize)
	if err != nil {
		return err
	}
	valf.SetUint(intVal)

	return nil
}
