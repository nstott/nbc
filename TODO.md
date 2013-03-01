TODO
====

+ add tests,
+ turn into a library, kill the main, change the package name 
+ rewrite genhash in ngram.go
+ change the storage backend to be an interface, and have mongo and memory be two implementations
+ these functions need to touch the storage backend
+ + func (n *nGram) exists() bool
+ + func (n *nGram) GetInstanceCount(class string) int
+ + func GetTotalNGrams(class string) int
+ + func CountDistinctNGrams() int
+ + func (d *Document)DumpToMongo()
+ + func (c *ClassData) Update()
+ + func GetClassProbabilities() map[string]float64


