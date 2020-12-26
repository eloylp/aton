// +build detector

package engine_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/internal/detector/engine"
)

func TestGoFaceDetector(t *testing.T) {
	faceDetector, err := engine.NewGoFaceDetector(ModelsDir)
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
