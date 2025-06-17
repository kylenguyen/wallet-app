package restjson_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kylenguyen/wallet-app/pkg/restjson"
)

func Test_ResponseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		code           int
		err            error
		expectedCode   int
		expectedResp   restjson.ErrorResponse
		expectedHeader string
	}{
		{
			name:         "Success With Error",
			code:         http.StatusBadRequest,
			err:          errors.New("bad request"),
			expectedCode: http.StatusBadRequest,
			expectedResp: restjson.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "bad request",
			},
			expectedHeader: "application/json; charset=utf-8",
		},
		{
			name:         "Success Without Error",
			code:         http.StatusNotFound,
			err:          nil,
			expectedCode: http.StatusNotFound,
			expectedResp: restjson.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "not found",
			},
			expectedHeader: "application/json; charset=utf-8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			restjson.ResponseError(c, tc.code, tc.err)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedHeader, w.Header().Get("Content-Type"))

			var resp restjson.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResp, resp)
		})
	}
}

func Test_ResponseData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		data           interface{}
		expectedCode   int
		expectedResp   restjson.Response
		expectedHeader string
	}{
		{
			name: "Success With Data",
			data: map[string]string{
				"foo": "bar",
			},
			expectedCode: http.StatusOK,
			expectedResp: restjson.Response{
				Code: http.StatusOK,
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			expectedHeader: "application/json; charset=utf-8",
		},
		{
			name:         "Success With Nil Data",
			data:         nil,
			expectedCode: http.StatusOK,
			expectedResp: restjson.Response{
				Code: http.StatusOK,
				Data: nil,
			},
			expectedHeader: "application/json; charset=utf-8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			restjson.ResponseData(c, tc.data)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedHeader, w.Header().Get("Content-Type"))

			var resp restjson.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResp, resp)
		})
	}
}
