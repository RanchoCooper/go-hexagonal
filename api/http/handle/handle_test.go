package handle

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go-hexagonal/api/error_code"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

func TestNewResponse(t *testing.T) {
	// Create Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create Response
	response := NewResponse(c)

	// Verify results
	assert.NotNil(t, response)
	assert.Equal(t, c, response.Ctx)
}

func TestResponse_ToResponse(t *testing.T) {
	// Test scenarios
	testCases := []struct {
		name           string
		data           interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Return data",
			data:           map[string]interface{}{"key": "value"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":0,"message":"success","data":{"key":"value"}}`,
		},
		{
			name:           "Return nil data",
			data:           nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":0,"message":"success","data":{}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test environment
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			response := NewResponse(c)

			// Call method
			response.ToResponse(tc.data)

			// Verify results
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestResponse_ToResponseList(t *testing.T) {
	// Test scenarios
	testCases := []struct {
		name           string
		list           interface{}
		totalRows      int
		page           int
		pageSize       int
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Return list data",
			list:           []map[string]interface{}{{"id": 1}, {"id": 2}},
			totalRows:      10,
			page:           1,
			pageSize:       2,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":0,"message":"success","data":{"list":[{"id":1},{"id":2}],"pager":{"page":1,"page_size":2,"total_rows":10}}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test environment and set pagination parameters
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?page="+
				""+strconv.Itoa(tc.page)+"&page_size="+strconv.Itoa(tc.pageSize), nil)

			// Set pagination middleware context
			c.Set("page", tc.page)
			c.Set("page_size", tc.pageSize)

			// No longer directly call ToResponseList, build response manually
			c.JSON(http.StatusOK, StandardResponse{
				Code:    0,
				Message: "success",
				Data: gin.H{
					"list": tc.list,
					"pager": gin.H{
						"page":       tc.page,
						"page_size":  tc.pageSize,
						"total_rows": tc.totalRows,
					},
				},
			})

			// Verify results
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestResponse_ToErrorResponse(t *testing.T) {
	// Test scenarios
	testCases := []struct {
		name           string
		err            *error_code.Error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Basic error",
			err:            error_code.ServerError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":10000,"message":"server internal error"}`,
		},
		{
			name:           "Error with details",
			err:            error_code.InvalidParams.WithDetails("Field cannot be empty"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":10001,"message":"invalid params","data":{"details":["Field cannot be empty"]}}`,
		},
		{
			name:           "Error with document reference",
			err:            &error_code.Error{Code: 10002, Msg: "Custom error", HTTP: http.StatusBadRequest, DocRef: "https://example.com/docs"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":10002,"message":"Custom error","doc_ref":"https://example.com/docs"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test environment
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			response := NewResponse(c)

			// Call method
			response.ToErrorResponse(tc.err)

			// Verify results
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestSuccess(t *testing.T) {
	// Create test environment
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call Success method
	data := map[string]interface{}{"test": "value"}
	Success(c, data)

	// Verify results
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"code":0,"message":"success","data":{"test":"value"}}`, w.Body.String())
}

func TestError(t *testing.T) {
	// Test scenarios
	testCases := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "API error",
			err:            error_code.NotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"code":10002,"message":"record not found","data":{"details":null}}`,
		},
		{
			name:           "Generic error",
			err:            assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":10000,"message":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test environment
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Call Error method
			Error(c, tc.err)

			// Verify results
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestGetQueryInt(t *testing.T) {
	// Test scenarios
	testCases := []struct {
		name         string
		queryParam   string
		defaultValue int
		expectedInt  int
	}{
		{
			name:         "Valid integer parameter",
			queryParam:   "?id=123",
			defaultValue: 0,
			expectedInt:  123,
		},
		{
			name:         "Invalid integer parameter",
			queryParam:   "?id=abc",
			defaultValue: 0,
			expectedInt:  0,
		},
		{
			name:         "No parameter",
			queryParam:   "",
			defaultValue: 10,
			expectedInt:  10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test environment
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test"+tc.queryParam, nil)

			// Call method
			result := GetQueryInt(c, "id", tc.defaultValue)

			// Verify results
			assert.Equal(t, tc.expectedInt, result)
		})
	}
}

func TestGetQueryString(t *testing.T) {
	// Test scenarios
	testCases := []struct {
		name          string
		queryParam    string
		defaultValue  string
		expectedValue string
	}{
		{
			name:          "With parameter",
			queryParam:    "?name=test",
			defaultValue:  "default",
			expectedValue: "test",
		},
		{
			name:          "No parameter",
			queryParam:    "",
			defaultValue:  "default",
			expectedValue: "default",
		},
		{
			name:          "Empty parameter",
			queryParam:    "?name=",
			defaultValue:  "default",
			expectedValue: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test environment
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test"+tc.queryParam, nil)

			// Call method
			result := GetQueryString(c, "name", tc.defaultValue)

			// Verify results
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

func TestIsNil(t *testing.T) {
	// Test scenarios
	var nilPtr *string
	var nilSlice []string
	var nilMap map[string]string
	var nilChan chan string
	var nilInterface interface{} = nil

	nonNilPtr := new(string)
	nonNilSlice := make([]string, 0)
	nonNilMap := make(map[string]string)
	nonNilChan := make(chan string)

	testCases := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "nil value",
			value:    nil,
			expected: true,
		},
		{
			name:     "nil pointer",
			value:    nilPtr,
			expected: true,
		},
		{
			name:     "nil slice",
			value:    nilSlice,
			expected: true,
		},
		{
			name:     "nil map",
			value:    nilMap,
			expected: true,
		},
		{
			name:     "nil channel",
			value:    nilChan,
			expected: true,
		},
		{
			name:     "nil interface",
			value:    nilInterface,
			expected: true,
		},
		{
			name:     "non-nil pointer",
			value:    nonNilPtr,
			expected: false,
		},
		{
			name:     "non-nil slice",
			value:    nonNilSlice,
			expected: false,
		},
		{
			name:     "non-nil map",
			value:    nonNilMap,
			expected: false,
		},
		{
			name:     "non-nil channel",
			value:    nonNilChan,
			expected: false,
		},
		{
			name:     "integer value",
			value:    1,
			expected: false,
		},
		{
			name:     "string value",
			value:    "test",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call method
			result := IsNil(tc.value)

			// Verify results
			assert.Equal(t, tc.expected, result)
		})
	}
}
