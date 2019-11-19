package ips

// CompareBytes returns -1 if b1 is less than b2, 0 if they are equal,
// or 1 if b1 is more than b2. The comparation is done from right to left, and
// if an input is shorter than the other, the shorter one is extended with zeros.
// If b1 or b2 are nil, nothing makes sense and we panic
// TODO: Change for builtin
func CompareBytes(b1, b2 []byte) int8 {
	if b1 == nil && b2 == nil {
		panic("cannot compare with nil array")
	}
	var i, j int
	if len(b1) < len(b2) {
		i = len(b1) - len(b2)
	} else if len(b1) > len(b2) {
		j = len(b2) - len(b1)
	}
	for {
		if i < 0 && j >= 0 && b2[j] > 0 {
			return -1
		} else if j < 0 && i >= 0 && b1[i] > 0 {
			return 1
		} else if j >= len(b1) && i >= len(b2) {
			return 0
		} else if b1[i] < b2[j] {
			return -1
		} else if b1[i] > b2[j] {
			return 1
		}
		i++
		j++
	}
}

// Operatebytes executes a byte operation between two arrays, from right to left.
// If length of both byte arrays is different, it is assumed that the shorter byte array
// is extended with zeroes.
func OperateBytes(b1, b2 []byte, f func(b1, b2 byte) byte) []byte {
	if b1 == nil || b2 == nil {
		return nil
	}
	i, j := len(b1), len(b2)
	var k int
	var maxLen int
	if i > j {
		maxLen = i
	} else {
		maxLen = j
	}
	k = maxLen
	result := make([]byte, maxLen)

L:
	for {
		i--
		j--
		k--
		switch {
		case i >= 0 && j >= 0:
			result[k] = f(b1[i], b2[j])
		case i < 0 && j >= 0:
			result[k] = 0 | b2[j]
		case j < 0 && i >= 0:
			result[k] = 0 | b1[i]
		default:
			break L
		}
	}
	return result
}
