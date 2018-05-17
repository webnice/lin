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

func isFloat64Valid(t *testing.T, f Float64, from string) {
	dv, _ := f.Value()
	if dv.(float64) != math.MaxFloat64 {
		t.Errorf("Bad %s float64: \"%f\" ≠ \"%f\"\n", from, dv, float64(math.MaxFloat64))
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func isFloat64Null(t *testing.T, f Float64, from string) {
	if f.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewFloat64(t *testing.T) {
	v1 := NewFloat64()
	isFloat64Null(t, v1, "NewFloat64()")
}

func TestNewFloat64Value(t *testing.T) {
	v1 := NewFloat64Value(math.MaxFloat64)
	isFloat64Valid(t, v1, "NewFloat64Value()")
}

func TestNewFloat64PointerValue(t *testing.T) {
	var bv = float64(math.MaxFloat64)

	v1 := NewFloat64PointerValue(&bv)
	isFloat64Valid(t, v1, "NewFloat64PointerValue()")

	v2 := NewFloat64PointerValue(nil)
	isFloat64Null(t, v2, "NewFloat64PointerValue()")
}

func TestFloat64SetValid(t *testing.T) {
	v1 := NewFloat64()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	v1.SetValid(float64(math.MaxFloat64))
	isFloat64Valid(t, v1, "SetValid()")
}

func TestFloat64Reset(t *testing.T) {
	v1 := NewFloat64()
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

func TestFloat64NullIfDefault(t *testing.T) {
	v1 := NewFloat64Value(float64(math.MaxFloat64))
	isFloat64Valid(t, v1, "NewFloat64Value()")
	v1.NullIfDefault()
	isFloat64Valid(t, v1, "NullIfDefault()")

	v1.SetValid(0)
	if !v1.Valid {
		t.Error("Valid property", "is false, but should be true")
	}
	v1.NullIfDefault()
	isFloat64Null(t, v1, "NullIfDefault()")
}

func TestFloat64MustValue(t *testing.T) {
	var buf interface{}

	v1 := NewFloat64()
	buf = v1.MustValue()
	if _, ok := buf.(float64); !ok {
		t.Error("MustValue()", "is nil, but should be not nil")
	}
	v1.SetValid(float64(math.MaxFloat64))
	bf2 := v1.MustValue()
	if bf2 != math.MaxFloat64 {
		t.Error("MustValue()", "is wrong")
	}
}

func TestFloat64Pointer(t *testing.T) {
	v1 := NewFloat64PointerValue(nil)
	isFloat64Null(t, v1, "NewFloat64PointerValue(nil)")
	if v1.Pointer() != nil {
		t.Error("Pointer()", "is not nil, but should be nil")
	}

	v2 := NewFloat64Value(float64(math.MaxFloat64))
	pb := v2.Pointer()
	if pb == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
	if *pb != float64(math.MaxFloat64) {
		t.Error("Pointer() value", "is wrong")
	}

	v3 := NewFloat64PointerValue(pb)
	isFloat64Valid(t, v3, "NewFloat64PointerValue(max)")
	if v3.Pointer() == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
}

func TestFloat64Scan(t *testing.T) {
	v1 := NewFloat64()
	errorPanic(v1.Scan(float64(math.MaxFloat64)))
	isFloat64Valid(t, v1, "Scan()")

	v2 := NewFloat64()
	errorPanic(v2.Scan("179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000"))
	isFloat64Valid(t, v2, "Scan()")

	v3 := NewFloat64()
	errorPanic(v3.Scan(nil))
	isFloat64Null(t, v3, "Scan()")

	v4 := NewFloat64()
	err := v4.Scan(false)
	if err == nil {
		t.Error("Scan()", "is nil, but should be not nil")
	}
}

func TestFloat64Value(t *testing.T) {
	var err error
	var dv driver.Value
	var ok bool
	var item float64

	v1 := NewFloat64Value(float64(math.MaxFloat64))
	dv, err = v1.Value()
	errorPanic(err)
	if dv == nil {
		t.Error("Value()", "returns nil, but should be not nil")
	}
	if item, ok = dv.(float64); !ok {
		t.Errorf("%s returns type %q, but should be %q", "Value()", reflect.TypeOf(dv).Name(), "float64")
	}
	if item != float64(math.MaxFloat64) {
		t.Error("Value() value", "is wrong")
	}

	v2 := NewFloat64()
	dv, err = v2.Value()
	errorPanic(err)
	if dv != nil {
		t.Error("Value()", "returns not nil, but should be nil")
	}
}

func TestFloat64UnmarshalJSON(t *testing.T) {
	var err error

	v1 := NewFloat64()
	err = json.Unmarshal(float64JSON, &v1)
	errorPanic(err)
	isFloat64Valid(t, v1, "UnmarshalJSON()")

	v2 := NewFloat64()
	err = json.Unmarshal(bytesNullJSON, &v2)
	errorPanic(err)
	isFloat64Null(t, v2, "UnmarshalJSON(null)")

	v3 := NewFloat64()
	err = v3.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("Expected json.SyntaxError, not %T", err)
	}

	v4 := NewFloat64()
	err = v4.UnmarshalJSON(uint64MaxValueValidJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	v5 := NewFloat64()
	err = v5.UnmarshalJSON(blankJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	v6 := NewFloat64()
	err = json.Unmarshal(float64StringJSON, &v6)
	errorPanic(err)
	isFloat64Valid(t, v6, "UnmarshalJSON()")

	v7 := NewFloat64()
	err = json.Unmarshal(float64BlankJSON, &v7)
	errorPanic(err)
	isFloat64Null(t, v7, "UnmarshalJSON()")

	v8 := NewFloat64()
	err = json.Unmarshal(boolFalseJSON, &v8)
	if err == nil {
		panic("err should not be nil")
	}
	isFloat64Null(t, v8, "UnmarshalJSON()")
}

func TestFloat64MarshalJSON(t *testing.T) {
	v1 := NewFloat64Value(float64(math.MaxFloat64))
	data, err := v1.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, string(float64JSON), "non-empty json marshal")

	v2 := NewFloat64()
	data, err = v2.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, "null", "null json marshal")

	v3 := NewFloat64Value(math.NaN())
	data, err = v3.MarshalJSON()
	if err == nil || len(data) > 0 {
		panic("err should not be nil")
	}
}

func TestFloat64UnmarshalText(t *testing.T) {
	var err error

	v1 := NewFloat64()
	err = v1.UnmarshalText(float64JSON)
	errorPanic(err)
	isFloat64Valid(t, v1, "UnmarshalText()")

	v2 := NewFloat64()
	err = v2.UnmarshalText([]byte(""))
	errorPanic(err)
	if v2.Float64 != 0 || !v2.Valid {
		t.Errorf("Value should be valid")
	}

	v3 := NewFloat64()
	err = v3.UnmarshalText(boolNullJSON)
	errorPanic(err)
	isFloat64Null(t, v3, "UnmarshalText()")
}

func TestFloat64MarshalText(t *testing.T) {
	v1 := NewFloat64Value(float64(math.MaxFloat64))
	data, err := v1.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, string(float64JSON), "Non-empty text marshal")

	v2 := NewFloat64()
	data, err = v2.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "null", "Null text marshal")
}

func TestFloat64UnmarshalBinary(t *testing.T) {
	var err error
	var btf []byte

	v1 := Float64{}
	btf, err = hex.DecodeString(floatNullInvalidGob)
	errorPanic(err)
	errorPanic(gob.NewDecoder(bytes.NewBuffer(btf)).Decode(&v1))
	isFloat64Null(t, v1, "UnmarshalBinary() invalid")

	v2 := Float64{}
	btf, err = hex.DecodeString(floatOkValidGob)
	errorPanic(err)
	errorPanic(gob.NewDecoder(bytes.NewBuffer(btf)).Decode(&v2))
	isFloat64Valid(t, v2, "UnmarshalBinary() ok")
}

func TestFloat64MarshalBinary(t *testing.T) {
	var err error
	var buf *bytes.Buffer
	var enc *gob.Encoder

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)

	v1 := NewFloat64()
	err = enc.Encode(&v1)
	errorPanic(err)
	v1h := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(v1h), floatNullInvalidGob, "NewBytes() -> MarshalBinary()")

	buf.Reset()
	enc = gob.NewEncoder(buf)
	v2 := NewFloat64Value(float64(math.MaxFloat64))
	err = enc.Encode(&v2)
	errorPanic(err)
	v2h := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(v2h), floatOkValidGob, "NewBoolValue(false) -> MarshalBinary()")
}
