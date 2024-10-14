package condition

import (
	"strconv"
	"time"
)

type stringValues []string

func (vs stringValues) intN(idx int) (int, bool) {
	if n, ok := vs.int64N(idx); ok {
		return int(n), true
	}
	return 0, false
}

func (vs stringValues) ints() []int {
	var ret []int
	for _, v := range vs.int64s() {
		ret = append(ret, int(v))
	}
	return ret
}

func (vs stringValues) int64N(idx int) (int64, bool) {
	str, ok := vs.getN(idx)
	if !ok {
		return 0, false
	}
	num, err := strconv.ParseInt(str, 10, 64)
	return num, err == nil
}

func (vs stringValues) int64s() []int64 {
	ret := make([]int64, 0, len(vs))
	for _, v := range vs {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			ret = append(ret, n)
		}
	}
	return ret
}

func (vs stringValues) timeN(idx int) (time.Time, bool) {
	var at time.Time
	str, ok := vs.getN(idx)
	if !ok {
		return at, false
	}
	err := at.UnmarshalText([]byte(str))

	return at, err == nil
}

func (vs stringValues) times() []time.Time {
	ret := make([]time.Time, 0, len(vs))
	for _, v := range vs {
		var at time.Time
		if err := at.UnmarshalText([]byte(v)); err == nil {
			ret = append(ret, at)
		}
	}
	return ret
}

func (vs stringValues) boolN(idx int) (bool, bool) {
	str, ok := vs.getN(idx)
	if !ok {
		return false, false
	}
	ret, err := strconv.ParseBool(str)

	return ret, err == nil
}

func (vs stringValues) getN(idx int) (string, bool) {
	sz := len(vs)
	if sz <= idx {
		return "", false
	}
	return vs[idx], true
}