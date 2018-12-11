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
	"time"

	"gopkg.in/webnice/lin.v1/wrapper"
)

// Time is an nullable time.Time object
type Time struct {
	Time  time.Time // Value of object
	Valid bool      // Valid is true if value is not NULL
}

// NewTime Создание нового объекта Time
func NewTime() Time {
	return Time{
		Time:  time.Time{},
		Valid: false,
	}
}

// NewTimeValue Создание нового действительного объекта Time из значения
func NewTimeValue(value time.Time) Time {
	return Time{
		Time:  value,
		Valid: true,
	}
}

// NewTimePointerValue Создание нового действительного объекта Time из ссылки на значение
func NewTimePointerValue(ptr *time.Time) Time {
	if ptr == nil {
		return NewTime()
	}
	return NewTimeValue(*ptr)
}

// SetValid Изменение значения и установка флага действительного значения
func (t *Time) SetValid(value time.Time) { t.Time, t.Valid = value, true }

// Reset Сброс значения и установка флага не действительного значения
func (t *Time) Reset() { t.Time, t.Valid = time.Time{}, false }

// NullIfDefault Выполняет сброс значения до null, если значение переменной явзяется дефолтовым
func (t *Time) NullIfDefault() Time {
	if t.Time.IsZero() {
		t.Reset()
	}
	return *t
}

// MustValue Возвращает значение в любом случае
func (t *Time) MustValue() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// Pointer Возвращает ссылку на значение
func (t *Time) Pointer() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// Scan Реализация интерфейса Scanner
func (t *Time) Scan(value interface{}) (err error) {
	switch x := value.(type) {
	case nil:
		t.Time, t.Valid = time.Time{}, false
		return
	case string:
		t.Time, err = time.Parse(time.RFC3339Nano, x)
		t.Valid = err == nil
	case time.Time:
		t.Time, t.Valid = x, true
	default:
		buf := asString(x)
		t.Time, err = time.Parse(time.RFC3339Nano, buf)
		t.Valid = err == nil
	}

	return
}

// Value Реализация интерфейса driver.Valuer
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// UnmarshalJSON Реализация интерфейса json.Unmarshaler
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	var v interface{}

	if err = json.Unmarshal(data, &v); err != nil {
		return
	}
	switch x := v.(type) {
	case nil:
		t.Time, t.Valid = time.Time{}, false
		return
	case string:
		err = t.Time.UnmarshalJSON(data)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			return fmt.Errorf(`Unmarshalling object into Go value of type nul.Time requires key "Time" to be of type `+
				`string and key "Valid" to be of type bool; `+
				`found %T and %T, respectively`, x["Time"], x["Valid"])
		}
		err = t.Time.UnmarshalText([]byte(ti))
		t.Valid = valid
		return
	default:
		err = fmt.Errorf("Can't unmarshal %s into Go value of type nul.Time", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil

	return
}

// MarshalJSON Реализация интерфейса json.Marshaler
func (t Time) MarshalJSON() (data []byte, err error) {
	const (
		nullString = "null"
	)

	if !t.Valid {
		data = []byte(nullString)
		return
	}
	data, err = t.Time.MarshalJSON()

	return
}

// UnmarshalText Реализация интерфейса encoding.TextUnmarshaler
func (t *Time) UnmarshalText(text []byte) (err error) {
	const (
		emptyString = ""
		nullString  = "null"
	)
	var str string

	switch str = string(text); str {
	case nullString:
		t.Time, t.Valid = time.Time{}, false
		return
	case emptyString:
		t.Time, t.Valid = time.Time{}, true
		return
	default:
		err = t.Time.UnmarshalText(text)
	}
	t.Valid = err == nil

	return
}

// MarshalText Реализация интерфейса encoding.TextMarshaler
func (t Time) MarshalText() (text []byte, err error) {
	const (
		emptyString = ""
		nullString  = "null"
	)

	if !t.Valid {
		text = []byte(nullString)
		return
	}
	if t.Time.IsZero() {
		text = []byte(emptyString)
		return
	}
	text, err = t.Time.MarshalText()

	return
}

// UnmarshalBinary Реализация интерфейса encoding.BinaryUnmarshaler
func (t *Time) UnmarshalBinary(data []byte) (err error) {
	var reader *bytes.Reader
	var dec *gob.Decoder
	var item *wrapper.TimeWrapper

	reader = bytes.NewReader(data)
	dec = gob.NewDecoder(reader)
	item = new(wrapper.TimeWrapper)
	if err = dec.Decode(item); err == nil {
		t.Time, t.Valid = item.Value, item.Valid
	}

	return
}

// MarshalBinary Реализация интерфейса encoding.BinaryMarshaler
func (t Time) MarshalBinary() (data []byte, err error) {
	var buf *bytes.Buffer
	var enc *gob.Encoder
	var item *wrapper.TimeWrapper

	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)
	item = &wrapper.TimeWrapper{Value: t.Time, Valid: t.Valid}
	err = enc.Encode(item)
	data = buf.Bytes()

	return
}
