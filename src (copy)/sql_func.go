package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func sql_MAX(array interface{}) ([]float64, error) {
	result := []float64{}
	var tmp float64
	switch v := array.(type) {
	case []int:

		tmp = math.Max(float64(v[0]), float64(v[1]))
		for i := 0; i < len(v); i++ {
			tmp = math.Max(tmp, float64(v[i]))
		}
		result = append(result, tmp)
		return result, nil

	case []float64:

		tmp = math.Max(float64(v[0]), float64(v[1]))
		for i := 0; i < len(v); i++ {
			tmp = math.Max(tmp, float64(v[i]))
		}
		result = append(result, tmp)
		return result, nil

	default:
		return nil, errors.New("sql: MAX function expected array")
	}
}

func sql_MIN(array interface{}) ([]float64, error) {
	result := []float64{}
	var tmp float64
	switch v := array.(type) {
	case []int:

		tmp = math.Min(float64(v[0]), float64(v[1]))
		for i := 0; i < len(v); i++ {
			tmp = math.Min(tmp, float64(v[i]))
		}
		result = append(result, tmp)
		return result, nil

	case []float64:

		tmp = math.Min(float64(v[0]), float64(v[1]))
		for i := 0; i < len(v); i++ {
			tmp = math.Min(tmp, float64(v[i]))
		}
		result = append(result, tmp)
		return result, nil

	default:
		return nil, errors.New("sql: MIN function expected array")
	}
}

func sql_AVG(array interface{}) ([]float64, error) {
	result := []float64{}
	var tmp float64 = 0
	switch v := array.(type) {
	case []int:
		for _, value := range v {
			tmp += float64(value)
		}

		tmp = tmp / float64(len(v))

	case []float64:
		for _, value := range v {
			tmp += value
		}
		tmp = tmp / float64(len(v))
	default:
		return nil, errors.New("sql: AVG function expected array")
	}

	result = append(result, tmp)

	return result, nil
}

func sql_SUM(array interface{}) ([]float64, error) {
	result := []float64{}
	var tmp float64 = 0
	switch v := array.(type) {
	case []int:
		for _, value := range v {
			tmp += float64(value)
		}

	case []float64:
		for _, value := range v {
			tmp += value
		}

	default:
		return nil, errors.New("sql: SUM function expected array")
	}

	result = append(result, tmp)

	return result, nil
}

// Start of sql string functions

func sql_CHAR_LENGTH(array []string) ([]int, error) {
	result := []int{}
	for _, v := range array {
		result = append(result, len(v))
	}

	return result, nil
}

func sql_CHARACTER_LENGTH(array []string) ([]int, error) {
	return sql_CHAR_LENGTH(array)
}

func sql_LCASE(array []string) ([]string, error) {
	result := []string{}
	for _, v := range array {
		result = append(result, strings.ToLower(v))

	}
	return result, nil
}
func sql_LOWER(array []string) ([]string, error) {
	return sql_LCASE(array)
}

func sql_UCASE(array []string) ([]string, error) {
	result := []string{}
	for _, v := range array {
		result = append(result, strings.ToUpper(v))

	}
	return result, nil
}

func sql_UPPER(array []string) ([]string, error) {
	return sql_UCASE(array)
}

func sql_LENGTH(array []string) ([]int, error) {
	result := []int{}
	for _, v := range array {
		result = append(result, len(v))
	}
	return result, nil
}

func sql_REVERSE(array []string) ([]string, error) {
	result := []string{}
	for _, v := range array {
		var str = ""
		for _, v := range strings.Split(v, "") {
			str = v + str
		}
		result = append(result, str)
	}
	return result, nil
}

