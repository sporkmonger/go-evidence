package evidence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombineConjunctive(t *testing.T) {
	const tolerance = 0.0025

	tcs := []struct {
		name        string
		mfns        func() []*MassFunction
		expectedMfn func() *MassFunction
	}{
		{
			name: "traffic light",
			mfns: func() []*MassFunction {
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

				return []*MassFunction{mf1, mf2}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K(), 0.0)
				cf.Set(K("red"), 0.32737)
				cf.Set(K("yellow"), 0.35403)
				cf.Set(K("green"), 0.18659)
				cf.Set(K("red", "yellow"), 0.03639)
				cf.Set(K("red", "green"), 0.02634)
				cf.Set(K("yellow", "green"), 0.02139)
				cf.Set(K("red", "yellow", "green"), 0.04789)
				return cf
			},
		},
		{
			name: "block decision factor out mass function",
			mfns: func() []*MassFunction {
				mf1 := &MassFunction{}
				mf1.Set(K(), 0.0)
				mf1.Set(K("allow"), 0.35)
				mf1.Set(K("deny"), 0.2)
				mf1.Set(K("allow", "deny"), 0.45)

				mf2 := &MassFunction{}
				mf2.Set(K(), 0.0)
				mf2.Set(K("allow"), 0.0)
				mf2.Set(K("deny"), 0.0)
				mf2.Set(K("allow", "deny"), 1.0)

				return []*MassFunction{mf1, mf2}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K(), 0.0)
				cf.Set(K("allow"), 0.35)
				cf.Set(K("deny"), 0.2)
				cf.Set(K("allow", "deny"), 0.45)
				return cf
			},
		},
		{
			name: "block decision high conflict",
			mfns: func() []*MassFunction {
				mf1 := &MassFunction{}
				mf1.Set(K(), 0.0)
				mf1.Set(K("allow"), 0.99) // NaN if 1.0
				mf1.Set(K("deny"), 0.01)
				mf1.Set(K("allow", "deny"), 0.0)

				mf2 := &MassFunction{}
				mf2.Set(K(), 0.0)
				mf2.Set(K("allow"), 0.01)
				mf2.Set(K("deny"), 0.99) // NaN if 1.0
				mf2.Set(K("allow", "deny"), 0.0)

				return []*MassFunction{mf1, mf2}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K(), 0.0)
				cf.Set(K("allow"), 0.5)
				cf.Set(K("deny"), 0.5)
				cf.Set(K("allow", "deny"), 0.0)
				return cf
			},
		},
		{
			// https://www.mdpi.com/1424-8220/18/5/1487/pdf
			name: "multi-sensor target recognition system 2x evidence",
			mfns: func() []*MassFunction {
				mf1 := &MassFunction{}
				mf1.Set(K(), 0.0)
				mf1.Set(K("a"), 0.30)
				mf1.Set(K("b"), 0.20)
				mf1.Set(K("c"), 0.10)
				mf1.Set(K("a", "b", "c"), 0.40)

				mf2 := &MassFunction{}
				mf2.Set(K("a"), 0.00)
				mf2.Set(K("b"), 0.90)
				mf2.Set(K("c"), 0.10)
				mf2.Set(K("a", "b", "c"), 0.00)

				return []*MassFunction{mf1, mf2}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K("a"), 0.00)
				cf.Set(K("b"), 0.9153)
				cf.Set(K("c"), 0.0847)
				cf.Set(K("a", "b", "c"), 0.00)
				return cf
			},
		},
		{
			// https://www.mdpi.com/1424-8220/18/5/1487/pdf
			name: "multi-sensor target recognition system 3x evidence",
			mfns: func() []*MassFunction {
				mf1 := &MassFunction{}
				mf1.Set(K(), 0.0)
				mf1.Set(K("a"), 0.30)
				mf1.Set(K("b"), 0.20)
				mf1.Set(K("c"), 0.10)
				mf1.Set(K("a", "b", "c"), 0.40)

				mf2 := &MassFunction{}
				mf2.Set(K("a"), 0.00)
				mf2.Set(K("b"), 0.90)
				mf2.Set(K("c"), 0.10)
				mf2.Set(K("a", "b", "c"), 0.00)

				mf3 := &MassFunction{}
				mf3.Set(K("a"), 0.60)
				mf3.Set(K("b"), 0.10)
				mf3.Set(K("c"), 0.10)
				mf3.Set(K("a", "b", "c"), 0.20)

				return []*MassFunction{mf1, mf2, mf3}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K("a"), 0.00)
				cf.Set(K("b"), 0.9153)
				cf.Set(K("c"), 0.0847)
				cf.Set(K("a", "b", "c"), 0.00)
				return cf
			},
		},
		{
			// https://www.mdpi.com/1424-8220/18/5/1487/pdf
			name: "multi-sensor target recognition system 4x evidence",
			mfns: func() []*MassFunction {
				mf1 := &MassFunction{}
				mf1.Set(K(), 0.0)
				mf1.Set(K("a"), 0.30)
				mf1.Set(K("b"), 0.20)
				mf1.Set(K("c"), 0.10)
				mf1.Set(K("a", "b", "c"), 0.40)

				mf2 := &MassFunction{}
				mf2.Set(K("a"), 0.00)
				mf2.Set(K("b"), 0.90)
				mf2.Set(K("c"), 0.10)
				mf2.Set(K("a", "b", "c"), 0.00)

				mf3 := &MassFunction{}
				mf3.Set(K("a"), 0.60)
				mf3.Set(K("b"), 0.10)
				mf3.Set(K("c"), 0.10)
				mf3.Set(K("a", "b", "c"), 0.20)

				mf4 := &MassFunction{}
				mf4.Set(K("a"), 0.70)
				mf4.Set(K("b"), 0.10)
				mf4.Set(K("c"), 0.10)
				mf4.Set(K("a", "b", "c"), 0.10)

				return []*MassFunction{mf1, mf2, mf3, mf4}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K("a"), 0.00)
				cf.Set(K("b"), 0.9153)
				cf.Set(K("c"), 0.0847)
				cf.Set(K("a", "b", "c"), 0.00)
				return cf
			},
		},
		{
			// https://www.mdpi.com/1424-8220/18/5/1487/pdf
			name: "multi-sensor target recognition system 5x evidence",
			mfns: func() []*MassFunction {
				mf1 := &MassFunction{}
				mf1.Set(K(), 0.0)
				mf1.Set(K("a"), 0.30)
				mf1.Set(K("b"), 0.20)
				mf1.Set(K("c"), 0.10)
				mf1.Set(K("a", "b", "c"), 0.40)

				mf2 := &MassFunction{}
				mf2.Set(K("a"), 0.00)
				mf2.Set(K("b"), 0.90)
				mf2.Set(K("c"), 0.10)
				mf2.Set(K("a", "b", "c"), 0.00)

				mf3 := &MassFunction{}
				mf3.Set(K("a"), 0.60)
				mf3.Set(K("b"), 0.10)
				mf3.Set(K("c"), 0.10)
				mf3.Set(K("a", "b", "c"), 0.20)

				mf4 := &MassFunction{}
				mf4.Set(K("a"), 0.70)
				mf4.Set(K("b"), 0.10)
				mf4.Set(K("c"), 0.10)
				mf4.Set(K("a", "b", "c"), 0.10)

				mf5 := &MassFunction{}
				mf5.Set(K("a"), 0.70)
				mf5.Set(K("b"), 0.10)
				mf5.Set(K("c"), 0.10)
				mf5.Set(K("a", "b", "c"), 0.10)

				return []*MassFunction{mf1, mf2, mf3, mf4, mf5}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K("a"), 0.00)
				cf.Set(K("b"), 0.9153)
				cf.Set(K("c"), 0.0847)
				cf.Set(K("a", "b", "c"), 0.00)
				return cf
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			mfns := tc.mfns()
			expectedMfn := tc.expectedMfn()
			cf := CombineConjunctive(mfns...)
			for _, possibility := range cf.Possibilities() {
				expectedValue := expectedMfn.Get(possibility)
				value := cf.Get(possibility)
				assert.InDelta(expectedValue, value, tolerance)
			}
			assert.True(cf.Valid())
		})
	}
}

