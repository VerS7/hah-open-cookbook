package storage

import (
	"time"

	"open-hah-cookbook/internal/auth"
)

type User struct {
	Id             int       `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	Name           string    `db:"name"`
	HashedPassword string    `db:"hashed_password"`
	AccessToken    string    `db:"access_token"`
	IsAdmin        bool      `db:"is_admin"`
}

type Session struct {
	AccessToken string `db:"access_token"`
	IsAdmin     bool   `db:"is_admin"`
}

func (st *Storage) AddUser(name string, password string, isAdmin bool) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	query := "INSERT INTO users (name, hashed_password, access_token, is_admin) VALUES (?, ?, ?, ?)"

	hashedPassword := auth.Hash(password)
	accessToken := auth.Hash(name + hashedPassword)

	_, err := st.db.Exec(query, name, hashedPassword, accessToken, isAdmin)
	return err
}

func (st *Storage) GetUserByName(name string) (*User, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var u User

	if err := st.db.Get(&u, "SELECT * FROM users WHERE name == ?", name); err != nil {
		return nil, err
	}

	return &u, nil
}

func (st *Storage) RemoveUser(id int) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	_, err := st.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err

}

func (st *Storage) GetSessions() ([]Session, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var sessions []Session

	err := st.db.Select(&sessions, "SELECT access_token, is_admin FROM users")
	return sessions, err
}
