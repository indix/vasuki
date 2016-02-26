package sets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptySet(t *testing.T) {
	emptySet := Empty()
	assert.Equal(t, 0, emptySet.Size())
}

func TestFromSlice(t *testing.T) {
	fromSlice := FromSlice([]string{
		"foo", "bar",
	})

	assert.Equal(t, 2, fromSlice.Size())
}

func TestSetHasNoDuplicates(t *testing.T) {
	fromSlice := FromSlice([]string{
		"foo", "foo", "foo",
	})

	assert.Equal(t, 1, fromSlice.Size())
}

func TestSetAdd(t *testing.T) {
	set := Empty()
	set.Add("foo")
	set.Add("bar")
	assert.Equal(t, 2, set.Size())
}

func TestSetContains(t *testing.T) {
	set := Empty()
	set.Add("foo")
	set.Add("bar")
	assert.True(t, set.Contains("foo"))
}

func TestSetUnion(t *testing.T) {
	set1 := FromSlice([]string{
		"foo", "bar",
	})

	set2 := FromSlice([]string{
		"baz",
	})

	union := set1.Union(set2)
	assert.Equal(t, 3, union.Size())
	assert.True(t, union.Contains("foo"))
	assert.True(t, union.Contains("bar"))
	assert.True(t, union.Contains("baz"))
}

func TestSetIntersect(t *testing.T) {
	set1 := FromSlice([]string{
		"foo", "bar",
	})

	set2 := FromSlice([]string{
		"foo", "baz",
	})

	intersect := set1.Intersect(set2)
	assert.Equal(t, 1, intersect.Size())
	assert.True(t, intersect.Contains("foo"))
}

func TestSetIsSuperSetOf(t *testing.T) {
	set1 := FromSlice([]string{
		"foo", "bar", "baz",
	})

	set2 := FromSlice([]string{
		"foo", "baz",
	})

	assert.True(t, set1.IsSupersetOf(set2))
	assert.True(t, set1.IsSupersetOf(Empty())) // ø is a subset of everything
}

func TestSetEqual(t *testing.T) {
	set1 := FromSlice([]string{
		"foo", "bar", "baz",
	})

	set2 := FromSlice([]string{
		"bar", "foo", "baz",
	})

	assert.True(t, set1.Equal(set2))
}
