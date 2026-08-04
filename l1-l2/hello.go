package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	input := "There once was a cat named Barry. He was a very good cat. This cat lived in Boston. He loved doing Boston-related activities (that were good for cats). He walked the esplanade. He shopped on Newbury. He ate at Tatte. He sometimes even went to TD Garden. Did you know that cats are not allowed in TD Garden?"

	for i := 0; i <= len(input); i++ {
		if input[i:i+3] == "cat" {
			fmt.Printf("found cat @ %v\n", i)
		}
	}
}
