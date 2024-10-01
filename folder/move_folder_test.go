package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

func Test_folder_MoveFolder(t *testing.T) {
	orgID1 := uuid.Must(uuid.NewV4())
	orgID2 := uuid.Must(uuid.NewV4())

	sampleFolders := []folder.Folder{
		{Name: "alpha", Paths: "alpha", OrgId: orgID1},
		{Name: "bravo", Paths: "alpha.bravo", OrgId: orgID1},
		{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
		{Name: "delta", Paths: "alpha.delta", OrgId: orgID1},
		{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgID1},
		{Name: "foxtrot", Paths: "foxtrot", OrgId: orgID2},
		{Name: "golf", Paths: "golf", OrgId: orgID1},
	}

	tests := []struct {
		name           string
		srcFolder      string
		dstFolder      string
		expectedOutput []folder.Folder
		expectError    bool
		errorMessage   string
	}{
		{
			name:      "Move bravo under delta",
			srcFolder: "bravo",
			dstFolder: "delta",
			expectedOutput: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgID1},
				{Name: "bravo", Paths: "alpha.delta.bravo", OrgId: orgID1},
				{Name: "charlie", Paths: "alpha.delta.bravo.charlie", OrgId: orgID1},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgID1},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgID1},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: orgID2},
				{Name: "golf", Paths: "golf", OrgId: orgID1},
			},
			expectError: false,
		},
		{
			name:      "Move bravo under golf",
			srcFolder: "bravo",
			dstFolder: "golf",
			expectedOutput: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgID1},
				{Name: "bravo", Paths: "golf.bravo", OrgId: orgID1},
				{Name: "charlie", Paths: "golf.bravo.charlie", OrgId: orgID1},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgID1},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgID1},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: orgID2},
				{Name: "golf", Paths: "golf", OrgId: orgID1},
			},
			expectError: false,
		},
		{
			name:         "Move folder to itself (bravo to bravo)",
			srcFolder:    "bravo",
			dstFolder:    "bravo",
			expectError:  true,
			errorMessage: "cannot move a folder to itself",
		},
		{
			name:         "Move folder to a child of itself (bravo to charlie)",
			srcFolder:    "bravo",
			dstFolder:    "charlie",
			expectError:  true,
			errorMessage: "cannot move a folder to its own descendant",
		},
		{
			name:         "Move folder between different organizations (bravo to foxtrot)",
			srcFolder:    "bravo",
			dstFolder:    "foxtrot",
			expectError:  true,
			errorMessage: "cannot move folder between different organizations",
		},
		{
			name:         "Move non-existent source folder",
			srcFolder:    "nonexistent",
			dstFolder:    "delta",
			expectError:  true,
			errorMessage: "source folder 'nonexistent' does not exist",
		},
		{
			name:         "Move folder to non-existent destination",
			srcFolder:    "bravo",
			dstFolder:    "nonexistent",
			expectError:  true,
			errorMessage: "destination folder 'nonexistent' does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := folder.NewDriver(sampleFolders)

			output, err := driver.MoveFolder(tt.srcFolder, tt.dstFolder)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMessage {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMessage, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(output) != len(tt.expectedOutput) {
				t.Errorf("Expected %d folders, got %d", len(tt.expectedOutput), len(output))
			}

			for _, expectedFolder := range tt.expectedOutput {
				found := false
				for _, actualFolder := range output {
					if actualFolder.Name == expectedFolder.Name && actualFolder.Paths == expectedFolder.Paths && actualFolder.OrgId == expectedFolder.OrgId {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected folder %+v not found in output", expectedFolder)
				}
			}
		})
	}
}
