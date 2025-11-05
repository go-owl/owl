package owl

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	// maxFieldLength is the maximum allowed length for a single field value (10KB)
	maxFieldLength = 10000
	// maxFileSize is the maximum size per uploaded file (50MB)
	maxFileSize = 50 << 20
)

// Binder handles different content type bindings.
type Binder struct {
	request *http.Request
}

// JSON binds request body as JSON.
// Go's json.Decoder automatically protects against deeply nested JSON (max depth ~10000).
func (b *Binder) JSON(dst interface{}) error {
	if b.request.Body == nil {
		return NewHTTPError(http.StatusBadRequest, "request body is empty")
	}
	defer b.request.Body.Close()

	dec := json.NewDecoder(b.request.Body)

	if err := dec.Decode(dst); err != nil {
		return NewHTTPError(http.StatusBadRequest, "invalid JSON: "+err.Error())
	}

	return nil
}

// XML binds request body as XML.
// Note: External entities are automatically disabled by Go's xml.Decoder for security.
func (b *Binder) XML(dst interface{}) error {
	if b.request.Body == nil {
		return NewHTTPError(http.StatusBadRequest, "request body is empty")
	}
	defer b.request.Body.Close()

	// Create decoder (Go's xml package is safe from XXE by default)
	decoder := xml.NewDecoder(b.request.Body)

	if err := decoder.Decode(dst); err != nil {
		return NewHTTPError(http.StatusBadRequest, "invalid XML: "+err.Error())
	}
	return nil
}

// Text binds request body as plain text string.
// Useful for webhooks or when you need raw body content.
// Note: Body size is automatically limited by App's BodyLimit config via MaxBytesReader.
func (b *Binder) Text(dst *string) error {
	data, err := b.readBodySafe()
	if err != nil {
		return err
	}
	*dst = string(data)
	return nil
}

// Bytes binds request body as raw bytes.
// Useful for binary data or when you need the raw payload.
// Note: Body size is automatically limited by App's BodyLimit config via MaxBytesReader.
func (b *Binder) Bytes(dst *[]byte) error {
	data, err := b.readBodySafe()
	if err != nil {
		return err
	}
	*dst = data
	return nil
}

// readBodySafe reads the request body safely (body limit handled by App-level MaxBytesReader)
func (b *Binder) readBodySafe() ([]byte, error) {
	if b.request.Body == nil {
		return nil, NewHTTPError(http.StatusBadRequest, "request body is empty")
	}
	defer b.request.Body.Close()

	// Read body - size limit is enforced by App's MaxBytesReader in wrapHandler
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(b.request.Body); err != nil {
		return nil, NewHTTPError(http.StatusBadRequest, "failed to read body: "+err.Error())
	}

	return buf.Bytes(), nil
}

// Query binds URL query parameters to dst struct.
// Supports string, int, int64, float64, bool types.
// Example: /users?name=John&age=25 -> struct{Name string; Age int}
func (b *Binder) Query(dst interface{}) error {
	values := b.request.URL.Query()
	return bindValues(values, dst)
}

// Form binds request form data (application/x-www-form-urlencoded) to dst struct.
// Supports string, int, int64, float64, bool types.
func (b *Binder) Form(dst interface{}) error {
	if err := b.request.ParseForm(); err != nil {
		return NewHTTPError(http.StatusBadRequest, "invalid form data: "+err.Error())
	}
	return bindValues(b.request.PostForm, dst)
}

// MultipartForm binds multipart form data (for file uploads) to dst struct.
// Use *multipart.FileHeader for file fields.
// Example: struct { Name string; Avatar *multipart.FileHeader }
func (b *Binder) MultipartForm(dst interface{}, maxMemory int64) error {
	if maxMemory == 0 {
		maxMemory = 32 << 20 // 32MB default
	}

	if err := b.request.ParseMultipartForm(maxMemory); err != nil {
		return NewHTTPError(http.StatusBadRequest, "invalid multipart form: "+err.Error())
	}

	// Bind form values
	if err := bindValues(b.request.MultipartForm.Value, dst); err != nil {
		return err
	}

	// Bind file uploads
	return bindFiles(b.request.MultipartForm.File, dst)
}

// File retrieves a single uploaded file by field name.
// Returns the file header and a reader.
func (b *Binder) File(name string) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := b.request.FormFile(name)
	if err != nil {
		return nil, nil, NewHTTPError(http.StatusBadRequest, "failed to get file: "+err.Error())
	}
	return file, header, nil
}

// Auto automatically detects the content type and binds accordingly.
// Provides excellent DX by eliminating manual content-type checking.
// Example: c.Bind().Auto(&data) - works with JSON, Form, Multipart, or XML
func (b *Binder) Auto(dst interface{}) error {
	ct := b.request.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "application/json"):
		return b.JSON(dst)
	case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
		return b.Form(dst)
	case strings.HasPrefix(ct, "multipart/form-data"):
		return b.MultipartForm(dst, 32<<20) // 32MB default
	case strings.HasPrefix(ct, "application/xml"), strings.HasPrefix(ct, "text/xml"):
		return b.XML(dst)
	default:
		return NewHTTPError(http.StatusUnsupportedMediaType, "unsupported content type: "+ct)
	}
}

