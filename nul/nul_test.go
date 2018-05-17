package nul // import "gopkg.in/webnice/nul.v1/nul"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"
)

var (
	invalidJSON             = []byte(`:{;-}`)
	blankJSON               = []byte(`{}`)
	boolTrueJSON            = []byte(`true`)
	boolFalseJSON           = []byte(`false`)
	boolNullJSON            = []byte(`null`)
	boolFalseValidJSON      = []byte(`{"Bool":false,"Valid":true}`)
	boolTrueValidGob        = `0aff81060102ff840000003aff8200362dff850301010b426f6f6c5772617070657201ff86000102010556616c7565010200010556616c6964010200000007ff860101010100`
	boolFalseValidGob       = `0aff81060102ff8400000038ff8200342dff850301010b426f6f6c5772617070657201ff86000102010556616c7565010200010556616c6964010200000005ff86020100`
	boolNullInvalidGob      = `0aff81060102ff8400000036ff8200322dff850301010b426f6f6c5772617070657201ff86000102010556616c7565010200010556616c6964010200000003ff8600`
	bytesTestValidJSON      = []byte(`{"Bytes":"VGVzdCBkYXRhIDFwSHVPeEFEWmtlaDhZOVd2TDc1","Valid":true}`)
	bytesTestTextBase64     = []byte(`VGVzdCBkYXRhIDFwSHVPeEFEWmtlaDhZOVd2TDc1`)
	bytesNullInvalidGob     = `0aff87060102ff8a00000012ff8b0301010642756666657201ff8c00000037ff8800332eff8d0301010c42797465735772617070657201ff8e000102010556616c7565010a00010556616c6964010200000003ff8e00`
	bytesOkValidGob         = `0aff87060102ff8a00000012ff8b0301010642756666657201ff8c00000059ff8800552eff8d0301010c42797465735772617070657201ff8e000102010556616c7565010a00010556616c6964010200000025ff8e011e54657374206461746120317048754f7841445a6b656838593957764c3735010100`
	bytesZeroValidGob       = `0aff87060102ff8a00000012ff8b0301010642756666657201ff8c00000039ff8800352eff8d0301010c42797465735772617070657201ff8e000102010556616c7565010a00010556616c6964010200000005ff8e020100`
	bytesNullJSON           = []byte(`null`)
	float64JSON             = []byte(`179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000`)
	float64StringJSON       = []byte(`"179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000"`)
	float64BlankJSON        = []byte(`""`)
	floatNullInvalidGob     = `0aff8f060102ff9200000039ff90003530ff930301010e466c6f617436345772617070657201ff94000102010556616c7565010800010556616c6964010200000003ff9400`
	floatOkValidGob         = `0aff8f060102ff9200000045ff90004130ff930301010e466c6f617436345772617070657201ff94000102010556616c7565010800010556616c696401020000000fff9401f8ffffffffffffef7f010100`
	int64JSON               = []byte(`9223372036854775807`)
	int64StringJSON         = []byte(`"9223372036854775807"`)
	int64BlankJSON          = []byte(`""`)
	intNullInvalidGob       = `0aff95060102ff9800000037ff9600332eff990301010c496e7436345772617070657201ff9a000102010556616c7565010400010556616c6964010200000003ff9a00`
	intOkValidGob           = `0aff95060102ff9800000043ff96003f2eff990301010c496e7436345772617070657201ff9a000102010556616c7565010400010556616c696401020000000fff9a01f8fffffffffffffffe010100`
	stringTestBody          = `3LbOVMltCjj1Mg6sSRYLzS5j64DDNEVax29ypIGxwEx9mnbFnT9FY0sZqP11`
	stringJSON              = []byte(`"3LbOVMltCjj1Mg6sSRYLzS5j64DDNEVax29ypIGxwEx9mnbFnT9FY0sZqP11"`)
	stringNullInvalidGob    = `0aff9b060102ff9e00000038ff9c00342fff9f0301010d537472696e675772617070657201ffa0000102010556616c7565010c00010556616c6964010200000003ffa000`
	stringOkValidGob        = `0aff9b060102ff9e00000078ff9c00742fff9f0301010d537472696e675772617070657201ffa0000102010556616c7565010c00010556616c6964010200000043ffa0013c334c624f564d6c74436a6a314d6736735352594c7a53356a363444444e4556617832397970494778774578396d6e62466e5439465930735a71503131010100`
	timeStringValue         = `2018-05-17T17:17:17.171717+03:00`
	timeStringValueJSON     = []byte(`"` + timeStringValue + `"`)
	timeOkValidValue, _     = time.Parse(time.RFC3339, timeStringValue)
	timeTestValidJSON       = []byte(`{"Time":"` + timeStringValue + `","Valid":true}`)
	timeNullInvalidGob      = `0affa1060102ffa400000010ffa50501010454696d6501ffa600000048ffa200442effa70301010b54696d655772617070657201ffa8000102010556616c756501ffa600010556616c6964010200000010ffa50501010454696d6501ffa600000003ffa800`
	timeOkValidGob          = `0affa1060102ffa400000010ffa50501010454696d6501ffa60000005bffa200572effa70301010b54696d655772617070657201ffa8000102010556616c756501ffa600010556616c6964010200000010ffa50501010454696d6501ffa600000016ffa8010f010000000ed28f85ed0a3c318800b4010100`
	uint64String            = fmt.Sprintf("%d", uint64(math.MaxUint64))
	uint64JSON              = []byte(uint64String)
	uint64StringJSON        = []byte(`"` + uint64String + `"`)
	uint64BlankJSON         = []byte(`""`)
	uint64MaxValueValidJSON = []byte(`{"Uint64":` + uint64String + `,"Valid":true}`)
	uint64NullInvalidGob    = `0affa9060102ffac00000038ffaa00342fffad0301010d55696e7436345772617070657201ffae000102010556616c7565010600010556616c6964010200000003ffae00`
	uint64OkValidGob        = `0affa9060102ffac00000044ffaa00402fffad0301010d55696e7436345772617070657201ffae000102010556616c7565010600010556616c696401020000000fffae01f8ffffffffffffffff010100`

	//	badObjectJSON           = []byte(`{"hello": "world"}`)
	//	int64JSON               = []byte(`12345`)
)

