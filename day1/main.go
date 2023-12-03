package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type color int

const (
	red color = iota
	blue
	green
)

var (
	stringToColor = map[string]color{
		"red":   red,
		"blue":  blue,
		"green": green,
	}
	colorToString = map[color]string{
		red:   "red",
		blue:  "blue",
		green: "green",
	}
	colorConstraints = map[color]int{
		red:   12,
		green: 13,
		blue:  14,
	}
)

type gameSession map[color]int

func (gs gameSession) String() string {
	out := ""
	for color, count := range gs {
		out += fmt.Sprintf("%s: %v\n", colorToString[color], count)
	}
	return out
}

type game struct {
	number   int
	sessions []gameSession
}

func (g *game) String() string {
	out := fmt.Sprintf("Game %v\n", g.number)
	for _, session := range g.sessions {
		out += session.String()
	}
	return out
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	games := map[int]game{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := strings.Split(scanner.Text(), ": ")
		gameStr := strings.Split(str[0], " ")
		cubeDrawStr := str[1]

		gameSessions := []gameSession{}
		for _, cubeDraws := range strings.Split(cubeDrawStr, "; ") {
			sesh := make(gameSession)
			for _, cubeDraw := range strings.Split(cubeDraws, ", ") {
				draw := strings.Split(cubeDraw, " ")
				numCubes, err := strconv.Atoi(draw[0])
				if err != nil {
					log.Fatal(err)
				}
				sesh[stringToColor[draw[1]]] = numCubes
			}
			gameSessions = append(gameSessions, sesh)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		gameNum, err := strconv.Atoi(gameStr[1])
		if err != nil {
			log.Fatal(err)
		}
		games[gameNum] = game{
			number:   gameNum,
			sessions: gameSessions,
		}
	}
	possilbeGames := map[int]bool{}
	cubeMins := map[int]map[color]int{}
	for _, game := range games {
		possible := true
		minCubes := make(map[color]int)
		for _, session := range game.sessions {
			for cubeColor, num := range session {
				if val := minCubes[cubeColor]; val < num {
					minCubes[cubeColor] = num
				}
				if constraint := colorConstraints[cubeColor]; num > constraint {
					possible = false
				}
			}
		}
		if possible {
			possilbeGames[game.number] = true
		}
		cubeMins[game.number] = minCubes
	}
	sum := 0
	for game := range possilbeGames {
		sum += game
	}
	powerSetSum := 0

	for _, mins := range cubeMins {
		powerSet := 1
		for _, min := range mins {
			powerSet *= min
		}
		powerSetSum += powerSet
	}
	fmt.Printf("Final Possible Game Sum: %v\nFinal Min Cube Powerset: %v\n", sum, powerSetSum)
}
