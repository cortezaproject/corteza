package sam

// Users
type User struct {
	ID       uint64
	Username string
	Password string `json:"-"`
}

func (User) new() *User {
	return &User{}
}

func (u *User) GetID() uint64 {
	return u.ID
}

func (u *User) SetID(value uint64) *User {
	u.ID = value
	return u
}
func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) SetUsername(value string) *User {
	u.Username = value
	return u
}
func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) SetPassword(value string) *User {
	u.Password = value
	return u
}
