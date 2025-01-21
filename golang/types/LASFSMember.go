package types

type LASFSMember struct {
	Name     string
	Password string
}

func (m *LASFSMember) GetName() string {
	return m.Name
}
