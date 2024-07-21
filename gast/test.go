package gast

import "fmt"

type User interface {
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

// GetAge 11
// GetAge 12
func (impl *userImpl) GetAge() int64 {
	return impl.Age
}

func (impl *userImpl) Print() {
	fmt.Println("Name: ", impl.Name, " Age: ", impl.Age)
}

func GetName() {}
