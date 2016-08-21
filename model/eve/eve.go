// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package oautheve

import (
	"encoding/json"
	"io"

	"github.com/mattermost/platform/einterfaces"
	"github.com/mattermost/platform/model"
)

type EveProvider struct {
}

type EveUser struct {
	CharacterId        int64  `json:"CharacterID"`
	CharacterName      string `json:"CharacterName"`
	CharacterOwnerHash string `json:"CharacterOwnerHash"`
}

func init() {
	provider := &EveProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_EVE, provider)
}

func userFromGitLabUser(glu *EveUser) *model.User {
	user := &model.User{}
	user.Username = model.CleanUsername(glu.CharacterName)
	user.FirstName = glu.CharacterName
	user.AuthData = &glu.CharacterOwnerHash
	user.AuthService = model.USER_AUTH_SERVICE_EVE
	user.Email = glu.CharacterOwnerHash + "@eveonline.com"

	return user
}

func gitLabUserFromJson(data io.Reader) *EveUser {
	decoder := json.NewDecoder(data)
	var glu EveUser
	err := decoder.Decode(&glu)
	if err == nil {
		return &glu
	}

	return nil
}

func (glu *EveUser) IsValid() bool {
	if glu.CharacterId == 0 {
		return false
	}

	return true
}

func (glu *EveUser) getAuthData() string {
	return glu.CharacterOwnerHash
}

func (m *EveProvider) GetIdentifier() string {
	return model.USER_AUTH_SERVICE_EVE
}

func (m *EveProvider) GetUserFromJson(data io.Reader) *model.User {
	glu := gitLabUserFromJson(data)
	if glu.IsValid() {
		return userFromGitLabUser(glu)
	}

	return &model.User{}
}

func (m *EveProvider) GetAuthDataFromJson(data io.Reader) string {
	glu := gitLabUserFromJson(data)

	if glu.IsValid() {
		return glu.getAuthData()
	}

	return ""
}
