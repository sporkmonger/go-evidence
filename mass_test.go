package evidence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessors(t *testing.T) {
	assert := assert.New(t)

	mf := &MassFunction{}

	// Unassigned keys are zero
	assert.Equal(0.0, mf.Get(K()))
	assert.False(mf.Valid())
	// Still zero if you assign zero
	mf.Set(K("a"), 0.0)
	assert.Equal(0.0, mf.Get(K("a")))
	assert.False(mf.Valid())
	assert.True(mf.focalSet.Contains("a"))
	mf.Set(K("b"), 0.1)
	assert.Equal(0.1, mf.Get(K("b")))
	assert.False(mf.Valid())
	assert.True(mf.focalSet.Contains("b"))
	// Set something non-zero
	mf.Set(K("a", "b", "c"), 0.4)
	// Adding something new shouldn't change what we already had
	assert.Equal(0.0, mf.Get(K("a")))
	assert.Equal(0.4, mf.Get(K("a", "b", "c")))
	assert.False(mf.Valid())
	assert.True(mf.focalSet.Contains("c"))
	// Setting an invalid value should give an error and leave value unchanged
	err := mf.Set(K("a", "b", "c"), 4.0)
	assert.NotNil(err)
	assert.Equal(0.4, mf.Get(K("a", "b", "c")))
	assert.False(mf.Valid())
	// Finally get a 1.0 sum and ensure it's now valid
	mf.Set(K("b", "c"), 0.5)
	assert.True(mf.Valid())
	// Make sure that checking for things that don't exist also works
	assert.False(mf.focalSet.Contains("d"))
}

func TestFocals(t *testing.T) {
	assert := assert.New(t)

	mf := &MassFunction{}

	// Still zero if you assign zero
	mf.Set(K("a"), 0.0)
	mf.Set(K("b"), 0.1)
	mf.Set(K("c"), 0.3)
	mf.Set(K("a", "b", "c"), 0.1)
	mf.Set(K("b", "c"), 0.5)
	fs := mf.Focals()
	assert.Equal(0.0, fs.Get(K("a")))
	assert.Equal(0.1, fs.Get(K("b")))
	assert.Equal(0.3, fs.Get(K("c")))
	assert.Equal(0.0, fs.Get(K("a", "b", "c")))
	assert.Equal(0.0, fs.Get(K("b", "c")))
	assert.Equal(0.0, fs.Get(K()))
}

func TestBelief(t *testing.T) {
	assert := assert.New(t)
	const tolerance = 0.00001

	mf := &MassFunction{}

	mf.Set(K("a"), 0.0)
	mf.Set(K("b"), 0.1)
	mf.Set(K("c"), 0.3)
	mf.Set(K("a", "b", "c"), 0.1)
	mf.Set(K("b", "c"), 0.5)
	bf := mf.Belief()
	assert.Equal(0.0, bf.Get(K()))
	assert.Equal(0.0, bf.Get(K("a")))
	assert.Equal(0.1, bf.Get(K("b")))
	assert.Equal(0.3, bf.Get(K("c")))
	assert.Equal(0.0+0.1+0.0, bf.Get(K("a", "b")))          // 0.1
	assert.Equal(0.0+0.3+0.0, bf.Get(K("a", "c")))          // 0.3
	assert.Equal(0.1+0.3+0.5, bf.Get(K("b", "c")))          // 0.9
	assert.Equal(0.1+0.3+0.5+0.1, bf.Get(K("a", "b", "c"))) // 1.0
	assert.True(bf.Valid())

	mf = &MassFunction{}
	mf.Set(K(), 0.0)
	mf.Set(K("alive"), 0.2)
	mf.Set(K("dead"), 0.5)
	mf.Set(K("alive", "dead"), 0.3)
	bf = mf.Belief()
	assert.Equal(0.0, bf.Get(K()))
	assert.Equal(0.2, bf.Get(K("alive")))
	assert.Equal(0.5, bf.Get(K("dead")))
	assert.Equal(0.2+0.5+0.3, bf.Get(K("alive", "dead"))) // 1.0
	assert.True(bf.Valid())

	mf = &MassFunction{}
	mf.Set(K(), 0.0)
	mf.Set(K("red"), 0.35)
	mf.Set(K("yellow"), 0.25)
	mf.Set(K("green"), 0.15)
	mf.Set(K("red", "yellow"), 0.06)
	mf.Set(K("red", "green"), 0.05)
	mf.Set(K("yellow", "green"), 0.04)
	mf.Set(K("red", "yellow", "green"), 0.1)
	bf = mf.Belief()
	assert.InDelta(0.0, bf.Get(K()), tolerance)
	assert.InDelta(0.35, bf.Get(K("red")), tolerance)
	assert.InDelta(0.25, bf.Get(K("yellow")), tolerance)
	assert.InDelta(0.15, bf.Get(K("green")), tolerance)
	assert.InDelta(0.35+0.25+0.06, bf.Get(K("red", "yellow")), tolerance)                             // 0.66
	assert.InDelta(0.35+0.15+0.05, bf.Get(K("red", "green")), tolerance)                              // 0.55
	assert.InDelta(0.25+0.15+0.04, bf.Get(K("yellow", "green")), tolerance)                           // 0.44
	assert.InDelta(0.35+0.25+0.15+0.06+0.05+0.04+0.1, bf.Get(K("red", "yellow", "green")), tolerance) // 1.0
	assert.True(bf.Valid())
}