func TestCombineDisjunctive(t *testing.T) {
	const tolerance = 0.0025

	tcs := []struct {
		name        string
		mfns        func() []*MassFunction
		expectedMfn func() *MassFunction
	}{
		{
			name: "traffic light",
			mfns: func() []*MassFunction {
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

				return []*MassFunction{mf1, mf2}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K(), 0.0)
				cf.Set(K("red"), 0.0525)
				cf.Set(K("yellow"), 0.075)
				cf.Set(K("green"), 0.03)
				cf.Set(K("red", "yellow"), 0.1893)
				cf.Set(K("red", "green"), 0.1155)
				cf.Set(K("yellow", "green"), 0.1194)
				cf.Set(K("red", "yellow", "green"), 0.4183)
				return cf
			},
		},
		{
			name: "block decision factor out mass function",
			mfns: func() []*MassFunction {
				mf1 := &MassFunction{}
				mf1.Set(K(), 0.0)
				mf1.Set(K("allow"), 0.35)
				mf1.Set(K("deny"), 0.2)
				mf1.Set(K("allow", "deny"), 0.45)

				mf2 := &MassFunction{}
				mf2.Set(K(), 0.0)
				mf2.Set(K("allow"), 0.0)
				mf2.Set(K("deny"), 0.0)
				mf2.Set(K("allow", "deny"), 1.0)

				return []*MassFunction{mf1, mf2}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K(), 0.0)
				cf.Set(K("allow"), 0.0)
				cf.Set(K("deny"), 0.0)
				cf.Set(K("allow", "deny"), 1.0)
				return cf
			},
		},
		{
			name: "block decision high conflict",
			mfns: func() []*MassFunction {
				mf1 := &MassFunction{}
				mf1.Set(K(), 0.0)
				mf1.Set(K("allow"), 1.0)
				mf1.Set(K("deny"), 0.0)
				mf1.Set(K("allow", "deny"), 0.0)

				mf2 := &MassFunction{}
				mf2.Set(K(), 0.0)
				mf2.Set(K("allow"), 0.0)
				mf2.Set(K("deny"), 1.0)
				mf2.Set(K("allow", "deny"), 0.0)

				return []*MassFunction{mf1, mf2}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K(), 0.0)
				cf.Set(K("allow"), 0.0)
				cf.Set(K("deny"), 0.0)
				cf.Set(K("allow", "deny"), 1.0)
				return cf
			},
		},
		{
			// https://www.mdpi.com/1424-8220/18/5/1487/pdf
			name: "multi-sensor target recognition system, 2x evidence",
			mfns: func() []*MassFunction {
				mf1 := &MassFunction{}
				mf1.Set(K(), 0.0)
				mf1.Set(K("a"), 0.30)
				mf1.Set(K("b"), 0.20)
				mf1.Set(K("c"), 0.10)
				mf1.Set(K("a", "b", "c"), 0.40)

				mf2 := &MassFunction{}
				mf2.Set(K("a"), 0.00)
				mf2.Set(K("b"), 0.90)
				mf2.Set(K("c"), 0.10)
				mf2.Set(K("a", "b", "c"), 0.00)

				return []*MassFunction{mf1, mf2}
			},
			expectedMfn: func() *MassFunction {
				cf := &MassFunction{}
				cf.Set(K("a"), 0.00)
				cf.Set(K("b"), 0.18)
				cf.Set(K("c"), 0.01)
				cf.Set(K("a", "b"), 0.27)
				cf.Set(K("b", "c"), 0.11)
				cf.Set(K("a", "c"), 0.03)
				cf.Set(K("a", "b", "c"), 0.40)
				return cf
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			mfns := tc.mfns()
			expectedMfn := tc.expectedMfn()
			cf := CombineDisjunctive(mfns...)
			for _, possibility := range cf.Possibilities() {
				expectedValue := expectedMfn.Get(possibility)
				value := cf.Get(possibility)
				assert.InDelta(expectedValue, value, tolerance)
			}
			assert.True(cf.Valid())
		})
	}
}

