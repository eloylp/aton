package engine

type Classifier interface {
	SaveCategories([]string, []byte) error
	FindCategories([]byte) ([]string, error)
}
