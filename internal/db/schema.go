package db

type UsersTable struct {
	Id int `gorm:"primaryKey,autoIncrement"`
	Email string
	PasswordHash string
}

type LanguagesTable struct {
	Id int `gorm:"primaryKey,autoIncrement"`
	UserId int
	Name string
	XP int
	LinesDel int
	LinesNew int
	Files int
}