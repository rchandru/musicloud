package metadata

import "testing"

func TestNewMetadataAndGetters(t *testing.T) {
	m := NewMetadata("Group1", "Teacher1", "virtual", []string{"Song1"}, []string{"Raga1"}, []string{"Tala1"}, []string{"Composer1"})
	if m.GetGroupName() != "Group1" {
		t.Errorf("expected Group1, got %s", m.GetGroupName())
	}
	if m.GetTeacher() != "Teacher1" {
		t.Errorf("expected Teacher1, got %s", m.GetTeacher())
	}
	if m.GetSessionType() != "virtual" {
		t.Errorf("expected virtual, got %s", m.GetSessionType())
	}
	if len(m.GetSongsTaught()) != 1 || m.GetSongsTaught()[0] != "Song1" {
		t.Errorf("unexpected SongsTaught")
	}
}
