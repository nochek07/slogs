package squeeze

import "sort"

type KeyValue struct {
	Key string
	Value uint
}

type KeyValueInterface interface {
	sort(kvArray []KeyValue) []KeyValue
	condition(kvFirst, kvSecond KeyValue) bool
}

type KeyValueByRepeat struct {}
type KeyValueByIndex struct {}

func (kvByType KeyValueByRepeat) condition(kvFirst, kvSecond KeyValue) bool {
	return kvFirst.Value > kvSecond.Value
}

func (kvByType KeyValueByIndex) condition(kvFirst, kvSecond KeyValue) bool {
	return kvFirst.Value < kvSecond.Value
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