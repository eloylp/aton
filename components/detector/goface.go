package detector

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Kagami/go-face"
)

type GoFace struct {
	rec *face.Recognizer
	cat map[int32]string
	r   *rand.Rand
}

func NewGoFace(modelsDir string) (*GoFace, error) {
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

func (d *GoFace) SaveCategories(names []string, bytes []byte) error {
	dup := duplicates(names)
	if len(dup) > 0 {
		return fmt.Errorf("gofacedetector: duplicated names: %s", strings.Join(dup, ","))
	}
	faces, err := d.rec.Recognize(bytes)
	if err != nil {
		return fmt.Errorf("gofacedetector: can't recognize samples: %w", err)
	}
	if len(faces) != len(names) {
		return fmt.Errorf("gofacedetector: passed faces number (%v) not match with recognized (%v)", len(names), len(faces))
	}
	descriptors := make([]face.Descriptor, len(faces))
	categories := make([]int32, len(faces))
	for i := 0; i < len(faces); i++ {
		descriptors[i] = faces[i].Descriptor
		categories[i] = d.categoryFromName(names[i])
	}
	d.rec.SetSamples(descriptors, categories)
	return nil
}

func duplicates(names []string) []string {
	m := map[string]struct{}{}
	var duplicates []string
	for i := 0; i < len(names); i++ {
		name := names[i]
		if _, ok := m[name]; ok {
			duplicates = append(duplicates, name)
		}
		m[name] = struct{}{}
	}
	return duplicates
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

func (d *GoFace) FindCategories(input []byte) ([]string, error) {
	faces, err := d.rec.Recognize(input)
	if err != nil {
		return nil, fmt.Errorf("gofacedetector: can't recognize input: %w", err)
	}
	var results []string
	done := map[string]struct{}{}
	for i := 0; i < len(faces); i++ {
		catN := d.rec.Classify(faces[i].Descriptor)
		catName, ok := d.cat[int32(catN)]
		_, duplicated := done[catName]
		if ok && !duplicated {
			results = append(results, catName)
			done[catName] = struct{}{}
		}
	}
	return results, nil
}
