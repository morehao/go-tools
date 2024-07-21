package gast

import "fmt"

type User interface {
	GetName() string
	GetAge(n int, n1 int) (int64, int)
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
func (impl *userImpl) GetAge(n int, n1 int) (int64, int) {
	return impl.Age, n
}

func (impl *userImpl) Print() {
	fmt.Println("Name: ", impl.Name, " Age: ", impl.Age)
}

func GetName() {}
