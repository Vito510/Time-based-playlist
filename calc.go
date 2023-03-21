package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"
)

type Song struct {
	Path     string  `json:"path"`
	Duration float64 `json:"duration"`
}

func shuffle(arr []Song) []Song {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func subset_sum(songs []Song, target float64, partial []Song) []Song {
	s := 0.0
	for _, v := range partial {
		s += v.Duration
	}
	s = roundFloat(s, 1)

	if s == target {
		for _, v := range partial {
			fmt.Printf("%s\n", v.Path)
		}
		return partial
	}
	if s >= target {
		return nil
	}

	for i := 0; i < len(songs); i++ {
		n := songs[i]
		remaining := songs[i+1:]

		r := subset_sum(remaining, target, append(partial, n))
		if r != nil {
			return r
		}
	}
	return nil
}

var MAX_LENGTH float64 = 60 * 5
var MIN_LENGTH float64 = 60 * 2

func main() {

	// Read json file
	jsonFile, err := os.Open("songs.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var songs []Song

	json.Unmarshal(byteValue, &songs)

	var temp []Song
	// Apply min, max song length
	for i := 0; i < len(songs); i++ {
		if songs[i].Duration < MAX_LENGTH && songs[i].Duration > MIN_LENGTH {
			temp = append(temp, songs[i])
		}
	}
	songs = temp

	songs = shuffle(songs)

	subset_sum(songs, 120*60, []Song{})
}
