package controllers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func mockGin() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// test request, must instantiate a request first
	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header), // if you need to test headers
	}
	// finally set the request to the gin context
	c.Request = req

	return c, w
}

func TestHomeHandler(t *testing.T) {
	c, r := mockGin()
	HomeHandler()(c)
	gotStatus := r.Result().Status
	expectedStatus := "200 OK"

	if gotStatus != expectedStatus {
		t.Error("got error on request homeHandler func")
	}

	gotBody := r.Body.String()
	expectedBody := "FeatWS Ruller Works!!!"

	if gotBody != expectedBody {
		t.Error("got error on request homeHandler func")
	}

}
