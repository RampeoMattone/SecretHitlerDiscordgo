package Utils

// Set type is a type alias of `map[interface{}]struct{}`
type Set map[interface{}]struct{}

// Add Adds an element to the set
func (s Set) Add(elem interface{}) {
	s[elem] = struct{}{}
}

// AddAll Adds a list of elements to the set
func (s Set) AddAll(elem ...interface{}) {
	for e := range elem {
		s[e] = struct{}{}
	}
}

// Remove Removes an element from the set
func (s Set) Remove(elem interface{}) {
	delete(s, elem)
}

// Clear Removes an element from the set
func (s Set) Clear() {
	for e := range s {
		s.Remove(e)
	}
}

// Has Returns a boolean value describing if the element exists in the set
func (s Set) Has(elem interface{}) bool {
	_, ok := s[elem]
	return ok
}
