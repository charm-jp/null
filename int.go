package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	grpc_shared_go "github.com/charm-jp/grpc-shared-go"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"reflect"
	"strconv"
)

// Int is an nullable int64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int struct {
	sql.NullInt64
}

// NewInt creates a new Int
func NewInt(i int64, valid bool) Int {
	return Int{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

// IntFrom creates a new Int that will always be valid.
func IntFrom(i int64) Int {
	return NewInt(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int64) Int {
	if i == nil {
		return NewInt(0, false)
	}
	return NewInt(*i, true)
}

// IntFromString creates a new Int from String
func IntFromString(s String) Int {
	val, err := strconv.Atoi(s.String)
	if err != nil {
		return NewInt(0, false)
	}
	return NewInt(int64(val), s.Valid)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int) ValueOrZero() int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

func IntFromPB(s *grpc_shared_go.NullInt64) Int {
	if s != nil {
		switch x := s.Kind.(type) {
		case *grpc_shared_go.NullInt64_Data:
			return IntFrom(x.Data)
		case *grpc_shared_go.NullInt64_Null:
			return Int{
				NullInt64: sql.NullInt64{
					Int64: 999999999,
					Valid: false,
				},
			}
		}
	}

	return Int{}
}

func (i Int) ToPB() *grpc_shared_go.NullInt64 {
	if i.Valid {
		return &grpc_shared_go.NullInt64{
			Kind: &grpc_shared_go.NullInt64_Data{Data: i.Int64},
		}
	} else if i.Int64 == 999999999 {
		return &grpc_shared_go.NullInt64{
			Kind: &grpc_shared_go.NullInt64_Null{Null: structpb.NullValue_NULL_VALUE},
		}
	}

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
// It also supports unmarshalling a sql.NullInt64.
func (i *Int) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch x := v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int64)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		i.Int64, err = strconv.ParseInt(str, 10, 64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt64)
	case nil:
		// Nulled
		i.Valid = false
		i.Int64 = 999999999
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
// A non-null Int with a 0 value will not be considered zero.
func (i Int) IsZero() bool {
	return !i.Valid
}
