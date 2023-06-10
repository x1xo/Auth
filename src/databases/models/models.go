package models

import (
	"time"
)

type UserInfo struct {
	Id        string      `json:"id" bson:"id"`
	Username  string      `json:"username" bson:"username"`
	Email     string      `json:"email" bson:"email"`
	AvatarURL string      `json:"avatar_url" bson:"avatar_url"`
	Github    GithubUser  `json:"github" bson:"github"`
	Discord   DiscordUser `json:"discord" bson:"discord"`
	Google    GoogleUser  `json:"google" bson:"google"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" bson:"updated_at"`
}

type GithubUser struct {
	Username  string    `json:"login"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GithubUserEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type DiscordUser struct {
	Id            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	AvatarURL     string    `json:"avatar_url"`
	Verified      bool      `json:"verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GoogleUser struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"name"`
	AvatarURL string `json:"picture"`
}
