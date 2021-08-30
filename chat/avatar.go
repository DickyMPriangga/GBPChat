package main

import (
	"errors"
	"io/ioutil"
	"path"

	gomniauthcommon "github.com/stretchr/gomniauth/common"
)

var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type Avatar interface {
	GetAvatarURL(u ChatUser) (string, error)
}

type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

type AvatarAuth struct {
}

type GravatarAvatar struct {
}

type FileSystemAvatar struct {
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

var UseAuthAvatar AvatarAuth
var UseGravatar GravatarAvatar
var UseFileSystemAvatar FileSystemAvatar

func (AvatarAuth) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()

	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}

	return url, nil
}

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.AvatarURL(), nil
}

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	files, err := ioutil.ReadDir("avatars_img")

	if err != nil {
		return "", ErrNoAvatarURL
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
			return "/avatars_img/" + file.Name(), nil
		}
	}

	return "", ErrNoAvatarURL
}
