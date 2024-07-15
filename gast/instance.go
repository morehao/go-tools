package gast

type User interface {
	GetName() string
}

type userImpl struct{}

func NewUser() User {
	return &userImpl{}
}

func (i *userImpl) GetName() string {
	return "user"
}

func (i *userImpl) GetAge() int64 {
	return 10
}
