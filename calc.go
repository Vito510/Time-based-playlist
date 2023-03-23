package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

type Song struct {
	Path     string  `json:"path"`
	Duration float64 `json:"duration"`
}

type Config struct {
	MIN_LENGTH  float64  `json:"MIN_LENGTH"`
	MAX_LENGTH  float64  `json:"MAX_LENGTH"`
	FINAL_SONGS []string `json:"FINAL_SONGS"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func random_element(arr []string) string {
	rand.Seed(time.Now().UnixNano())
	return arr[rand.Intn(len(arr))]
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
		// for _, v := range partial {
		// 	fmt.Printf("%s\n", v.Path)
		// }
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

var PRECISION uint = 1

func test_range(songs []Song) int {

	var songs_len float64 = 0
	var avg [10]int
	var playlist []Song
	for _, v := range songs {
		songs_len += v.Duration
	}
	var s int = 0

	color.Green("Starting range test")

	for i := 5; i <= 24*60; i += 5 {

		if i*60 > int(songs_len*0.98) {
			color.Red("Failed: playlist duration might be unreachable")
			return 0
		}

		fmt.Printf("%4d | %4.1fhrs: ", i, float64(i)/60)
		for j := 0; j < 10; j++ {

			songs = shuffle(songs)

			start_time := time.Now()
			playlist = subset_sum(songs, float64(i*60), []Song{})
			end_time := time.Since(start_time)
			avg[j] = int(end_time.Milliseconds())
			fmt.Printf("%3d (%3d) ", end_time.Milliseconds(), len(playlist))
		}
		s = 0
		for j := 0; j < 10; j++ {
			s += avg[j]
		}
		fmt.Printf("| %5.1fms\n", float64(s)/float64(len(avg)))

		if float64(s)/float64(len(avg)) > 500 {
			color.Red("Failed: average to high")
			return 0
		}

	}

	return 1
}

func test_precise(songs []Song) int {

	var songs_len float64 = 0
	var avg [10]int
	var playlist []Song
	for _, v := range songs {
		songs_len += v.Duration
	}
	var s int = 0

	color.Green("Starting precise test")

	for i := 15 * 60; i <= 60*60; i += 1 {

		if i > int(songs_len*0.99) {
			color.Red("Failed: playlist duration might be unreachable")
			return 0
		}

		fmt.Printf("%4.2fmin : ", float64(i)/60)
		for j := 0; j < 10; j++ {

			songs = shuffle(songs)

			start_time := time.Now()
			playlist = subset_sum(songs, float64(i), []Song{})
			end_time := time.Since(start_time)
			avg[j] = int(end_time.Milliseconds())
			fmt.Printf("%3d (%3d) ", end_time.Milliseconds(), len(playlist))
		}
		s = 0
		for j := 0; j < 10; j++ {
			s += avg[j]
		}
		fmt.Printf("| %5.1fms\n", float64(s)/float64(len(avg)))

		if float64(s)/float64(len(avg)) > 500 {
			color.Red("Failed: average to high")
			return 0
		}

	}
	return 1
}

func test(songs []Song) {

	passed := 0

	passed += test_range(songs)
	passed += test_precise(songs)

	if passed == 2 {
		color.Green("Passed")
	} else if passed == 0 {
		color.Red("Failed")
	} else {
		color.Yellow("Partial pass")
	}

}

func get_song_index(path string, songs []Song) int {
	for i, v := range songs {
		if path == v.Path {
			return i
		}
	}
	return -1
}

func main() {

	// Read config file

	jsonFile, err := os.Open("config.json")
	check(err)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	var FINAL_SONG_INDEX int
	var FINAL_SONG string = ""

	json.Unmarshal(byteValue, &config)

	// Read json file
	jsonFile, err = os.Open("songs.json")
	check(err)
	defer jsonFile.Close()

	byteValue, _ = ioutil.ReadAll(jsonFile)
	var songs []Song

	json.Unmarshal(byteValue, &songs)

	var temp []Song
	// Apply min, max song length
	for i := 0; i < len(songs); i++ {
		if songs[i].Duration < config.MAX_LENGTH && songs[i].Duration > config.MIN_LENGTH {
			temp = append(temp, songs[i])
		}
	}
	songs = temp

	if len(config.FINAL_SONGS) > 0 {
		FINAL_SONG_INDEX = get_song_index(random_element(config.FINAL_SONGS), songs)
	}

	if len(os.Args) > 1 && os.Args[1] == "test" {
		test(songs)
		return
	}

	var hour int = 0
	var minute int = 0

	for hour+minute == 0 {
		fmt.Printf("Enter end time [h:m]: ")
		fmt.Scanf("%d:%d", &hour, &minute)
	}

	now := time.Now()
	end_target_time := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

	target := roundFloat(end_target_time.Sub(time.Now()).Seconds(), PRECISION)

	if FINAL_SONG_INDEX != -1 {
		if target-songs[FINAL_SONG_INDEX].Duration >= 0 {
			target -= songs[FINAL_SONG_INDEX].Duration
			fmt.Printf("\n%f\n", target)
			target = roundFloat(target, PRECISION)
			FINAL_SONG = songs[FINAL_SONG_INDEX].Path
		} else {
			color.Yellow("Final song will not fit, skipping")
			FINAL_SONG_INDEX = -1
		}
	}

	songs = shuffle(songs)
	subset_time_start := time.Now()
	playlist := subset_sum(songs, float64(target), []Song{})
	subset_time_end := time.Now()

	if len(playlist) == 0 {
		color.Red("Couldn't fit songs into target time")
		time.Sleep(time.Duration(time.Second * 2))
		os.Exit(0)
	}

	// write to file
	write_time_start := time.Now()
	f, err := os.Create("playlist.m3u")

	check(err)
	f.WriteString("#EXTM3U\n")
	for _, v := range playlist {
		f.WriteString(v.Path + "\n")
	}

	if FINAL_SONG != "" {
		f.WriteString(FINAL_SONG)
		playlist = append(playlist, Song{})
	}

	f.Sync()
	f.Close()
	write_time_end := time.Now()

	cmd := exec.Command("cmd", "/C start playlist.m3u")
	cmd.Run()

	fmt.Printf("\nCreated playlist with %d items\n\nParameters:\n\tMin song length: %f\n\tMax song length: %f", len(playlist), config.MIN_LENGTH, config.MAX_LENGTH)
	fmt.Printf("\nTime:\n\tSubset: %dms\n\tWrite: %dms", subset_time_end.Sub(subset_time_start).Milliseconds(), write_time_end.Sub(write_time_start).Milliseconds())

	time.Sleep(time.Duration(time.Second * 3))

}
