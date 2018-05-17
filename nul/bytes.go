package nul // import "gopkg.in/webnice/nul.v1/nul"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"gopkg.in/webnice/nul.v1/wrapper"
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
	return Bytes{
		Bytes: bytes.NewBuffer(value),
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
func (bt *Bytes) SetValid(value []byte) { bt.Bytes, bt.Valid = bytes.NewBuffer(value), true }

// Reset Сброс значения и установка флага не действительного значения
func (bt *Bytes) Reset() { bt.Bytes.Reset(); bt.Valid = false }

// NullIfDefault Выполняет сброс значения до null, если значение переменной явзяется дефолтовым
func (bt *Bytes) NullIfDefault() {
	if bt.Bytes.Len() == 0 {
		bt.Reset()
	}
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
		bt.Bytes = bytes.NewBuffer(value.([]byte))
	case string:
		bt.Bytes = bytes.NewBufferString(value.(string))
	default:
		err = fmt.Errorf("Can't scan type %T into nul.Bytes: %v", x, value)
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
	var v interface{}

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
			return fmt.Errorf("Unmarshalling object into Go value of type nul.Bytes requires key "+
				`"Bytes" to be of type string and key "Valid" to be of type string; `+
				"found %T and %T, respectively", x["Bytes"], x["Valid"])
		}
		var buf []byte
		buf, err = base64.StdEncoding.DecodeString(value)
		bt.Bytes = bytes.NewBuffer(buf)
		bt.Valid = valid
	}
	bt.Valid = err == nil

	return
}

// MarshalJSON Реализация интерфейса json.Marshaler
func (bt Bytes) MarshalJSON() (data []byte, err error) {
	const (
		nullString = "null"
	)

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
	var str string
	var buf []byte

	switch str = string(text); str {
	case nullString:
		bt.Bytes, bt.Valid = &bytes.Buffer{}, false
		return
	case emptyString:
		bt.Bytes, bt.Valid = &bytes.Buffer{}, true
		return
	default:
		buf, err = base64.StdEncoding.DecodeString(str)
		bt.Bytes, bt.Valid = bytes.NewBuffer(buf), err == nil
	}

	return
}

// MarshalText Реализация интерфейса encoding.TextMarshaler
func (bt Bytes) MarshalText() (text []byte, err error) {
	const (
		nullString = "null"
	)
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
	var reader *bytes.Reader
	var dec *gob.Decoder
	var item *wrapper.BytesWrapper

	reader = bytes.NewReader(data)
	dec = gob.NewDecoder(reader)
	item = new(wrapper.BytesWrapper)
	if err = dec.Decode(item); err == nil {
		bt.Bytes, bt.Valid = bytes.NewBuffer(item.Value), item.Valid
	}

	return
}

// MarshalBinary Реализация интерфейса encoding.BinaryMarshaler
func (bt Bytes) MarshalBinary() (data []byte, err error) {
	var buf *bytes.Buffer
	var enc *gob.Encoder
	var item *wrapper.BytesWrapper

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)
	item = &wrapper.BytesWrapper{Value: bt.Bytes.Bytes(), Valid: bt.Valid}
	err = enc.Encode(item)
	data = buf.Bytes()

	return
}
