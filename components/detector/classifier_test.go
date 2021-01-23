// +build detector

package detector_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/components/detector"
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
			resp, err := d.FindCategories(helper.ReadFile(t, c))
			assert.NoError(t, err)
			assert.Equal(t, 1, resp.TotalEntities)
			assert.Equal(t, []string{"bona"}, resp.Matches)
		}
	}
}

func AssertSingleFaceDetectionInGroup(d detector.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona"}, helper.ReadFile(t, faceBona1))
		assert.NoError(t, err)
		resp, err := d.FindCategories(helper.ReadFile(t, groupFaces))
		assert.NoError(t, err)
		assert.Equal(t, 10, resp.TotalEntities)
		assert.Equal(t, []string{"bona"}, resp.Matches)
	}
}

func AssertMultipleFacesDetectionInGroup(d detector.Classifier) func(t *testing.T) {
	return func(t *testing.T) {
		err := d.SaveCategories([]string{"bona", "luda"}, helper.ReadFile(t, groupBonaAndLuda))
		assert.NoError(t, err)
		resp, err := d.FindCategories(helper.ReadFile(t, groupFaces))
		assert.NoError(t, err)
		assert.Equal(t, 10, resp.TotalEntities)
		assert.Equal(t, []string{"luda", "bona"}, resp.Matches)
	}
}
