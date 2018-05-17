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
	"time"
)

func isTimeValid(t *testing.T, tm Time, from string) {
	dv, _ := tm.Value()
	if dv.(time.Time) != timeOkValidValue {
		t.Errorf("Bad %s time.Time: \"%s\" ≠ \"%s\"\n", from, dv, timeOkValidValue)
	}
	if !tm.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func isTimeNull(t *testing.T, tm Time, from string) {
	if tm.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewTime(t *testing.T) {
	v1 := NewTime()
	isTimeNull(t, v1, "NewTime()")
}

func TestNewTimeValue(t *testing.T) {
	v1 := NewTimeValue(timeOkValidValue)
	isTimeValid(t, v1, "NewTimeValue()")
}

func TestNewTimePointerValue(t *testing.T) {
	var bv = timeOkValidValue

	v1 := NewTimePointerValue(&bv)
	isTimeValid(t, v1, "NewTimePointerValue()")

	v2 := NewTimePointerValue(nil)
	isTimeNull(t, v2, "NewTimePointerValue()")
}

func TestTimeSetValid(t *testing.T) {
	v1 := NewTime()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	v1.SetValid(timeOkValidValue)
	isTimeValid(t, v1, "SetValid()")
}

func TestTimeReset(t *testing.T) {
	v1 := NewTime()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	v1.SetValid(time.Time{})
	if !v1.Valid {
		t.Error("Valid property", "is false, but should be true")
	}
	v1.Reset()
	if v1.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
}

func TestTimeMustValue(t *testing.T) {
	var buf interface{}

	v1 := NewTime()
	buf = v1.MustValue()
	if _, ok := buf.(time.Time); !ok {
		t.Error("MustValue()", "is nil, but should be not nil")
	}
	v1.SetValid(timeOkValidValue)
	bf2 := v1.MustValue()
	if bf2 != timeOkValidValue {
		t.Error("MustValue()", "is wrong")
	}
}

func TestTimePointer(t *testing.T) {
	v1 := NewTimePointerValue(nil)
	isTimeNull(t, v1, "NewTimePointerValue(nil)")
	if v1.Pointer() != nil {
		t.Error("Pointer()", "is not nil, but should be nil")
	}

	v2 := NewTimeValue(timeOkValidValue)
	pb := v2.Pointer()
	if pb == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
	if *pb != timeOkValidValue {
		t.Error("Pointer() value", "is wrong")
	}

	v3 := NewTimePointerValue(pb)
	isTimeValid(t, v3, "NewTimePointerValue()")
	if v3.Pointer() == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
}

func TestTimeScan(t *testing.T) {
	v1 := NewTime()
	errorPanic(v1.Scan(timeOkValidValue))
	isTimeValid(t, v1, "Scan()")

	v2 := NewTime()
	errorPanic(v2.Scan(timeStringValue))
	isTimeValid(t, v2, "Scan()")

	v3 := NewTime()
	errorPanic(v3.Scan(nil))
	isTimeNull(t, v3, "Scan()")

	v4 := NewTime()
	err := v4.Scan(false)
	if err == nil {
		t.Error("Scan()", "is nil, but should be not nil")
	}
}

func TestTimeValue(t *testing.T) {
	var err error
	var dv driver.Value
	var ok bool
	var item time.Time

	v1 := NewTimeValue(timeOkValidValue)
	dv, err = v1.Value()
	errorPanic(err)
	if dv == nil {
		t.Error("Value()", "returns nil, but should be not nil")
	}
	if item, ok = dv.(time.Time); !ok {
		t.Errorf("%s returns type %q, but should be %q", "Value()", reflect.TypeOf(dv).Name(), "time.Time")
	}
	if item != timeOkValidValue {
		t.Error("Value() value", "is wrong")
	}

	v2 := NewTime()
	dv, err = v2.Value()
	errorPanic(err)
	if dv != nil {
		t.Error("Value()", "returns not nil, but should be nil")
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	var err error

	v1 := NewTime()
	err = json.Unmarshal(timeStringValueJSON, &v1)
	errorPanic(err)
	isTimeValid(t, v1, "UnmarshalJSON()")

	v2 := NewTime()
	err = json.Unmarshal(bytesNullJSON, &v2)
	errorPanic(err)
	isTimeNull(t, v2, "UnmarshalJSON(null)")

	v3 := NewTime()
	err = v3.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("Expected json.SyntaxError, not %T", err)
	}

	v4 := NewTime()
	err = v4.UnmarshalJSON(uint64MaxValueValidJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	v5 := NewTime()
	err = v5.UnmarshalJSON(blankJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	v6 := NewTime()
	err = json.Unmarshal(boolFalseJSON, &v6)
	if err == nil {
		panic("err should not be nil")
	}
	isTimeNull(t, v6, "UnmarshalJSON()")

	v7 := NewTime()
	err = json.Unmarshal(timeTestValidJSON, &v7)
	errorPanic(err)
	isTimeValid(t, v7, "UnmarshalJSON()")
}

func TestTimeMarshalJSON(t *testing.T) {
	v1 := NewTimeValue(timeOkValidValue)
	data, err := v1.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, string(timeStringValueJSON), "non-empty json marshal")

	v2 := NewTime()
	data, err = v2.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, "null", "null json marshal")
}

func TestTimeUnmarshalText(t *testing.T) {
	var err error

	v1 := NewTime()
	err = v1.UnmarshalText([]byte(timeStringValue))
	errorPanic(err)
	isTimeValid(t, v1, "UnmarshalText()")

	v2 := NewTime()
	err = v2.UnmarshalText([]byte(""))
	errorPanic(err)
	if !v2.Time.IsZero() || !v2.Valid {
		t.Errorf("Value should be valid")
	}

	v3 := NewTime()
	err = v3.UnmarshalText(boolNullJSON)
	errorPanic(err)
	isTimeNull(t, v3, "UnmarshalText()")
}

func TestTimeMarshalText(t *testing.T) {
	v1 := NewTimeValue(timeOkValidValue)
	data, err := v1.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, timeStringValue, "Non-empty text marshal")

	v2 := NewTime()
	data, err = v2.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "null", "Null text marshal")

	v3 := NewTimeValue(time.Time{})
	data, err = v3.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "", "Null text marshal")
}

func TestTimeUnmarshalBinary(t *testing.T) {
	var err error
	var btf []byte

	v1 := Time{}
	btf, err = hex.DecodeString(timeNullInvalidGob)
	errorPanic(err)
	errorPanic(gob.NewDecoder(bytes.NewBuffer(btf)).Decode(&v1))
	isTimeNull(t, v1, "UnmarshalBinary() invalid")

	v2 := Time{}
	btf, err = hex.DecodeString(timeOkValidGob)
	errorPanic(err)
	errorPanic(gob.NewDecoder(bytes.NewBuffer(btf)).Decode(&v2))
	isTimeValid(t, v2, "UnmarshalBinary() ok")
}

func TestTimeMarshalBinary(t *testing.T) {
	var err error
	var buf *bytes.Buffer
	var enc *gob.Encoder

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)

	v1 := NewTime()
	err = enc.Encode(&v1)
	errorPanic(err)
	v1h := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(v1h), timeNullInvalidGob, "NewBytes() -> MarshalBinary()")

	buf.Reset()
	enc = gob.NewEncoder(buf)
	v2 := NewTimeValue(timeOkValidValue)
	err = enc.Encode(&v2)
	errorPanic(err)
	v2h := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(v2h), timeOkValidGob, "NewBoolValue(false) -> MarshalBinary()")
}
