package main

import (
	"bufio"
	"log"
	"os"
)

type color int

const (
	red color = iota
	blue
	green
)

type session struct {
	green int
	blue  int
	red   int
}

type game struct {
	number   int
	sessions []session
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

}
