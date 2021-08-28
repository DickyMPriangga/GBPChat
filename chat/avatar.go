package main

import (
	"errors"
)

var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AvatarAuth struct {
}

type GravatarAvatar struct {
}

var UseAuthAvatar AvatarAuth
var UseGravatar GravatarAvatar

func (AvatarAuth) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]

	if !ok {
		return "", ErrNoAvatarURL
	}

	urlStr, ok := url.(string)

	if !ok {
		return "", ErrNoAvatarURL
	}

	return urlStr, nil
}

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	userid, ok := c.userData["userid"]

	if !ok {
		return "", ErrNoAvatarURL
	}

	useridStr, ok := userid.(string)

	if !ok {
		return "", ErrNoAvatarURL
	}

	return "//www.gravatar.com/avatar/" + useridStr, nil
}
