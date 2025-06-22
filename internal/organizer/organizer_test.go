package organizer

import (
	"testing"
)

type dummyService struct{}

func TestOrganizeFiles_InvalidService(t *testing.T) {
	meta := Metadata{GroupName: "G"}
	err := OrganizeFiles(nil, "fileid", meta)
	if err == nil {
		t.Error("expected error with nil service")
	}
}
