package evidence

// DempsterCombine takes two MassFunctions and returns a new MassFunction
// according to Dempster's rule of combination.
func DempsterCombine(mf1 *MassFunction, mf2 *MassFunction) (cf *MassFunction) {
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
