package folder

import (
	"fmt"
	"strings"
)

func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	var srcFolder, dstFolder Folder
	srcFound, dstFound := false, false

	for _, folder := range f.folders {
		if folder.Name == name {
			srcFolder = folder
			srcFound = true
		}
		if folder.Name == dst {
			dstFolder = folder
			dstFound = true
		}
	}

	if !srcFound {
		return nil, fmt.Errorf("source folder '%s' does not exist", name)
	}
	if !dstFound {
		return nil, fmt.Errorf("destination folder '%s' does not exist", dst)
	}

	if srcFolder.OrgId != dstFolder.OrgId {
		return nil, fmt.Errorf("cannot move folder between different organizations")
	}

	// edge case: moving folder to itself
	if srcFolder.Name == dstFolder.Name {
		return nil, fmt.Errorf("cannot move a folder to itself")
	}

	// edge case: moving folder into its child
	if strings.HasPrefix(dstFolder.Paths+".", srcFolder.Paths+".") {
		return nil, fmt.Errorf("cannot move a folder to its own descendant")
	}

	oldPathPrefix := srcFolder.Paths
	newPathPrefix := dstFolder.Paths + "." + srcFolder.Name

	for i, folder := range f.folders {
		if folder.OrgId == srcFolder.OrgId {
			if strings.HasPrefix(folder.Paths, oldPathPrefix) {

				relativePath := strings.TrimPrefix(folder.Paths, oldPathPrefix)
				if len(relativePath) > 0 && relativePath[0] == '.' {
					relativePath = relativePath[1:]
				}
				newPath := newPathPrefix
				if relativePath != "" {
					newPath += "." + relativePath
				}
				f.folders[i].Paths = newPath
			}
		}
	}

	return f.folders, nil
}
