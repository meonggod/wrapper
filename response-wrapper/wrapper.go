package response_wrapper

import (
	"context"
	"encoding/json"
	"fmt"
	error_wrapper "github.com/meonggod/wrapper/error-wrapper"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Header Header      `json:"header,omitempty"`
	Body   interface{} `json:"body,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
}

type Header struct {
	ProcessTime string `json:"process_time"`
	IsSuccess   bool   `json:"is_success"`
}

var (
	defaultCategory   = error_wrapper.NewCategory(500, "Errors at New Errors")
	defaultDefinition = error_wrapper.NewDefinition(0, "an Errors", true, defaultCategory)
	defaultErr        = error_wrapper.New(defaultDefinition, "an Errors")
)

func New(w *http.ResponseWriter, context context.Context, isSuccess bool, body interface{}, errW *error_wrapper.ErrorWrapper) {
	var (
		errors interface{}
	)
	now := time.Now().UnixNano() / int64(time.Millisecond)

	if errW == nil {
		writeHeader(*w, 200)
	} else {
		writeHeader(*w, errW.StatusCode())
	}

	startTime, err := strconv.Atoi(fmt.Sprint(context.Value("start_time")))
	requestId := fmt.Sprint(context.Value("request_id"))
	if err != nil {
		New(w, context, false, nil, defaultErr)
		return
	}
	start := int64(startTime)
	start /= int64(time.Millisecond)

	header := Header{
		ProcessTime: fmt.Sprintf(`%v ms`, now-start),
		IsSuccess:   isSuccess,
	}

	if errW != nil {
		errors = [3]string{
			fmt.Sprintf(`%s (%d) (%s)`, errW.Error(), errW.Code(), requestId),
			fmt.Sprintf(`%s (%d)`, errW.Error(), errW.Code()),
			fmt.Sprintf(`%s`, errW.ActualError()),
		}
	}

	json.NewEncoder(*w).Encode(Response{
		Header: header,
		Body:   body,
		Errors: errors,
	})

	return
}

func writeHeader(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}
