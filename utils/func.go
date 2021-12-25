package utils

import (
	cryptoRand "crypto/rand"
	"fmt"
	"io"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/maja42/goval"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ParseIntToPointer is func convert int to *int
func ParseIntToPointer(val int) *int {
	return &val
}

// ParseInt32ToPointer is func convert int to *int
func ParseInt32ToPointer(val int32) *int32 {
	return &val
}

// ParseInt64ToPointer is func convert int to *int
func ParseInt64ToPointer(val int64) *int64 {
	return &val
}

// ParseFloat64ToPointer is func convert int to *int
func ParseFloat64ToPointer(val float64) *float64 {
	return &val
}

// ParseStringToPointer is func convert string to *string
func ParseStringToPointer(val string) *string {
	return &val
}

// ParseBoolToPointer is func convert bool to *bool
func ParseBoolToPointer(val bool) *bool {
	return &val
}

// ParseTimeToPointer is func convert time.Time to *time.Time
func ParseTimeToPointer(val time.Time) *time.Time {
	return &val
}

// PointerToInt ...
func PointerToInt(val *int) int {
	return *val
}

func GetIntArrayFromString(str string) ([]int32, error) {
	numList := strings.Split(str, ",")
	var arr = []int32{}
	for _, v := range numList {
		if v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				arr = append(arr, int32(n))
			} else {
				return nil, err
			}
		}
	}

	return arr, nil
}

func GetDatesString(times []time.Time) []string {
	dates := []string{}
	layout := "2006-01-02"
	for _, v := range times {
		dates = append(dates, v.Format(layout))
	}
	return dates
}

func GetNumberRange(from, to int32) []int32 {
	var numbers []int32
	for from <= to {
		numbers = append(numbers, from)
		from++
	}
	return numbers
}

func GetDatesRange(from time.Time, to time.Time) []time.Time {
	var ans []time.Time
	t := from
	ans = append(ans, t)
	for t.Before(to) {
		t = t.Add(24 * time.Hour)
		ans = append(ans, t)
	}
	return ans
}

// CalcValWithFormula ...
func CalcValWithFormula(formula string, params map[string]interface{}) (float64, error) {
	result, err := parseFormulaToInterface(formula, params)
	if err != nil {
		return 0, err
	}

	val, ok := result.(float64)
	if !ok {
		return 0, err
	}

	return val, nil
}

