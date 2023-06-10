package models

import (
	"time"
)

type UserInfo struct {
	Id        string      `json:"id" bson:"id"`
	Username  string      `json:"username" bson:"username"`
	Email     string      `json:"email" bson:"email"`
	AvatarURL string      `json:"avatar_url" bson:"avatar_url"`
	Github    GithubUser  `json:"github,omitempty" bson:"github"`
	Discord   DiscordUser `json:"discord,omitempty" bson:"discord"`
	Google    GoogleUser  `json:"google,omitempty" bson:"google"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" bson:"updated_at"`
}

type GithubUser struct {
	Username  string    `json:"login,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type GithubUserEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type DiscordUser struct {
	Id            string    `json:"id,omitempty"`
	Email         string    `json:"email,omitempty"`
	Username      string    `json:"username,omitempty"`
	Discriminator string    `json:"discriminator,omitempty"`
	AvatarURL     string    `json:"avatar_url,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	Verified      bool      `json:"verified,omitempty"`
}

type GoogleUser struct {
	Id        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	Username  string `json:"name,omitempty"`
	AvatarURL string `json:"picture,omitempty"`
}
