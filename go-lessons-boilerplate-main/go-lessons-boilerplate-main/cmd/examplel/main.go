package main

import (
	"errors"
	"fmt"
)

func counter(lim int) func() bool {
	i := 0
	return func() bool {
		i++
		fmt.Println(i)
		return i < lim
	}
}

type Status int

const (
	Pending  Status = 1
	Approved Status = 2
	Rejected Status = 3
)

func main2() {
	var v []int

	v2 := make([]int, 0, 7)

	fmt.Println(v == nil)
	fmt.Println(v2 == nil)

	v2 = append(v2, 7)
	v2 = append(v2, 8)
	v2 = append(v2, 9)
	fmt.Println(v2)

	a := map[int]string{
		2: "untoeheant",
	}

	a[1] = "nutohan"

	a[135] = "altceva"
	v7, ok := a[135]
	if !ok {
		fmt.Println("Not found")
	} else {
		fmt.Println("Found", v7)
	}
	fmt.Println(a)

	s := map[int]bool{}

	s[7] = true
	fmt.Println(s[7])

	c := counter(7)
	for c() {
	}

	for i := range 7 {
		fmt.Println(i)
	}

	for k, v := range a {
		fmt.Println(k, "->", v)
	}

	for _, e := range v2 {
		// fmt.Println("at", idx, "->", e)
		fmt.Println(e)
	}
}

func f() {
	fmt.Println("unethu")
}

type X string
func (x X) Error() string {
	return string(x)
}

func g() (int, error) {
	return 7, X("am o eroare custom")
	return 7, errors.New("Error")
}


func main() {
	f()

	// g := f
	// g()
	//
	// g = func() {
	// 	fmt.Println("123")
	// }
	// g()

	v, err := g()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
}

