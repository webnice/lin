package nul // import "gopkg.in/webnice/lin.v1/nl"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"gopkg.in/webnice/lin.v1/wrapper"
)

// Bytes is an nullable []byte object
type Bytes struct {
	Bytes *bytes.Buffer // Value of object
	Valid bool          // Valid is true if value is not NULL
}

// NewBytes Создание нового объекта []byte
func NewBytes() Bytes {
	return Bytes{
		Bytes: &bytes.Buffer{},
		Valid: false,
	}
}

// NewBytesValue Создание нового действительного объекта Bytes из значения
func NewBytesValue(value []byte) Bytes {
	var buf = make([]byte, len(value))
	_ = copy(buf, value)
	return Bytes{
		Bytes: bytes.NewBuffer(buf),
		Valid: true,
	}
}

// NewBytesPointerValue Создание нового действительного объекта Bytes из ссылки на значение
func NewBytesPointerValue(ptr *[]byte) Bytes {
	if ptr == nil {
		return NewBytes()
	}
	return NewBytesValue(*ptr)
}

// SetValid Изменение значения и установка флага действительного значения
func (bt *Bytes) SetValid(value []byte) {
	var buf = make([]byte, len(value))
	_ = copy(buf, value)
	bt.Bytes, bt.Valid = bytes.NewBuffer(buf), true
}

// Reset Сброс значения и установка флага не действительного значения
func (bt *Bytes) Reset() { bt.Bytes.Reset(); bt.Valid = false }

// NullIfDefault Выполняет сброс значения до null, если значение переменной явзяется дефолтовым
func (bt *Bytes) NullIfDefault() Bytes {
	if bt.Bytes != nil && bt.Bytes.Len() == 0 {
		bt.Reset()
	}
	return *bt
}

// MustValue Возвращает значение в любом случае
func (bt *Bytes) MustValue() []byte {
	if !bt.Valid {
		return []byte{}
	}
	return bt.Bytes.Bytes()
}

// Pointer Возвращает ссылку на значение
func (bt *Bytes) Pointer() *[]byte {
	var ret []byte
	if !bt.Valid {
		return nil
	}
	ret = bt.Bytes.Bytes()
	return &ret
}

// Scan Реализация интерфейса Scanner
func (bt *Bytes) Scan(value interface{}) (err error) {
	switch x := value.(type) {
	case nil:
		bt.Valid = false
		return
	case []byte:
		bt.SetValid(value.([]byte))
	case string:
		bt.SetValid([]byte(value.(string)))
	default:
		err = fmt.Errorf("can't scan type %T into nul.Bytes: %v", x, value)
	}
	bt.Valid = err == nil

	return
}

// Value Реализация интерфейса driver.Valuer
func (bt Bytes) Value() (driver.Value, error) {
	if !bt.Valid {
		return nil, nil
	}
	return bt.Bytes.Bytes(), nil
}

// UnmarshalJSON Реализация интерфейса json.Unmarshaler
func (bt *Bytes) UnmarshalJSON(data []byte) (err error) {
	var (
		v   interface{}
		buf []byte
	)

	if err = json.Unmarshal(data, &v); err != nil {
		return
	}
	switch x := v.(type) {
	case nil:
		bt.Valid = false
		return
	case map[string]interface{}:
		value, okValue := x["Bytes"].(string)
		valid, okValid := x["Valid"].(bool)
		if !okValue || !okValid {
			return fmt.Errorf("unmarshalling object into Go value of type nul.Bytes requires key "+
				`"Bytes" to be of type string and key "Valid" to be of type string; `+
				"found %T and %T, respectively", x["Bytes"], x["Valid"])
		}
		if buf, err = base64.StdEncoding.DecodeString(value); err == nil {
			bt.SetValid(buf)
		}
		bt.Valid = valid
	}
	bt.Valid = err == nil

	return
}

// MarshalJSON Реализация интерфейса json.Marshaler
func (bt Bytes) MarshalJSON() (data []byte, err error) {
	const nullString = "null"

	if !bt.Valid {
		data = []byte(nullString)
		return
	}
	data, err = json.Marshal(bt.Bytes.Bytes())

	return
}

// UnmarshalText Реализация интерфейса encoding.TextUnmarshaler
func (bt *Bytes) UnmarshalText(text []byte) (err error) {
	const (
		emptyString = ""
		nullString  = "null"
	)
	var (
		str string
		buf []byte
	)

	switch str = string(text); str {
	case nullString:
		bt.Bytes, bt.Valid = &bytes.Buffer{}, false
		return
	case emptyString:
		bt.Bytes, bt.Valid = &bytes.Buffer{}, true
		return
	default:
		if buf, err = base64.StdEncoding.DecodeString(str); err == nil {
			bt.SetValid(buf)
		}
	}

	return
}

// MarshalText Реализация интерфейса encoding.TextMarshaler
func (bt Bytes) MarshalText() (text []byte, err error) {
	const nullString = "null"

	if !bt.Valid {
		text = []byte(nullString)
		return
	}
	if bt.Bytes.Len() == 0 {
		text = make([]byte, 0)
		return
	}
	text = []byte(base64.StdEncoding.EncodeToString(bt.Bytes.Bytes()))

	return
}

// UnmarshalBinary Реализация интерфейса encoding.BinaryUnmarshaler
func (bt *Bytes) UnmarshalBinary(data []byte) (err error) {
	var (
		reader *bytes.Reader
		dec    *gob.Decoder
		item   *wrapper.BytesWrapper
	)

	reader = bytes.NewReader(data)
	dec = gob.NewDecoder(reader)
	item = new(wrapper.BytesWrapper)
	if err = dec.Decode(item); err == nil {
		bt.SetValid(item.Value)
		bt.Valid = item.Valid
	}

	return
}

// MarshalBinary Реализация интерфейса encoding.BinaryMarshaler
func (bt Bytes) MarshalBinary() (data []byte, err error) {
	var (
		buf  *bytes.Buffer
		enc  *gob.Encoder
		item *wrapper.BytesWrapper
	)

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)
	item = &wrapper.BytesWrapper{Value: bt.Bytes.Bytes(), Valid: bt.Valid}
	err = enc.Encode(item)
	data = buf.Bytes()

	return
}
