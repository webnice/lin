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

const bytesTestString = `Test data 1pHuOxADZkeh8Y9WvL75`

func isBytesValid(t *testing.T, bt Bytes, from string) {
	if !bytes.Equal(bt.Bytes.Bytes(), []byte(bytesTestString)) {
		t.Errorf("Bad %s bytes: %q ≠ %q\n", from, bt.Bytes.String(), bytesTestString)
	}
	if !bt.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func isNullBytes(t *testing.T, bt Bytes, from string) {
	if !bytes.Equal(bt.Bytes.Bytes(), []byte{}) {
		t.Errorf("Bad %s bytes: %q ≠ %q\n", from, bt.Bytes.String(), "")
	}
	if bt.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewBytes(t *testing.T) {
	bt := NewBytes()
	isNullBytes(t, bt, "NewBytes()")
}

func TestNewBytesValue(t *testing.T) {
	bt := NewBytesValue([]byte(bytesTestString))
	isBytesValid(t, bt, "NewBytesValue()")
}

func TestNewBytesPointerValue(t *testing.T) {
	var bv = []byte(bytesTestString)

	bt := NewBytesPointerValue(&bv)
	isBytesValid(t, bt, "NewBytesPointerValue()")

	bt = NewBytesPointerValue(nil)
	isNullBytes(t, bt, "NewBytesPointerValue()")
}

func TestBytesSetValid(t *testing.T) {
	bt := NewBytes()
	if bt.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	bt.SetValid([]byte(bytesTestString))
	isBytesValid(t, bt, "SetValid()")
}

func TestBytesInvalidate(t *testing.T) {
	bt := NewBytes()
	if bt.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
	bt.SetValid([]byte("Test"))
	if !bt.Valid {
		t.Error("Valid property", "is false, but should be true")
	}
	bt.Invalidate()
	if bt.Valid {
		t.Error("Valid property", "is true, but should be false")
	}
}

func TestBytesMustValue(t *testing.T) {
	var buf interface{}

	bt := NewBytes()
	buf = bt.MustValue()
	if _, ok := buf.([]byte); !ok {
		t.Error("MustValue()", "is nil, but should be not nil")
	}
	bt.SetValid([]byte(bytesTestString))
	bf2 := bt.MustValue()
	if !bytes.Equal(bf2, []byte(bytesTestString)) {
		t.Error("MustValue()", "is wrong")
	}
}

func TestBytesPointer(t *testing.T) {
	nb := NewBytesPointerValue(nil)
	isNullBytes(t, nb, "NewBytesPointerValue(nil)")
	if nb.Pointer() != nil {
		t.Error("Pointer()", "is not nil, but should be nil")
	}

	fb := NewBytesValue([]byte(bytesTestString))
	pb := fb.Pointer()
	if pb == nil {
		t.Error("Pointer()", "is nil, but should be not nil")
	}
	if !bytes.Equal(*pb, []byte(bytesTestString)) {
		t.Error("Pointer() value", "is wrong")
	}
}

func TestBytesScan(t *testing.T) {
	bt := NewBytes()
	errorPanic(bt.Scan([]byte(bytesTestString)))
	isBytesValid(t, bt, "Scan()")

	bs := NewBytes()
	errorPanic(bs.Scan(bytesTestString))
	isBytesValid(t, bt, "Scan()")

	bn := NewBytes()
	errorPanic(bn.Scan(nil))
	isNullBytes(t, bn, "Scan()")

	be := NewBytes()
	err := be.Scan(false)
	if err == nil {
		t.Error("Scan()", "is nil, but should be not nil")
	}
}

func TestBytesValue(t *testing.T) {
	var err error
	var dv driver.Value
	var ok bool
	var item []byte

	b := NewBytesValue([]byte(bytesTestString))
	dv, err = b.Value()
	errorPanic(err)
	if dv == nil {
		t.Error("Value()", "returns nil, but should be not nil")
	}
	if item, ok = dv.([]byte); !ok {
		t.Errorf("%s returns type %q, but should be %q", "Value()", reflect.TypeOf(dv).Name(), "[]byte")
	}
	if !bytes.Equal(item, []byte(bytesTestString)) {
		t.Error("Value() value", "is wrong")
	}

	nb := NewBytes()
	dv, err = nb.Value()
	errorPanic(err)
	if dv != nil {
		t.Error("Value()", "returns not nil, but should be nil")
	}
}

func TestBytesUnmarshalJSON(t *testing.T) {
	var err error

	b := NewBytes()
	err = json.Unmarshal(bytesTestValidJSON, &b)
	errorPanic(err)
	isBytesValid(t, b, "UnmarshalJSON()")

	nb := NewBytes()
	err = json.Unmarshal(bytesNullJSON, &nb)
	errorPanic(err)
	isNullBytes(t, nb, "UnmarshalJSON(null)")

	invalid := NewBytes()
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("Expected json.SyntaxError, not %T", err)
	}

	badType := NewBytes()
	err = badType.UnmarshalJSON(uint64MaxValueValidJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	badType = NewBytes()
	err = badType.UnmarshalJSON(blankJSON)
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestBytesMarshalJSON(t *testing.T) {
	bv := NewBytesValue([]byte(bytesTestString))
	data, err := bv.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, `"VGVzdCBkYXRhIDFwSHVPeEFEWmtlaDhZOVd2TDc1"`, "non-empty json marshal")

	null := NewBytes()
	data, err = null.MarshalJSON()
	errorPanic(err)
	jsonEquals(t, data, "null", "null json marshal")
}

func TestBytesUnmarshalText(t *testing.T) {
	var err error

	bt := NewBytes()
	err = bt.UnmarshalText(bytesTestTextBase64)
	errorPanic(err)
	isBytesValid(t, bt, "UnmarshalText()")

	null := NewBytes()
	err = null.UnmarshalText([]byte("null"))
	errorPanic(err)
	isNullBytes(t, null, "UnmarshalText(`null`)")

	zero := NewBytes()
	err = zero.UnmarshalText([]byte(""))
	errorPanic(err)
	if !bytes.Equal(zero.MustValue(), []byte{}) || !zero.Valid {
		t.Errorf("Value should be valid")
	}
}

func TestBytesMarshalText(t *testing.T) {
	bv := NewBytesValue([]byte(bytesTestString))
	data, err := bv.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "VGVzdCBkYXRhIDFwSHVPeEFEWmtlaDhZOVd2TDc1", "Non-empty text marshal")

	zero := NewBytesValue(nil)
	data, err = zero.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "", "Zero text marshal")

	null := NewBytes()
	data, err = null.MarshalText()
	errorPanic(err)
	jsonEquals(t, data, "null", "Null text marshal")
}

