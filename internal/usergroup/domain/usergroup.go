package domain

import "errors"

type UserGroup struct {
	GroupId   int    `json:"id"`
	GroupName string `json:"name"`
	GroupDesc string `json:"description"`
	Type      string `json:"type"`
	Thumbnail string `json:"thumbnailurl"`
	IsActive  bool   `json:"is_active"`
}

func NewUserGroup(name, desc, groupType string) (*UserGroup, error) {
	if name == "" {
		return nil, errors.New("GroupName field cannot be empty")
	}
	if groupType == "" {
		return nil, errors.New("GroupType field cannot be empty")
	}
	return &UserGroup{
		GroupName: name,
		GroupDesc: desc,
		Type:      groupType,
		IsActive:  true,
		Thumbnail: "",
	}, nil
}
