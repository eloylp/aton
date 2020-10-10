package detector

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Kagami/go-face"
)

type Facial interface {
	SaveFaces([]string, []byte) error
	FindFaces([]byte) ([]string, error)
}

type GoFace struct {
	rec *face.Recognizer
	cat map[int32]string
	r   *rand.Rand
}

func NewGoFaceDetector(modelsDir string) (*GoFace, error) {
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("gofacedetector: can't init face recognizer: %w", err)
	}
	d := &GoFace{
		rec: rec,
		cat: map[int32]string{},
		r:   rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
	}
	return d, nil
}

func (d *GoFace) SaveFaces(names []string, bytes []byte) error {
	faces, err := d.rec.Recognize(bytes)
	if err != nil {
		return fmt.Errorf("gofacedetector: can't recognize samples: %w", err)
	}
	if len(faces) != len(names) {
		return fmt.Errorf("gofacedetector: passed faces number (%v) not match with recognized (%v)", len(names), len(faces))
	}
	descriptors := make([]face.Descriptor, len(faces))
	categories := make([]int32, len(faces))
	for i, f := range faces {
		descriptors[i] = f.Descriptor
		categories[i] = d.categoryFromName(names[i])
	}
	d.rec.SetSamples(descriptors, categories)
	return nil
}

func (d *GoFace) categoryFromName(name string) int32 {
	var cat int32
	for cat == 0 || d.catExists(cat) {
		cat = d.r.Int31()
	}
	d.cat[cat] = name
	return cat
}

func (d *GoFace) catExists(cat int32) bool {
	_, ok := d.cat[cat]
	return ok
}

func (d *GoFace) FindFaces(input []byte) ([]string, error) {
	faces, err := d.rec.Recognize(input)
	if err != nil {
		return nil, fmt.Errorf("gofacedetector: can't recognize input: %w", err)
	}
	var results []string
	done := map[string]bool{}
	for _, f := range faces {
		catN := d.rec.Classify(f.Descriptor)
		catName, ok := d.cat[int32(catN)]
		_, duplicated := done[catName]
		if ok && !duplicated {
			results = append(results, catName)
			done[catName] = true
		}
	}
	return results, nil
}
