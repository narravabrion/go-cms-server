package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
}

type Password struct {
	Pass *string
	Hash []byte
}

func (p *Password) Set(pass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Pass = &pass
	p.Hash = hash

	return nil

}


func (p *Password) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(password))
}
