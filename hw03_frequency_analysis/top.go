package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type keyValue struct {
	key   string
	value int
}

func Top10(str string) []string {
	mp := make(map[string]int)
	str = strings.ReplaceAll(str, "\n", " ")
	initSlice := strings.Split(str, " ")
	for _, word := range initSlice {
		word = strings.TrimSpace(word)
		if len(word) > 0 {
			mp[strings.TrimSpace(word)]++
		}
	}

	// Формируем слайс структур keyValue - готовим к сортировке
	kvSlice := make([]keyValue, 0, len(mp))
	for k, v := range mp {
		kvSlice = append(kvSlice, keyValue{key: k, value: v})
	}

	sort.Slice(kvSlice, func(i, j int) bool {
		// Сначала сравниваем по значению (по убыванию)
		if kvSlice[i].value != kvSlice[j].value {
			return kvSlice[i].value > kvSlice[j].value
		}
		// Если значения равны, сравниваем по ключу (по возрастанию)
		return kvSlice[i].key < kvSlice[j].key
	})

	returnStr := make([]string, 0, 10)
	for idx, kv := range kvSlice {
		if idx > 9 {
			break
		}
		returnStr = append(returnStr, kv.key)
	}
	return returnStr
}
