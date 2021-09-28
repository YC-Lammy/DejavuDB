package datastore

import (
	"errors"
	"src/types"
	"src/types/int128"
	"src/types/uint128"
	"strconv"

	"rogchap.com/v8go"
)

func JsGet(ctx *v8go.Context, key string) (*v8go.Value, error) {
	dtype, p := Get(key)
	if p == nil && dtype == 0x00 {
		return nil, errors.New("undefined key " + key)
	}
	switch dtype {

	case types.String:
		return ctx.RunScript("'"+*(*string)(p)+"'", "string.js")

	case types.Int, types.Int64:
		return ctx.RunScript(strconv.FormatInt(*(*int64)(p), 10), "string.js")
	case types.Int32:
		return ctx.RunScript(strconv.FormatInt(int64(*(*int32)(p)), 10), "string.js")
	case types.Int16:
		return ctx.RunScript(strconv.FormatInt(int64(*(*int16)(p)), 10), "string.js")
	case types.Int8:
		return ctx.RunScript(strconv.FormatInt(int64(*(*int8)(p)), 10), "string.js")
	case types.Int128:
		return ctx.RunScript("BigInt('"+(*(*int128.Int128)(p)).String()+"')", "string.js")
	case types.Uint, types.Uint64:
		return ctx.RunScript(strconv.FormatUint(*(*uint64)(p), 10), "string.js")
	case types.Uint32:
		return ctx.RunScript(strconv.FormatUint(uint64(*(*uint32)(p)), 10), "string.js")
	case types.Uint16:
		return ctx.RunScript(strconv.FormatUint(uint64(*(*uint16)(p)), 10), "string.js")
	case types.Uint8:
		return ctx.RunScript(strconv.FormatUint(uint64(*(*uint8)(p)), 10), "string.js")
	case types.Uint128:
		return ctx.RunScript("BigInt('"+(*(*uint128.Uint128)(p)).String()+"')", "string.js")
	case types.Decimal, types.Decimal64:
	case types.Decimal32:
	case types.Decimal128:
	case types.Float, types.Float64:
		return ctx.RunScript(strconv.FormatFloat(*(*float64)(p), 'f', 15, 64), "string.js")
	case types.Float32:
		return ctx.RunScript(strconv.FormatFloat(float64(*(*float32)(p)), 'f', 15, 32), "string.js")
	case types.Float128:

	case types.Byte:
	case types.Byte_arr:
		return ctx.RunScript("(new TextEncoder).encode('"+string(*(*[]byte)(p))+"')", "string.js")
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
		return ctx.RunScript("null", "null.js")
	}
	return nil, nil
}
