package evidence

import (
	"fmt"
	"math"
	"strings"
)

// A MassFunction is a mapping of possibilities to probabilities.
type MassFunction struct {
	Function
}

func (mf *MassFunction) String() string {
	var sb strings.Builder
	bf := mf.Belief()
	pf := mf.Plausibility()
	for _, p := range mf.Possibilities() {
		sb.WriteString(fmt.Sprintf("%s\t%f\t%f\t%f\n",
			p, mf.Get(p), bf.Get(p), pf.Get(p)))
	}
	return sb.String()
}

// Valid verifies that a given MassFunction meets the defined requirements for
// one. All probabilities must be in the range 0.0 >= p >= 1.0, and all
// probabilities must ultimately sum to 1.0.
func (mf *MassFunction) Valid() bool {
	mf.mux.Lock()
	mf.init()
	var sum float64
	for _, probability := range mf.possibilities {
		if probability < 0.0 || probability > 1.0 {
			// This probably shouldn't happen
			mf.mux.Unlock()
			return false
		}
		sum += probability
	}
	if !floatEq(sum, 1.0) {
		mf.mux.Unlock()
		return false
	}
	mf.mux.Unlock()
	return true
}

// FocalKeys returns a slice containing just the focal keys
func (mf *MassFunction) FocalKeys() (fks []functionKey) {
	for k := range mf.focalSet {
		fks = append(fks, functionKey(k))
	}
	return fks
}

// Focals returns a function containing just the focal set
func (mf *MassFunction) Focals() (f *Function) {
	return mf.Select(mf.FocalKeys())
}

// Belief converts a MassFunction into a BeliefFunction
func (mf *MassFunction) Belief() (bf *BeliefFunction) {
	fks := mf.Powerset()
	mf.mux.Lock()
	bf = &BeliefFunction{}
	bf.init()
	// Iterate over all keys in the mass function's powerset
	for _, p := range fks {
		value := 0.0
		// Then iterate over all keys in that key's powerset, summing their values
		for _, k := range p.Powerset() {
			value += mf.getUnsafe(k)
		}
		bf.Set(p, value)
	}
	mf.mux.Unlock()
	return
}

// Plausibility converts a MassFunction into a PlausibilityFunction
func (mf *MassFunction) Plausibility() (pf *PlausibilityFunction) {
	fks := mf.Powerset()
	bf := mf.Belief()
	mf.mux.Lock()
	pf = &PlausibilityFunction{}
	pf.init()
	// Iterate over all keys in the mass function's powerset
	for _, p := range fks {
		// ~p is the focals that don't make up proposition p
		// i.e. if p = a,b and the any-set is a,b,c,d, then ~p is c,d
		notP := make([]string, 0)
		for _, k := range mf.FocalKeys() {
			inP := false
			for _, k2 := range p.FocalElements() {
				if k == k2 {
					inP = true
					break
				}
			}
			if !inP {
				notP = append(notP, string(k))
			}
		}
		fk := K(notP...)
		pf.Set(p, 1.0-bf.getUnsafe(fk))
	}
	mf.mux.Unlock()
	return
}

// Commonality converts a MassFunction into a CommonalityFunction
func (mf *MassFunction) Commonality() (cf *CommonalityFunction) {
	fks := mf.Powerset()
	mf.mux.Lock()
	cf = &CommonalityFunction{}
	cf.init()
	// Iterate over all keys in the mass function's powerset
	for _, p := range fks {
		value := 0.0
		for _, k := range fks {
			if p.IsSubset(k) {
				value += mf.getUnsafe(k)
			}
		}
		cf.Set(p, value)
	}
	mf.mux.Unlock()
	return
}

// Pignistic returns a new MassFunction after application of the pignistic
// transformation containing only singletons.
func (mf *MassFunction) Pignistic() (nmf *MassFunction) {
	fks := mf.Powerset()
	mf.mux.Lock()
	nmf = &MassFunction{}
	for _, p := range fks {
		v := mf.getUnsafe(p)
		if v > 0.0 {
			pfe := p.FocalElements()
			size := float64(len(pfe))
			for _, s := range pfe {
				nmf.Set(s, nmf.getUnsafe(s)+(v/size))
			}
		}
	}
	mf.mux.Unlock()
	return
}

// Entropy returns the Deng entropy for the MassFunction.
func (mf *MassFunction) Entropy() float64 {
	entropy := 0.0
	fks := mf.Powerset()
	mf.mux.Lock()
	for _, p := range fks {
		v := mf.getUnsafe(p)
		if floatEq(v, 0.0) {
			continue
		}
		n := len(p.FocalElements())
		entropy -= v * math.Log2(v/(math.Pow(2.0, float64(n))-1.0))
	}
	mf.mux.Unlock()
	return entropy
}
