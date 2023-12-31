package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("this is a demo file\n")
	fmt.Printf("// --- begin %s ---\n", os.Args[1])
	fmt.Printf("%s\n", os.Args[2])
	fmt.Printf("// --- end %s ---\n", os.Args[1])
	fmt.Printf("this is line 5 of the demo output\n")
	fmt.Printf("and some other text.\n")
}
