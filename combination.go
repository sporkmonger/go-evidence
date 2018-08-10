package evidence

// CombineConjunctive takes two MassFunctions and returns a new MassFunction
// according to Dempster's rule of combination.
func CombineConjunctive(mf1 *MassFunction, mf2 *MassFunction) (cf *MassFunction) {
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

// CombineDisjunctive takes two MassFunctions and returns a new MassFunction
// according to Dempster's rule of combination.
func CombineDisjunctive(mf1 *MassFunction, mf2 *MassFunction) (cf *MassFunction) {
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
