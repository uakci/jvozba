package main

import (
	"bufio"
	"fmt"
	"os"

	"jvozba"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in := scanner.Text()
		res, err := jvozba.Jvozba(in)
		if err != nil {
			fmt.Printf("got error: %v\n", err)
		} else {
			fmt.Printf("%s → %s (%d)\n", in, res, jvozba.Score(res))
		}
	}
}
