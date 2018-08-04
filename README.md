# :syringe: syringe
A simple but functional dependency injector for Golang

## Installation
```
go get github.com/thomazfeitoza/syringe
```

## Usage

Basic container creation:
```
package main

import (
	"fmt"

	"github.com/thomazfeitoza/syringe"
)

type Printer struct{}

func (p *Printer) Print(message string) {
	fmt.Println(message)
}

type Interactor struct {
	pr *Printer `di:"printer"`
}

func (i *Interactor) DoSomething() {
	i.pr.Print("Hello World")
}

func main() {
	container := syringe.NewContainer()

	container.RegisterType("printer", new(Printer), syringe.Transient)
	container.RegisterType("interactor", new(Interactor), syringe.Transient)

	interactor := container.Resolve("interactor").(*Interactor)
	interactor.DoSomething()
}
```

Resolving dependencies from scopes:
```
package main

import (
	"fmt"

	"github.com/thomazfeitoza/syringe"
)

type Printer struct {
	requestId string `di:"reqId"`
}

func (p *Printer) Print(message string) {
	fmt.Printf("%s: %s\n", p.requestId, message)
}

type Interactor struct {
	pr *Printer `di:"printer"`
}

func (i *Interactor) DoSomething() {
	i.pr.Print("Hello World")
}

func main() {
	container := syringe.NewContainer()
	container.RegisterType("printer", new(Printer), syringe.Transient)
	container.RegisterType("interactor", new(Interactor), syringe.Transient)

	scope := container.Scope()
	scope.RegisterValue("reqId", "request-123")

	interactor := scope.Resolve("interactor").(*Interactor)
	interactor.DoSomething()
}
```