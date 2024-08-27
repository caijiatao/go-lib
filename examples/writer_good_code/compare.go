package main

import "fmt"

type User struct {
	Name string
}

func main() {
	u1 := &User{Name: "user1"}
	u2 := &User{Name: "user1"}
	fmt.Println(*u1 == *u2)
}
