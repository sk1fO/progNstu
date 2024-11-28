package main

import (
	"fmt"
	"math/rand"
)

func main() {
	alpha := "АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯABCDEFGHIJKLMNOPQRSTUVWXYZ"

	alphaRune := []rune(alpha)

	rand.Shuffle(len(alphaRune), func(i, j int) {
		alphaRune[i], alphaRune[j] = alphaRune[j], alphaRune[i]
	})

	reverse := Reverse("TDLGЕЖKRНИUЛЪБФQДJВVEРАЦЗЮXХFЯNОYУOBПЩZHЬЭСIAWШЫSMТPЙКГЧCМ")
	fmt.Println(string(reverse))
}

func Reverse(s string) string {

	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {

		runes[i], runes[j] = runes[j], runes[i]

	}

	return string(runes)

}
