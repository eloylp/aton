// +build integration

package detector_test

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/internal/detector"
)

var (
	ModelsDir = "../../models"
	imagesDir = "../../images"
	faceBona1 = filepath.Join(imagesDir, "bona.jpg")
	faceBona2 = filepath.Join(imagesDir, "bona2.jpg")
	faceBona3 = filepath.Join(imagesDir, "bona3.jpg")
	faceBona4 = filepath.Join(imagesDir, "bona4.jpg")
)

func TestFaceDetectors(t *testing.T) {
	faceDetector, err := detector.NewDLIBFaceDetector(ModelsDir)
	assert.NoError(t, err)
	t.Run("Testing DLIB face detector",
		AssertFacialDetection(faceDetector))
}

func AssertFacialDetection(d detector.Facial) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveFace("bona", readFile(t, faceBona1))
		assert.NoError(t, err)
		for _, c := range []string{faceBona2, faceBona3, faceBona4} {
			face, err := d.FindFace(readFile(t, c))
			assert.NoError(t, err)
			assert.Equal(t, "bona", face)
		}
	}
}

func readFile(t *testing.T, file string) []byte {
	t.Helper()
	face1, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return face1
}
