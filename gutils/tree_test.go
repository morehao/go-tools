package gutils

import (
	"testing"
)

func TestBuildTree(t *testing.T) {
	type Data struct {
		Description string
	}
	nodes := []TreeNode{
		{ID: 1, ParentID: 0, TreeNodeData: Data{Description: "Root 1"}},
		{ID: 2, ParentID: 1, TreeNodeData: Data{Description: "Child 1.1"}},
		{ID: 3, ParentID: 1, TreeNodeData: Data{Description: "Child 1.2"}},
		{ID: 4, ParentID: 2, TreeNodeData: Data{Description: "Child 1.1.1"}},
		{ID: 5, ParentID: 2, TreeNodeData: Data{Description: "Child 1.1.2"}},
		{ID: 6, ParentID: 3, TreeNodeData: Data{Description: "Child 1.2.1"}},
		{ID: 7, ParentID: 0, TreeNodeData: Data{Description: "Root 2"}},
		{ID: 8, ParentID: 7, TreeNodeData: Data{Description: "Child 2.1"}},
		{ID: 9, ParentID: 7, TreeNodeData: Data{Description: "Child 2.2"}},
		{ID: 10, ParentID: 9, TreeNodeData: Data{Description: "Child 2.2.1"}},
		{ID: 11, ParentID: 6, TreeNodeData: Data{Description: "Child 1.2.1.1"}},
	}
	roots := BuildTree(nodes)
	printTree(roots, 0)
}
