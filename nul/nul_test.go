package nul // import "gopkg.in/webnice/nul.v1/nul"

import (
	"testing"
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
