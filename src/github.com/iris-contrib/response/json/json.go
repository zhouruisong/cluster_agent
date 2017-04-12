package json

import (
	"bytes"
	"encoding/json"

	"github.com/valyala/bytebufferpool"
)

const (
	// ContentType the key for the engine, the user/dev can still use its own
	ContentType = "application/json"
)

var buffer bytebufferpool.Pool

// Engine the response engine which renders a JSON 'object'
type Engine struct {
	config Config
}

// New returns a new json response engine
func New(cfg ...Config) *Engine {
	c := DefaultConfig().Merge(cfg)
	return &Engine{config: c}
}

var (
	newLineB = []byte("\n")
	ltHex    = []byte("\\u003c")
	lt       = []byte("<")

	gtHex = []byte("\\u003e")
	gt    = []byte(">")

	andHex = []byte("\\u0026")
	and    = []byte("&")
)

// Response accepts the 'object' value and converts it to bytes in order to be 'renderable'
// implements the iris.ResponseEngine
func (e *Engine) Response(val interface{}, options ...map[string]interface{}) ([]byte, error) {
	if e.config.StreamingJSON {
		w := buffer.Get()
		if len(e.config.Prefix) > 0 {
			w.Write(e.config.Prefix)
		}
		err := json.NewEncoder(w).Encode(val)
		result := w.Bytes()
		buffer.Put(w)
		return result, err
	}

	var result []byte
	var err error

	if e.config.Indent {
		result, err = json.MarshalIndent(val, "", "  ")
		result = append(result, newLineB...)
	} else {
		result, err = json.Marshal(val)
	}
	if err != nil {
		return nil, err
	}

	if e.config.UnEscapeHTML {
		result = bytes.Replace(result, ltHex, lt, -1)
		result = bytes.Replace(result, gtHex, gt, -1)
		result = bytes.Replace(result, andHex, and, -1)
	}
	if len(e.config.Prefix) > 0 {
		result = append(e.config.Prefix, result...)
	}
	return result, nil
}
