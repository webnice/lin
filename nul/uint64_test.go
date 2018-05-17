package nul // import "gopkg.in/webnice/nul.v1/nul"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func isUint64Valid(t *testing.T, u Uint64, from string) {
	dv := interface{}(u.Uint64)
	if dv.(uint64) != uint64(math.MaxUint64) {
		t.Errorf("Bad %s uint64: \"%d\" ≠ \"%d\"\n", from, dv, uint64(math.MaxUint64))
	}
	if !u.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func isUint64Null(t *testing.T, u Uint64, from string) {
	if u.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewUint64(t *testing.T) {
	v1 := NewUint64()
	isUint64Null(t, v1, "NewUint64()")
}

func TestNewUint64Value(t *testing.T) {
	v1 := NewUint64Value(math.MaxUint64)
	isUint64Valid(t, v1, "NewUint64Value()")
}

func TestNewUint64PointerValue(t *testing.T) {
	var bv = uint64(math.MaxUint64)

	v1 := NewUint64PointerValue(&bv)
	isUint64Valid(t, v1, "NewUint64PointerValue()")

	v2 := NewUint64PointerValue(nil)
	isUint64Null(t, v2, "NewUint64PointerValue()")
}

func TestUint64SetValid(t *testing.T) {
	v1 := NewUint64()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	v1.SetValid(uint64(math.MaxUint64))
	isUint64Valid(t, v1, "SetValid()")
}

func TestUint64Reset(t *testing.T) {
	v1 := NewUint64()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	v1.SetValid(0)
	if !v1.Valid {
		t.Error("Valid property", "is false, but should be true")
	}
	v1.Reset()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
}

func TestUint64NullIfDefault(t *testing.T) {
	v1 := NewUint64Value(uint64(math.MaxUint64))
	isUint64Valid(t, v1, "NewUint64Value()")
	v1.NullIfDefault()
	isUint64Valid(t, v1, "NullIfDefault()")

	v1.SetValid(0)
	if !v1.Valid {
		t.Error("Valid property", "is false, but should be true")
	}
	v1.NullIfDefault()
	isUint64Null(t, v1, "NullIfDefault()")
}

func TestUint64MustValue(t *testing.T) {
	var buf interface{}

	v1 := NewUint64()
	buf = v1.MustValue()
	if _, ok := buf.(uint64); !ok {
		t.Error("MustValue()", "is nil, but should be not nil")
	}
	v1.SetValid(uint64(math.MaxUint64))
	bf2 := v1.MustValue()
	if bf2 != math.MaxUint64 {
		t.Error("MustValue()", "is wrong")
	}
}

func TestUint64Pointer(t *testing.T) {
	v1 := NewUint64PointerValue(nil)
	isUint64Null(t, v1, "NewUint64PointerValue(nil)")
	if v1.Pointer() != nil {
		t.Error("Pointer()", "is not nil, but should be nil")
	}

	v2 := NewUint64Value(uint64(math.MaxUint64))
	pb := v2.Pointer()
	if pb == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
	if *pb != uint64(math.MaxUint64) {
		t.Error("Pointer() value", "is wrong")
	}

	v3 := NewUint64PointerValue(pb)
	isUint64Valid(t, v3, "NewUint64PointerValue(max)")
	if v3.Pointer() == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
}

func TestUint64Scan(t *testing.T) {
	v1 := NewUint64()
	errorPanic(v1.Scan(uint64(math.MaxUint64)))
	isUint64Valid(t, v1, "Scan()")

	v2 := NewUint64()
	errorPanic(v2.Scan(uint64String))
	isUint64Valid(t, v2, "Scan()")

	v3 := NewUint64()
	errorPanic(v3.Scan(nil))
	isUint64Null(t, v3, "Scan()")

	v4 := NewUint64()
	err := v4.Scan(false)
	if err == nil {
		t.Error("Scan()", "is nil, but should be not nil")
	}
}

func TestUint64Value(t *testing.T) {
	var err error
	var dv driver.Value
	var ok bool
	var item []byte

	v1 := NewUint64Value(uint64(math.MaxUint64))
	dv, err = v1.Value()
	errorPanic(err)
	if dv == nil {
		t.Error("Value()", "returns nil, but should be not nil")
	}
	if item, ok = dv.([]byte); !ok {
		t.Errorf("%s returns type %q, but should be %q", "Value()", reflect.TypeOf(dv).Name(), "string")
	}
	if !bytes.Equal(item, uint64JSON) {
		t.Error("Value() value", "is wrong")
	}

	v2 := NewUint64()
	dv, err = v2.Value()
	errorPanic(err)
	if dv != nil {
		t.Error("Value()", "returns not nil, but should be nil")
	}
}

func TestUint64UnmarshalJSON(t *testing.T) {
	var err error

	v1 := NewUint64()
	err = json.Unmarshal(uint64JSON, &v1)
	errorPanic(err)
	isUint64Valid(t, v1, "UnmarshalJSON()")

	v2 := NewUint64()
	err = json.Unmarshal(bytesNullJSON, &v2)
	errorPanic(err)
	isUint64Null(t, v2, "UnmarshalJSON(null)")

	v3 := NewUint64()
	err = v3.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("Expected json.SyntaxError, not %T", err)
	}

	v4 := NewUint64()
	err = v4.UnmarshalJSON(uint64MaxValueValidJSON)
	errorPanic(err)
	isUint64Valid(t, v1, "UnmarshalJSON()")

	v5 := NewUint64()
	err = v5.UnmarshalJSON(blankJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	v6 := NewUint64()
	err = json.Unmarshal(uint64StringJSON, &v6)
	errorPanic(err)
	isUint64Valid(t, v6, "UnmarshalJSON()")

	v7 := NewUint64()
	err = json.Unmarshal(uint64BlankJSON, &v7)
	errorPanic(err)
	isUint64Null(t, v7, "UnmarshalJSON()")

	v8 := NewUint64()
	err = json.Unmarshal(boolFalseJSON, &v8)
	if err == nil {
		panic("err should not be nil")
	}
	isUint64Null(t, v8, "UnmarshalJSON()")
}

func TestUint64MarshalJSON(t *testing.T) {
	v1 := NewUint64Value(uint64(math.MaxUint64))
	data, err := v1.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, string(uint64JSON), "non-empty json marshal")

	v2 := NewUint64()
	data, err = v2.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, "null", "null json marshal")
}

