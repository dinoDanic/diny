package users

type User struct {
	ID   int
	Name string
}

func FindUser(id int) (*User, bool) {
	for _, u := range seedUsers {
		if u.ID == id {
			return &u, true
		}
	}
	return nil, false
}

var seedUsers = []User{
	{ID: 1, Name: "alice"},
	{ID: 2, Name: "bob"},
}
