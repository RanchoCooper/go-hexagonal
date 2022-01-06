package handle

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "go-hexagonal/api/http/dto"
)

/**
 * @author Rancho
 * @date 2022/1/7
 */

func TestCreateExample(t *testing.T) {
    var w = httptest.NewRecorder()
    var response map[string]interface{}
    body := dto.CreateExampleReq{
        Name:  "RanchoCooper",
        Alias: "Rancho",
    }
    b, err := json.Marshal(&body)
    require.NoError(t, err)
    req, _ := http.NewRequest(http.MethodPost, "/example", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")

    NewServerRoute().ServeHTTP(w, req)

    // verify
    assert.Equal(t, http.StatusOK, w.Code)
    err = json.Unmarshal(w.Body.Bytes(), &response)
    assert.Nil(t, err)
    assert.Equal(t, "RanchoCooper", response["name"])
    assert.Equal(t, "Rancho", response["alias"])
}
