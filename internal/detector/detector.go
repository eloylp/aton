package detector

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Kagami/go-face"
)

type Facial interface {
	SaveFace(string, []byte) error
	FindFaces([]byte) ([]string, error)
}

type GoFace struct {
	rec *face.Recognizer
	cat map[int32]string
	r   *rand.Rand
}

func NewDLIBFaceDetector(modelsDir string) (*GoFace, error) {
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("dlibfacerecognizer: can't init face recognizer: %w", err)
	}
	d := &GoFace{
		rec: rec,
		cat: map[int32]string{},
		r:   rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
	}
	return d, nil
}

func (d *GoFace) SaveFace(name string, bytes []byte) error {
	f, err := d.rec.Recognize(bytes)
	if err != nil {
		return err
	}
	descriptors := []face.Descriptor{f[0].Descriptor}
	d.rec.SetSamples(descriptors, []int32{d.categoryFromName(name)})
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
		return nil, err
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
