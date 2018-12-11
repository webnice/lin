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
	"strconv"

	"gopkg.in/webnice/lin.v1/wrapper"
)

// Int64 is an nullable int64 object
type Int64 struct {
	Int64 int64 // Value of object
	Valid bool  // Valid is true if value is not NULL
}

// NewInt64 Создание нового объекта Int64
func NewInt64() Int64 {
	return Int64{
		Int64: 0,
		Valid: false,
	}
}

// NewInt64Value Создание нового действительного объекта Int64 из значения
func NewInt64Value(value int64) Int64 {
	return Int64{
		Int64: value,
		Valid: true,
	}
}

// NewInt64PointerValue Создание нового действительного объекта Int64 из ссылки на значение
func NewInt64PointerValue(ptr *int64) Int64 {
	if ptr == nil {
		return NewInt64()
	}
	return NewInt64Value(*ptr)
}

// SetValid Изменение значения и установка флага действительного значения
func (i *Int64) SetValid(value int64) { i.Int64, i.Valid = value, true }

// Reset Сброс значения и установка флага не действительного значения
func (i *Int64) Reset() { i.Int64, i.Valid = 0, false }

// NullIfDefault Выполняет сброс значения до null, если значение переменной явзяется дефолтовым
func (i *Int64) NullIfDefault() Int64 {
	if i.Int64 == 0 {
		i.Reset()
	}
	return *i
}

// MustValue Возвращает значение в любом случае
func (i *Int64) MustValue() int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

// Pointer Возвращает ссылку на значение
func (i *Int64) Pointer() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// Scan Реализация интерфейса Scanner
func (i *Int64) Scan(value interface{}) (err error) {
	switch x := value.(type) {
	case nil:
		i.Int64, i.Valid = 0, false
		return
	case int64:
		i.Int64 = value.(int64)
	default:
		buf := asString(x)
		i.Int64, err = strconv.ParseInt(buf, 10, 64)
	}
	i.Valid = err == nil

	return
}

// Value Реализация интерфейса driver.Valuer
func (i Int64) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.Int64, nil
}

// UnmarshalJSON Реализация интерфейса json.Unmarshaler
func (i *Int64) UnmarshalJSON(data []byte) (err error) {
	var v interface{}

	if err = json.Unmarshal(data, &v); err != nil {
		return
	}
	switch x := v.(type) {
	case float64:
		i.Int64, err = strconv.ParseInt(string(data), 10, 64)
	case string:
		var str = x
		if len(str) == 0 {
			i.Valid = false
			return
		}
		i.Int64, err = strconv.ParseInt(str, 10, 64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.Int64)
	case nil:
		i.Valid = false
		return
	default:
		err = fmt.Errorf("Can't unmarshal %v into Go value of type nul.Int64", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil

	return
}

// MarshalJSON Реализация интерфейса json.Marshaler
func (i Int64) MarshalJSON() (data []byte, err error) {
	const (
		nullString = "null"
	)

	if !i.Valid {
		data = []byte(nullString)
		return
	}
	data = []byte(strconv.FormatInt(i.Int64, 10))

	return
}

// UnmarshalText Реализация интерфейса encoding.TextUnmarshaler
func (i *Int64) UnmarshalText(text []byte) (err error) {
	const (
		emptyString = ""
		nullString  = "null"
	)
	var str string

	switch str = string(text); str {
	case nullString:
		i.Int64, i.Valid = 0, false
		return
	case emptyString:
		i.Int64, i.Valid = 0, true
		return
	default:
		i.Int64, err = strconv.ParseInt(str, 10, 64)
	}
	i.Valid = err == nil

	return
}

// MarshalText Реализация интерфейса encoding.TextMarshaler
func (i Int64) MarshalText() (text []byte, err error) {
	const (
		nullString = "null"
	)

	if !i.Valid {
		text = []byte(nullString)
		return
	}
	text = []byte(strconv.FormatInt(i.Int64, 10))

	return
}

// UnmarshalBinary Реализация интерфейса encoding.BinaryUnmarshaler
func (i *Int64) UnmarshalBinary(data []byte) (err error) {
	var reader *bytes.Reader
	var dec *gob.Decoder
	var item *wrapper.Int64Wrapper

	reader = bytes.NewReader(data)
	dec = gob.NewDecoder(reader)
	item = new(wrapper.Int64Wrapper)
	if err = dec.Decode(item); err == nil {
		i.Int64, i.Valid = item.Value, item.Valid
	}

	return
}

// MarshalBinary Реализация интерфейса encoding.BinaryMarshaler
func (i Int64) MarshalBinary() (data []byte, err error) {
	var buf *bytes.Buffer
	var enc *gob.Encoder
	var item *wrapper.Int64Wrapper

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)
	item = &wrapper.Int64Wrapper{Value: i.Int64, Valid: i.Valid}
	err = enc.Encode(item)
	data = buf.Bytes()

	return
}
