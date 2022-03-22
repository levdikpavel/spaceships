package vector

type Vector []int

func New(values []int) Vector {
	v := make(Vector, len(values))
	copy(v, values)
	return v
}

func Add(v1, v2 Vector) Vector {
	length := min(len(v1), len(v2))
	result := make(Vector, length)
	for i := 0; i < length; i++ {
		result[i] = v1[i] + v2[i]
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
