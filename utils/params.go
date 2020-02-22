package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
	"time"
)

// Params represents a map with string params for the task.
type Params map[string]string

// FormatArgs represents the union between a Params struct and the current date.
// TODO: use Params to set the date?
type FormatArgs struct {
	Params Params    // A list of params
	Date   time.Time // Current date
}

// Returns a new FormatArgs object.
func NewFormatArgs(params Params) *FormatArgs {
	return &FormatArgs{
		Params: params,
		Date:   time.Now(),
	}
}

// transforms a list of strings in the format key:value to a Params struct.
func ListToParams(paramsList []string) Params {
	params := make(Params)
	for _, param := range paramsList {
		kv := strings.Split(param, ":")
		if len(kv) > 1 {
			key := kv[0]
			value := strings.Join(kv[1:], ":")
			params[key] = value
		}
	}
	return params
}

// FormatString formats a string using the defined params.
// If it fails, it returns the input argument.
func (params Params) FormatString(str string) string {
	tmpl, err := template.New(GenerateRandomHex(6)).Parse(str)
	if err != nil {
		return str
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, NewFormatArgs(params)); err != nil {
		return str
	}
	return buf.String()
}

// FormatStringMapString formats a map with string keys and string values, and returns
// a formatted map of the same types.
func (params Params) FormatStringMapString(m map[string]string) map[string]string {
	if m == nil {
		return m
	}
	newMap := make(map[string]string)
	for k, v := range m {
		newMap[params.FormatString(k)] = params.FormatString(v)
	}
	return newMap
}

// FormatStringArray formats an array with string values, and returns
// a formatted array of the same types.
func (params Params) FormatStringArray(arr []string) []string {
	if arr == nil {
		return arr
	}
	newArray := make([]string, len(arr))
	for i, v := range arr {
		newArray[i] = params.FormatString(v)
	}
	return newArray
}

// FormatReader formats a reader object using the defined params.
// If it fails, it returns the input argument.
func (params Params) FormatReader(reader io.Reader) io.Reader {
	if reader == nil {
		return reader
	}
	str, err := ioutil.ReadAll(reader)
	if err != nil {
		return reader
	}
	tmpl, err := template.New(GenerateRandomHex(6)).Parse(string(str))
	if err != nil {
		return reader
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, NewFormatArgs(params)); err != nil {
		return reader
	}
	return &buf
}

// Returns a new Params map with params from both sources.
// The preference of params is for params2.
func (params Params) Join(params2 Params) Params {
	newParams := make(Params)
	for k, v := range params {
		newParams[k] = v
	}
	for k, v := range params2 {
		newParams[k] = v
	}
	return newParams
}

func (params Params) Get(key string, defVal string) string {
	if val, ok := params[key]; ok {
		return val
	}
	return defVal
}

func (formatArgs *FormatArgs) Get(key, defVal string) string {
	return formatArgs.Params.Get(key, defVal)
}

func (formatArgs *FormatArgs) Today() string {
	return formatArgs.Date.Format("02-01-2006")
}
