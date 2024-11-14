package folder

import (
	"strings"

	"github.com/gofrs/uuid"
)

type IDriver interface {
	// GetFoldersByOrgID returns all folders that belong to a specific orgID.
	GetFoldersByOrgID(orgID uuid.UUID) []Folder
	// component 1
	// Implement the following methods:
	// GetAllChildFolders returns all child folders of a specific folder.
	GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error)

	// component 2
	// Implement the following methods:
	// MoveFolder moves a folder to a new destination.
	MoveFolder(name string, dst string) ([]Folder, error)
}

type driver struct {
	// define attributes here
	// data structure to store folders
	// or preprocessed data

	// example: feel free to change the data structure, if slice is not what you want
	folders []Folder

	// updated solution attributes
	folderMap     map[uuid.UUID]map[string]*Folder
	folderNodes   map[string]*FolderNode
	rootNodes     map[uuid.UUID][]*FolderNode
	nameToFolders map[string][]*Folder
}

type FolderNode struct {
	Folder   *Folder
	Children []*FolderNode
}

func NewDriver(folders []Folder) IDriver {
	d := &driver{
		folders:       folders,
		folderMap:     make(map[uuid.UUID]map[string]*Folder),
		folderNodes:   make(map[string]*FolderNode),
		rootNodes:     make(map[uuid.UUID][]*FolderNode),
		nameToFolders: make(map[string][]*Folder),
	}

	for i := range folders {
		f := &folders[i]
		if _, ok := d.folderMap[f.OrgId]; !ok {
			d.folderMap[f.OrgId] = make(map[string]*Folder)
		}
		d.folderMap[f.OrgId][f.Name] = f
		d.nameToFolders[f.Name] = append(d.nameToFolders[f.Name], f)
	}

	for i := range folders {
		f := &folders[i]
		node := &FolderNode{
			Folder: f,
		}
		d.folderNodes[f.Paths] = node
	}

	for _, node := range d.folderNodes {
		pathParts := strings.Split(node.Folder.Paths, ".")
		if len(pathParts) == 1 {
			// root folder
			d.rootNodes[node.Folder.OrgId] = append(d.rootNodes[node.Folder.OrgId], node)
		} else {
			// has a parent
			parentPath := strings.Join(pathParts[:len(pathParts)-1], ".")
			parentNode, ok := d.folderNodes[parentPath]
			if ok {
				parentNode.Children = append(parentNode.Children, node)
			}
		}
	}

	return d
}
