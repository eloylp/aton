package detector_test

import (
	"io/ioutil"
	"log"
	"testing"
)

func readFile(t *testing.T, file string) []byte {
	t.Helper()
	face1, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return face1
}
