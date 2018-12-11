package nul // import "gopkg.in/webnice/lin.v1/nl"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"

	"gopkg.in/webnice/lin.v1/wrapper"
)

// String is an nullable string object
type String struct {
	String string // Value of object
	Valid  bool   // Valid is true if value is not NULL
}

// NewString Создание нового объекта String
func NewString() String {
	return String{
		String: ``,
		Valid:  false,
	}
}

// NewStringValue Создание нового действительного объекта String из значения
func NewStringValue(value string) String {
	return String{
		String: value,
		Valid:  true,
	}
}

// NewStringPointerValue Создание нового действительного объекта String из ссылки на значение
func NewStringPointerValue(ptr *string) String {
	if ptr == nil {
		return NewString()
	}
	return NewStringValue(*ptr)
}

// SetValid Изменение значения и установка флага действительного значения
func (s *String) SetValid(value string) { s.String, s.Valid = value, true }

// Reset Сброс значения и установка флага не действительного значения
func (s *String) Reset() {
	const emptyString = ""
	s.String, s.Valid = emptyString, false
}

// NullIfDefault Выполняет сброс значения до null, если значение переменной явзяется дефолтовым
func (s *String) NullIfDefault() String {
	const emptyString = ""
	if s.String == emptyString {
		s.Reset()
	}
	return *s
}

// MustValue Возвращает значение в любом случае
func (s *String) MustValue() string {
	const emptyString = ""
	if !s.Valid {
		return emptyString
	}
	return s.String
}

// Pointer Возвращает ссылку на значение
func (s *String) Pointer() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// Scan Реализация интерфейса Scanner
func (s *String) Scan(value interface{}) (err error) {
	const emptyString = ""
	if value == nil {
		s.String, s.Valid = emptyString, false
		return nil
	}
	s.String = asString(value)
	s.Valid = err == nil

	return
}

// Value Реализация интерфейса driver.Valuer
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}

// UnmarshalJSON Реализация интерфейса json.Unmarshaler
func (s *String) UnmarshalJSON(data []byte) (err error) {
	var v interface{}

	if err = json.Unmarshal(data, &v); err != nil {
		return
	}
	switch x := v.(type) {
	case nil:
		s.Valid = false
		return
	case string:
		s.String = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &s.String)
	default:
		err = fmt.Errorf("Can't unmarshal %s into Go value of type nul.String", reflect.TypeOf(v).Name())
	}
	s.Valid = err == nil

	return
}

// MarshalJSON Реализация интерфейса json.Marshaler
func (s String) MarshalJSON() (data []byte, err error) {
	const (
		nullString = "null"
	)

	if !s.Valid {
		data = []byte(nullString)
		return
	}
	data, err = json.Marshal(s.String)

	return
}

// UnmarshalText Реализация интерфейса encoding.TextUnmarshaler
func (s *String) UnmarshalText(text []byte) (err error) {
	const (
		emptyString = ""
		nullString  = "null"
	)
	var str string

	switch str = string(text); str {
	case nullString:
		s.String, s.Valid = emptyString, false
		return
	case emptyString:
		s.String, s.Valid = emptyString, true
		return
	default:
		s.String = str
		s.Valid = err == nil
	}

	return
}

// MarshalText Реализация интерфейса encoding.TextMarshaler
func (s String) MarshalText() (text []byte, err error) {
	const (
		nullString = "null"
	)

	if !s.Valid {
		text = []byte(nullString)
		return
	}
	text = []byte(s.String)

	return
}

// UnmarshalBinary Реализация интерфейса encoding.BinaryUnmarshaler
func (s *String) UnmarshalBinary(data []byte) (err error) {
	var reader *bytes.Reader
	var dec *gob.Decoder
	var item *wrapper.StringWrapper

	reader = bytes.NewReader(data)
	dec = gob.NewDecoder(reader)
	item = new(wrapper.StringWrapper)
	if err = dec.Decode(item); err == nil {
		s.String, s.Valid = item.Value, item.Valid
	}

	return
}

// MarshalBinary Реализация интерфейса encoding.BinaryMarshaler
func (s String) MarshalBinary() (data []byte, err error) {
	var buf *bytes.Buffer
	var enc *gob.Encoder
	var item *wrapper.StringWrapper

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)
	item = &wrapper.StringWrapper{Value: s.String, Valid: s.Valid}
	err = enc.Encode(item)
	data = buf.Bytes()

	return
}
