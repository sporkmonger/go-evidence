package evidence

// A CommonalityFunction is a mapping of possibilities to levels of commonality.
type CommonalityFunction struct {
	Function
}

// Valid verifies that a given CommonalityFunction meets the defined
// requirements for one. All probabilities must be in the range 0.0 >= p >= 1.0.
func (pf *CommonalityFunction) Valid() bool {
	pf.mux.Lock()
	pf.init()
	var sum float64
	for _, probability := range pf.possibilities {
		if probability < 0.0 || probability > 1.0 {
			// This probably shouldn't happen
			pf.mux.Unlock()
			return false
		}
		sum += probability
	}
	// Unlike MassFunction, a CommonalityFunction does not enforce that all
	// values sum to 1.0.
	pf.mux.Unlock()
	return true
}
