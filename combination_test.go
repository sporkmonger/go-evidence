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

	cf1 := CombineConjunctive(mf1, mf2)
	assert.InDelta(0.0, cf1.Get(K()), tolerance)
	assert.InDelta(0.32737, cf1.Get(K("red")), tolerance)
	assert.InDelta(0.35403, cf1.Get(K("yellow")), tolerance)
	assert.InDelta(0.18659, cf1.Get(K("green")), tolerance)
	assert.InDelta(0.03639, cf1.Get(K("red", "yellow")), tolerance)
	assert.InDelta(0.02634, cf1.Get(K("red", "green")), tolerance)
	assert.InDelta(0.02139, cf1.Get(K("yellow", "green")), tolerance)
	assert.InDelta(0.04789, cf1.Get(K("red", "yellow", "green")), tolerance)
	assert.True(cf1.Valid())

	mf3 := &MassFunction{}
	mf3.Set(K(), 0.0)
	mf3.Set(K("allow"), 0.35)
	mf3.Set(K("deny"), 0.2)
	mf3.Set(K("allow", "deny"), 0.45)

	mf4 := &MassFunction{}
	mf4.Set(K(), 0.0)
	mf4.Set(K("allow"), 0.0)
	mf4.Set(K("deny"), 0.0)
	mf4.Set(K("allow", "deny"), 1.0)

	cf2 := CombineConjunctive(mf3, mf4)
	assert.InDelta(0.0, cf2.Get(K()), tolerance)
	assert.InDelta(0.35, cf2.Get(K("allow")), tolerance)
	assert.InDelta(0.2, cf2.Get(K("deny")), tolerance)
	assert.InDelta(0.45, cf2.Get(K("allow", "deny")), tolerance)
	assert.True(cf2.Valid())

	mf5 := &MassFunction{}
	mf5.Set(K(), 0.0)
	mf5.Set(K("allow"), 0.99) // NaN if 1.0
	mf5.Set(K("deny"), 0.01)
	mf5.Set(K("allow", "deny"), 0.0)

	mf6 := &MassFunction{}
	mf6.Set(K(), 0.0)
	mf6.Set(K("allow"), 0.01)
	mf6.Set(K("deny"), 0.99) // NaN if 1.0
	mf6.Set(K("allow", "deny"), 0.0)

	cf3 := CombineConjunctive(mf5, mf6)
	assert.InDelta(0.0, cf3.Get(K()), tolerance)
	assert.InDelta(0.5, cf3.Get(K("allow")), tolerance)
	assert.InDelta(0.5, cf3.Get(K("deny")), tolerance)
	assert.InDelta(0.0, cf3.Get(K("allow", "deny")), tolerance)
	assert.True(cf3.Valid())
}

func TestCombineDisjunctive(t *testing.T) {
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

	cf1 := CombineDisjunctive(mf1, mf2)
	assert.InDelta(0.0, cf1.Get(K()), tolerance)
	assert.InDelta(0.0525, cf1.Get(K("red")), tolerance)
	assert.InDelta(0.075, cf1.Get(K("yellow")), tolerance)
	assert.InDelta(0.03, cf1.Get(K("green")), tolerance)
	assert.InDelta(0.1893, cf1.Get(K("red", "yellow")), tolerance)
	assert.InDelta(0.1155, cf1.Get(K("red", "green")), tolerance)
	assert.InDelta(0.1194, cf1.Get(K("yellow", "green")), tolerance)
	assert.InDelta(0.4183, cf1.Get(K("red", "yellow", "green")), tolerance)
	assert.True(cf1.Valid())

	mf3 := &MassFunction{}
	mf3.Set(K(), 0.0)
	mf3.Set(K("allow"), 0.35)
	mf3.Set(K("deny"), 0.2)
	mf3.Set(K("allow", "deny"), 0.45)

	mf4 := &MassFunction{}
	mf4.Set(K(), 0.0)
	mf4.Set(K("allow"), 0.0)
	mf4.Set(K("deny"), 0.0)
	mf4.Set(K("allow", "deny"), 1.0)

	cf2 := CombineDisjunctive(mf3, mf4)
	assert.InDelta(0.0, cf2.Get(K()), tolerance)
	assert.InDelta(0.0, cf2.Get(K("allow")), tolerance)
	assert.InDelta(0.0, cf2.Get(K("deny")), tolerance)
	assert.InDelta(1.0, cf2.Get(K("allow", "deny")), tolerance)
	assert.True(cf2.Valid())

	mf5 := &MassFunction{}
	mf5.Set(K(), 0.0)
	mf5.Set(K("allow"), 1.0)
	mf5.Set(K("deny"), 0.0)
	mf5.Set(K("allow", "deny"), 0.0)

	mf6 := &MassFunction{}
	mf6.Set(K(), 0.0)
	mf6.Set(K("allow"), 0.0)
	mf6.Set(K("deny"), 1.0)
	mf6.Set(K("allow", "deny"), 0.0)

	cf3 := CombineDisjunctive(mf5, mf6)
	assert.InDelta(0.0, cf3.Get(K()), tolerance)
	assert.InDelta(0.0, cf3.Get(K("allow")), tolerance)
	assert.InDelta(0.0, cf3.Get(K("deny")), tolerance)
	assert.InDelta(1.0, cf3.Get(K("allow", "deny")), tolerance)
	assert.True(cf3.Valid())
}
