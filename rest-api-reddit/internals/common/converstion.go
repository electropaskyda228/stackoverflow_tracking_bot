package common

import "strconv"

func StringToUint(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(u64), nil
}

func UintToString(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}