// Основной интерфейс который должны удовлетворять все типы
type mainInterface interface {
	// Reset Сброс значения и установка флага не действительного значения
	Reset()
}

func errorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func jsonEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("Bad %q data: %q ≠ %q\n", from, data, cmp)
	}
}

func TestAsString(t *testing.T) {
	jsonEquals(t, []byte(asString(string(boolTrueValidGob))), boolTrueValidGob, "asString(string)")
	jsonEquals(t, []byte(asString([]byte(bytesTestTextBase64))), string(bytesTestTextBase64), "asString([]byte)")
	jsonEquals(t, []byte(asString(int8(math.MaxInt8))), "127", "asString(int8)")
	jsonEquals(t, []byte(asString(int8(math.MinInt8))), "-128", "asString(int8)")
	jsonEquals(t, []byte(asString(int16(math.MaxInt16))), "32767", "asString(int16)")
	jsonEquals(t, []byte(asString(int16(math.MinInt16))), "-32768", "asString(int16)")
	jsonEquals(t, []byte(asString(int32(math.MaxInt32))), "2147483647", "asString(int32)")
	jsonEquals(t, []byte(asString(int32(math.MinInt32))), "-2147483648", "asString(int32)")
	jsonEquals(t, []byte(asString(int64(math.MaxInt64))), "9223372036854775807", "asString(int64)")
	jsonEquals(t, []byte(asString(int64(math.MinInt64))), "-9223372036854775808", "asString(int64)")
	jsonEquals(t, []byte(asString(uint8(math.MaxUint8))), "255", "asString(uint8)")
	jsonEquals(t, []byte(asString(uint16(math.MaxUint16))), "65535", "asString(uint16)")
	jsonEquals(t, []byte(asString(uint32(math.MaxUint32))), "4294967295", "asString(uint32)")
	jsonEquals(t, []byte(asString(uint64(math.MaxUint64))), "18446744073709551615", "asString(uint64)")
	jsonEquals(t, []byte(asString(float32(math.MaxFloat32))), "3.4028235e+38", "asString(float32)")
	jsonEquals(t, []byte(asString(float64(math.MaxFloat64))), "1.7976931348623157e+308", "asString(float64)")
	jsonEquals(t, []byte(asString(bool(true))), "true", "asString(bool)")
	jsonEquals(t, []byte(asString(bool(false))), "false", "asString(bool)")
	jsonEquals(t, []byte(asString(struct {
		a string
		b int8
	}{a: "custom", b: 120})), "{custom 120}", "asString(bool)")
}

