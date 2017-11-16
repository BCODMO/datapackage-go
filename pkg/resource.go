package pkg

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/frictionlessdata/datapackage-go/clone"
)

type pathType byte

const (
	urlPath      pathType = 0
	relativePath pathType = 1
)

const (
	schemaProp    = "schema"
	nameProp      = "name"
	formatProp    = "format"
	mediaTypeProp = "mediatype"
	pathProp      = "path"
	dataProp      = "data"
	jsonFormat    = "json"
)

// Resource describes a data resource such as an individual file or table.
type Resource struct {
	descriptor map[string]interface{}
	Path       []string    `json:"path,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Name       string      `json:"name,omitempty"`
}

// MarshalJSON returns the JSON encoding of the resource.
func (r *Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.descriptor)
}

// UnmarshalJSON parses the JSON-encoded data and stores the result in the resource descriptor.
func (r *Resource) UnmarshalJSON(b []byte) error {
	var descriptor map[string]interface{}
	if err := json.Unmarshal(b, &descriptor); err != nil {
		return err
	}
	aux, err := NewResource(descriptor)
	if err != nil {
		return err
	}
	*r = *aux
	return nil
}

// Descriptor returns a copy of the underlying descriptor which describes the resource.
func (r *Resource) Descriptor() (map[string]interface{}, error) {
	return clone.Descriptor(r.descriptor)
}

// Valid checks whether the resource is valid.
func (r *Resource) Valid() bool {
	_, err := NewResource(r.descriptor)
	return err == nil
}

// NewResource creates a new Resource from the passed-in descriptor.
func NewResource(d map[string]interface{}) (*Resource, error) {
	if d[pathProp] != nil && d[dataProp] != nil {
		return nil, fmt.Errorf("either path or data properties MUST be set (only one of them). Descriptor:%v", d)
	}
	var err error
	r := Resource{
		descriptor: d,
	}
	r.Name, err = parseName(d[nameProp])
	if err != nil {
		return nil, err
	}
	schemaI := d[schemaProp]
	if schemaI != nil {
		if err := validateSchema(schemaI, d); err != nil {
			return nil, err
		}
	}
	pathI := d[pathProp]
	if pathI != nil {
		p, err := parsePath(pathI, d)
		if err != nil {
			return nil, err
		}
		r.Path = append(r.Path, p...)
		return &r, nil
	}
	dataI := d[dataProp]
	if dataI != nil {
		data, err := parseData(dataI, d)
		if err != nil {
			return nil, err
		}
		r.Data = data
		return &r, nil
	}
	return nil, fmt.Errorf("either path or data properties MUST be set  (only one of them). Descriptor:%v", d)
}

func validateSchema(schI interface{}, d map[string]interface{}) error {
	switch schI.(type) {
	case string:
		if _, err := parsePath(schI, d); err != nil {
			return err
		}
		return nil
	case map[string]interface{}:
		return nil
	}
	return fmt.Errorf("resource schema MUST be a string or a JSON schema object: %v", schI)
}

var nameRegexp = regexp.MustCompile(`^[a-z\._]+$`)

func parseName(name interface{}) (string, error) {
	if name == nil {
		return "", fmt.Errorf("resource MUST contain a name property. ")
	}
	n, ok := name.(string)
	if !ok {
		return "", fmt.Errorf("resource names MUST be strings: %v", name)
	}
	if !nameRegexp.MatchString(n) {
		return "", fmt.Errorf("resource names MUST consist only of lowercase alphanumeric characters plus \".\", \"-\" and \"_\":%v", n)
	}
	return n, nil
}

func parseData(dataI interface{}, d map[string]interface{}) (interface{}, error) {
	if dataI != nil {
		switch dataI.(type) {
		case string:
			if d[formatProp] == nil && d[mediaTypeProp] == nil {
				return nil, fmt.Errorf("format or mediatype properties MUST be provided for JSON data strings. Descriptor:%v", d)
			}
			return dataI, nil
		case []map[string]interface{}, map[string]interface{}:
			return dataI, nil
		}
	}
	return nil, fmt.Errorf("data property must be either a JSON array/object OR a JSON string. Descriptor:%v", d)
}

func parsePath(pathI interface{}, d map[string]interface{}) ([]string, error) {
	var returned []string
	// Parse.
	switch pathI.(type) {
	default:
		return nil, fmt.Errorf("path MUST be a string or an array of strings. Descriptor:%v", d)
	case string:
		if p, ok := pathI.(string); ok {
			returned = append(returned, p)
		}
	case []string:
		returned = append(returned, pathI.([]string)...)
	}
	var lastType, currType pathType
	// Validation.
	for index, p := range returned {
		// Check if it is a relative path.
		u, err := url.Parse(p)
		if err != nil || u.Scheme == "" {
			if path.IsAbs(p) || strings.HasPrefix(path.Clean(p), "..") {
				return nil, fmt.Errorf("absolute paths (/) and relative parent paths (../) MUST NOT be used. Descriptor:%v", d)
			}
			currType = relativePath
		} else { // Check if it is a valid URL.
			if u.Scheme != "http" && u.Scheme != "https" {
				return nil, fmt.Errorf("URLs MUST be fully qualified. MUST be using either http or https scheme. Descriptor:%v", d)
			}
			currType = urlPath
		}
		if index > 0 {
			if currType != lastType {
				return nil, fmt.Errorf("it is NOT permitted to mix fully qualified URLs and relative paths in a single resource. Descriptor:%v", d)
			}
			lastType = currType
		}
	}
	return returned, nil
}

// NewUncheckedResource returns an Resource instance based on the descriptor without any verification. The returned Resource might
// not be valid.
func NewUncheckedResource(d map[string]interface{}) (*Resource, error) {
	r := &Resource{descriptor: d}
	nI, ok := d["name"]
	if ok {
		nStr, ok := nI.(string)
		if ok {
			r.Name = nStr
		}
	}
	return r, nil
}
