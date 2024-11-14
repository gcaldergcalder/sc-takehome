package folder

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (d *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	var res []Folder

	// makes use of thenew rootNodes to traverse all folders under the given OrgID
	roots, ok := d.rootNodes[orgID]
	if !ok {
		return res
	}

	var dfs func(*FolderNode)
	dfs = func(n *FolderNode) {
		res = append(res, *n.Folder)
		for _, child := range n.Children {
			dfs(child)
		}
	}

	for _, root := range roots {
		dfs(root)
	}

	return res
}

func (d *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	// find the folder with the given name and OrgID
	foldersWithName, ok := d.nameToFolders[name]
	if !ok {
		return nil, fmt.Errorf("folder '%s' does not exist", name)
	}

	var folder *Folder
	for _, f := range foldersWithName {
		if f.OrgId == orgID {
			folder = f
			break
		}
	}

	if folder == nil {
		return nil, fmt.Errorf("folder '%s' does not exist in the specified organization", name)
	}

	// get the FolderNode
	node, ok := d.folderNodes[folder.Paths]
	if !ok {
		return nil, fmt.Errorf("folder node not found")
	}

	// collect all child folders using DFS
	var childFolders []Folder
	var dfs func(*FolderNode)
	dfs = func(n *FolderNode) {
		for _, child := range n.Children {
			childFolders = append(childFolders, *child.Folder)
			dfs(child)
		}
	}
	dfs(node)

	return childFolders, nil
}
