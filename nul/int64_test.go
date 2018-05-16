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

func isInt64Valid(t *testing.T, i Int64, from string) {
	dv, _ := i.Value()
	if dv.(int64) != math.MaxInt64 {
		t.Errorf("Bad %s int64: \"%d\" ≠ \"%d\"\n", from, dv, int64(math.MaxInt64))
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func isInt64Null(t *testing.T, i Int64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewInt64(t *testing.T) {
	v1 := NewInt64()
	isInt64Null(t, v1, "NewInt64()")
}

func TestNewInt64Value(t *testing.T) {
	v1 := NewInt64Value(math.MaxInt64)
	isInt64Valid(t, v1, "NewInt64Value()")
}

func TestNewInt64PointerValue(t *testing.T) {
	var bv = int64(math.MaxInt64)

	v1 := NewInt64PointerValue(&bv)
	isInt64Valid(t, v1, "NewInt64PointerValue()")

	v2 := NewInt64PointerValue(nil)
	isInt64Null(t, v2, "NewInt64PointerValue()")
}

func TestInt64SetValid(t *testing.T) {
	v1 := NewInt64()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	v1.SetValid(int64(math.MaxInt64))
	isInt64Valid(t, v1, "SetValid()")
}

func TestInt64Invalidate(t *testing.T) {
	v1 := NewInt64()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	v1.SetValid(0)
	if !v1.Valid {
		t.Error("Valid property", "is false, but should be true")
	}
	v1.Invalidate()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
}

func TestInt64MustValue(t *testing.T) {
	var buf interface{}

	v1 := NewInt64()
	buf = v1.MustValue()
	if _, ok := buf.(int64); !ok {
		t.Error("MustValue()", "is nil, but should be not nil")
	}
	v1.SetValid(int64(math.MaxInt64))
	bf2 := v1.MustValue()
	if bf2 != math.MaxInt64 {
		t.Error("MustValue()", "is wrong")
	}
}

func TestInt64Pointer(t *testing.T) {
	v1 := NewInt64PointerValue(nil)
	isInt64Null(t, v1, "NewInt64PointerValue(nil)")
	if v1.Pointer() != nil {
		t.Error("Pointer()", "is not nil, but should be nil")
	}

	v2 := NewInt64Value(int64(math.MaxInt64))
	pb := v2.Pointer()
	if pb == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
	if *pb != int64(math.MaxInt64) {
		t.Error("Pointer() value", "is wrong")
	}

	v3 := NewInt64PointerValue(pb)
	isInt64Valid(t, v3, "NewInt64PointerValue(max)")
	if v3.Pointer() == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
}

func TestInt64Scan(t *testing.T) {
	v1 := NewInt64()
	errorPanic(v1.Scan(int64(math.MaxInt64)))
	isInt64Valid(t, v1, "Scan()")

	v2 := NewInt64()
	errorPanic(v2.Scan("9223372036854775807"))
	isInt64Valid(t, v2, "Scan()")

	v3 := NewInt64()
	errorPanic(v3.Scan(nil))
	isInt64Null(t, v3, "Scan()")

	v4 := NewInt64()
	err := v4.Scan(false)
	if err == nil {
		t.Error("Scan()", "is nil, but should be not nil")
	}
}

func TestInt64Value(t *testing.T) {
	var err error
	var dv driver.Value
	var ok bool
	var item int64

	v1 := NewInt64Value(int64(math.MaxInt64))
	dv, err = v1.Value()
	errorPanic(err)
	if dv == nil {
		t.Error("Value()", "returns nil, but should be not nil")
	}
	if item, ok = dv.(int64); !ok {
		t.Errorf("%s returns type %q, but should be %q", "Value()", reflect.TypeOf(dv).Name(), "int64")
	}
	if item != int64(math.MaxInt64) {
		t.Error("Value() value", "is wrong")
	}

	v2 := NewInt64()
	dv, err = v2.Value()
	errorPanic(err)
	if dv != nil {
		t.Error("Value()", "returns not nil, but should be nil")
	}
}

func TestInt64UnmarshalJSON(t *testing.T) {
	var err error

	v1 := NewInt64()
	err = json.Unmarshal(int64JSON, &v1)
	errorPanic(err)
	isInt64Valid(t, v1, "UnmarshalJSON()")

	v2 := NewInt64()
	err = json.Unmarshal(bytesNullJSON, &v2)
	errorPanic(err)
	isInt64Null(t, v2, "UnmarshalJSON(null)")

	v3 := NewInt64()
	err = v3.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("Expected json.SyntaxError, not %T", err)
	}

	v4 := NewInt64()
	err = v4.UnmarshalJSON(uint64MaxValueValidJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	v5 := NewInt64()
	err = v5.UnmarshalJSON(blankJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	v6 := NewInt64()
	err = json.Unmarshal(int64StringJSON, &v6)
	errorPanic(err)
	isInt64Valid(t, v6, "UnmarshalJSON()")

	v7 := NewInt64()
	err = json.Unmarshal(int64BlankJSON, &v7)
	errorPanic(err)
	isInt64Null(t, v7, "UnmarshalJSON()")

	v8 := NewInt64()
	err = json.Unmarshal(boolFalseJSON, &v8)
	if err == nil {
		panic("err should not be nil")
	}
	isInt64Null(t, v8, "UnmarshalJSON()")
}

func TestInt64MarshalJSON(t *testing.T) {
	v1 := NewInt64Value(int64(math.MaxInt64))
	data, err := v1.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, "9223372036854775807", "non-empty json marshal")

	v2 := NewInt64()
	data, err = v2.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, "null", "null json marshal")
}

func TestInt64UnmarshalText(t *testing.T) {
	var err error

	v1 := NewInt64()
	err = v1.UnmarshalText(int64JSON)
	errorPanic(err)
	isInt64Valid(t, v1, "UnmarshalText()")

	v2 := NewInt64()
	err = v2.UnmarshalText([]byte(""))
	errorPanic(err)
	if v2.Int64 != 0 || !v2.Valid {
		t.Errorf("Value should be valid")
	}

	v3 := NewInt64()
	err = v3.UnmarshalText(boolNullJSON)
	errorPanic(err)
	isInt64Null(t, v3, "UnmarshalText()")
}

func TestInt64MarshalText(t *testing.T) {
	v1 := NewInt64Value(int64(math.MaxInt64))
	data, err := v1.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, string(int64JSON), "Non-empty text marshal")

	v2 := NewInt64()
	data, err = v2.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "null", "Null text marshal")
}

