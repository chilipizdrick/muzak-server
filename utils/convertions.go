package utils

import "github.com/lib/pq"

func Int64SliceToUIntSlice(arr []int64) []uint {
	if arr == nil {
		return nil
	}
	res := make([]uint, len(arr))
	for i, v := range arr {
		res[i] = uint(v)
	}
	return res
}

func PQInt64ArrayPtrToInt64Slice(arr pq.Int64Array) []int64 {
	if arr == nil {
		return nil
	} else {
		return []int64(arr)
	}
}

func PQInt64ArrayPtrToUIntSlice(arr pq.Int64Array) []uint {
	tmpRes := PQInt64ArrayPtrToInt64Slice(arr)
	if tmpRes == nil {
		return nil
	}
	return Int64SliceToUIntSlice(tmpRes)
}
