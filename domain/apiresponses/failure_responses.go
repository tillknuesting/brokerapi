package apiresponses

import (
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"
)

// FailureResponse can be returned from any of the `ServiceBroker` interface methods
// which allow an error to be returned. Doing so will provide greater control over
// the HTTP response.
type FailureResponse struct {
	error
	statusCode    int
	loggerAction  string
	emptyResponse bool
	errorKey      string
	errorKeyBool  bool // needed for instance_usable and update_repeatable
}

// NewFailureResponse returns a pointer to a new instance of FailureResponse.
// err will by default be used as both a logging message and HTTP response description.
// statusCode is the HTTP status code to be returned, must be 4xx or 5xx
// loggerAction is a short description which will be used as the action if the error is logged.
func NewFailureResponse(err error, statusCode int, loggerAction string) *FailureResponse {
	return &FailureResponse{
		error:        err,
		statusCode:   statusCode,
		loggerAction: loggerAction,
	}
}

// ErrorResponse returns an interface{} which will be JSON encoded and form the body
// of the HTTP response
func (f *FailureResponse) ErrorResponse() interface{} {
	if f.emptyResponse {
		return EmptyResponse{}
	}

	return ErrorResponse{
		Description:      f.error.Error(),
		Error:            f.errorKey,
		InstanceUsable:   f.errorKeyBool,
		UpdateRepeatable: f.errorKeyBool,
	}
}

// ValidatedStatusCode returns the HTTP response status code. If the code is not 4xx
// or 5xx, an InternalServerError will be returned instead.
func (f *FailureResponse) ValidatedStatusCode(logger lager.Logger) int {
	if f.statusCode < 400 || 600 <= f.statusCode {
		if logger != nil {
			logger.Error("validating-status-code", fmt.Errorf("Invalid failure http response code: 600, expected 4xx or 5xx, returning internal server error: 500."))
		}
		return http.StatusInternalServerError
	}
	return f.statusCode
}

// LoggerAction returns the loggerAction, used as the action when logging
func (f *FailureResponse) LoggerAction() string {
	return f.loggerAction
}

// AppendErrorMessage returns an error with the message updated. All other properties are preserved.
func (f *FailureResponse) AppendErrorMessage(msg string) *FailureResponse {
	return &FailureResponse{
		error:         fmt.Errorf("%s %s", f.Error(), msg),
		statusCode:    f.statusCode,
		loggerAction:  f.loggerAction,
		emptyResponse: f.emptyResponse,
		errorKey:      f.errorKey,
	}
}

// FailureResponseBuilder provides a fluent set of methods to build a *FailureResponse.
type FailureResponseBuilder struct {
	error
	statusCode    int
	loggerAction  string
	emptyResponse bool
	errorKey      string
	errorKeyBool  bool // needed for instance_usable and update_repeatable
}

// NewFailureResponseBuilder returns a pointer to a newly instantiated FailureResponseBuilder
// Accepts required arguments to create a FailureResponse.
func NewFailureResponseBuilder(err error, statusCode int, loggerAction string) *FailureResponseBuilder {
	return &FailureResponseBuilder{
		error:         err,
		statusCode:    statusCode,
		loggerAction:  loggerAction,
		emptyResponse: false,
	}
}

// WithErrorKey adds a custom ErrorKey which will be used in FailureResponse to add an `Error`
// field to the JSON HTTP response body
func (f *FailureResponseBuilder) WithErrorKey(errorKey string) *FailureResponseBuilder {
	f.errorKey = errorKey
	return f
}

// WithEmptyResponse will cause the built FailureResponse to return an empty JSON object as the
// HTTP response body
func (f *FailureResponseBuilder) WithEmptyResponse() *FailureResponseBuilder {
	f.emptyResponse = true
	return f
}

// Build returns the generated FailureResponse built using previously configured variables.
func (f *FailureResponseBuilder) Build() *FailureResponse {
	return &FailureResponse{
		error:         f.error,
		statusCode:    f.statusCode,
		loggerAction:  f.loggerAction,
		emptyResponse: f.emptyResponse,
		errorKey:      f.errorKey,
		errorKeyBool:  f.errorKeyBool,
	}
}
