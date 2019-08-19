package main

import (
	"github.com/kataras/iris/httptest"
	"testing"
)

// terminal run: go test
func TestMVCHelloWorld(t *testing.T) {
	e := httptest.New(t, newApp())
	e.GET("/").Expect().Status(httptest.StatusOK).ContentType("text/html", "utf-8").Body().Equal("<h1>Welcome</h1>")
}
