package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

var (
	nums = map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
	}
)

type numTracker struct {
	value int
	index int
}

func (n *numTracker) String() string {
	return fmt.Sprintf("[V: %v, I: %v]", n.value, n.index)
}
func (n *numTracker) Empty() bool {
	return n.value == 0 && n.index == 0
}
func findDigit(input string) []numTracker {
	foundNums := []numTracker{}
	for i, char := range input {
		val := int(char)
		if val >= 48 && val <= 57 {
			foundNums = append(foundNums, numTracker{
				value: val - 48,
				index: i,
			})
		}
	}
	return foundNums
}
func findNumSubStr(input string) []numTracker {
	var matches = make(map[string][]int)
	foundNums := []numTracker{}
	for num := range nums {
		pattern := regexp.MustCompile(fmt.Sprintf(`(?i)%s`, num))
		for _, foundMatches := range pattern.FindAllStringIndex(input, -1) {
			sort.Ints(foundMatches)
			matches[num] = append(matches[num], foundMatches[0])
		}
	}
	for num, indexes := range matches {
		for _, index := range indexes {
			foundNums = append(foundNums, numTracker{
				value: nums[num],
				index: index,
			})
		}
	}
	return foundNums
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	strs := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strs = append(strs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	sum := 0
	for _, str := range strs {
		wg.Add(1)
		str := str
		go func() {
			defer wg.Done()
			foundNumbers := findDigit(str)
			foundNumbers = append(foundNumbers, findNumSubStr(str)...)
			indexedDigits := make(map[int]numTracker)
			indexKeys := []int{}
			for _, digit := range foundNumbers {
				if digit.Empty() {
					continue
				}
				index := digit.index
				indexedDigits[index] = digit
				indexKeys = append(indexKeys, index)
			}
			sort.Ints(indexKeys)
			firstIndex, lastIndex := indexKeys[0], indexKeys[len(indexKeys)-1]
			first, last := indexedDigits[firstIndex], indexedDigits[lastIndex]
			digitStr := fmt.Sprintf("%d%d", first.value, last.value)
			digit, err := strconv.Atoi(digitStr)
			if err != nil {
				log.Fatal(err)
			}
			sum += digit
		}()
	}
	wg.Wait()
	fmt.Printf("Final Sum: %v", sum)
}
