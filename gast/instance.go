package gast

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

func (impl *userImpl) GetAge() int64 {
	return impl.Age
}
func GetName() string {
	return "test"
}
