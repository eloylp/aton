// +build integration

package detector_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/internal/detector"
)

var (
	ModelsDir        = "../../models"
	imagesDir        = "../../images"
	groupFaces       = filepath.Join(imagesDir, "pristin.jpg")
	faceBona1        = filepath.Join(imagesDir, "bona.jpg")
	faceBona2        = filepath.Join(imagesDir, "bona2.jpg")
	faceBona3        = filepath.Join(imagesDir, "bona3.jpg")
	faceBona4        = filepath.Join(imagesDir, "bona4.jpg")
	groupBonaAndLuda = filepath.Join(imagesDir, "bonaAndLuda.jpg")
)

func TestFaceDetectors(t *testing.T) {
	faceDetector, err := detector.NewGoFaceDetector(ModelsDir)
	assert.NoError(t, err)
	t.Run("Testing GoFace face detector",
		AssertSingleFaceDetection(faceDetector))
	t.Run("Testing GoFace face detector with group",
		AssertSingleFaceDetectionInGroup(faceDetector))
	t.Run("Testing GoFace multiple face detector with group",
		AssertMultipleFacesDetectionInGroup(faceDetector))
}

func AssertSingleFaceDetection(d detector.Facial) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveFaces([]string{"bona"}, readFile(t, faceBona1))
		assert.NoError(t, err)
		for _, c := range []string{faceBona2, faceBona3, faceBona4} {
			faces, err := d.FindFaces(readFile(t, c))
			assert.NoError(t, err)
			assert.Equal(t, []string{"bona"}, faces)
		}
	}
}

func AssertSingleFaceDetectionInGroup(d detector.Facial) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveFaces([]string{"bona"}, readFile(t, faceBona1))
		assert.NoError(t, err)
		faces, err := d.FindFaces(readFile(t, groupFaces))
		assert.NoError(t, err)
		assert.Equal(t, []string{"bona"}, faces)
	}
}

func AssertMultipleFacesDetectionInGroup(d detector.Facial) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveFaces([]string{"bona", "luda"}, readFile(t, groupBonaAndLuda))
		assert.NoError(t, err)
		faces, err := d.FindFaces(readFile(t, groupFaces))
		assert.NoError(t, err)
		assert.Equal(t, []string{"luda", "bona"}, faces)
	}
}
