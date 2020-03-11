package squeeze

import (
	"sort"
	"testing"
)

var kv = []KeyValue{
	{"a", 4}, {"b", 1}, {"c", 17}, {"n", 3}, {"f", 15},
}

func TestSortRepeat(t *testing.T) {
	sortedStructByRepeat := KeyValueByRepeat{}.sort(kv)
	data := getDataIntSlice(sortedStructByRepeat)

	if sort.IsSorted(data) {
		t.Fail()
	}
}

func TestSortIndex(t *testing.T) {
	sortedStructByIndex := KeyValueByIndex{}.sort(kv)
	data := getDataIntSlice(sortedStructByIndex)

	if sort.IsSorted(sort.Reverse(data)) {
		t.Fail()
	}
}

func getDataIntSlice(dataKeyValue []KeyValue) sort.IntSlice  {
	var data sort.IntSlice
	for _, keyValue := range dataKeyValue {
		data = append(data, int(keyValue.Value))
	}
	return data
}