package nul // import "gopkg.in/webnice/nul.v1/nul"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"math"
	"testing"
)

var (
	invalidJSON         = []byte(`:{;-}`)
	blankJSON           = []byte(`{}`)
	boolTrueJSON        = []byte(`true`)
	boolFalseJSON       = []byte(`false`)
	boolNullJSON        = []byte(`null`)
	boolFalseValidJSON  = []byte(`{"Bool":false,"Valid":true}`)
	boolTrueValidGob    = `0aff81060102ff840000003aff8200362dff850301010b426f6f6c5772617070657201ff86000102010556616c7565010200010556616c6964010200000007ff860101010100`
	boolFalseValidGob   = `0aff81060102ff8400000038ff8200342dff850301010b426f6f6c5772617070657201ff86000102010556616c7565010200010556616c6964010200000005ff86020100`
	boolNullInvalidGob  = `0aff81060102ff8400000036ff8200322dff850301010b426f6f6c5772617070657201ff86000102010556616c7565010200010556616c6964010200000003ff8600`
	bytesTestValidJSON  = []byte(`{"Bytes":"VGVzdCBkYXRhIDFwSHVPeEFEWmtlaDhZOVd2TDc1","Valid":true}`)
	bytesTestTextBase64 = []byte(`VGVzdCBkYXRhIDFwSHVPeEFEWmtlaDhZOVd2TDc1`)
	bytesNullInvalidGob = `0aff87060102ff8a00000012ff8b0301010642756666657201ff8c00000037ff8800332eff8d0301010c42797465735772617070657201ff8e000102010556616c7565010a00010556616c6964010200000003ff8e00`
	bytesOkValidGob     = `0aff87060102ff8a00000012ff8b0301010642756666657201ff8c00000059ff8800552eff8d0301010c42797465735772617070657201ff8e000102010556616c7565010a00010556616c6964010200000025ff8e011e54657374206461746120317048754f7841445a6b656838593957764c3735010100`
	bytesZeroValidGob   = `0aff87060102ff8a00000012ff8b0301010642756666657201ff8c00000039ff8800352eff8d0301010c42797465735772617070657201ff8e000102010556616c7565010a00010556616c6964010200000005ff8e020100`
	bytesNullJSON       = []byte(`null`)
	float64JSON         = []byte(`179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000`)
	float64StringJSON   = []byte(`"179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000"`)
	float64BlankJSON    = []byte(`""`)
	floatNullInvalidGob = `0aff8f060102ff9200000039ff90003530ff930301010e466c6f617436345772617070657201ff94000102010556616c7565010800010556616c6964010200000003ff9400`
	floatOkValidGob     = `0aff8f060102ff9200000045ff90004130ff930301010e466c6f617436345772617070657201ff94000102010556616c7565010800010556616c696401020000000fff9401f8ffffffffffffef7f010100`
	int64JSON           = []byte(`9223372036854775807`)
	int64StringJSON     = []byte(`"9223372036854775807"`)
	int64BlankJSON      = []byte(`""`)
	intNullInvalidGob   = `0aff95060102ff9800000037ff9600332eff990301010c496e7436345772617070657201ff9a000102010556616c7565010400010556616c6964010200000003ff9a00`
	intOkValidGob       = `0aff95060102ff9800000043ff96003f2eff990301010c496e7436345772617070657201ff9a000102010556616c7565010400010556616c696401020000000fff9a01f8fffffffffffffffe010100`

	uint64MaxValueValidJSON = []byte(`{"Uint64":18446744073709551615,"Valid":true}`)

	//	badObjectJSON           = []byte(`{"hello": "world"}`)
	//	int64JSON               = []byte(`12345`)
)

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
