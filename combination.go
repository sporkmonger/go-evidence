package evidence

// combinePairwise takes a pairwise combination function and two or more
// MassFunctions and returns a new MassFunction according to the rule of
// combination given by the combination function. Returns nil if no
// MassFunctions are provided.
func combinePairwise(combiner func(*MassFunction, *MassFunction) *MassFunction,
	mfns ...*MassFunction) *MassFunction {
	if len(mfns) == 0 {
		return nil
	}
	accumulator := mfns[0]
	for _, mf := range mfns[1:] {
		accumulator = combiner(accumulator, mf)
	}
	return accumulator
}

// CombineConjunctive takes two or more MassFunctions and returns a new
// MassFunction according to Dempster's rule of combination. Returns
// nil if no MassFunctions are provided.
func CombineConjunctive(mfns ...*MassFunction) *MassFunction {
	return combinePairwise(pairwiseCombineConjunctive, mfns...)
}

// pairwiseCombineConjunctive takes two MassFunctions and returns a new
// MassFunction according to Dempster's rule of combination.
func pairwiseCombineConjunctive(mf1 *MassFunction, mf2 *MassFunction) (cf *MassFunction) {
	cf = &MassFunction{}
	cf.init()
	for _, p1 := range mf1.Powerset() {
		for _, p2 := range mf2.Powerset() {
			intersect := p1.Intersect(p2)
			cf.Set(intersect, cf.getUnsafe(intersect)+(mf1.Get(p1)*mf2.Get(p2)))
		}
	}
	for _, p := range mf1.Powerset() {
		if p != K() {
			newP := cf.getUnsafe(p) / (1.0 - cf.getUnsafe(K()))
			cf.Set(p, newP)
		}
	}
	cf.Set(K(), 0.0)
	return cf
}

// CombineDisjunctive takes two or more MassFunctions and returns a new
// MassFunction according to the disjunctive rule of combination. Returns
// nil if no MassFunctions are provided.
func CombineDisjunctive(mfns ...*MassFunction) (cf *MassFunction) {
	return combinePairwise(pairwiseCombineDisjunctive, mfns...)
}

// pairwiseCombineDisjunctive takes two MassFunctions and returns a new
// MassFunction according to the disjunctive rule of combination.
func pairwiseCombineDisjunctive(mf1 *MassFunction, mf2 *MassFunction) (cf *MassFunction) {
	cf = &MassFunction{}
	cf.init()
	for _, p1 := range mf1.Powerset() {
		for _, p2 := range mf2.Powerset() {
			union := p1.Union(p2)
			cf.Set(union, cf.getUnsafe(union)+(mf1.Get(p1)*mf2.Get(p2)))
		}
	}
	return cf
}

// CombineMurphyAverage takes two or more MassFunctions and returns a new
// MassFunction according to Murphy's rule of combination, first averaging
// the masses and then performing a conjunctive combination. Returns
// nil if no MassFunctions are provided.
func CombineMurphyAverage(mfns ...*MassFunction) (cf *MassFunction) {
	if len(mfns) == 0 {
		return nil
	}
	count := len(mfns)
	cf = &MassFunction{}
	cf.init()
	for _, p1 := range mfns[0].Powerset() {
		sum := 0.0
		for _, mf := range mfns {
			sum += mf.getUnsafe(p1)
		}
		cf.Set(p1, sum/float64(count))
	}
	cfRepeat := make([]*MassFunction, count)
	for i := 0; i < count; i++ {
		cfRepeat[i] = cf
	}
	return CombineConjunctive(cfRepeat...)
}