func TestPlausibility(t *testing.T) {
	assert := assert.New(t)
	const tolerance = 0.00001

	mf := &MassFunction{}

	mf.Set(K("a"), 0.0)
	mf.Set(K("b"), 0.1)
	mf.Set(K("c"), 0.3)
	mf.Set(K("a", "b", "c"), 0.1)
	mf.Set(K("b", "c"), 0.5)
	pf := mf.Plausibility()
	assert.InDelta(0.0, pf.Get(K()), tolerance)
	assert.InDelta(1.0-0.9, pf.Get(K("a")), tolerance)           // 0.1
	assert.InDelta(1.0-0.3, pf.Get(K("b")), tolerance)           // 0.7
	assert.InDelta(1.0-0.1, pf.Get(K("c")), tolerance)           // 0.9
	assert.InDelta(1.0-0.3, pf.Get(K("a", "b")), tolerance)      // 0.7
	assert.InDelta(1.0-0.1, pf.Get(K("a", "c")), tolerance)      // 0.9
	assert.InDelta(1.0-0.0, pf.Get(K("b", "c")), tolerance)      // 1.0
	assert.InDelta(1.0-0.0, pf.Get(K("a", "b", "c")), tolerance) // 1.0
	assert.True(pf.Valid())

	mf = &MassFunction{}
	mf.Set(K(), 0.0)
	mf.Set(K("alive"), 0.2)
	mf.Set(K("dead"), 0.5)
	mf.Set(K("alive", "dead"), 0.3)
	pf = mf.Plausibility()
	assert.Equal(0.0, pf.Get(K()))
	assert.Equal(1.0-0.5, pf.Get(K("alive")))         // 0.5
	assert.Equal(1.0-0.2, pf.Get(K("dead")))          // 0.8
	assert.Equal(1.0-0.0, pf.Get(K("alive", "dead"))) // 1.0
	assert.True(pf.Valid())
}

func TestCommonality(t *testing.T) {
	assert := assert.New(t)
	const tolerance = 0.00001

	mf := &MassFunction{}
	mf.Set(K(), 0.0)
	mf.Set(K("red"), 0.35)
	mf.Set(K("yellow"), 0.25)
	mf.Set(K("green"), 0.15)
	mf.Set(K("red", "yellow"), 0.06)
	mf.Set(K("red", "green"), 0.05)
	mf.Set(K("yellow", "green"), 0.04)
	mf.Set(K("red", "yellow", "green"), 0.1)
	cf := mf.Commonality()
	assert.InDelta(1.0, cf.Get(K()), tolerance)
	assert.InDelta(0.35+0.06+0.05+0.1, cf.Get(K("red")), tolerance)    // 0.56
	assert.InDelta(0.25+0.06+0.04+0.1, cf.Get(K("yellow")), tolerance) // 0.45
	assert.InDelta(0.15+0.05+0.04+0.1, cf.Get(K("green")), tolerance)  // 0.34
	assert.InDelta(0.06+0.1, cf.Get(K("red", "yellow")), tolerance)    // 0.16
	assert.InDelta(0.05+0.1, cf.Get(K("red", "green")), tolerance)     // 0.15
	assert.InDelta(0.04+0.1, cf.Get(K("yellow", "green")), tolerance)  // 0.14
	assert.InDelta(0.1, cf.Get(K("red", "yellow", "green")), tolerance)
	assert.True(cf.Valid())
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	mf := &MassFunction{}

	mf.Set(K("a"), 0.0)
	mf.Set(K("b"), 0.1)
	mf.Set(K("c"), 0.3)
	mf.Set(K("a", "b", "c"), 0.1)
	mf.Set(K("b", "c"), 0.5)
	str := mf.String()
	assert.Contains(str, "{a}\t0.000000\t0.000000\t0.100000")
	assert.Contains(str, "{b}\t0.100000\t0.100000\t0.700000")
	assert.Contains(str, "{c}\t0.300000\t0.300000\t0.900000")
	assert.Contains(str, "{b,c}\t0.500000\t0.900000\t1.000000")
	assert.Contains(str, "{a,b,c}\t0.100000\t1.000000\t1.000000")
}
