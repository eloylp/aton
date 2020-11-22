package helper

import (
	"io/ioutil"
	"log"
	"testing"
)

func ReadFile(t *testing.T, file string) []byte {
	t.Helper()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
