package sort

type Method uint8

const (
	SelectionSort Method = iota
	InsertionSort
	MergeSort
	endMethod
)

var methodNames = []string{
	"SelectionSort",
	"InsertionSort",
	"MergeSort",
}

func (m Method) String() string {
	if m < endMethod {
		return methodNames[m]
	}

	panic("Method is unknown")
}
