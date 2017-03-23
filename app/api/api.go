package api

import (
	"encoding/json"
	"net/http"

	"github.com/jloup/utils"
	"github.com/labstack/echo"
)

var log = utils.StandardL().WithField("module", "api")

// add go generate for stringify
type ApiStatus uint32
type ApiErrorCode uint32

const (
	OK ApiStatus = iota
	ERROR

	INTERNAL_ERROR ApiErrorCode = iota
	INVALID_REQUEST
	NOT_FOUND
)

func (a ApiErrorCode) String() string {
	switch a {
	case INTERNAL_ERROR:
		return "INTERNAL_ERROR"
	case INVALID_REQUEST:
		return "INVALID_REQUEST"
	case NOT_FOUND:
		return "NOT_FOUND"
	}

	return ""
}

func ApiErrorToHttpStatus(code ApiErrorCode) int {
	switch code {
	case INVALID_REQUEST, NOT_FOUND:
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

type Error struct {
	errType ApiErrorCode
	err     error
}

func (e Error) Err() error {
	if e.err != nil {
		return e
	}

	return nil
}

func (e Error) ErrType() ApiErrorCode {
	return e.errType
}

func (e Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return ""
}

type ApiResponse struct {
	Status   ApiStatus       `json:"status"`
	Response json.RawMessage `json:"res"`
}

type ApiErrorResponse struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

func ApiError(code ApiErrorCode) *ApiResponse {
	msgResp := ApiErrorResponse{code.String(), ""}
	j, err := json.Marshal(msgResp)
	if err != nil {
		panic(err.Error())
	}

	return &ApiResponse{
		Status:   ERROR,
		Response: j,
	}
}

func ApiSuccess(data []byte) *ApiResponse {
	return &ApiResponse{
		Status:   OK,
		Response: data,
	}
}

func WriteResponse(httpCode int, res *ApiResponse, c echo.Context) error {
	resp := c.Response()

	resp.Header().Set("Content-Type", "application/json;charset=utf-8")
	resp.WriteHeader(httpCode)
	j, err := json.Marshal(res)
	if err != nil {
		return err
	}

	resp.Write(j)
	return nil
}

func WriteResponseRaw(httpCode int, res []byte, contentType string, c echo.Context) error {
	resp := c.Response()

	resp.Header().Set("Content-Type", contentType)
	resp.WriteHeader(httpCode)
	resp.Write(res)
	return nil
}

type ApiQuery interface {
	Run()
	Err() error
	GetRawOut() interface{}
	GetJSONOut() json.RawMessage
}

func RunApiQuery(c echo.Context, query ApiQuery) error {
	query.Run()
	b := query.GetJSONOut()

	if query.Err() != nil {
		log.Errorf("%v", query.Err())
		if err, ok := query.Err().(Error); ok {
			return WriteResponse(ApiErrorToHttpStatus(err.ErrType()), ApiError(err.ErrType()), c)
		} else {
			return WriteResponse(http.StatusInternalServerError, ApiError(INTERNAL_ERROR), c)
		}
	}

	return WriteResponse(http.StatusOK, ApiSuccess(b), c)
}

func RunApiQueryResponseRaw(c echo.Context, query ApiQuery) error {
	query.Run()

	if query.Err() != nil {
		log.Errorf("%v", query.Err())
		if err, ok := query.Err().(Error); ok {
			return WriteResponse(ApiErrorToHttpStatus(err.ErrType()), ApiError(err.ErrType()), c)
		} else {
			return WriteResponse(http.StatusInternalServerError, ApiError(INTERNAL_ERROR), c)
		}
	}

	return WriteResponseRaw(http.StatusOK, []byte(query.GetRawOut().(string)), "text/plain;charset=utf-8", c)
}
