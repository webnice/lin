package nul // import "gopkg.in/webnice/lin.v1/nl"

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

func isTrueBool(t *testing.T, b Bool, from string) {
	if !b.Bool {
		t.Errorf("Bad %s bool: %t ≠ %t\n", from, b.Bool, true)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func isFalseBool(t *testing.T, b Bool, from string) {
	if b.Bool {
		t.Errorf("Bad %v bool: %v ≠ %v\n", from, b.Bool, false)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func isNullBool(t *testing.T, b Bool, from string) {
	if b.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewBool(t *testing.T) {
	b := NewBool()
	isNullBool(t, b, "NewBool()")
}

func TestNewBoolValue(t *testing.T) {
	fb := NewBoolValue(false)
	isFalseBool(t, fb, "NewBoolValue(false)")
	tb := NewBoolValue(true)
	isTrueBool(t, tb, "NewBoolValue(true)")
}

func TestNewBoolPointerValue(t *testing.T) {
	var bv bool

	fb := NewBoolPointerValue(&bv)
	isFalseBool(t, fb, "NewBoolPointerValue(*false)")

	bv = true
	tb := NewBoolPointerValue(&bv)
	isTrueBool(t, tb, "NewBoolPointerValue(*true)")

	nb := NewBoolPointerValue(nil)
	isNullBool(t, nb, "NewBoolPointerValue(nil)")
}

func TestBoolSetValid(t *testing.T) {
	tb := NewBool()
	tb.SetValid(false)
	isFalseBool(t, tb, "SetValid(false)")

	fb := NewBool()
	fb.SetValid(true)
	isTrueBool(t, fb, "SetValid(true)")
}

func TestBoolReset(t *testing.T) {
	ib := NewBool()
	isNullBool(t, ib, "NewBool()")

	ib.SetValid(false)
	isFalseBool(t, ib, "SetValid(false)")
	if !ib.Valid {
		t.Error("Valid property", "is false, but should be true")
	}

	ib.Reset()
	isNullBool(t, ib, "Invalidate()")
	if ib.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
}

func TestBoolNullIfDefault(t *testing.T) {
	v1 := NewBoolValue(true)
	isTrueBool(t, v1, "NewBoolValue(true)")
	v1.NullIfDefault()
	isTrueBool(t, v1, "NullIfDefault()")

	v1.SetValid(false)
	if !v1.Valid {
		t.Error("Valid property", "is false, but should be true")
	}
	v1.NullIfDefault()
	isNullBool(t, v1, "NullIfDefault()")
}

func TestBoolMustValue(t *testing.T) {
	mb := NewBool()
	isNullBool(t, mb, "NewBool()")
	if mb.MustValue() {
		t.Error("MustValue()", "is true, but should be false")
	}

	mb.SetValid(false)
	if mb.MustValue() {
		t.Error("MustValue()", "is true, but should be false")
	}

	mb.SetValid(true)
	if !mb.MustValue() {
		t.Error("MustValue()", "is false, but should be true")
	}
}

func TestBoolPointer(t *testing.T) {
	nb := NewBoolPointerValue(nil)
	isNullBool(t, nb, "NewBoolPointerValue(nil)")
	if nb.Pointer() != nil {
		t.Error("Pointer()", "is not nil, but should be nil")
	}

	fb := NewBoolValue(false)
	pb := fb.Pointer()
	if pb == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
	if *pb {
		t.Error("Pointer() value", "is true, but should be false")
	}

	tb := NewBoolValue(true)
	pb = tb.Pointer()
	if pb == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
	if !*pb {
		t.Error("Pointer() value", "is false, but should be true")
	}
}

func TestBoolScan(t *testing.T) {
	fb := NewBoolValue(false)
	errorPanic(fb.Scan(true))
	isTrueBool(t, fb, "Scan(true)")

	nb := NewBoolValue(false)
	errorPanic(nb.Scan(nil))
	isNullBool(t, nb, "Scan(nil)")

	eb := NewBoolValue(true)
	if err := eb.Scan(-1); err == nil {
		t.Error("Scan()", "is returns nil, but should be not nil")
	}
	isNullBool(t, eb, "Scan(nil)")
}

func TestBoolValue(t *testing.T) {
	var err error
	var dv driver.Value
	var item, ok bool

	b := NewBoolValue(true)
	dv, err = b.Value()
	errorPanic(err)
	if dv == nil {
		t.Error("Value()", "returns nil, but should be not nil")
	}
	if item, ok = dv.(bool); !ok {
		t.Errorf("%s returns type %q, but should be %q", "Value()", reflect.TypeOf(dv).Name(), "bool")
	}
	if !item {
		t.Error("Value() value", "is false, but should be true")
	}

	nb := NewBool()
	dv, err = nb.Value()
	errorPanic(err)
	if dv != nil {
		t.Error("Value()", "returns not nil, but should be nil")
	}
}

func TestBoolUnmarshalJSON(t *testing.T) {
	var err error
	var b, fb, nb, badType, valid, invalid Bool

	b = NewBool()
	err = json.Unmarshal(boolTrueJSON, &b)
	errorPanic(err)
	isTrueBool(t, b, "UnmarshalJSON(true)")

	fb = NewBool()
	err = json.Unmarshal(boolFalseJSON, &fb)
	errorPanic(err)
	isFalseBool(t, fb, "UnmarshalJSON(false)")

	nb = NewBool()
	err = json.Unmarshal(boolNullJSON, &nb)
	errorPanic(err)
	isNullBool(t, nb, "UnmarshalJSON(null)")

	valid = NewBool()
	err = json.Unmarshal(boolFalseValidJSON, &valid)
	errorPanic(err)

	invalid = NewBool()
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("Expected json.SyntaxError, not %T", err)
	}

	badType = NewBool()
	err = badType.UnmarshalJSON(uint64MaxValueValidJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	badType = NewBool()
	err = badType.UnmarshalJSON(blankJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestBoolMarshalJSON(t *testing.T) {
	b := NewBoolValue(true)
	data, err := b.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, "true", "non-empty json marshal")

	empty := NewBoolValue(false)
	data, err = empty.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, "false", "empty json marshal")

	null := NewBool()
	data, err = null.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, "null", "null json marshal")

}

func TestBoolUnmarshalText(t *testing.T) {
	var err error
	var b, zero, blank, null, invalid Bool

	b = NewBool()
	err = b.UnmarshalText([]byte("true"))
	errorPanic(err)
	isTrueBool(t, b, "UnmarshalText(`true`) bool")

	zero = NewBool()
	err = zero.UnmarshalText([]byte("false"))
	errorPanic(err)
	isFalseBool(t, zero, "UnmarshalText(`false`) bool")

	blank = NewBool()
	err = blank.UnmarshalText([]byte(""))
	errorPanic(err)
	isNullBool(t, blank, "UnmarshalText(``) bool")

	null = NewBool()
	err = null.UnmarshalText([]byte("null"))
	errorPanic(err)
	isNullBool(t, null, "UnmarshalText(`null`) bool")

	invalid = NewBool()
	err = invalid.UnmarshalText([]byte(":D"))
	if err == nil {
		panic("Error should not be nil")
	}
	isNullBool(t, invalid, "Invalid json")
}

func TestBoolMarshalText(t *testing.T) {
	b := NewBoolValue(true)
	data, err := b.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "true", "Non-empty text marshal")

	zero := NewBoolValue(false)
	data, err = zero.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "false", "Zero text marshal")

	null := NewBool()
	data, err = null.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "null", "Null text marshal")
}

func TestBoolUnmarshalBinary(t *testing.T) {
	var err error
	var buf *bytes.Buffer
	var dec *gob.Decoder
	var btf []byte
	var trueBool, falseBool, nullBool Bool

	btf, err = hex.DecodeString(boolTrueValidGob)
	errorPanic(err)
	buf = bytes.NewBuffer(btf)
	dec = gob.NewDecoder(buf)
	err = dec.Decode(&trueBool)
	errorPanic(err)
	isTrueBool(t, trueBool, "UnmarshalBinary() -> Bool(true)")

	btf, err = hex.DecodeString(boolFalseValidGob)
	errorPanic(err)
	buf = bytes.NewBuffer(btf)
	dec = gob.NewDecoder(buf)
	err = dec.Decode(&falseBool)
	errorPanic(err)
	isFalseBool(t, falseBool, "UnmarshalBinary() -> Bool(false)")

	btf, err = hex.DecodeString(boolNullInvalidGob)
	errorPanic(err)
	buf = bytes.NewBuffer(btf)
	dec = gob.NewDecoder(buf)
	err = dec.Decode(&nullBool)
	errorPanic(err)
	isNullBool(t, nullBool, "UnmarshalBinary() -> Bool(null)")
}

func TestBoolMarshalBinary(t *testing.T) {
	var err error
	var buf *bytes.Buffer
	var enc *gob.Encoder
	var trueBool, falseBool, nullBool Bool

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)

	trueBool = NewBoolValue(true)
	err = enc.Encode(&trueBool)
	errorPanic(err)
	tbh := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(tbh), boolTrueValidGob, "NewBoolValue(true) -> MarshalBinary()")

	buf.Reset()
	enc = gob.NewEncoder(buf)
	falseBool = NewBoolValue(false)
	err = enc.Encode(&falseBool)
	errorPanic(err)
	fbh := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(fbh), boolFalseValidGob, "NewBoolValue(false) -> MarshalBinary()")

	buf.Reset()
	enc = gob.NewEncoder(buf)
	nullBool = NewBool()
	err = enc.Encode(&nullBool)
	errorPanic(err)
	nbh := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(nbh), boolNullInvalidGob, "NewBool() -> MarshalBinary()")
}
