package evidence

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/google/go-cmp/cmp"
)

// functionKey is an unexported hashable map key for functions
type functionKey string

func (fk functionKey) FocalElements() (fks []functionKey) {
	// Special-case the empty set
	if string(fk) == "" {
		return
	}
	fs := strings.Split(string(fk), ",")
	fks = make([]functionKey, 0, len(fs))
	for i := range fs {
		fks = append(fks, functionKey(fs[i]))
	}
	return fks
}

// IsSubset returns true if a possibility is a subset of another possibility
func (fk functionKey) IsSubset(ofk functionKey) bool {
	count := 0
	fkfe := fk.FocalElements()
	for _, a := range fkfe {
		for _, b := range ofk.FocalElements() {
			if a == b {
				count++
				break
			}
		}
	}
	return count == len(fkfe)
}

// IsSuperset returns true if a possibility is a superset of another possibility
func (fk functionKey) IsSuperset(ofk functionKey) bool {
	count := 0
	ofkfe := ofk.FocalElements()
	for _, a := range ofkfe {
		for _, b := range fk.FocalElements() {
			if a == b {
				count++
				break
			}
		}
	}
	return count == len(ofkfe)
}

// Intersect returns the functionKey that would be the intersection of the two
// keys.
func (fk functionKey) Intersect(ofk functionKey) (ifk functionKey) {
	var intersectingKeys []string
	ofkfe := ofk.FocalElements()
	for _, a := range ofkfe {
		for _, b := range fk.FocalElements() {
			if a == b {
				intersectingKeys = append(intersectingKeys, string(a))
				break
			}
		}
	}
	return K(intersectingKeys...)
}

// Powerset returns all combinations of function keys within this function key.
func (fk functionKey) Powerset() (fks []functionKey) {
	focals := fk.FocalElements()
	n := len(focals)
	for num := 0; num < (1 << uint(n)); num++ {
		combination := []string{}
		for i := 0; i < n; i++ {
			// bit set
			if num&(1<<uint(i)) != 0 {
				// append to the combination
				combination = append(combination, string(focals[i]))
			}
		}
		fks = append(fks, K(combination...))
	}
	return
}

// String presents a human-readable version of the function key
func (fk functionKey) String() string {
	return fmt.Sprintf("{%s}", string(fk))
}

var keyValidator = regexp.MustCompile(`^[a-z0-9-]+$`)

// K generates a mass key from a set of string values.
func K(focals ...string) functionKey {
	// This function can't reasonably return an error due to the context in which
	// it's intended to be used, so remove dupes and sort to generate our key.
	focusSet := make(stringSet)
	for _, focus := range focals {
		if !keyValidator.MatchString(focus) {
			panic(fmt.Sprintf(
				"invalid focus key name (%q), must be lowercase alphanumeric or hyphen",
				focus))
		}
		focusSet[focus] = struct{}{}
	}
	keys := make([]string, 0, len(focusSet))
	for k := range focusSet {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return functionKey(strings.Join(keys, ","))
}

// A Function is a mapping of possibilities to values in the 0.0 to 1.0 range,
// usually probabilities.
type Function struct {
	focalSet      stringSet
	possibilities map[functionKey]float64
	mux           sync.Mutex
}

func (f *Function) init() {
	if f.possibilities == nil {
		f.possibilities = make(map[functionKey]float64)
	}
	if f.focalSet == nil {
		f.focalSet = make(stringSet)
	}
}

// Set assigns a probability to a given possibility.
func (f *Function) Set(key functionKey, probability float64) (err error) {
	f.mux.Lock()
	f.init()
	// If we don't truncate, floating point errors may cause values to go over 1.0
	probability = floatFixed(probability, 5)
	if probability < 0.0 || probability > 1.0 {
		f.mux.Unlock()
		return errors.New("probability out of range")
	}
	// Don't validate further as this can lead to difficulties when changing a
	// mass function in-place. We're only validating the input in isolation.
	f.possibilities[key] = probability
	for _, focus := range key.FocalElements() {
		f.focalSet[string(focus)] = exists
	}
	f.mux.Unlock()
	return nil
}

// Get assigns a probability to a given possibility.
func (f *Function) Get(key functionKey) (probability float64) {
	f.mux.Lock()
	probability = f.getUnsafe(key)
	f.mux.Unlock()
	return
}

// Get assigns a probability to a given possibility.
func (f *Function) getUnsafe(key functionKey) (probability float64) {
	f.init()
	var ok bool
	if probability, ok = f.possibilities[key]; !ok {
		// If the key is missing, the probability is zero.
		probability = 0.0
	}
	return
}

// sort.Interface for function key lists
type fkList struct {
	fks []functionKey
}

func (fkl fkList) Len() int {
	return len(fkl.fks)
}
func (fkl fkList) Less(i, j int) bool {
	ifk := fkl.fks[i]
	jfk := fkl.fks[j]
	if len(ifk) == len(jfk) {
		return strings.Compare(string(ifk), string(jfk)) < 0
	}
	return len(ifk) < len(jfk)
}
func (fkl fkList) Swap(i, j int) {
	ifk := fkl.fks[i]
	jfk := fkl.fks[j]
	fkl.fks[i] = jfk
	fkl.fks[j] = ifk
}

// Possibilities returns a lexically sorted slice of possibilities
func (f *Function) Possibilities() (fks []functionKey) {
	f.mux.Lock()
	f.init()
	for p := range f.possibilities {
		fks = append(fks, p)
	}
	fkl := fkList{
		fks: fks,
	}
	sort.Sort(fkl)
	f.mux.Unlock()
	return fkl.fks
}

// Select returns a new function containing just the specified subset.
func (f *Function) Select(fks []functionKey) (nf *Function) {
	f.mux.Lock()
	f.init()
	nf = &Function{}
	nf.init()
	for _, fk := range fks {
		nf.Set(fk, f.getUnsafe(fk))
	}
	f.mux.Unlock()
	return nf
}

// Powerset returns all combinations of function keys for this Function.
func (f *Function) Powerset() (fks []functionKey) {
	f.mux.Lock()
	f.init()
	focals := make([]string, 0, len(f.focalSet))
	for focal := range f.focalSet {
		focals = append(focals, focal)
	}
	n := len(focals)
	for num := 0; num < (1 << uint(n)); num++ {
		combination := []string{}
		for i := 0; i < n; i++ {
			// bit set
			if num&(1<<uint(i)) != 0 {
				// append to the combination
				combination = append(combination, focals[i])
			}
		}
		fks = append(fks, K(combination...))
	}
	f.mux.Unlock()
	return
}

func floatEq(a float64, b float64) bool {
	const tolerance = 0.00001
	opt := cmp.Comparer(func(x, y float64) bool {
		diff := math.Abs(x - y)
		mean := math.Abs(x+y) / 2.0
		// https://www.juliaferraioli.com/blog/2018/06/golang-testing-floats/
		if math.IsNaN(diff / mean) {
			return true
		}
		return (diff / mean) < tolerance
	})
	return cmp.Equal(a, b, opt)
}

func floatFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(math.Round(num*output)) / output
}
