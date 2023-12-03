package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type numIndex struct {
	number     int
	start, end int
}

type numberIndexes []numIndex

func (i numberIndexes) match(index int) numberIndexes {
	indexes := numberIndexes{}
	for _, number := range i {
		if inRange := (index >= number.start) && (index <= number.end); inRange {
			indexes = append(indexes, number)
		}
	}
	return indexes
}

func extractNumberIndexes(input string) numberIndexes {
	indexes := numberIndexes{}
	pattern := `\d+`
	re := regexp.MustCompile(pattern)
	for _, match := range re.FindAllStringIndex(input, -1) {
		start, end := match[0], match[1]
		num, err := strconv.Atoi(input[start:end])
		if err != nil {
			log.Fatal(err)
		}
		index := numIndex{
			number: num,
			start:  start,
			end:    end,
		}
		indexes = append(indexes, index)
	}
	return indexes
}

type symbolIndex struct {
	symbol string
	index  int
}

type symbolIndexes []symbolIndex

func extractSymbolIndexes(input string) symbolIndexes {
	indexes := symbolIndexes{}
	for index, character := range input {
		if charValue := int(character); charValue == 46 || (charValue <= 57 && charValue >= 48) {
			continue
		}
		indexes = append(indexes, symbolIndex{symbol: string(character), index: index})
	}
	return indexes
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	strs := []string{}
	partIndexes := make(map[int]numberIndexes)
	symbolIndexes := map[int]symbolIndexes{}
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		str := scanner.Text()
		partIndexes[lineCount] = extractNumberIndexes(str)
		symbolIndexes[lineCount] = extractSymbolIndexes(str)
		strs = append(strs, str)
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	width, height := len(strs[0]), len(strs)
	for line, symbols := range symbolIndexes {
		fmt.Printf("Line %v:\n", line+1)
		for _, symbol := range symbols {
			// fmt.Printf("Symbol %v [%v]\n", symbol.symbol, symbol.index)
			// Check on the same line
			if parts := partIndexes[line]; len(parts) > 0 {
				for _, part := range parts.match(symbol.index) {
					fmt.Printf("Adjacent Part %v on %v for %v\n", part.number, line+1, symbol.symbol)
				}
			}
			leftDiagIndex, rightDiagIndex := symbol.index-1, symbol.index+1
			nextLine, prevLine := line+1, line-1
			for _, newLine := range []int{nextLine, prevLine} {
				if newLine <= height {
					for _, diagIndex := range []int{leftDiagIndex, rightDiagIndex} {
						if diagIndex <= width {
							if parts := partIndexes[newLine]; len(parts) > 0 {
								for _, part := range parts.match(diagIndex) {
									fmt.Printf("Diagonal Part %v on %v for %v\n", part.number, newLine+1, symbol.symbol)
								}
							}
						}
					}
				}
			}
		}
	}
}
