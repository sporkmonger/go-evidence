package evidence

// A BeliefFunction is a mapping of possibilities to levels of support/belief.
type BeliefFunction struct {
	Function
}

// Valid verifies that a given BeliefFunction meets the defined requirements for
// one. All probabilities must be in the range 0.0 >= p >= 1.0.
func (bf *BeliefFunction) Valid() bool {
	bf.mux.Lock()
	bf.init()
	var sum float64
	for _, probability := range bf.possibilities {
		if probability < 0.0 || probability > 1.0 {
			// This probably shouldn't happen
			bf.mux.Unlock()
			return false
		}
		sum += probability
	}
	// Unlike MassFunction, a BeliefFunction does not enforce that all values sum
	// to 1.0.
	bf.mux.Unlock()
	return true
}
