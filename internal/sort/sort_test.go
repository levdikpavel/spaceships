package sort

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSort(t *testing.T) {
	suite.Run(t, new(SortSuite))
}

type SortSuite struct {
	suite.Suite

	input   []int
	output  []int
	methods []string
}

func (s *SortSuite) SetupSuite() {
	s.input = []int{6, 5, 3, 1, 8, 7, 2, 4}
	s.output = []int{1, 2, 3, 4, 5, 6, 7, 8}
	s.methods = []string{
		"SelectionSort",
		"InsertionSort",
		"MergeSort",
	}
}

func (s *SortSuite) TestIntSort() {
	for _, method := range s.methods {
		fabric := NewFabric(method)
		sort := fabric.CreateIntSort()

		var input []int
		input = append(input, s.input...)
		sort.Sort(input)
		s.Require().Equal(s.output, input, method)
	}
}