func BenchmarkCombineConjunctive(b *testing.B) {
	mf1 := &MassFunction{}
	mf1.Set(K(), 0.0)
	mf1.Set(K("a"), 0.30)
	mf1.Set(K("b"), 0.20)
	mf1.Set(K("c"), 0.10)
	mf1.Set(K("a", "b", "c"), 0.40)

	mf2 := &MassFunction{}
	mf2.Set(K("a"), 0.00)
	mf2.Set(K("b"), 0.90)
	mf2.Set(K("c"), 0.10)
	mf2.Set(K("a", "b", "c"), 0.00)

	mf3 := &MassFunction{}
	mf3.Set(K("a"), 0.60)
	mf3.Set(K("b"), 0.10)
	mf3.Set(K("c"), 0.10)
	mf3.Set(K("a", "b", "c"), 0.20)

	mf4 := &MassFunction{}
	mf4.Set(K("a"), 0.70)
	mf4.Set(K("b"), 0.10)
	mf4.Set(K("c"), 0.10)
	mf4.Set(K("a", "b", "c"), 0.10)

	mf5 := &MassFunction{}
	mf5.Set(K("a"), 0.70)
	mf5.Set(K("b"), 0.10)
	mf5.Set(K("c"), 0.10)
	mf5.Set(K("a", "b", "c"), 0.10)

	mfns := []*MassFunction{mf1, mf2, mf3, mf4, mf5}

	for n := 0; n < b.N; n++ {
		CombineConjunctive(mfns...)
	}
}

func BenchmarkCombineDisjunctive(b *testing.B) {
	mf1 := &MassFunction{}
	mf1.Set(K(), 0.0)
	mf1.Set(K("a"), 0.30)
	mf1.Set(K("b"), 0.20)
	mf1.Set(K("c"), 0.10)
	mf1.Set(K("a", "b", "c"), 0.40)

	mf2 := &MassFunction{}
	mf2.Set(K("a"), 0.00)
	mf2.Set(K("b"), 0.90)
	mf2.Set(K("c"), 0.10)
	mf2.Set(K("a", "b", "c"), 0.00)

	mf3 := &MassFunction{}
	mf3.Set(K("a"), 0.60)
	mf3.Set(K("b"), 0.10)
	mf3.Set(K("c"), 0.10)
	mf3.Set(K("a", "b", "c"), 0.20)

	mf4 := &MassFunction{}
	mf4.Set(K("a"), 0.70)
	mf4.Set(K("b"), 0.10)
	mf4.Set(K("c"), 0.10)
	mf4.Set(K("a", "b", "c"), 0.10)

	mf5 := &MassFunction{}
	mf5.Set(K("a"), 0.70)
	mf5.Set(K("b"), 0.10)
	mf5.Set(K("c"), 0.10)
	mf5.Set(K("a", "b", "c"), 0.10)

	mfns := []*MassFunction{mf1, mf2, mf3, mf4, mf5}

	for n := 0; n < b.N; n++ {
		CombineDisjunctive(mfns...)
	}
}
