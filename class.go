package nbc

type Classification struct {
	Name string
	Count int
}

func NewClassification(name string) Classification {
	return Classification{name, 0}
}