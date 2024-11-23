package main

import (
	"cow-bot/cows"
	"fmt"
)

func main() {
	fortunes := cows.ReadFortuneFiles("data/fortunes/mathematics.txt")

	animal := cows.GetCowAsset(cows.GetAnimal())
	s := string(animal)
	formatedAnimal := cows.FormatAnimal(s)

	// fmt.Println("fortunes\n", getFortune(fortunes))
	fmt.Println(cows.MakeSpeechBubble(cows.GetFortune(fortunes)))
	fmt.Println(formatedAnimal)
}
