package types

import "election/utils"

type LASFSMember struct {
	ID       string
	Name     string
	Password string
}

func (m *LASFSMember) GetID() string {
	return m.ID
}

func (m *LASFSMember) GetName() string {
	return m.Name
}

func (m *LASFSMember) GenerateID(string) string {
	return utils.LASFSMemberIDFromString(m.Name)
}

func NewLASFSMember(username string, password string) *LASFSMember {
	m := LASFSMember{
		Name:     username,
		Password: password,
	}
	m.ID = m.GenerateID(m.Name)
	return &m
}
