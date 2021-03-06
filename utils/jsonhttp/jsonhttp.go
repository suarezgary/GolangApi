package jsonhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

// APIResponse contains the attributes found in an API response
type APIResponse struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Debug   string      `json:"debug,omitempty"`
}

// CheckableRequest defines an interface for request payloads that can be checked with the jsonhttp checker. See JSONDecodeAndCatchForAPI for the usage
type CheckableRequest interface {
	Parameters() error
}

// JSONSuccess returns a successful APIResponse on the http response with the provided parameters
func JSONSuccess(w http.ResponseWriter, data interface{}, message string) {
	if message == "" {
		message = "ok"
	}
	resp := APIResponse{
		Message: message,
		Success: true,
		Data:    data,
	}
	JSONWriter(w, resp, http.StatusOK)
}

// JSONSuccessNoContent returns a successful APIResponse
func JSONSuccessNoContent(w http.ResponseWriter) {
	message := "Record Not Found"
	resp := APIResponse{
		Message: message,
		Success: true,
	}
	JSONWriter(w, resp, http.StatusNoContent)
}

// JSONError returns an APIResponse on the http response with the provided parameters and status code
func JSONError(w http.ResponseWriter, data interface{}, message string, debug string, statusCode int) {
	if message == "" {
		message = "error"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    data,
		Debug:   debug,
	}
	JSONWriter(w, resp, statusCode)
}

// JSONInternalError returns an internal server error APIResponse on the http response with the provided parameters
func JSONInternalError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "error"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	JSONWriter(w, resp, http.StatusInternalServerError)
}

// JSONBadRequestError returns a bad request error APIResponse on the http response with the provided parameters
func JSONBadRequestError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "bad_request"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	JSONWriter(w, resp, http.StatusBadRequest)
}

// JSONNotFoundError returns a not found error APIResponse on the http response with the provided parameters
func JSONNotFoundError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "not_found"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	JSONWriter(w, resp, http.StatusNotFound)
}

// JSONForbiddenError returns a forbidden error APIResponse on the http response with the provided parameters
func JSONForbiddenError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "forbidden"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	JSONWriter(w, resp, http.StatusForbidden)
}

// JSONDetailed returns the provided APIResponse on the http response with the provided HTTP status code
func JSONDetailed(w http.ResponseWriter, resp APIResponse, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	JSONWriter(w, resp, statusCode)
}

// JSONWriter provides a wrapper function to marshal an interface{} type to JSON and then send the bytes back over an http.ResponseWriter
func JSONWriter(w http.ResponseWriter, payload interface{}, statusCode int) {
	//dj, err := json.MarshalIndent(payload, "", "  ")
	dj, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s", dj)
}

// JSONDecodeAndCatchForAPI is the primary function for decoding checkable (and non-checkable) payloads into structs. If the struct passed into `outStruct` satisfied the `CheckableRequest` interface, the check will also be run after decoding the JSON
func JSONDecodeAndCatchForAPI(w http.ResponseWriter, r *http.Request, outStruct interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&outStruct)
	if err != nil {
		JSONBadRequestError(w, "Invalid JSON", "")
		return err
	}
	if !isCheckableRequest(outStruct) {
		return nil
	}
	method := reflect.ValueOf(outStruct).MethodByName("Parameters").Interface().(func() error)
	err = method()
	if err != nil {
		JSONBadRequestError(w, "", err.Error())
		return err
	}
	return nil
}

func isCheckableRequest(checkAgainst interface{}) bool {
	reader := reflect.TypeOf((*CheckableRequest)(nil)).Elem()
	return reflect.TypeOf(checkAgainst).Implements(reader)
}

//DecodeJSONBody DecodeJSONBody
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			return errors.New("Content-Type header is not application/json")
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return errors.New(msg)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return errors.New(msg)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return errors.New(msg)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return errors.New(msg)

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return errors.New(msg)

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return errors.New(msg)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return errors.New(msg)
	}

	return nil
}

//DecodeJSONFromFormValue DecodeJSONFromFormValue
func DecodeJSONFromFormValue(w http.ResponseWriter, r *http.Request, dst interface{}, formKey string) error {
	keyValue := r.FormValue(formKey)
	if keyValue == "" {
		return errors.New("Error Getting Value from Form")
	}

	err := json.Unmarshal([]byte(keyValue), &dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return errors.New(msg)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return errors.New(msg)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return errors.New(msg)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return errors.New(msg)

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return errors.New(msg)

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return errors.New(msg)

		default:
			return err
		}
	}

	return nil
}
