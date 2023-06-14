package api

import (
	"reflect"
	"sort"
)

// SliceSort slice排序
func SliceSort[S ~[]E, E any](slice S, less func(a, b E) bool) {
	sort.Slice(slice, func(i, j int) bool {
		a := slice[i]
		b := slice[j]
		return less(a, b)
	})
}

// SliceUnique sorts the slice pointed by the provided pointer given the provided
// less function and removes repeated elements.
// The function panics if the provided interface is not a pointer to a slice.
func SliceUnique[S ~[]E, E any](slicePtr *S, less func(i, j int) bool) {
	v := reflect.ValueOf(slicePtr).Elem()
	if v.Len() <= 1 {
		return
	}
	sort.Slice(v.Interface(), less)

	i := 0
	for j := 1; j < v.Len(); j++ {
		if !less(i, j) {
			continue
		}
		i++
		v.Index(i).Set(v.Index(j))
	}
	i++
	v.SetLen(i)
}