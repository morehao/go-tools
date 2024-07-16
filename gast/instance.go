package gast

import "fmt"

type User interface {
	GetAge() int64
}

type userImpl struct {
	Name string
	Age  int64
}

func NewUser() User {
	return &userImpl{}
}

func (impl *userImpl) GetName() string {
	return impl.Name
}

func (impl *userImpl) GetAge() int64 {
	return impl.Age
}

func (impl *userImpl) Print() {
	fmt.Println("Name: ", impl.Name, " Age: ", impl.Age)
}

func GetName() string {
	return "test"
}
