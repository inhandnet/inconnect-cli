package api

import (
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
)

// BackendError is the typed form of the "HTTP 200 + body error" pattern, so
// callers can branch on the specific error_code (e.g. 10001, which the backend
// returns even after a resource was already persisted).
type BackendError struct {
	Message string
	Code    string
}

func (e *BackendError) Error() string {
	return fmt.Sprintf("%s (error_code %s)", e.Message, e.Code)
}

// ErrorCode returns the backend error_code carried by err, or "" if err is not
// a BackendError.
func ErrorCode(err error) string {
	var be *BackendError
	if errors.As(err, &be) {
		return be.Code
	}
	return ""
}

// BodyError detects the backend's "HTTP 200 + body error" pattern. The shared
// common-lib exception handler wraps business/internal ErrorCodeException into a
// 200 response with a top-level {"error","error_code"} body (only auth errors map
// to real 4xx). We treat a top-level non-empty `error` accompanied by `error_code`
// as a failure. Task-level errors live under result.error (nested) and are left
// alone, since they carry expected semantics like timeouts.
func BodyError(body []byte) error {
	msg := gjson.GetBytes(body, "error")
	code := gjson.GetBytes(body, "error_code")
	if msg.Exists() && msg.String() != "" && code.Exists() {
		return &BackendError{Message: msg.String(), Code: code.String()}
	}
	return nil
}

func ResultID(body []byte) string {
	return firstNonEmpty(body, "result._id", "result.id", "_id", "id")
}

// ResultIDName extracts the id and name from an API response, tolerating both
// result-wrapped objects ({"result":{...}}, e.g. update) and bare objects
// ({...}, e.g. create), and both `_id` and `id` field names.
func ResultIDName(body []byte) (id, name string) {
	id = firstNonEmpty(body, "result._id", "result.id", "_id", "id")
	name = firstNonEmpty(body, "result.name", "name")
	return id, name
}

func firstNonEmpty(body []byte, paths ...string) string {
	for _, p := range paths {
		if v := gjson.GetBytes(body, p).String(); v != "" {
			return v
		}
	}
	return ""
}
