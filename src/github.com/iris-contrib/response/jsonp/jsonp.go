package jsonp

import (
	"encoding/json"
)

const (
	// ContentType the key for the engine, the user/dev can still use its own
	ContentType = "application/javascript"
)

// Engine the response engine which renders a JSONP 'object' with its callback
type Engine struct {
	config Config
}

// New returns a new jsonp response engine
func New(cfg ...Config) *Engine {
	c := DefaultConfig().Merge(cfg)
	return &Engine{config: c}
}

func (e *Engine) getCallbackOption(options map[string]interface{}) string {
	callbackOpt := options["callback"]
	if s, isString := callbackOpt.(string); isString {
		return s
	}
	return e.config.Callback
}

var (
	finishCallbackB = []byte(");")
	newLineB        = []byte("\n")
)

// Response accepts the 'object' value and converts it to bytes in order to be 'renderable'
// implements the iris.ResponseEngine
func (e *Engine) Response(val interface{}, options ...map[string]interface{}) ([]byte, error) {
	var result []byte
	var err error
	if e.config.Indent {
		result, err = json.MarshalIndent(val, "", "  ")
	} else {
		result, err = json.Marshal(val)
	}

	if err != nil {
		return nil, err
	}

	// the config's callback can be overriden with the options
	callback := e.config.Callback
	if len(options) > 0 {
		callback = e.getCallbackOption(options[0])
	}

	if callback != "" {
		result = append([]byte(callback+"("), result...)
		result = append(result, finishCallbackB...)
	}

	if e.config.Indent {
		result = append(result, newLineB...)
	}
	return result, nil
}
