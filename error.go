package gocsv

// Error represents any error that can be returned by the gocsv package.
type Error string

func (err Error) Error() string {
	return string(err)
}

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
