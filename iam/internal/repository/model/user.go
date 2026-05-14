package model

type User struct {
	UserUUID            string `db:"user_id"`
	Login               string `db:"login"`
	Email               string `db:"email"`
	PasswordHash        string `db:"password_hash"`
	NotificationMethods []byte `db:"notification_methods"`
}
