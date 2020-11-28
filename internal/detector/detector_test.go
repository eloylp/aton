// +build detector

package detector_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/internal/detector"
	"github.com/eloylp/aton/pkg/test/helper"
)

var (
	ModelsDir        = "../../models"
	imagesDir        = "../../testdata/images"
	groupFaces       = filepath.Join(imagesDir, "pristin.jpg")
	faceBona1        = filepath.Join(imagesDir, "bona.jpg")
	faceBona2        = filepath.Join(imagesDir, "bona2.jpg")
	faceBona3        = filepath.Join(imagesDir, "bona3.jpg")
	faceBona4        = filepath.Join(imagesDir, "bona4.jpg")
	groupBonaAndLuda = filepath.Join(imagesDir, "bonaAndLuda.jpg")
)

func TestGoFaceDetector(t *testing.T) {
	faceDetector, err := detector.NewGoFaceDetector(ModelsDir)
	assert.NoError(t, err)
	t.Run("Error if duplicated names",
		AssertErrorIfDuplicatedNames(faceDetector))
	t.Run("Error if initial samples and names number mismatch",
		AssertErrorIfNotAllFacesRecognized(faceDetector))
	t.Run("Detect one face",
		AssertSingleFaceDetection(faceDetector))
	t.Run("Detect single face in group",
		AssertSingleFaceDetectionInGroup(faceDetector))
	t.Run("Detect multiple faces within group",
		AssertMultipleFacesDetectionInGroup(faceDetector))
}

func AssertErrorIfDuplicatedNames(d detector.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona", "luda", "bona_dep2", "luda", "bona"}, helper.ReadFile(t, faceBona1))
		fmt.Println(err)
		assert.EqualError(t, err, "gofacedetector: duplicated names: luda,bona")
	}
}

func AssertErrorIfNotAllFacesRecognized(d detector.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona", "EXTRA_NON_EXISTENT_FACE"}, helper.ReadFile(t, faceBona1))
		assert.EqualError(t, err, "gofacedetector: passed faces number (2) not match with recognized (1)")
	}
}

func AssertSingleFaceDetection(d detector.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona"}, helper.ReadFile(t, faceBona1))
		assert.NoError(t, err)
		for _, c := range []string{faceBona2, faceBona3, faceBona4} {
			faces, err := d.FindCategories(helper.ReadFile(t, c))
			assert.NoError(t, err)
			assert.Equal(t, []string{"bona"}, faces)
		}
	}
}

func AssertSingleFaceDetectionInGroup(d detector.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona"}, helper.ReadFile(t, faceBona1))
		assert.NoError(t, err)
		faces, err := d.FindCategories(helper.ReadFile(t, groupFaces))
		assert.NoError(t, err)
		assert.Equal(t, []string{"bona"}, faces)
	}
}

func AssertMultipleFacesDetectionInGroup(d detector.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona", "luda"}, helper.ReadFile(t, groupBonaAndLuda))
		assert.NoError(t, err)
		faces, err := d.FindCategories(helper.ReadFile(t, groupFaces))
		assert.NoError(t, err)
		assert.Equal(t, []string{"luda", "bona"}, faces)
	}
}
