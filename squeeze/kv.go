package squeeze

import "sort"

type KeyValue struct {
	Key string
	Value uint
}

type KeyValueInterface interface {
	sort(kvArray []KeyValue) []KeyValue
	condition(kv1, kv2 KeyValue) bool
}

type KeyValueByRepeat struct {}
type KeyValueByIndex struct {}

func (kvByType KeyValueByRepeat) condition(kv1, kv2 KeyValue) bool {
	return kv1.Value > kv2.Value
}

func (kvByType KeyValueByIndex) condition(kv1, kv2 KeyValue) bool {
	return kv1.Value < kv2.Value
}

func (kvByType KeyValueByRepeat) sort(kvArray []KeyValue) []KeyValue {
	return sortByType(kvByType, kvArray)
}

func (kvByType KeyValueByIndex) sort(kvArray []KeyValue) []KeyValue {
	return sortByType(kvByType, kvArray)
}

func sortByType(kvInterface KeyValueInterface, kvArray []KeyValue) []KeyValue {
	sort.Slice(kvArray, func(i, j int) bool {
		return kvInterface.condition(kvArray[i], kvArray[j])
	})
	return kvArray
}