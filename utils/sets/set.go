package sets

// Set datastructure for strings
type Set interface {
	Contains(elem string) bool
	Add(elem string)
	Union(another Set) Set
	Intersect(another Set) Set
	Size() int
	Values() []string
}

// Map based Set implementation
type MapSet struct {
	_data map[string]*struct{}
}

// Empty set of strings
func Empty() Set {
	return &MapSet{
		_data: make(map[string]*struct{}),
	}
}

// Creates a new Set from a slice of strings
func FromSlice(slice []string) Set {
	set := Empty()
	for _, elem := range slice {
		set.Add(elem)
	}

	return set
}

// Checks for an existence of an element
func (s *MapSet) Contains(elem string) bool {
	_, present := s._data[elem]
	return present
}

// Add an element to the Set
func (s *MapSet) Add(elem string) {
	s._data[elem] = nil
}

// Union another Set to this set and returns that
func (s *MapSet) Union(another Set) Set {
	union := FromSlice(s.Values())
	for _, value := range another.Values() {
		union.Add(value)
	}
	return union
}

// Intersect another Set to this Set and returns that
func (s *MapSet) Intersect(another Set) Set {
	intersection := Empty()
	for _, elem := range another.Values() {
		if s.Contains(elem) {
			intersection.Add(elem)
		}
	}
	return intersection
}

// Values of the underlying set
func (s *MapSet) Values() []string {
	var values []string
	for key := range s._data {
		values = append(values, key)
	}

	return values
}

// Size of the set
func (s *MapSet) Size() int {
	return len(s._data)
}
