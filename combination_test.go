package evidence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombineConjunctive(t *testing.T) {
	assert := assert.New(t)
	const tolerance = 0.00001

	mf1 := &MassFunction{}
	mf1.Set(K(), 0.0)
	mf1.Set(K("red"), 0.35)
	mf1.Set(K("yellow"), 0.25)
	mf1.Set(K("green"), 0.15)
	mf1.Set(K("red", "yellow"), 0.06)
	mf1.Set(K("red", "green"), 0.05)
	mf1.Set(K("yellow", "green"), 0.04)
	mf1.Set(K("red", "yellow", "green"), 0.1)

	mf2 := &MassFunction{}
	mf2.Set(K(), 0.0)
	mf2.Set(K("red"), 0.15)
	mf2.Set(K("yellow"), 0.3)
	mf2.Set(K("green"), 0.2)
	mf2.Set(K("red", "yellow"), 0.03)
	mf2.Set(K("red", "green"), 0.01)
	mf2.Set(K("yellow", "green"), 0.01)
	mf2.Set(K("red", "yellow", "green"), 0.3)

	cf := CombineConjunctive(mf1, mf2)
	assert.InDelta(0.0, cf.Get(K()), tolerance)
	assert.InDelta(0.32737, cf.Get(K("red")), tolerance)
	assert.InDelta(0.35403, cf.Get(K("yellow")), tolerance)
	assert.InDelta(0.18659, cf.Get(K("green")), tolerance)
	assert.InDelta(0.03639, cf.Get(K("red", "yellow")), tolerance)
	assert.InDelta(0.02634, cf.Get(K("red", "green")), tolerance)
	assert.InDelta(0.02139, cf.Get(K("yellow", "green")), tolerance)
	assert.InDelta(0.04789, cf.Get(K("red", "yellow", "green")), tolerance)
	assert.True(cf.Valid())
}
