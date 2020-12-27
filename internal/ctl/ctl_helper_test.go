// +build integration

package ctl_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
)

func testLogger(output io.Writer) *logrus.Logger {
	l := logrus.New()
	l.SetOutput(output)
	return l
}

func fetchResource(t *testing.T, s string) []byte {
	t.Helper()
	resp, err := http.Get(s)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
