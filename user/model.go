package user

import "github.com/hhru/meetup/jira"

type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`
}

func fromJira(user *jira.User) *User {
	return &User{
		user.Name,
		user.DisplayName,
		user.AvatarUrls.Big,
	}
}
