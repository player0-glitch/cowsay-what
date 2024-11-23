package cows

import (
	"bufio"
	"cow-bot/asset"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

// defining what a speech bubble type should have and look like
type bubble struct {
	width      int
	startL     rune
	startR     rune
	bigSideL   rune
	bigSideR   rune
	endL       rune
	endR       rune
	smallSideL rune
	smallSideR rune
	speech     string
}

var (
	say = &bubble{
		0,
		'/', '\\',
		'|', '|',
		'\\', '/',
		'<', '>',
		`\`,
	}
	think = &bubble{
		0,
		'(', ')',
		'(', ')',
		'(', ')',
		'(', ')',
		`o`,
	}
)

// read the file with all the fortunes
func ReadFortuneFiles(path string) []string {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	allQuotes := string(file)
	result := strings.FieldsFunc(allQuotes, func(r rune) bool {
		return r == '%'
	})
	return result
}

// function to load different cows by name
// this code was stolen from https://github.com/msmith491/go-cowsay/blob/master/src/go-cowsay/main.go
func GetCowAsset(cowName string) []byte {
	if strings.HasSuffix(cowName, ".cow") {
		// getting the cow data
		data, err := assets.Asset(cowName)
		// error checking, we don't work with broken files
		if err != nil {
			fmt.Println("Couldn't Load Asset ", err)
		}
		return data
	} else {
		// here we load the cowName we want
		name := fmt.Sprintf("data/cows/%s.cow", cowName)
		data, err := os.ReadFile(name)
		if err != nil {
			fmt.Println("Couldn't Access Asset", err)
		}
		return data
	}
}

// Random Number Generator [min,max)
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// This function returns a random pair of eyes as a string
func GetEyes() string {
	eyes := [9]string{"==", "$$", "@@", "**", "--", "OO", "oO", "..", "xx"}

	return eyes[RandomInt(0, len(eyes))]
}

func GetAnimal() string {
	animals := [30]string{
		"default", "bunny", "cower", "daemon", "kitty",
		"moose", "moofasa", "moose", "sheep", "skeleton",
		"www", "tux", "vader", "meow", "satanic",
		"hellokitty", "ghostbusters", "dragon", "flaming-sheep", "elephant",
		"bunny", "squirrel", "sodomized", "stegosaurus", "eyes",
		"ren", "mutilated", "small", "telebears", "three-eyes",
	}

	return animals[RandomInt(0, len(animals))]
}

// Returns a random qoute from our fortunes
func GetFortune(fortunes []string) string {
	size := len(fortunes)
	return fortunes[RandomInt(0, size)]
}

func FormatAnimal(str string) string {
	var animal string
	animal = strings.Replace(str, "$eyes", GetEyes(), -1)
	animal = strings.Replace(animal, "${eyes}", GetEyes(), -1)
	animal = strings.Replace(animal, "$thoughts", "o", -1)
	animal = strings.Replace(animal, "$tongue", " ", -1)
	animal = strings.Replace(animal, "\\\\", "\\", -1)
	animal = strings.Replace(animal, "\\@", "@", -1)

	// there are dangling EOC characters
	// some animals have unnecessary hashtags, we'll just use the function
	// to remove EOC to remove custom delim we pass to it
	animal = RemoveDanglingEOC(animal, "#")
	return RemoveDanglingEOC(animal, "EOC")
}

// This is a great little recursive word wrapping algo, courtesy of Peter
// Mortensen: https://stackoverflow.com/a/857770
func Wrap(text string, width int) []string {
	text = strings.TrimSpace(text)
	if len(text) <= width {
		return []string{text}
	} else {
		isplit := width
		for i := width; i > 0; i-- {
			if text[i] == ' ' {
				isplit = i
				break
			}
		}
		before := strings.TrimRight(text[:isplit], " ")
		after := strings.TrimLeft(text[isplit:], " ")
		return append([]string{before}, Wrap(after, width)...)
	}
}

// this function takes in our fortune and wraps it around a speech bubble
func MakeSpeechBubble(str string) string {
	var (
		longestLine int
		speech      = make([]string, 0)
		bubble      *bubble // this is a struct type
	)
	// i'll change this one day to allow for a thinking cow
	bubble = say
	// padding makes the speech bubble look neat
	pad := " "
	bubble.width = 40
	hspace := regexp.MustCompile(`[ \t]+`)
	vspace := regexp.MustCompile(`[r\n]{2,}`)
	str = hspace.ReplaceAllString(str, " ")
	str = vspace.ReplaceAllString(str, "\n\n")

	// fortunes are broken into new lines
	fortuneLines := strings.Split(str, "\n")
	innerTextWidth := bubble.width - 1
	// implement word wrapping and get the longest line in the wrapped text
	for _, line := range fortuneLines {
		// if len(line) > longestLine {
		// 	longestLine = len(line)
		// }
		// speech = append(speech, line)
		wordLines := Wrap(line, innerTextWidth)
		for _, wordLine := range wordLines {
			if len(wordLine) > longestLine {
				longestLine = len(wordLine)
			}
			speech = append(speech, wordLine)
		}
	}
	// padding each line to the length of the longest line and the bubble sides
	if len(speech) == 1 {
		wordLine := speech[0]
		left, right := bubble.smallSideL, bubble.smallSideR

		// we'll be using this to build our speech bubble
		var builder strings.Builder
		builder.WriteRune(left)
		builder.WriteString(pad) // adding padding in our bubble
		builder.WriteString(wordLine)
		builder.WriteString(pad)
		builder.WriteRune(right)
		speech[0] = builder.String()
	} else {
		for i, wordLine := range speech {
			// runes are like chars (int32 type)
			var left, right rune
			switch {
			case i == 0:
				left, right = bubble.startL, bubble.startR
			case i == len(speech)-1:
				left, right = bubble.endL, bubble.endR
			default:
				left, right = bubble.bigSideL, bubble.bigSideR
			}
			// creating room 'in' the bubble
			room := strings.Repeat(pad, longestLine-len(wordLine))

			var builder strings.Builder
			builder.WriteRune(left)
			builder.WriteString(pad)
			builder.WriteString(wordLine)
			builder.WriteString(room)
			builder.WriteString(pad)
			builder.WriteRune(right)
			speech[i] = builder.String()
		}
	}

	// addding padding on the top and bottom of the speech bubble
	var b strings.Builder
	w := longestLine + 2*len(pad)

	b.WriteString(pad)
	b.WriteString(strings.Repeat("_", w))
	b.WriteByte('\n')
	b.WriteString(strings.Join(speech, "\n"))
	b.WriteByte('\n')
	b.WriteString(pad)
	b.WriteString(strings.Repeat("-", w))
	return b.String()
}

func RemoveDanglingEOC(str string, delim string) string {
	var result strings.Builder
	eoc := delim
	scanner := bufio.NewScanner(strings.NewReader(str))

	// iterate through the scanner class
	for scanner.Scan() {
		line := scanner.Text()
		// ignore anyline with EOC
		// this function is designed to ONLY be used after formatAnimal function is ran
		if !strings.Contains(line, eoc) {
			result.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning to remove EOC ", err)
	}
	return result.String()
}