func TestUint64UnmarshalText(t *testing.T) {
	var err error

	v1 := NewUint64()
	err = v1.UnmarshalText(uint64JSON)
	errorPanic(err)
	isUint64Valid(t, v1, "UnmarshalText()")

	v2 := NewUint64()
	err = v2.UnmarshalText([]byte(""))
	errorPanic(err)
	if v2.Uint64 != 0 || !v2.Valid {
		t.Errorf("Value should be valid")
	}

	v3 := NewUint64()
	err = v3.UnmarshalText(boolNullJSON)
	errorPanic(err)
	isUint64Null(t, v3, "UnmarshalText()")
}

func TestUint64MarshalText(t *testing.T) {
	v1 := NewUint64Value(uint64(math.MaxUint64))
	data, err := v1.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, string(uint64JSON), "Non-empty text marshal")

	v2 := NewUint64()
	data, err = v2.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "null", "Null text marshal")
}

func TestUint64UnmarshalBinary(t *testing.T) {
	var err error
	var btf []byte

	v1 := Uint64{}
	btf, err = hex.DecodeString(uint64NullInvalidGob)
	errorPanic(err)
	errorPanic(gob.NewDecoder(bytes.NewBuffer(btf)).Decode(&v1))
	isUint64Null(t, v1, "UnmarshalBinary() invalid")

	v2 := Uint64{}
	btf, err = hex.DecodeString(uint64OkValidGob)
	errorPanic(err)
	errorPanic(gob.NewDecoder(bytes.NewBuffer(btf)).Decode(&v2))
	isUint64Valid(t, v2, "UnmarshalBinary() ok")
}

func TestUint64MarshalBinary(t *testing.T) {
	var err error
	var buf *bytes.Buffer
	var enc *gob.Encoder

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)

	v1 := NewUint64()
	err = enc.Encode(&v1)
	errorPanic(err)
	v1h := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(v1h), uint64NullInvalidGob, "NewBytes() -> MarshalBinary()")

	buf.Reset()
	enc = gob.NewEncoder(buf)
	v2 := NewUint64Value(uint64(math.MaxUint64))
	err = enc.Encode(&v2)
	errorPanic(err)
	v2h := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(v2h), uint64OkValidGob, "NewBoolValue(false) -> MarshalBinary()")
}
