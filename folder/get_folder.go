package folder

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	folders := f.folders

	res := []Folder{}
	for _, f := range folders {
		if f.OrgId == orgID {
			res = append(res, f)
		}
	}

	return res

}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	var parentPath string
	found := false

	//TODO:  bit of a clumsy way of accessing the parent folder, look for more efficient search
	// iterates over folders attribute in driver struct
	for _, folder := range f.folders {
		// breaks loop when finding the folder
		if folder.OrgId == orgID && folder.Name == name {
			parentPath = folder.Paths
			found = true
			break
		}
	}

	if !found {
		// notifies file not found
		return nil, fmt.Errorf("Folder '%s' does not exist in the specified organization", name)
	}

	var childFolders []Folder
	parentPathWithDot := parentPath + "."

	for _, folder := range f.folders {
		if folder.OrgId == orgID {
			if strings.HasPrefix(folder.Paths, parentPathWithDot) {
				childFolders = append(childFolders, folder)
			}
		}
	}

	return childFolders, nil
}
