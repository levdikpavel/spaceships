package sort

func NewFabric(sortMethod string) Fabric {
	return fabric{
		method: sortMethod,
	}
}

type fabric struct {
	method string
}

func (f fabric) CreateIntSort() IntSort {
	switch f.method {
	case "SelectionSort":
		return selectionSort{}
	case "InsertionSort":
		return insertionSort{}
	case "MergeSort":
		return mergeSort{}
	default:
		panic("unknown sort method")
	}
}
