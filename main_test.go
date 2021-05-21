package main

import (
	"fmt"
	"healthBridge/assert"
	"testing"
)

type ibase interface {
	DoSomething() string
}

type myBase struct {
	first  string
	second string
}

type firstSub struct {
	myBase
}

type secondSub struct {
	myBase
}

func (o firstSub) DoSomething() string {
	return fmt.Sprintf("%s-%s", o.first, o.second)
}

func (o secondSub) DoSomething() string {
	return fmt.Sprintf("%s-%s", o.second, o.first)
}

func Factory(b bool) ibase {
	if b {
		return firstSub{
			myBase{
				first:  "one",
				second: "two",
			},
		}
	} else {
		return secondSub{
			myBase{
				first:  "one",
				second: "two",
			},
		}
	}
}

func TestInheritance(t *testing.T) {

	i1 := Factory(true)
	fmt.Printf("Result i1: %s\r\n", i1.DoSomething())

	i2 := Factory(false)
	fmt.Printf("Result i1: %s\r\n", i2.DoSomething())

	list := make(map[string]*ibase)
	list["a"] = &i1
	list["b"] = &i2

	for i, v := range list {
		fmt.Printf("Result %v: %s\r\n", i, (*v).DoSomething())
	}
	assert.StrEqual("one-two", i1.DoSomething(), "", t)
	assert.StrEqual("two-one", i2.DoSomething(), "", t)
}
