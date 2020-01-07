package squeeze

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Stat struct {
	repeat, index uint
	date string
}

type keyValue struct {
	Key string
	Value uint
}

func GetMapStat(nameFile string, dateLength int, datePattern string) (map[string]Stat, error) {
	file, err := os.Open(nameFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	mapStat := make(map[string]Stat)

	scanner := bufio.NewScanner(file)
	var index uint = 0
	for scanner.Scan() {
		str := scanner.Text()
		if len(str)>dateLength {
			strData := str[:dateLength]

			_, err = time.Parse(datePattern, strData)
			if err == nil {
				str = strings.TrimPrefix(str, strData)
				str = strings.TrimSpace(str)

				if value, ok := mapStat[str]; ok {
					value.repeat++
					value.index = index
					value.date = strData
					mapStat[str] = value
				} else {
					mapStat[str] = Stat{1, index, strData}
				}
			}
		}
		index++
	}
	return mapStat, nil
}

func ReturnResult(nameFile string, mapValues map[string]Stat) error {

	file, err := os.OpenFile(nameFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)

	if err != nil {
		return err
	}
	defer file.Close()

	var sortedStructByRepeat, sortedStructByIndex []keyValue

	for key, mapValue := range mapValues {
		sortedStructByRepeat = append(sortedStructByRepeat, keyValue{key, mapValue.repeat})
		sortedStructByIndex  = append(sortedStructByIndex, keyValue{key, mapValue.index})
	}

	sort.Slice(sortedStructByRepeat, func(i, j int) bool {
		return sortedStructByRepeat[i].Value > sortedStructByRepeat[j].Value
	})

	sort.Slice(sortedStructByIndex, func(i, j int) bool {
		return sortedStructByIndex[i].Value < sortedStructByIndex[j].Value
	})

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

	_, err = file.WriteString(strResultStart + strResultEnd)

	return err
}