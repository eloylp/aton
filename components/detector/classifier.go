package detector

type Classifier interface {
	SaveCategories([]string, []byte) error
	FindCategories([]byte) (*FoundCategories, error)
}

type FoundCategories struct {
	Matches       []string
	DetectedCount int
}
