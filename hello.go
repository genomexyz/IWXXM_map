package main

import "fmt"
import "strconv"

func main() {
	var a [4]int
	fmt.Printf("hello, world\n")
	coba := strconv.Itoa(-42)
	fmt.Printf(coba+"ssd\n")
	
	sum := 0
	for i := 0; i < 4; i++ {
		sum += i
		a[i] = sum
		fmt.Printf(strconv.Itoa(sum)+"\n")
		fmt.Println(a)
	}
}