func TestBytesUnmarshalBinary(t *testing.T) {
	var err error
	var buf *bytes.Buffer
	var dec *gob.Decoder
	var btf []byte
	var invalidBytes, okBytes, zeroBytes Bytes

	btf, err = hex.DecodeString(bytesNullInvalidGob)
	errorPanic(err)
	buf = bytes.NewBuffer(btf)
	dec = gob.NewDecoder(buf)
	err = dec.Decode(&invalidBytes)
	errorPanic(err)
	isNullBytes(t, invalidBytes, "UnmarshalBinary() -> Bytes(invalid)")

	btf, err = hex.DecodeString(bytesOkValidGob)
	errorPanic(err)
	buf = bytes.NewBuffer(btf)
	dec = gob.NewDecoder(buf)
	err = dec.Decode(&okBytes)
	errorPanic(err)
	isBytesValid(t, okBytes, "UnmarshalBinary() -> Bool(ok)")

	btf, err = hex.DecodeString(bytesZeroValidGob)
	errorPanic(err)
	buf = bytes.NewBuffer(btf)
	dec = gob.NewDecoder(buf)
	err = dec.Decode(&zeroBytes)
	errorPanic(err)
	if !bytes.Equal(zeroBytes.MustValue(), []byte{}) || !zeroBytes.Valid {
		t.Errorf("Value should be valid")
	}
}

func TestBytesMarshalBinary(t *testing.T) {
	var err error
	var buf *bytes.Buffer
	var enc *gob.Encoder

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)

	bt := NewBytes()
	err = enc.Encode(&bt)
	errorPanic(err)
	bth := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(bth), bytesNullInvalidGob, "NewBytes() -> MarshalBinary()")

	buf.Reset()
	enc = gob.NewEncoder(buf)
	btok := NewBytesValue([]byte(bytesTestString))
	err = enc.Encode(&btok)
	errorPanic(err)
	bto := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(bto), bytesOkValidGob, "NewBytesValue(max) -> MarshalBinary()")

	buf.Reset()
	enc = gob.NewEncoder(buf)
	zero := NewBytesValue(nil)
	err = enc.Encode(&zero)
	errorPanic(err)
	btz := hex.EncodeToString(buf.Bytes())
	jsonEquals(t, []byte(btz), bytesZeroValidGob, "NewBytesValue(nil) -> MarshalBinary()")
}
