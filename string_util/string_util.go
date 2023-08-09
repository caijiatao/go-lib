package string_util

import "strings"

// SplitAndTransferStrSliceToInterfaceSlice
// @Description: 切分字符串数组，size 为每个切片的
func SplitAndTransferStrSliceToInterfaceSlice(stringArray []string, maxSize int) []interface{} {
	var (
		res            = make([]interface{}, 0)
		start          = 0
		end            = 0
		splitNum       int
		stringArrayLen = len(stringArray)
	)
	if len(stringArray) <= maxSize {
		return []interface{}{stringArray}
	}
	if stringArrayLen%maxSize == 0 {
		splitNum = stringArrayLen / maxSize
	} else {
		splitNum = stringArrayLen/maxSize + 1
	}
	for i := 1; i <= splitNum; i++ {
		end = i * maxSize
		if i != splitNum {
			res = append(res, stringArray[start:end])
		} else {
			res = append(res, stringArray[start:])
		}
		start = i * maxSize
	}
	return res
}

// SplitStringSlice
// @Description: 切分字符串数组，size 为每个切片的
func SplitStringSlice(stringArray []string, chunkSize int) [][]string {
	var (
		res            = make([][]string, 0)
		start          = 0
		end            = 0
		splitNum       int
		stringArrayLen = len(stringArray)
	)
	if len(stringArray) <= chunkSize {
		return [][]string{stringArray}
	}
	splitNum = (stringArrayLen + chunkSize - 1) / chunkSize
	for i := 1; i <= splitNum; i++ {
		end = i * chunkSize
		if i != splitNum {
			res = append(res, stringArray[start:end])
		} else {
			res = append(res, stringArray[start:])
		}
		start = i * chunkSize
	}
	return res
}

func FilterStrings(originStrings []string, filterStrings []string) []string {
	res := make([]string, 0)
	m := make(map[string]bool)
	for _, s := range filterStrings {
		m[s] = true
	}

	for _, s := range originStrings {
		if m[s] {
			continue
		}
		res = append(res, s)
	}
	return res
}

func GetPrefix(s, delimiter string, count int) string {
	parts := strings.Split(s, delimiter)
	if len(parts) < count+1 {
		return s
	}
	return strings.Join(parts[:count+1], delimiter)
}
