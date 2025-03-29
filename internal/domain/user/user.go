package user

import "strings"

type IUser interface {
	Nickname() string
	SetNickname(nickname string)
}

type User struct {
	nickname string
}

func (u *User) Nickname() string {
	return u.nickname
}

func (u *User) SetNickname(nickname string) {
	nickname = strings.ReplaceAll(nickname, "\n", "")
	nickname = strings.ReplaceAll(nickname, "\r", "")
	nickname = strings.TrimSpace(nickname)
	u.nickname = nickname
}

func NewUser(nickname string) *User {
	return &User{nickname: nickname}
}
