package wrapper // import "gopkg.in/webnice/nul.v1/wrapper"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"encoding/gob"
)

func init() {
	// Register the concrete type for the encoder and decoder
	gob.Register(BoolWrapper{})
	gob.Register(BytesWrapper{})
	gob.Register(Float64Wrapper{})
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
