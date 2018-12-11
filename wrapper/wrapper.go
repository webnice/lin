package wrapper // import "gopkg.in/webnice/lin.v1/wrapper"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"encoding/gob"
	"time"
)

func init() {
	// Register the concrete type for the encoder and decoder
	gob.Register(BoolWrapper{})
	gob.Register(BytesWrapper{})
	gob.Register(Float64Wrapper{})
	gob.Register(Int64Wrapper{})
	gob.Register(StringWrapper{})
	gob.Register(TimeWrapper{})
	gob.Register(Uint64Wrapper{})
}

// BoolWrapper Обёртка для Bool
type BoolWrapper struct {
	Value bool
	Valid bool
}

// BytesWrapper Обёртка для Bytes
type BytesWrapper struct {
	Value []byte
	Valid bool
}

// Float64Wrapper Обёртка для Float64
type Float64Wrapper struct {
	Value float64
	Valid bool
}

// Int64Wrapper Обёртка для Int64
type Int64Wrapper struct {
	Value int64
	Valid bool
}

// StringWrapper Обёртка для String
type StringWrapper struct {
	Value string
	Valid bool
}

// TimeWrapper Обёртка для Time
type TimeWrapper struct {
	Value time.Time
	Valid bool
}

// Uint64Wrapper Обёртка для Uint64
type Uint64Wrapper struct {
	Value uint64
	Valid bool
}
