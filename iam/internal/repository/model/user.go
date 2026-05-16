package model

type User struct {
	UserUUID            string `db:"user_id"`
	Login               string `db:"login"`
	Email               string `db:"email"`
	PasswordHash        string `db:"password_hash"`
	NotificationMethods []byte `db:"notification_methods"`
}

type LoginCredentials struct {
	Id string `db:"user_id"`
	Pw string `db:"password_hash"`
}

type UserRedisView struct {
	UserUUID            string `redis:"user_id"`
	Login               string `redis:"login"`
	Email               string `redis:"email"`
	NotificationMethods []byte `redis:"notification_methods"`
}
