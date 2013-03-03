package nbc

type Classification struct {
	Name string
	ngrams map[string]*nGram
	Count int
}

func NewClassification(name string) Classification {
	return Classification{name, make(map[string]*nGram), 0}
}