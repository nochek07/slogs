package squeeze

import (
	"sort"
	"testing"
)

var kv = []KeyValue {
	{"a", 4},
	{"b", 1},
	{"c", 17},
	{"n", 3},
	{"f", 15},
}

func TestSortRepeat(t *testing.T) {
	sortedStructByRepeat := KeyValueByRepeat{}.sort(kv)
	data := getIntSlice(sortedStructByRepeat)

	if sort.IsSorted(data) {
		t.Errorf("Not sorted")
	}
}

func TestSortIndex(t *testing.T) {
	sortedStructByIndex := KeyValueByIndex{}.sort(kv)
	data := getIntSlice(sortedStructByIndex)

	if sort.IsSorted(sort.Reverse(data)) {
		t.Errorf("Not sorted")
	}
}

func getIntSlice(kvArray []KeyValue) sort.IntSlice  {
	var data sort.IntSlice
	for _, keyValue := range kvArray {
		data = append(data, int(keyValue.Value))
	}
	return data
}