// tagName extracts the field name from struct tags, handling options like "name,omitempty"
func tagName(field reflect.StructField, keys ...string) string {
	for _, key := range keys {
		if raw := field.Tag.Get(key); raw != "" && raw != "-" {
			// Split by comma to handle options like "name,omitempty"
			name := strings.Split(raw, ",")[0]
			if name != "" && name != "-" {
				return name
			}
		}
	}
	// Fallback to lowercase field name
	return strings.ToLower(field.Name)
}

// bindValues binds url.Values to a struct using reflection
func bindValues(values url.Values, dst interface{}) (err error) {
	// Panic recovery for reflection errors
	defer func() {
		if r := recover(); r != nil {
			err = NewHTTPError(http.StatusBadRequest, "binding panic: reflection error")
		}
	}()

	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return NewHTTPError(http.StatusBadRequest, "dst must be a pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return NewHTTPError(http.StatusBadRequest, "dst must be a pointer to struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		fieldType := t.Field(i)

		// Get tag name, handling options like "name,omitempty"
		tag := tagName(fieldType, "form", "query", "json")

		// Handle pointer fields by dereferencing
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
		}

		// Handle array fields
		if field.Kind() == reflect.Array {
			vals := values[tag]
			if len(vals) == 0 {
				continue
			}
			n := field.Len()
			if len(vals) < n {
				n = len(vals)
			}
			for i := 0; i < n; i++ {
				if len(vals[i]) > maxFieldLength {
					return NewHTTPError(http.StatusBadRequest, "field value too long: "+fieldType.Name)
				}
				if err := setField(field.Index(i), vals[i]); err != nil {
					return NewHTTPError(http.StatusBadRequest, "invalid value for field "+fieldType.Name+": "+err.Error())
				}
			}
			continue
		}

		// Handle slices for multiple values (?tag=a&tag=b&score=1&score=2)
		if field.Kind() == reflect.Slice {
			vals := values[tag]
			if len(vals) == 0 {
				continue
			}

			elem := field.Type().Elem()
			out := reflect.MakeSlice(field.Type(), 0, len(vals))

			for _, sv := range vals {
				// Check value length for security
				if len(sv) > maxFieldLength {
					return NewHTTPError(http.StatusBadRequest, "field value too long: "+fieldType.Name)
				}

				ev := reflect.New(elem).Elem()
				if err := setField(ev, sv); err != nil {
					return NewHTTPError(http.StatusBadRequest, "invalid value for field "+fieldType.Name+": "+err.Error())
				}
				out = reflect.Append(out, ev)
			}
			field.Set(out)
			continue
		}

		// Single value
		valueStr := values.Get(tag)
		if valueStr == "" {
			continue
		}

		// Limit string length to prevent memory exhaustion
		if len(valueStr) > maxFieldLength {
			return NewHTTPError(http.StatusBadRequest, "field value too long: "+fieldType.Name)
		}

		// Set field based on type
		if err := setField(field, valueStr); err != nil {
			return NewHTTPError(http.StatusBadRequest, "invalid value for field "+fieldType.Name+": "+err.Error())
		}
	}

	return nil
}

// bindFiles binds uploaded files to struct fields with security checks
func bindFiles(files map[string][]*multipart.FileHeader, dst interface{}) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return nil
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		fieldType := t.Field(i)

		// Get tag name using helper function
		tag := tagName(fieldType, "form")

		fileHeaders, exists := files[tag]
		if !exists || len(fileHeaders) == 0 {
			continue
		}

		// Security: Check file size to prevent DoS attacks
		for _, header := range fileHeaders {
			if header.Size > maxFileSize {
				return NewHTTPError(http.StatusBadRequest, "file too large: "+header.Filename)
			}
		}

		// Handle *multipart.FileHeader
		if field.Type() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
			field.Set(reflect.ValueOf(fileHeaders[0]))
		}
		// Handle []*multipart.FileHeader for multiple files
		if field.Type() == reflect.TypeOf([]*multipart.FileHeader{}) {
			field.Set(reflect.ValueOf(fileHeaders))
		}
	}

	return nil
}

// setField sets a reflect.Value based on string input
func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		// Check overflow for smaller int types
		if field.OverflowInt(intVal) {
			return NewHTTPError(http.StatusBadRequest, "integer overflow")
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		// Check overflow for smaller uint types
		if field.OverflowUint(uintVal) {
			return NewHTTPError(http.StatusBadRequest, "unsigned integer overflow")
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		// Check overflow for float32
		if field.OverflowFloat(floatVal) {
			return NewHTTPError(http.StatusBadRequest, "float overflow")
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return NewHTTPError(http.StatusBadRequest, "unsupported field type: "+field.Kind().String())
	}
	return nil
}
