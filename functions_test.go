package evidence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestK(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(functionKey("a,b,c"), K("a", "b", "c"))
	assert.Equal(functionKey("a,b,c"), K("c", "b", "a"))
	assert.Equal(functionKey("a,b,c"), K("a", "c", "a", "b", "c"))
	assert.Equal(functionKey(""), K())

	assert.Panics(func() {
		K("A")
	})
	assert.Panics(func() {
		K("")
	})
	assert.Panics(func() {
		K("_")
	})
	assert.Panics(func() {
		K("a,b,c")
	})
}

func TestFocalElements(t *testing.T) {
	assert := assert.New(t)

	fk := K()
	fe := fk.FocalElements()
	assert.Len(fe, 0)

	fk = K("a")
	fe = fk.FocalElements()
	assert.Len(fe, 1)
	assert.Contains(fe, functionKey("a"))
}

func TestIsSubset(t *testing.T) {
	assert := assert.New(t)

	assert.True(K("a").IsSubset(K("a", "b", "c")))
	assert.False(K("a", "b", "c").IsSubset(K("a")))
	assert.True(K("a", "b").IsSubset(K("a", "b", "c")))
	assert.True(K("b", "c").IsSubset(K("a", "b", "c")))
	assert.True(K("a", "b").IsSubset(K("a", "b")))
	assert.False(K("d").IsSubset(K("a", "b")))
	assert.False(K("d", "b").IsSubset(K("a", "b")))
	assert.True(K().IsSubset(K("a", "b", "c")))
	assert.True(K().IsSubset(K()))
	assert.False(K("a", "b", "c").IsSubset(K()))
	assert.False(K("a").IsSubset(K()))

	assert.True(K("yellow", "green").IsSubset(K("red", "yellow", "green")))
	assert.False(K("red", "yellow", "green").IsSubset(K("yellow", "green")))
}

func TestIsSuperset(t *testing.T) {
	assert := assert.New(t)

	assert.False(K("a").IsSuperset(K("a", "b", "c")))
	assert.True(K("a", "b", "c").IsSuperset(K("a")))
	assert.False(K("a", "b").IsSuperset(K("a", "b", "c")))
	assert.False(K("b", "c").IsSuperset(K("a", "b", "c")))
	assert.True(K("a", "b").IsSuperset(K("a", "b")))
	assert.False(K("d").IsSuperset(K("a", "b")))
	assert.False(K("d", "b").IsSuperset(K("a", "b")))
	assert.False(K().IsSuperset(K("a", "b", "c")))
	assert.True(K().IsSuperset(K()))
	assert.True(K("a", "b", "c").IsSuperset(K()))
	assert.True(K("a").IsSuperset(K()))

	assert.False(K("yellow", "green").IsSuperset(K("red", "yellow", "green")))
	assert.True(K("red", "yellow", "green").IsSuperset(K("yellow", "green")))
}

func TestPowerset(t *testing.T) {
	assert := assert.New(t)

	f := &Function{}

	f.Set(K("a"), 0.0)
	f.Set(K("b"), 0.1)
	f.Set(K("c"), 0.3)
	f.Set(K("a", "b", "c"), 0.1)
	f.Set(K("b", "c"), 0.5)

	// [[] [a] [b] [a b] [c] [a c] [b c] [a b c]]
	ps := f.Powerset()
	assert.Len(ps, 8)
	assert.Contains(ps, K())
	assert.Contains(ps, K("a"))
	assert.Contains(ps, K("b"))
	assert.Contains(ps, K("c"))
	assert.Contains(ps, K("a", "b"))
	assert.Contains(ps, K("a", "c"))
	assert.Contains(ps, K("b", "c"))
	assert.Contains(ps, K("a", "b", "c"))
	// Shouldn't contain things outside the set
	assert.NotContains(ps, K("d"))
	// Shouldn't contain any invalid keys
	assert.NotContains(ps, functionKey(","))
	assert.NotContains(ps, functionKey("c,b"))
	assert.NotContains(ps, functionKey("c,a"))
	assert.NotContains(ps, functionKey("b,a"))
	assert.NotContains(ps, functionKey("b,c,a"))
	assert.NotContains(ps, functionKey("c,b,a"))
}
