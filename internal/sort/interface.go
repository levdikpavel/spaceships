package sort

type IntSort interface {
	Sort([]int)
	SortMethod() Method
}

type Fabric interface {
	CreateIntSort() IntSort
}
