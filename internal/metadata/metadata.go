package metadata

type Metadata struct {
	GroupName   string
	Teacher     string
	SessionType string // in-person or virtual
	SongsTaught []string
	Ragas       []string
	Talas       []string
	Composers   []string
}

func NewMetadata(groupName, teacher, sessionType string, songsTaught, ragas, talas, composers []string) *Metadata {
	return &Metadata{
		GroupName:   groupName,
		Teacher:     teacher,
		SessionType: sessionType,
		SongsTaught: songsTaught,
		Ragas:       ragas,
		Talas:       talas,
		Composers:   composers,
	}
}

func (m *Metadata) GetGroupName() string {
	return m.GroupName
}

func (m *Metadata) GetTeacher() string {
	return m.Teacher
}

func (m *Metadata) GetSessionType() string {
	return m.SessionType
}

func (m *Metadata) GetSongsTaught() []string {
	return m.SongsTaught
}

func (m *Metadata) GetRagas() []string {
	return m.Ragas
}

func (m *Metadata) GetTalas() []string {
	return m.Talas
}

func (m *Metadata) GetComposers() []string {
	return m.Composers
}