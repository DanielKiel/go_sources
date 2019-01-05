package main

import (
  "fmt"
  "time"
)

func main() {
  fmt.Println("Welcome to the playground!")

	fmt.Println("The time is", time.Now())

  fmt.Println("The result of 4 + 2 is", add(2, 4))

  fmt.Println(split(10))
}

func add(x, y int) int {
  return x + y
}

//ein nacktes return, interessantes pattern
func split(sum int) (x, y int) {
  x = sum * 4 / 10
  y = sum + 20

  return
}
