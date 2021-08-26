package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
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
	email, ok := c.userData["email"]

	if !ok {
		return "", ErrNoAvatarURL
	}

	emailStr, ok := email.(string)

	if !ok {
		return "", ErrNoAvatarURL
	}

	m := md5.New()
	io.WriteString(m, strings.ToLower(emailStr))
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
}
