package folder

import (
	"fmt"
	"strings"
)

func (d *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	// find src
	srcFolders, ok := d.nameToFolders[name]
	if !ok || len(srcFolders) == 0 {
		return nil, fmt.Errorf("source folder '%s' does not exist", name)
	}
	srcFolder := srcFolders[0]

	// find dest
	dstFolders, ok := d.nameToFolders[dst]
	if !ok || len(dstFolders) == 0 {
		return nil, fmt.Errorf("destination folder '%s' does not exist", dst)
	}
	dstFolder := dstFolders[0]

	if srcFolder.OrgId != dstFolder.OrgId {
		return nil, fmt.Errorf("cannot move folder between different organizations")
	}

	// edge cases for test criteria
	if srcFolder.Paths == dstFolder.Paths {
		return nil, fmt.Errorf("cannot move a folder to itself")
	}
	if strings.HasPrefix(dstFolder.Paths+".", srcFolder.Paths+".") {
		return nil, fmt.Errorf("cannot move a folder to its own descendant")
	}

	srcNode := d.folderNodes[srcFolder.Paths]
	dstNode := d.folderNodes[dstFolder.Paths]

	oldPathPrefix := srcFolder.Paths
	newPathPrefix := dstFolder.Paths + "." + srcFolder.Name

	var updatePaths func(*FolderNode)
	updatePaths = func(n *FolderNode) {
		oldPath := n.Folder.Paths
		n.Folder.Paths = strings.Replace(n.Folder.Paths, oldPathPrefix, newPathPrefix, 1)

		delete(d.folderNodes, oldPath)
		d.folderNodes[n.Folder.Paths] = n

		for _, child := range n.Children {
			updatePaths(child)
		}
	}
	updatePaths(srcNode)

	// update the tree structure + remove srcNode from old parent's children
	pathParts := strings.Split(oldPathPrefix, ".")
	if len(pathParts) > 1 {
		oldParentPath := strings.Join(pathParts[:len(pathParts)-1], ".")
		oldParentNode := d.folderNodes[oldParentPath]
		var newChildren []*FolderNode
		for _, child := range oldParentNode.Children {
			if child != srcNode {
				newChildren = append(newChildren, child)
			}
		}
		oldParentNode.Children = newChildren
	} else {
		// root node
		var newRoots []*FolderNode
		for _, root := range d.rootNodes[srcFolder.OrgId] {
			if root != srcNode {
				newRoots = append(newRoots, root)
			}
		}
		d.rootNodes[srcFolder.OrgId] = newRoots
	}

	// add srcNode to new parent
	dstNode.Children = append(dstNode.Children, srcNode)

	return d.folders, nil
}
