package nul // import "gopkg.in/webnice/nul.v1/nul"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"testing"
)

func isStringValid(t *testing.T, s String, from string) {
	dv, _ := s.Value()
	if dv.(string) != stringTestBody {
		t.Errorf("Bad %s string: \"%s\" ≠ \"%s\"\n", from, dv, stringTestBody)
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func isStringNull(t *testing.T, s String, from string) {
	if s.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewString(t *testing.T) {
	v1 := NewString()
	isStringNull(t, v1, "NewString()")
}

func TestNewStringValue(t *testing.T) {
	v1 := NewStringValue(stringTestBody)
	isStringValid(t, v1, "NewStringValue()")
}

func TestNewStringPointerValue(t *testing.T) {
	var bv = stringTestBody

	v1 := NewStringPointerValue(&bv)
	isStringValid(t, v1, "NewStringPointerValue()")

	v2 := NewStringPointerValue(nil)
	isStringNull(t, v2, "NewStringPointerValue()")
}

func TestStringSetValid(t *testing.T) {
	v1 := NewString()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	v1.SetValid(stringTestBody)
	isStringValid(t, v1, "SetValid()")
}

func TestStringReset(t *testing.T) {
	v1 := NewString()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	v1.SetValid(``)
	if !v1.Valid {
		t.Error("Valid property", "is false, but should be true")
	}
	v1.Reset()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
}

func TestStringNullIfDefault(t *testing.T) {
	v1 := NewStringValue(stringTestBody)
	isStringValid(t, v1, "NewStringValue()")
	v1.NullIfDefault()
	isStringValid(t, v1, "NullIfDefault()")

	v1.SetValid(``)
	if !v1.Valid {
		t.Error("Valid property", "is false, but should be true")
	}
	v1.NullIfDefault()
	isStringNull(t, v1, "NullIfDefault()")
}

func TestStringMustValue(t *testing.T) {
	var buf interface{}

	v1 := NewString()
	buf = v1.MustValue()
	if _, ok := buf.(string); !ok {
		t.Error("MustValue()", "is nil, but should be not nil")
	}
	v1.SetValid(stringTestBody)
	bf2 := v1.MustValue()
	if bf2 != stringTestBody {
		t.Error("MustValue()", "is wrong")
	}
}

func TestStringPointer(t *testing.T) {
	v1 := NewStringPointerValue(nil)
	isStringNull(t, v1, "NewStringPointerValue(nil)")
	if v1.Pointer() != nil {
		t.Error("Pointer()", "is not nil, but should be nil")
	}

	v2 := NewStringValue(stringTestBody)
	pb := v2.Pointer()
	if pb == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
	if *pb != stringTestBody {
		t.Error("Pointer() value", "is wrong")
	}

	v3 := NewStringPointerValue(pb)
	isStringValid(t, v3, "NewStringPointerValue(max)")
	if v3.Pointer() == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
}

func TestStringScan(t *testing.T) {
	v1 := NewString()
	errorPanic(v1.Scan(stringTestBody))
	isStringValid(t, v1, "Scan()")

	v2 := NewString()
	errorPanic(v2.Scan(nil))
	isStringNull(t, v2, "Scan()")
}

func TestStringValue(t *testing.T) {
	var err error
	var dv driver.Value
	var ok bool
	var item string

	v1 := NewStringValue(stringTestBody)
	dv, err = v1.Value()
	errorPanic(err)
	if dv == nil {
		t.Error("Value()", "returns nil, but should be not nil")
	}
	if item, ok = dv.(string); !ok {
		t.Errorf("%s returns type %q, but should be %q", "Value()", reflect.TypeOf(dv).Name(), "string")
	}
	if item != stringTestBody {
		t.Error("Value() value", "is wrong")
	}

	v2 := NewString()
	dv, err = v2.Value()
	errorPanic(err)
	if dv != nil {
		t.Error("Value()", "returns not nil, but should be nil")
	}
}

func TestStringUnmarshalJSON(t *testing.T) {
	var err error

	v1 := NewString()
	err = json.Unmarshal(stringJSON, &v1)
	errorPanic(err)
	isStringValid(t, v1, "UnmarshalJSON()")

	v2 := NewString()
	err = json.Unmarshal(bytesNullJSON, &v2)
	errorPanic(err)
	isStringNull(t, v2, "UnmarshalJSON(null)")

	v3 := NewString()
	err = v3.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("Expected json.SyntaxError, not %T", err)
	}

	v4 := NewString()
	err = v4.UnmarshalJSON(uint64MaxValueValidJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	v5 := NewString()
	err = v5.UnmarshalJSON(blankJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	v6 := NewString()
	err = json.Unmarshal(int64JSON, &v6)
	if err == nil {
		panic("error should not be nil")
	}
	isStringNull(t, v6, "UnmarshalJSON()")
}

func TestStringMarshalJSON(t *testing.T) {
	v1 := NewStringValue(stringTestBody)
	data, err := v1.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, `"3LbOVMltCjj1Mg6sSRYLzS5j64DDNEVax29ypIGxwEx9mnbFnT9FY0sZqP11"`, "non-empty json marshal")

	v2 := NewString()
	data, err = v2.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, "null", "null json marshal")
}

func TestStringUnmarshalText(t *testing.T) {
	var err error

	v1 := NewString()
	err = v1.UnmarshalText([]byte(stringTestBody))
	errorPanic(err)
	isStringValid(t, v1, "UnmarshalText()")

	v2 := NewString()
	err = v2.UnmarshalText([]byte(""))
	errorPanic(err)
	if v2.String != "" || !v2.Valid {
		t.Errorf("Value should be valid")
	}

	v3 := NewString()
	err = v3.UnmarshalText(boolNullJSON)
	errorPanic(err)
	isStringNull(t, v3, "UnmarshalText()")
}

func TestStringMarshalText(t *testing.T) {
	v1 := NewStringValue(stringTestBody)
	data, err := v1.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, stringTestBody, "Non-empty text marshal")

	v2 := NewString()
	data, err = v2.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "null", "Null text marshal")
}

func TestStringUnmarshalBinary(t *testing.T) {
	var err error
	var btf []byte

	v1 := String{}
	btf, err = hex.DecodeString(stringNullInvalidGob)
	errorPanic(err)
	errorPanic(gob.NewDecoder(bytes.NewBuffer(btf)).Decode(&v1))
	isStringNull(t, v1, "UnmarshalBinary() invalid")

	v2 := String{}
	btf, err = hex.DecodeString(stringOkValidGob)
	errorPanic(err)
	errorPanic(gob.NewDecoder(bytes.NewBuffer(btf)).Decode(&v2))
	isStringValid(t, v2, "UnmarshalBinary() ok")
}

func TestStringMarshalBinary(t *testing.T) {
	var err error
	var buf *bytes.Buffer
	var enc *gob.Encoder

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)

	v1 := NewString()
	err = enc.Encode(&v1)
	errorPanic(err)
	v1h := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(v1h), stringNullInvalidGob, "NewBytes() -> MarshalBinary()")

	buf.Reset()
	enc = gob.NewEncoder(buf)
	v2 := NewStringValue(stringTestBody)
	err = enc.Encode(&v2)
	errorPanic(err)
	v2h := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(v2h), stringOkValidGob, "NewBoolValue(false) -> MarshalBinary()")
}
