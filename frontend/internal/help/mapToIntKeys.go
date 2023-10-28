package help

import "strconv"

func MapToIntKeys(inputMap map[string]interface{}) map[int]interface{} {
	out := make(map[int]interface{})
	for i, v := range inputMap {
		idx, _ := strconv.Atoi(i)
		out[idx] = v
	}
	return out
}
