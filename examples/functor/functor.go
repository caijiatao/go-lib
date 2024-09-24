package main

import (
	"fmt"
	"strings"
)

type User struct {
	ID   int
	Name string
}

type Functor[T any] []T

func (f Functor[T]) Map(fn func(T) T) Functor[T] {
	var result Functor[T]
	for _, v := range f {
		result = append(result, fn(v))
	}
	return result
}

func main() {
	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}

	upperUserName := func(u User) User {
		u.Name = strings.ToUpper(u.Name)
		return u
	}

	us := Functor[User](users)
	us = us.Map(upperUserName)

	// 打印转换后的用户数据
	for _, u := range us {
		fmt.Printf("ID: %d, Name: %s\n", u.ID, u.Name)
	}
}
