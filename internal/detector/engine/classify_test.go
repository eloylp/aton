// +build detector

package engine_test

import (
	"fmt"
	"github.com/eloylp/aton/internal/detector/engine"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/pkg/test/helper"
)

var (
	ModelsDir        = "../../models"
	imagesDir        = "../../samples/images"
	groupFaces       = filepath.Join(imagesDir, "pristin.jpg")
	faceBona1        = filepath.Join(imagesDir, "bona.jpg")
	faceBona2        = filepath.Join(imagesDir, "bona2.jpg")
	faceBona3        = filepath.Join(imagesDir, "bona3.jpg")
	faceBona4        = filepath.Join(imagesDir, "bona4.jpg")
	groupBonaAndLuda = filepath.Join(imagesDir, "bonaAndLuda.jpg")
)

func AssertErrorIfDuplicatedNames(d engine.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona", "luda", "bona_dep2", "luda", "bona"}, helper.ReadFile(t, faceBona1))
		fmt.Println(err)
		assert.EqualError(t, err, "gofacedetector: duplicated names: luda,bona")
	}
}

func AssertErrorIfNotAllFacesRecognized(d engine.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona", "EXTRA_NON_EXISTENT_FACE"}, helper.ReadFile(t, faceBona1))
		assert.EqualError(t, err, "gofacedetector: passed faces number (2) not match with recognized (1)")
	}
}

func AssertSingleFaceDetection(d engine.Classifier) func(t *testing.T) {
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

func AssertSingleFaceDetectionInGroup(d engine.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona"}, helper.ReadFile(t, faceBona1))
		assert.NoError(t, err)
		faces, err := d.FindCategories(helper.ReadFile(t, groupFaces))
		assert.NoError(t, err)
		assert.Equal(t, []string{"bona"}, faces)
	}
}

func AssertMultipleFacesDetectionInGroup(d engine.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona", "luda"}, helper.ReadFile(t, groupBonaAndLuda))
		assert.NoError(t, err)
		faces, err := d.FindCategories(helper.ReadFile(t, groupFaces))
		assert.NoError(t, err)
		assert.Equal(t, []string{"luda", "bona"}, faces)
	}
}
