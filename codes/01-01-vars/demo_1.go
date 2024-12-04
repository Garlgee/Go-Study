package main

import "fmt"

func main() {
	const NAME = "TOM"
	const AGE int = 20
	fmt.Println(NAME, AGE)

	const NAME_1, NAME_2 = "Tom", "Jerry"
	fmt.Println(NAME_1, NAME_2)

	const NAME_3, AGE_3 = "Tom", 20
	fmt.Println(NAME_3, AGE_3)

}
