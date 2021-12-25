package utils

import (
	"fmt"
	"strconv"
)

func GetIntArrayFromInt32Array(arr []int32) []int {
	var ans []int
	for _, a := range arr {
		ans = append(ans, int(a))
	}
	return ans
}

func GetIntArrayFromInt64Array(arr []int64) []int {
	var ans []int
	for _, a := range arr {
		ans = append(ans, int(a))
	}
	return ans
}

func GetInt32ArrayFromIntArray(arr []int) []int32 {
	var ans []int32
	for _, a := range arr {
		ans = append(ans, int32(a))
	}
	return ans
}

func GetInt64ArrayFromIntArray(arr []int) []int64 {
	var ans []int64
	for _, a := range arr {
		ans = append(ans, int64(a))
	}
	return ans
}

func GetIntArrayFromStringArray(arr []string) ([]int, error) {
	var ans []int
	for _, a := range arr {
		i, err := strconv.Atoi(a)
		if err != nil {
			return nil, err
		}
		ans = append(ans, int(i))
	}
	return ans, nil
}

func GetStringArrayFromIntArray(arr []int) []string {
	var ans []string
	for _, a := range arr {
		ans = append(ans, fmt.Sprint(a))
	}
	return ans
}
