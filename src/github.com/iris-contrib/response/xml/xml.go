package xml

import (
	"encoding/xml"
)

const (
	// ContentType the key for the engine, the user/dev can still use its own
	ContentType = "text/xml"
)

// Engine the response engine which renders an XML 'object'
type Engine struct {
	config Config
}

// New returns a new xml response engine
func New(cfg ...Config) *Engine {
	c := DefaultConfig().Merge(cfg)
	return &Engine{config: c}
}

// Response accepts the 'object' value and converts it to bytes in order to be 'renderable'
// implements the iris.ResponseEngine
func (e *Engine) Response(val interface{}, options ...map[string]interface{}) ([]byte, error) {

	var result []byte
	var err error

	if e.config.Indent {
		result, err = xml.MarshalIndent(val, "", "  ")
		result = append(result, '\n')
	} else {
		result, err = xml.Marshal(val)
	}
	if err != nil {
		return nil, err
	}
	if len(e.config.Prefix) > 0 {
		result = append(e.config.Prefix, result...)
	}
	return result, nil
}
