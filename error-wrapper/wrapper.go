package error_wrapper

import (
	"fmt"
	"runtime"
)

type ErrorWrapper struct {
	definition *ErrorDefinition
	stacktrace string

	args []interface{}
}

func New(definition *ErrorDefinition, err interface{}) *ErrorWrapper {
	var msg string

	_, file, no, ok := runtime.Caller(1)
	if ok {
		msg = fmt.Sprintf(`%s:%d => [%s] %s`, file, no, msg, err)
	}

	return &ErrorWrapper{
		definition: definition,
		stacktrace: msg,
	}
}

func (e *ErrorWrapper) Error() string {
	return e.definition.Error(e.args)
}

func (e *ErrorWrapper) Is(definition *ErrorDefinition) bool {
	if e.definition == definition {
		return true
	}
	return false
}

func (e *ErrorWrapper) StackTrace() string {
	return e.stacktrace
}

func (e *ErrorWrapper) ActualError() string {
	return e.definition.ActualError(e.args)
}

func (e *ErrorWrapper) Code() int {
	return e.definition.Code()
}

func (e *ErrorWrapper) StatusCode() int {
	return e.definition.category.StatusCode()
}

func (e *ErrorWrapper) IsMasked() bool {
	return e.definition.IsMasked()
}

func (e *ErrorWrapper) With(data string) *ErrorWrapper {
	e.args = append(e.args, data)
	return e
}

func (e *ErrorWrapper) Wrap(msgs ...string) {
	var msg string

	for _, m := range msgs {
		msg = fmt.Sprintf("%s %s", msg, m)
	}

	_, file, no, ok := runtime.Caller(1)
	if ok {
		msg = fmt.Sprintf(`%s:%d => [%s] %s`, file, no, msg, e.stacktrace)
	}

	e.stacktrace = msg
}
