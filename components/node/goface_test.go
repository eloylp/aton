// +build node

package node_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/components/node"
)

func TestGoFaceDetector(t *testing.T) {
	faceDetector, err := node.NewGoFace(ModelsDir)
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
