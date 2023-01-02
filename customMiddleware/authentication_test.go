package customMiddleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	h := middleware.BasicAuth(Authentication)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	t.Run("should return nil when username = admin & password = admin", func(t *testing.T) {
		auth := "Basic" + " " + base64.StdEncoding.EncodeToString([]byte("admin:admin"))
		req.Header.Set(echo.HeaderAuthorization, auth)
		err := h(c)

		assert.Nil(t, err)
	})

	t.Run("should return ", func(t *testing.T) {
		auth := "Basic" + " " + base64.StdEncoding.EncodeToString([]byte("admin:wrongpassword"))
		req.Header.Set(echo.HeaderAuthorization, auth)
		he := h(c).(*echo.HTTPError)

		assert.Equal(t, http.StatusUnauthorized, he.Code)
		assert.Equal(t, "code=401, message=Unauthorized", he.Error())
	})
}
