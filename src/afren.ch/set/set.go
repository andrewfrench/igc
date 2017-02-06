package set

import (
	"encoding/json"
	"afren.ch/db"
)

type CorrelationSet struct {
	base     string
	Count    int `json:"count"`
	AssocMap map[string]int `json:"set"`
}

func DataSet(base string, data []byte) (*CorrelationSet, error) {
	set := EmptySet()
	set.SetBase(base)

	var err error
	if len(data) > 0 {
		err = json.Unmarshal(data, set)
	}

	return set, err
}

func EmptySet() *CorrelationSet {
	set := new(CorrelationSet)
	set.AssocMap = map[string]int{}

	return set
}

func (s *CorrelationSet) SetBase(base string) {
	s.base = base
}

func (s *CorrelationSet) GetBase() string {
	return s.base
}

func (s *CorrelationSet) IncrementAssoc(assoc string) {
	s.AssocMap[assoc]++
	s.Count++
}

func (s *CorrelationSet) Map() map[string]int {
	return s.AssocMap
}

func (s *CorrelationSet) Serialize() ([]byte, error) {
	data, err := json.Marshal(s)

	return data, err
}

func (s *CorrelationSet) Save() error {
	serialized, err := s.Serialize()
	if err != nil { return err }

	if db.AssociationBaseExists(s.GetBase()) {
		db.UpdateAssociationSet(s.GetBase(), serialized)
	} else {
		db.InsertAssociationSet(s.GetBase(), serialized)
	}

	return err
}

func Merge(a, b *CorrelationSet) *CorrelationSet {
	new := EmptySet()
	new.SetBase(a.GetBase())
	new.Count = a.Count + b.Count

	for a, c := range a.AssocMap {
		new.increaseAssocBy(a, c)
	}

	for a, c := range b.AssocMap {
		new.increaseAssocBy(a, c)
	}

	return new
}

func (s *CorrelationSet) increaseAssocBy(assoc string, count int) {
	s.AssocMap[assoc] += count
}

func PullNew() *CorrelationSet {
	base, data := db.QueryIncomingSet()

	set := EmptySet()
	set.SetBase(base)

	for _, assoc := range data {
		set.IncrementAssoc(assoc)
	}

	return set
}

func PullExisting(base string) (*CorrelationSet, error) {
	data := db.QueryAssociationSet(base)

	return DataSet(base, data)
}