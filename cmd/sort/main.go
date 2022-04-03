package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"

	"modules/internal/sort"
)

var (
	inputFile  string
	outputFile string
	sortMethod string
)

func init() {
	flag.StringVar(&inputFile, "i", "input.txt", "input file")
	flag.StringVar(&outputFile, "o", "output.txt", "output file")
	flag.StringVar(&sortMethod, "m", "SelectionSort", "sort method")
}

func main() {
	fabric := sort.NewFabric(sortMethod)

	nums, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	sort := fabric.CreateIntSort()
	sort.Sort(nums)

	err = dumpResult(nums)
	if err != nil {
		log.Fatal(err)
	}
}

func readInput() (result []int, err error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		num, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, err
		}

		result = append(result, int(num))
	}

	err = scanner.Err()
	return
}

func dumpResult(nums []int) (err error) {
	f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}

	_, err = f.WriteString(sortMethod + "\n")
	for _, num := range nums {
		numStr := strconv.Itoa(num)
		_, err = f.WriteString(numStr + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