func TestInt64UnmarshalBinary(t *testing.T) {
	var err error
	var btf []byte

	v1 := Int64{}
	btf, err = hex.DecodeString(intNullInvalidGob)
	errorPanic(err)
	errorPanic(gob.NewDecoder(bytes.NewBuffer(btf)).Decode(&v1))
	isInt64Null(t, v1, "UnmarshalBinary() invalid")

	v2 := Int64{}
	btf, err = hex.DecodeString(intOkValidGob)
	errorPanic(err)
	errorPanic(gob.NewDecoder(bytes.NewBuffer(btf)).Decode(&v2))
	isInt64Valid(t, v2, "UnmarshalBinary() ok")
}

func TestInt64MarshalBinary(t *testing.T) {
	var err error
	var buf *bytes.Buffer
	var enc *gob.Encoder

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)

	v1 := NewInt64()
	err = enc.Encode(&v1)
	errorPanic(err)
	v1h := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(v1h), intNullInvalidGob, "NewBytes() -> MarshalBinary()")

	buf.Reset()
	enc = gob.NewEncoder(buf)
	v2 := NewInt64Value(int64(math.MaxInt64))
	err = enc.Encode(&v2)
	errorPanic(err)
	v2h := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(v2h), intOkValidGob, "NewBoolValue(false) -> MarshalBinary()")
}
