package main

/* an ngram */
type nGram struct {
	length int
	tokens []string
	hash string
	count map[string]int
}

func NewNGram(n int, tokens []string, class string) nGram  {
	return nGram{n, tokens, genhash(tokens), map[string]int{class: 1}}
}

func genhash(in []string) string {
	var ret string
	for _, v := range in {
		ret += " " + v
	}
	return ret
}

func GenerateNGrams(in []string, n int, class string) []nGram {
	out := make([]nGram, 0)
	for i := 0; i <= len(in) - n; i += 1 {
		out = append(out, NewNGram(n, in[i:i+n], class))
	}
	return out
}


func AggregateNGrams(ngrams []nGram, class string) map[string]nGram {
	ret := make(map[string]nGram)
	var mng nGram // declare these here
	var ok bool	 

	for _, v := range ngrams {
		mng, ok = ret[v.hash]
		if ok {
			mng.count[class]++
		} else {
			v.count[class] = 1
			ret[v.hash] = v
		}
	}
	return ret
}