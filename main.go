// Copyright © 2022 Rak Laptudirm <raklaptudirm@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/qeesung/image2ascii/ascii"
	"github.com/qeesung/image2ascii/convert"
)

func main() {
	// itp picture.jpg
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: itp <image path>")
		os.Exit(1)
	}

	fileName := os.Args[1]

	// only use numbers in the image
	// 7 is darkest and 8 is brightest
	// assumed black background with white text
	ascii.DefaultOptions.Pixels = []byte("7772299408")

	converter := convert.NewImageConverter()
	options := convert.DefaultOptions
	options.Colored = false // no ansi escapes

	formatted := converter.ImageFile2ASCIIString(fileName, &options)

	fmt.Println(formatted) // show preview

	img, width := formatImage(formatted) // remove word wrapping
	prime := searchPrime(img)            // search for similar prime

	// display
	fmt.Println(formatToImage(prime, width)) // number with wrapping
	fmt.Println(prime)                       // raw number
}

// formatImage removes word wrapping from the image and returns it as a
// []byte, along with the width of the word wrapping. It changes the last
// digit so that the number can be a prime, and also replaces any leading
// 0s with an 8.
func formatImage(img string) ([]byte, int) {
	// remove word wrapping line breaks
	raw := []byte(strings.ReplaceAll(img, "\n", ""))

	// prime numbers can't have 0, 2, 4, 5, 6, or 8 as their last digit,
	// therefore replace it with a digit of similar brightness if so
	lastDigitReplace := map[byte]byte{
		'0': '3',
		'2': '3',
		'4': '9',
		'6': '9',
		'8': '9',
		'5': '3',
	}

	lastIndex := len(raw) - 1
	lastDigit := raw[lastIndex]

	// selectively replace last digit
	if sub, notPrime := lastDigitReplace[lastDigit]; notPrime {
		raw[lastIndex] = sub
	}

	// make sure number does not have leading zeros
	if raw[0] == '0' {
		raw[0] = '8' // use digit with similar brightness
	}

	// get width from index of wrapping line break
	width := strings.Index(img, "\n")

	return raw, width
}

// formatToImage adds word wrapping to the given prime at the given width,
// essentially undoing what formatImage did.
func formatToImage(prime []byte, width int) string {
	var img string

	length := len(prime)

	// add word wrapping line breaks
	for i := 0; i < length; i += width {
		img += string(prime[i:i+width]) + "\n"
	}

	return img
}

// searchPrime searches for a prime number which "looks similar" to the
// given number. The []byte contains the ascii values of each digit of the
// number. The returned number is also in the same format.
func searchPrime(raw []byte) []byte {
	// seed the prng
	rand.Seed(time.Now().Unix())

	// iterations of Miller-Rabin to use for testing
	testCount := 100

	// create a copy to prevent manipulating the source slice
	img := make([]byte, len(raw))
	copy(img, raw)

	var n big.Int
	for i := 0; ; i++ {
		// print index of number and it's sha256 sum
		fmt.Printf("prime %5d: %x\n", i, sha256.Sum256(img))

		// check if number is prime
		if n.SetString(string(img), 10); n.ProbablyPrime(testCount) {
			return img
		}

		// slightly change number
		replaceRandomDigit(img)
	}
}

// replaceRandomDigit replaces a random digit of the given number with a
// digit of similar brightness. This excludes the digits at indexes 0 and
// len(number)-1, to prevent leading 0s and non-prime last digits.
func replaceRandomDigit(img []byte) {
retry:
	// replacement table with digits of similar brightness
	replacementTable := map[byte][]byte{
		'0': {'8', '9', '5'},
		'1': {'7'},
		'2': {'6'},
		'5': {'0'},
		'8': {'0', '9', '5'},
		'6': {'2'},
		'9': {'4', '3'},
		'4': {'9'},
	}

	// generate random index excluding 0 and len-1
	index := rand.Int()%(len(img)-2) + 1

	// if replacement found replace digit, or try again
	if subList, ok := replacementTable[img[index]]; ok {
		// select random entry
		subN := rand.Int() % len(subList)

		// replace digit
		sub := subList[subN]
		img[index] = sub
	} else {
		goto retry
	}
}
