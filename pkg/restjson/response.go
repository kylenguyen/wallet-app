// Package rest provides functions to format API responses, including data and errors, according to company conventions.
package restjson

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// ResponseError formats an error response with the specified HTTP status code and error message.
func ResponseError(c *gin.Context, code int, err error) {
	resp := ErrorResponse{
		Code:    code,
		Message: strings.ToLower(http.StatusText(code)),
	}

	if err != nil {
		resp.Message = err.Error()
	}

	c.JSON(code, resp)
}

// ResponseData formats a successful response with the default HTTP code (200 OK) and the provided data.
func ResponseData(c *gin.Context, data interface{}) {
	resp := Response{
		Code: http.StatusOK,
		Data: data,
	}

	c.JSON(resp.Code, resp)
}
