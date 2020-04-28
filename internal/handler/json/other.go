package json

import (
	"errors"
	"log"
	"regexp"
	"smarthome-home/internal"
	"smarthome-home/internal/domain/errormessage"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type Answer struct {
	CodeToError    map[string]errormessage.Error
	CodeToCode     map[string]int
	jsonStreamPool jsoniter.StreamPool
}

type StatusOKAnswer struct {
	Data interface{} `json:"data"`
}

//nolint will be long bc off error messages
func InitAnswer(jsp jsoniter.StreamPool) *Answer {
	err := Answer{
		jsonStreamPool: jsp,
		CodeToError:    map[string]errormessage.Error{},
		CodeToCode:     map[string]int{},
	}
	return &err
}

func (a *Answer) SetAnswer(ctx *fasthttp.RequestCtx, data ...interface{}) {
	if err, ok := data[0].(error); ok {
		a.setValidationErrors(ctx, err)
		return
	}

	if data[0] == nil {
		setStatusNoContent(ctx)
		return
	}
	a.setStatusOK(ctx, data...)
}

func setStatusNoContent(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

func (a *Answer) setStatusOK(ctx *fasthttp.RequestCtx, data ...interface{}) {
	jsonStream := a.jsonStreamPool.BorrowStream(nil)
	defer a.jsonStreamPool.ReturnStream(jsonStream)
	status := StatusOKAnswer{
		Data: data[0],
	}
	jsonStream.WriteVal(status)

	ctx.SetBody(jsonStream.Buffer())
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (a *Answer) setValidationErrors(ctx *fasthttp.RequestCtx, err error) {
	requestID := getRequestID(ctx)

	errors := &errormessage.Errors{
		ID: requestID,
	}

	errors.Errors, errors.StatusCode = a.makeErrors(err.Error())

	if errors.StatusCode == fasthttp.StatusUnprocessableEntity {
		jsonStream := a.jsonStreamPool.BorrowStream(nil)
		defer a.jsonStreamPool.ReturnStream(jsonStream)
		jsonStream.WriteVal(errors)

		ctx.SetBody(jsonStream.Buffer())
	}
	if errors.StatusCode >= 500 {
		log.Println(err.Error())
	}
	ctx.SetStatusCode(errors.StatusCode)
}

// TODO fix this mess
func (a *Answer) makeErrors(err string) (e *[]errormessage.Error, statusCode int) {
	splittedError := SplitErr(err)
	errors := make([]errormessage.Error, len(splittedError))
	statusCode = fasthttp.StatusInternalServerError
	for k, v := range splittedError {
		key, value := getKeyValue(v)
		errors[k] = a.getError(value, errormessage.Source{
			Pointer: key,
		})
		statusCode = a.getStatusCode(value)
	}
	return &errors, statusCode
}

func (a *Answer) getStatusCode(errCode string) int {
	// Check if code exists otherwise return 500
	if code, ok := a.CodeToCode[errCode]; ok {
		return code
	}
	return fasthttp.StatusInternalServerError
}

//nolint is global otherwise it will get calculated every time
var regexpForSplitting = regexp.MustCompile(`(\(|\)|\;)`)

func SplitErr(err string) []string {
	s := err
	s = strings.ReplaceAll(s, ".", ";")
	s = strings.ReplaceAll(s, " ", "")

	a := []string{}
	head := ""
	for _, value := range regexpForSplitting.Split(s, -1) {
		v3 := strings.Split(value, ":")
		if len(v3) > 1 {
			if v3[1] == "" {
				head = v3[0]
			} else {
				if head == "" {
					a = append(a, strings.ToLower(internal.Concat(v3[0], ":", v3[1])))
				} else {
					a = append(a, strings.ToLower(internal.Concat(head, "/", v3[0], ":", v3[1])))
				}
			}
		} else {
			head = ""
		}
	}
	return a
}

func (a *Answer) getError(code string, source errormessage.Source) errormessage.Error {
	err := a.CodeToError[code]
	return errormessage.Error{
		Title:  err.Title,
		Detail: err.Detail,
		Source: source,
	}
}

func getKeyValue(v string) (key, value string) {
	data := strings.Split(v, ":")

	// removing artifacts from err.error()
	if len(data) < 2 {
		return data[0], ""
	}
	key = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(data[0], ".", ""), " ", ""))
	value = strings.ReplaceAll(strings.ReplaceAll(data[1], ".", ""), " ", "")
	return key, value
}

func getRequestID(ctx *fasthttp.RequestCtx) []byte {
	var requestID []byte
	var ok bool

	if requestID, ok = ctx.UserValue("requestId").([]byte); !ok {
		requestID = []byte("\"request ID not found\"")
	}
	return requestID
}

func getID(ctx *fasthttp.RequestCtx, parameter string) (id int, err error) {
	parameterID, ok := ctx.UserValue(parameter).(string)
	if !ok {
		return 0, errors.New("interface error")
	}
	if id, err = strconv.Atoi(parameterID); err != nil {
		log.Println(err)
		return 0, errors.New("not a number")
	}
	return
}
