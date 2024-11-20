package helper

import "strconv"

func StringToUInt(str string) uint {
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {

	}
	return uint(num)
}

func UIntToString(num uint) string {
	return strconv.FormatUint(uint64(num), 10)
}