// parseFormulaToInterface ...
func parseFormulaToInterface(formula string, params map[string]interface{}) (interface{}, error) {
	expression, err := govaluate.NewEvaluableExpression(formula)
	if err != nil {
		return nil, err
	}

	result, err := expression.Evaluate(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CalcValWithFnsFormula
func CalcValWithFnsFormula(formula string, variables map[string]interface{}, fns map[string]goval.ExpressionFunction) (result float64, err error) {
	defer func() {
		if a := recover(); a != nil {
			err = fmt.Errorf("panic")
		}
	}()
	eval := goval.NewEvaluator()
	res, err := eval.Evaluate(formula, variables, fns) // Returns <36, nil>
	if err != nil {
		result = 0
	}
	switch res.(type) {
	case int:
		result = float64(res.(int))
	case float64:
		result = res.(float64)
	}
	return
}

func GetDatesByTypeDate(typeDate int, dates []string) ([]string, error) {
	switch typeDate {
	case 0:
		return dates, nil
	case 1:
		mapWeek := make(map[string]bool)
		var weeks []string
		for _, t := range dates {
			f, _ := time.Parse("2006-01-02", t)
			y, w := f.ISOWeek()
			var week string
			if w < 10 {
				week = fmt.Sprintf("0%v", w)
			} else {
				week = fmt.Sprint(w)
			}
			if _, ok := mapWeek[week]; !ok {
				weeks = append(weeks, fmt.Sprintf("%vW%v", y, week))
				mapWeek[week] = true
			}
		}
		return weeks, nil
	case 2:
		mapMonth := make(map[string]bool)
		var months []string
		for _, t := range dates {
			arr := strings.Split(t, "-")
			if len(arr) != 3 {
				return nil, fmt.Errorf("bad request: invalid date %v", t)
			}
			date := fmt.Sprintf("%v-%v-01", arr[0], arr[1])
			if _, ok := mapMonth[date]; !ok {
				months = append(months, date)
				mapMonth[date] = true
			}
		}
		return months, nil
	case 3:
		mapYear := make(map[string]bool)
		var years []string
		for _, t := range dates {
			f, _ := time.Parse("2006-01-02", t)
			y := f.Year()
			year := fmt.Sprint(y)
			if _, ok := mapYear[year]; !ok {
				years = append(years, year)
				mapYear[year] = true
			}
		}
		return years, nil
	default:
		return nil, fmt.Errorf("bad request: not support type_date = %v", typeDate)
	}
}

func UniqueSlice(intSlice []int32) []int32 {
	keys := make(map[int32]bool)
	list := []int32{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func ParseStandardFormula(formula string) string {
	newFormula := ""
	arr := strings.Split(formula, "")
	for idx, k := range arr {
		if k == "+" && idx == 0 || k == "+" && ((arr[idx-1] > "9" || arr[idx-1] < "0") && arr[idx-1] != ")") {
			continue
		}
		newFormula = newFormula + k
	}
	return newFormula
}

// round up with particular precision
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func GenOTPCode(length int) string {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(cryptoRand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RemoveUnicode(str string) string {
	m1 := regexp.MustCompile(`/à|á|ạ|ả|ã|â|ầ|ấ|ậ|ẩ|ẫ|ă|ằ|ắ|ặ|ẳ|ẵ/g`)
	str = m1.ReplaceAllString(str, "a")
	m1 = regexp.MustCompile(`/è|é|ẹ|ẻ|ẽ|ê|ề|ế|ệ|ể|ễ/g`)
	str = m1.ReplaceAllString(str, "e")
	m1 = regexp.MustCompile(`/ì|í|ị|ỉ|ĩ/g`)
	str = m1.ReplaceAllString(str, "i")
	m1 = regexp.MustCompile(`/ò|ó|ọ|ỏ|õ|ô|ồ|ố|ộ|ổ|ỗ|ơ|ờ|ớ|ợ|ở|ỡ/g`)
	str = m1.ReplaceAllString(str, "o")
	m1 = regexp.MustCompile(`/ù|ú|ụ|ủ|ũ|ư|ừ|ứ|ự|ử|ữ/g`)
	str = m1.ReplaceAllString(str, "u")
	m1 = regexp.MustCompile(`/ỳ|ý|ỵ|ỷ|ỹ/g`)
	str = m1.ReplaceAllString(str, "y")
	m1 = regexp.MustCompile(`/đ/g`)
	str = m1.ReplaceAllString(str, "d")
	m1 = regexp.MustCompile(`/À|Á|Ạ|Ả|Ã|Â|Ầ|Ấ|Ậ|Ẩ|Ẫ|Ă|Ằ|Ắ|Ặ|Ẳ|Ẵ/g`)
	str = m1.ReplaceAllString(str, "A")
	m1 = regexp.MustCompile(`/È|É|Ẹ|Ẻ|Ẽ|Ê|Ề|Ế|Ệ|Ể|Ễ/g`)
	str = m1.ReplaceAllString(str, "E")
	m1 = regexp.MustCompile(`/Ì|Í|Ị|Ỉ|Ĩ/g`)
	str = m1.ReplaceAllString(str, "I")
	m1 = regexp.MustCompile(`/Ò|Ó|Ọ|Ỏ|Õ|Ô|Ồ|Ố|Ộ|Ổ|Ỗ|Ơ|Ờ|Ớ|Ợ|Ở|Ỡ/g`)
	str = m1.ReplaceAllString(str, "O")
	m1 = regexp.MustCompile(`/Ù|Ú|Ụ|Ủ|Ũ|Ư|Ừ|Ứ|Ự|Ử|Ữ/g`)
	str = m1.ReplaceAllString(str, "U")
	m1 = regexp.MustCompile(`/Ỳ|Ý|Ỵ|Ỷ|Ỹ/g`)
	str = m1.ReplaceAllString(str, "Y")
	m1 = regexp.MustCompile(`/Đ/g`)
	str = m1.ReplaceAllString(str, "D")

	return str
}
