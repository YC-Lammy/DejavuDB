package yaegiAPI

import (
	"errors"
	"src/types"
	"src/types/int128"
	"strconv"
	"unsafe"
)

type value struct {
	Ptr   unsafe.Pointer
	Dtype byte
}

func (v value) String() string {
	switch v.Dtype {
	case types.String:
		return *(*string)(v.Ptr)
	case types.Int, types.Int64:
		return strconv.FormatInt(*(*int64)(v.Ptr), 10)
	case types.Int32:
		return strconv.FormatInt(int64(*(*int32)(v.Ptr)), 10)
	case types.Int16:
		return strconv.FormatInt(int64(*(*int16)(v.Ptr)), 10)
	case types.Int8:
		return strconv.FormatInt(int64(*(*int8)(v.Ptr)), 10)
	case types.Int128:
	case types.Uint, types.Uint64:
		return strconv.FormatUint(*(*uint64)(v.Ptr), 10)
	case types.Uint32:
		return strconv.FormatUint(uint64(*(*uint32)(v.Ptr)), 10)
	case types.Uint16:
		return strconv.FormatUint(uint64(*(*uint16)(v.Ptr)), 10)
	case types.Uint8:
		return strconv.FormatUint(uint64(*(*uint8)(v.Ptr)), 10)
	case types.Uint128:
	case types.Decimal, types.Decimal64:
	case types.Decimal32:
	case types.Decimal128:
	case types.Float, types.Float64:
		return strconv.FormatFloat(*(*float64)(v.Ptr), 'f', 15, 64)
	case types.Float32:
	case types.Float128:
	case types.Byte:
	case types.Byte_arr:
	case types.Bool:
	case types.Graph:
	case types.Table:
	case types.Json:
	case types.SmartContract:
	case types.Contract:
	case types.Money:
	case types.SmallMoney:
	case types.Time:
	case types.Date:
	case types.Datetime:
	case types.Smalldatetime:
	case types.Null:
	}
	return ""
}

func (v value) Add(substance interface{}) error {

	switch val := substance.(type) {
	case string:
		if v.Dtype == types.String {
			*(*string)(v.Ptr) += val
		} else {
			return errors.New("Add: type mismatch, expected string")
		}
	case int:
		switch v.Dtype {
		case types.Int, types.Int64:
			*(*int64)(v.Ptr) += int64(val)
		case types.Int32:
			*(*int32)(v.Ptr) += int32(val)
		case types.Int16:
			*(*int16)(v.Ptr) += int16(val)
		case types.Int8:
			*(*int8)(v.Ptr) += int8(val)
		case types.Int128:
			(*int128.Int128)(v.Ptr).Add(val)
		}
	case int64:
	case int32:
	case int16:
	case int8:
	case uint:
	case uint64:
	case uint32:
	case uint16:
	case uint8:
	}
	switch v.Dtype {
	case types.String:
		if s, ok := substance.(string); ok {
			*(*string)(v.Ptr) += s
		} else {
			return errors.New("Add: type mismatch, expected string")
		}
	case types.Int, types.Int64, types.Int32, types.Int16, types.Int8, types.Int128:
		switch substance.(type) {
		case int:
		}
		if s, ok := substance.(int64); ok {
			*(*int64)(v.Ptr) += s
		} else {
			return errors.New("type mismatch, expected int64")
		}
	case types.Uint, types.Uint64:
	case types.Uint32:

	case types.Uint16:

	case types.Uint8:
	case types.Uint128:
	case types.Decimal, types.Decimal64:
	case types.Decimal32:
	case types.Decimal128:
	case types.Float, types.Float64:
	case types.Float32:
	case types.Float128:
	case types.Byte:
	case types.Byte_arr:
	case types.Bool:
	case types.Graph:
	case types.Table:
	case types.Json:
	case types.SmartContract:
	case types.Contract:
	case types.Money:
	case types.SmallMoney:
	case types.Time:
	case types.Date:
	case types.Datetime:
	case types.Smalldatetime:
	case types.Null:
	}
	return nil
}
