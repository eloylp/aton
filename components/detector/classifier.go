package detector

type Classifier interface {
	SaveCategories([]string, []byte) error
	FindCategories([]byte) ([]string, error)
}
