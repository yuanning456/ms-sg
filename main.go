package main

import (
	"fmt"
	"ms-sg/config"
)

func main() {
	host := config.File.MustValue("login_server", "host", "123")
	fmt.Println(host)
	config.A()

	g := &A{}

	g.a()
	g.b()
}


type A struct {}

func (a *A) a() {

}

func (a A) b() {

}