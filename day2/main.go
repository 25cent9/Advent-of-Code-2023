package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type numIndex struct {
	number     int
	start, end int
}

type numberIndexes []numIndex

var debugStrs = make(map[string][]string)

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
	// fmt.Printf("Matches: %v\n", input)
	re := regexp.MustCompile(pattern)
	for _, match := range re.FindAllStringIndex(input, -1) {
		start, end := match[0], match[1]
		num, err := strconv.Atoi(input[start:end])
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Mach: {%v: %v - %v}\n", num, start, end)
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

type lineTracker map[string]int

func (l lineTracker) at(line int) string {
	for text, num := range l {
		if num == line {
			return text
		}
	}
	return ""
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	strs := make(lineTracker)
	partIndexes := make(map[int]numberIndexes)
	symbolIndexes := map[int]symbolIndexes{}
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		str := scanner.Text()
		partIndexes[lineCount] = extractNumberIndexes(str)
		symbolIndexes[lineCount] = extractSymbolIndexes(str)
		strs[str] = lineCount
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	width, height := len(strs), len(strs)
	// fmt.Printf("WxH: %vx%v\n", width, height)
	allParts := []int{}
	for line, symbols := range symbolIndexes {
		// fmt.Printf("Line %v:\n", line+1)
		for _, symbol := range symbols {
			// fmt.Printf("Symbol %v [%v]\n", symbol.symbol, symbol.index)
			// Check on the same line
			if parts := partIndexes[line]; len(parts) > 0 {
				for _, part := range parts.match(symbol.index) {
					msg := fmt.Sprintf("Adjacent Part Match: {%v: %v - %v}", part.number, part.start, part.end)
					allParts = append(allParts, part.number)
					debugStrs[strs.at(line)] = append(debugStrs[strs.at(line)], msg)
					// fmt.Printf("Adjacent Part %v on %v for %v\n", part.number, line+1, symbol.symbol)
				}
			}
			leftDiagIndex, rightDiagIndex := symbol.index-1, symbol.index+1
			nextLine, prevLine := line+1, line-1
			for _, newLine := range []int{nextLine, prevLine} {
				if newLine != -1 && newLine < height {
					for _, diagIndex := range []int{leftDiagIndex, rightDiagIndex} {
						if diagIndex < width && diagIndex >= 0 {
							if parts := partIndexes[newLine]; len(parts) > 0 {
								matchingParts := map[int]bool{}
								for _, part := range parts.match(diagIndex) {
									if _, ok := matchingParts[part.number]; !ok {
										matchingParts[part.number] = true
										matchLine := strs.at(newLine)
										msg := fmt.Sprintf("Diagonal Part Match for %v on line %v: {%v: %v - %v}\n%s", symbol.symbol, newLine+1, part.number, part.start, part.end, matchLine)
										// allParts = append(allParts, part.number)
										debugStrs[strs.at(line)] = append(debugStrs[strs.at(line)], msg)
									}
								}
								for part := range matchingParts {
									allParts = append(allParts, part)
								}
							}
						}
					}
				}
			}
		}
	}
	partSum := 0
	for _, part := range allParts {
		partSum += part
	}
	for line, msgs := range debugStrs {
		fmt.Printf("\n--------\n%v\n--------\n%v", line, strings.Join(msgs, "\n"))
	}
	fmt.Printf("All Parts: %v\nPart Sum: %v", allParts, partSum)

}
