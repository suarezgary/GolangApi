package models

//User - User Model
type User struct {
	ID       uint64 `json:"-"`
	FullName string `json:"-"`
}

//FindByID - Mocked Find by Id
func (user *User) FindByID() error {
	user.ID = 1
	user.FullName = "Demo User"
	return nil
}