/*
END of sql string functions
*/
func sql_CONCAT(columns []interface{}) ([]string, error) {
	result := []string{}
	var num int = 0
	switch v := columns[0].(type) {
	case []int:
		num = len(v)
	case []string:
		num = len(v)
	case []bool:
		num = len(v)
	case [][]byte:
		num = len(v)
	case []float64:
		num = len(v)
	}
	for i := 0; i < num; i++ {
		result = append(result, "")
	}
	for _, y := range columns {
		switch v := y.(type) {
		case []int:
			for i, a := range v {
				result[i] += strconv.FormatInt(int64(a), 10)
			}

		case []string:
			for i, a := range v {
				result[i] += a
			}
		case []bool:
			for i, a := range v {
				result[i] += strconv.FormatBool(a)
			}
		case [][]byte:
			for i, a := range v {
				result[i] += string(a)
			}
		case []float64:
			for i, a := range v {
				result[i] += strconv.FormatFloat(a, 'g', -1, 64)
			}
		}
	}
	return result, nil

}
func sql_CONCAT_WS(columns []interface{}, ws string) ([]string, error) {

	result := []string{}
	var num int = 0
	switch v := columns[0].(type) {
	case []int:
		num = len(v)
	case []string:
		num = len(v)
	case []bool:
		num = len(v)
	case [][]byte:
		num = len(v)
	case []float64:
		num = len(v)
	}
	for i := 0; i < num; i++ {
		result = append(result, "")
	}
	for _, y := range columns {
		switch v := y.(type) {
		case []int:
			for i, a := range v {
				result[i] += strconv.FormatInt(int64(a), 10)
				if i < num-1 {
					result[i] += ws
				}
			}

		case []string:
			for i, a := range v {
				result[i] += a
				if i < num-1 {
					result[i] += ws
				}
			}
		case []bool:
			for i, a := range v {
				result[i] += strconv.FormatBool(a)
				if i < num-1 {
					result[i] += ws
				}
			}
		case [][]byte:
			for i, a := range v {
				result[i] += string(a)
				if i < num-1 {
					result[i] += ws
				}

			}
		case []float64:
			for i, a := range v {
				result[i] += strconv.FormatFloat(a, 'g', -1, 64)
				if i < num-1 {
					result[i] += ws
				}
			}
		}
	}
	return result, nil
}

// start of sql Numeric functions

func sql_ABS(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Abs(v))

	}
	return result, nil
}

func sql_ACOS(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Acos(v))
	}
	return result, nil
}

func sql_ACOSH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Acosh(v))
	}
	return result, nil
}

func sql_ASIN(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Asin(v))

	}
	return result, nil
}

func sql_ASINH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Asinh(v))

	}
	return result, nil
}
func sql_ATAN(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Atan(v))

	}
	return result, nil
}
func sql_ATANH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Atanh(v))

	}
	return result, nil
}

func sql_CEIL(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Ceil(v))

	}
	return result, nil
}

func sql_COS(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Cos(v))

	}
	return result, nil
}

func sql_COSH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Cosh(v))

	}
	return result, nil
}

func sql_CBRT(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Cbrt(v))

	}
	return result, nil
}

func sql_EXP(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Exp(v))

	}
	return result, nil
}

func sql_EXP2(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Exp2(v))

	}
	return result, nil
}

func sql_EXMP1(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Expm1(v))

	}
	return result, nil
}

func sql_FLOOR(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Floor(v))

	}
	return result, nil
}

func sql_LN(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Log(v))

	}
	return result, nil
}

func sql_LOG(array []float64) ([]float64, error) {
	return sql_LN(array)
}

func sql_LOG10(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Log10(v))

	}
	return result, nil
}

func sql_LOG2(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Log2(v))

	}
	return result, nil
}

func sql_SIGNBIT(array []float64) ([]bool, error) {
	result := []bool{}
	for _, v := range array {
		result = append(result, math.Signbit(v))

	}
	return result, nil
}

func sql_SIN(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Sin(v))

	}
	return result, nil
}

func sql_SINH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Sinh(v))

	}
	return result, nil
}

func sql_SQRT(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Sqrt(v))

	}
	return result, nil
}

func sql_TAN(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Tan(v))

	}
	return result, nil
}

func sql_TANH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Tanh(v))

	}
	return result, nil
}
