package squeeze

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Stat struct {
	repeat, index uint
	date string
}

type MapStat map[string]Stat

func GetMapStat(nameFile string, dateLength int, datePattern string) (MapStat, error) {
	file, err := os.Open(nameFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mapStat := make(MapStat)

	scanner := bufio.NewScanner(file)
	var index uint = 0
	for scanner.Scan() {
		str := scanner.Text()
		if len(str) > dateLength {
			strData := str[:dateLength]

			_, err := time.Parse(datePattern, strData)
			if err == nil {
				key := strings.TrimSpace(strings.TrimPrefix(str, strData))

				if value, ok := mapStat[key]; ok {
					value.repeat++
					value.index = index
					value.date = strData
					mapStat[key] = value
				} else {
					mapStat[key] = Stat{1, index, strData}
				}
			}
		}
		index++
	}

	return mapStat, nil
}

func ReturnResult(nameFile string, mapValues MapStat) error {

	file, err := os.OpenFile(nameFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)

	if err != nil {
		return err
	}
	defer file.Close()

	var sortedStructByRepeat, sortedStructByIndex []KeyValue

	for key, mapValue := range mapValues {
		sortedStructByRepeat = append(sortedStructByRepeat, KeyValue{key, mapValue.repeat})
		sortedStructByIndex  = append(sortedStructByIndex, KeyValue{key, mapValue.index})
	}

	sortedStructByRepeat = KeyValueByRepeat{}.sort(sortedStructByRepeat)
	sortedStructByIndex = KeyValueByIndex{}.sort(sortedStructByIndex)

	var strResultStart, strResultEnd string
	for _, keyValue := range sortedStructByIndex {
		strResultStart += mapValues[keyValue.Key].date + " " + keyValue.Key + "\n"
	}

	for _, keyValue := range sortedStructByRepeat {
		if keyValue.Value != 1 {
			strResultEnd += fmt.Sprint(keyValue.Value) + "\t" + keyValue.Key + "\n"
		}
	}

	if len(strResultEnd) > 0 {
		strResultEnd = "\n==================== Lines more than 2 ====================\n\n" + strResultEnd
	}

	result := strResultStart + strResultEnd
	if len(result) > 0 {
		_, err = file.WriteString(result)
	} else {
		_ = os.Remove(nameFile)
		err = errors.New("empty result file")
	}

	return err
}