package main

import "fmt"

type Animal interface {
	Speak()
}
func SpeakWithAnimal(a Animal) {
	a.Speak()
}


type Cat struct { }
func (c Cat) Miau() {
	fmt.Println("miau")
}
func (c Cat) Speak() {
	c.Miau()
}

type Dog struct { 
	weight int
}
func (d *Dog) Speak() {
	fmt.Println("ham")
}

type X struct {
	Dog
}

type F func() string
func (f F) Speak() {
	fmt.Println(f())
}

func ceva() string {
	return "cevaaa"
}

func main() {
	SpeakWithAnimal(F(ceva))
	return 

	c := Cat{}
	c.Miau()

	SpeakWithAnimal(c)

	d := Dog{}
	SpeakWithAnimal(&d)

	x := X{}
	x.Speak()
	fmt.Println(x.weight)
	SpeakWithAnimal(&x)
}
