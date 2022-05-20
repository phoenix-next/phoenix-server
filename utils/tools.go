package utils

func GetKeys(m map[uint64]int) []uint64 {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	keys := make([]uint64, len(m))
	for k := range m {
		keys[j] = k
		j++
	}
	return keys

}
