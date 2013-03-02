package nbc

import (
	"fmt"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
)

type MongoEngine struct {
	host string
	db string
	collection string
	session *mgo.Session
}

func newMongoEngine(host, db, collection string) MongoEngine {
	return MongoEngine{host, db, collection, nil}
}

func (m *MongoEngine) Setup() {
	var err error
	m.session, err = mgo.Dial(m.host)
    if err != nil {
        panic(err)
    }
    m.session.SetSafe(&mgo.Safe{})
    m.indexCollection()
}


func (m *MongoEngine) TearDown() {
	m.session.Close()
}

func (m *MongoEngine) getCollection() *mgo.Collection {
	return m.session.DB(m.db).C(m.collection)
}

func (m *MongoEngine) getClassCollection() *mgo.Collection {
	return m.session.DB(m.db).C(m.collection + "_classes")
}

// indexCollection ensures that the mongo collection for the data is properly indexed
func (m *MongoEngine) indexCollection() {
	col := m.getCollection()
	index := mgo.Index{
		Key: []string{"hash"},
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: true,
	}
	_ = col.EnsureIndex(index)
}


// forgetData clears the data out of the mongo collection
func (m *MongoEngine) forgetData() {
	c := m.getCollection()
	c.RemoveAll(bson.M{})

	c = m.getClassCollection()
	c.RemoveAll(bson.M{})
}

// DumpDocument commits the Document to mongo
func (m *MongoEngine) DumpDocument(d *Document) {
	collection := m.getCollection()	
	field := "cound." + d.class.Name
		
	for _, ngram := range d.ngrams {
		if m.nGramExists(&ngram) {
			q := bson.M{"hash": ngram.Hash}
			err := collection.Update(q, bson.M{"$inc": bson.M{field: 1}})
			if err != nil {
				panic(err)
			}
		} else {
			// straight up insert
			collection.Insert(ngram)
		}
	}

	// update Classification stats
	fmt.Printf("class: %s\n", d.class.Name)
	err := collection.Update(bson.M{"name": d.class.Name}, bson.M{"$inc": bson.M{"count": 1}})
	if (err != nil) {
		err = collection.Insert(&Classification{d.class.Name, 1})
	}

}


/* exists performs a query against mongo to determine if a specific ngram exists in the database or not.
 * this is part of an 'insert on duplicate key update' type situation
 */
func (m *MongoEngine) nGramExists(n *nGram) bool {
	collection := m.getCollection()

	if n.Hash == "" {
		panic("NGram has unitialized Hash") 
	} 

	c, err := collection.Find(bson.M{"hash": n.Hash}).Count()
	if err != nil || c == 0 {
		return false
	}
	return true
}

// GetInstanceCount returns the number of times an ngram has been seen in a class
func (m *MongoEngine) GetInstanceCount(n *nGram, class string) int {
	collection := m.getCollection()
	var ngram nGram
	err := collection.Find(bson.M{"hash": n.Hash}).One(&ngram)
	if err != nil {
		return 0
	}
	return ngram.Count[class]
}

/* 
 * CountDistinctNGrams returns the size of our vocabulary.  
 * i.e. the number of distinct ngrams across all documents and classes
 */ 
func (m *MongoEngine) CountDistinctNGrams() int {
	collection := m.getCollection()
	count, err := collection.Find(bson.M{}).Count()
	if err != nil {
		panic(err)
	}
	return count
}

// getTotalNGrams returns the total number of times a specific ngram has been seen in a class across all documents
func (m *MongoEngine) GetTotalNGrams(class string) int {
	collection := m.getCollection()
	var field = "count." + class
	job := mgo.MapReduce{
        Map:      "function() { emit(\"total\", this.count."+class+")}",
        Reduce:   "function(key, values) { var t = 0; values.forEach(function (i) {t += i});return t; }",
	}
	var result []struct { Id string "_id"; Value int }
	q := collection.Find(bson.M{field: bson.M{"$gt": 0}})
	_, err := q.MapReduce(job, &result)
	if err != nil  {
	    panic(err)
	}
	if len(result) > 0 {
		return result[0].Value
	}
	return 0;
}

func (m *MongoEngine) GetClassProbabilities() map[string]float64 {
	collection := m.getClassCollection()
	var result Classification

	counts := make(map[string]int)
	var total int

	iter := collection.Find(bson.M{}).Limit(100).Iter()
	for iter.Next(&result) {
		total += result.Count
		counts[result.Name] = result.Count
    }	

	classCount := len(counts)
	probabilities := make(map[string]float64)

	for k, v := range counts {
		probabilities[k] = laplaceSmoothing(v, total, classCount)
	} 
	return probabilities
}