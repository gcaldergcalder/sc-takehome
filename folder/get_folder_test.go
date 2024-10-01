package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	// "github.com/stretchr/testify/assert"
)

// feel free to change how the unit test is structured
func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()

	orgID1 := uuid.Must(uuid.NewV4())
	orgID2 := uuid.Must(uuid.NewV4())
	orgID3 := uuid.Must(uuid.NewV4()) // For testing an orgID with no folders

	// Prepare sample folders
	sampleFolders := []folder.Folder{
		{Name: "alpha", Paths: "alpha", OrgId: orgID1},
		{Name: "bravo", Paths: "alpha.bravo", OrgId: orgID1},
		{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
		{Name: "delta", Paths: "delta", OrgId: orgID2},
		{Name: "echo", Paths: "delta.echo", OrgId: orgID2},
		{Name: "foxtrot", Paths: "foxtrot", OrgId: orgID2},
	}

	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		want    []folder.Folder
	}{
		{
			name:    "Get folders for orgID1",
			orgID:   orgID1,
			folders: sampleFolders,
			want: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgID1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgID1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
			},
		},
		{
			name:    "Get folders for orgID2",
			orgID:   orgID2,
			folders: sampleFolders,
			want: []folder.Folder{
				{Name: "delta", Paths: "delta", OrgId: orgID2},
				{Name: "echo", Paths: "delta.echo", OrgId: orgID2},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: orgID2},
			},
		},
		{
			name:    "No folders for orgID3",
			orgID:   orgID3,
			folders: sampleFolders,
			want:    []folder.Folder{},
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)

			got := f.GetFoldersByOrgID(tt.orgID)

			if len(got) != len(tt.want) {
				t.Errorf("Expected %d folders, got %d", len(tt.want), len(got))
			}

			// Optionally, compare the contents of the folders
			for _, expectedFolder := range tt.want {
				found := false
				for _, actualFolder := range got {
					if actualFolder.Name == expectedFolder.Name &&
						actualFolder.Paths == expectedFolder.Paths &&
						actualFolder.OrgId == expectedFolder.OrgId {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected folder %+v not found in result", expectedFolder)
				}
			}
		})
	}
}
