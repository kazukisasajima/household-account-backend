package entity

type User struct {
	ID       int    `json:"id"`
	Email	 string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Credentials struct {
	Email    string
	Password string
}