func TestMainInterface(t *testing.T) {
	_ = mainInterface(&Bool{})
	_ = mainInterface(&Bytes{})
	_ = mainInterface(&Float64{})
	_ = mainInterface(&Int64{})
	_ = mainInterface(&String{})
	_ = mainInterface(&Time{})
	_ = mainInterface(&Uint64{})
}

func TestEncodingBinaryInterface(t *testing.T) {
	_ = encoding.BinaryMarshaler(&Bool{})
	_ = encoding.BinaryMarshaler(&Bytes{})
	_ = encoding.BinaryMarshaler(&Float64{})
	_ = encoding.BinaryMarshaler(&Int64{})
	_ = encoding.BinaryMarshaler(&String{})
	_ = encoding.BinaryMarshaler(&Time{})
	_ = encoding.BinaryMarshaler(&Uint64{})

	_ = encoding.BinaryUnmarshaler(&Bool{})
	_ = encoding.BinaryUnmarshaler(&Bytes{})
	_ = encoding.BinaryUnmarshaler(&Float64{})
	_ = encoding.BinaryUnmarshaler(&Int64{})
	_ = encoding.BinaryUnmarshaler(&String{})
	_ = encoding.BinaryUnmarshaler(&Time{})
	_ = encoding.BinaryUnmarshaler(&Uint64{})
}

func TestEncodingTextInterface(t *testing.T) {
	_ = encoding.TextMarshaler(&Bool{})
	_ = encoding.TextMarshaler(&Bytes{})
	_ = encoding.TextMarshaler(&Float64{})
	_ = encoding.TextMarshaler(&Int64{})
	_ = encoding.TextMarshaler(&String{})
	_ = encoding.TextMarshaler(&Time{})
	_ = encoding.TextMarshaler(&Uint64{})

	_ = encoding.TextUnmarshaler(&Bool{})
	_ = encoding.TextUnmarshaler(&Bytes{})
	_ = encoding.TextUnmarshaler(&Float64{})
	_ = encoding.TextUnmarshaler(&Int64{})
	_ = encoding.TextUnmarshaler(&String{})
	_ = encoding.TextUnmarshaler(&Time{})
	_ = encoding.TextUnmarshaler(&Uint64{})
}

func TestEncodingJsonInterface(t *testing.T) {
	_ = json.Marshaler(&Bool{})
	_ = json.Marshaler(&Bytes{})
	_ = json.Marshaler(&Float64{})
	_ = json.Marshaler(&Int64{})
	_ = json.Marshaler(&String{})
	_ = json.Marshaler(&Time{})
	_ = json.Marshaler(&Uint64{})

	_ = json.Unmarshaler(&Bool{})
	_ = json.Unmarshaler(&Bytes{})
	_ = json.Unmarshaler(&Float64{})
	_ = json.Unmarshaler(&Int64{})
	_ = json.Unmarshaler(&String{})
	_ = json.Unmarshaler(&Time{})
	_ = json.Unmarshaler(&Uint64{})
}

func TestSqlDriverValuerInterface(t *testing.T) {
	_ = driver.Valuer(&Bool{})
	_ = driver.Valuer(&Bytes{})
	_ = driver.Valuer(&Float64{})
	_ = driver.Valuer(&Int64{})
	_ = driver.Valuer(&String{})
	_ = driver.Valuer(&Time{})
	_ = driver.Valuer(&Uint64{})
}

func TestSqlScannerInterface(t *testing.T) {
	_ = sql.Scanner(&Bool{})
	_ = sql.Scanner(&Bytes{})
	_ = sql.Scanner(&Float64{})
	_ = sql.Scanner(&Int64{})
	_ = sql.Scanner(&String{})
	_ = sql.Scanner(&Time{})
	_ = sql.Scanner(&Uint64{})
}
