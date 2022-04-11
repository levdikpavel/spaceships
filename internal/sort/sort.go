package sort

type selectionSort struct {}

func (selectionSort) SortMethod() Method {
	return SelectionSort
}

func (selectionSort) Sort(nums []int) {
	n := len(nums)

	for i := 0; i < n; i++ {
		min := i
		for j := i + 1; j < n; j++ {
			if nums[j] < nums[min] {
				min = j
			}
		}
		nums[i], nums[min] = nums[min], nums[i]
	}
}

type insertionSort struct {}

func (insertionSort) SortMethod() Method {
	return InsertionSort
}

func (insertionSort) Sort(nums []int) {
	for i := 1; i < len(nums); i++ {
		current := nums[i]
		j := i
		for ; j > 0 && nums[j-1] > current; j-- {
			nums[j] = nums[j-1]
		}
		nums[j] = current
	}
}

type mergeSort struct {}

func (mergeSort) SortMethod() Method {
	return MergeSort
}

func (s mergeSort) Sort(list []int) {
	count := len(list)

	switch {
	case count > 2:
		s.Sort(list[:count/2])
		s.Sort(list[count/2:])
		//lb := s.Sort(list[:count/2])
		//rb := mergeSort(list[count/2:])
		//list = append(lb, rb...)
		lastIndex := len(list) - 1

		for i := 0; i < lastIndex; i++ {
			mv := list[i]
			mi := i

			for j := i + 1; j < lastIndex+1; j++ {
				if mv > list[j] {
					mv = list[j]
					mi = j
				}
			}

			if mi != i {
				list[i], list[mi] = list[mi], list[i]
			}
		}

	case len(list) > 1 && list[0] > list[1]:
		list[0], list[1] = list[1], list[0]
	}
}
