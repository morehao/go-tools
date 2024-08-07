package gutils

import "fmt"

// TreeNode 结构体表示树的一个节点
type TreeNode struct {
	ID       int // 节点的唯一标识符
	ParentID int // 节点的父节点的唯一标识符
	TreeNodeData
	Children []*TreeNode // 节点的子节点
}

type TreeNodeData any

// BuildTree 从节点列表构建树
func BuildTree(nodes []TreeNode) []*TreeNode {
	nodeMap := make(map[int]*TreeNode) // 创建一个映射，便于查找节点
	var roots []*TreeNode              // 存储所有的根节点

	// 创建节点映射，方便查找
	for i := range nodes {
		node := &nodes[i]
		node.Children = []*TreeNode{} // 初始化子节点列表
		nodeMap[node.ID] = node       // 将节点添加到映射中
	}

	// 递归构建树
	for i := range nodes {
		node := &nodes[i]
		if node.ParentID == 0 {
			// 该节点是根节点
			roots = append(roots, node)
		} else {
			// 该节点有父节点，将其添加到父节点的子节点列表中
			if parent, ok := nodeMap[node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			} else {
				// 处理父节点不存在的情况
				fmt.Printf("警告: 节点 ID %d 的父节点 ID %d 未找到\n", node.ID, node.ParentID)
			}
		}
	}

	return roots
}

// printTree 递归打印树的结构
func printTree(nodes []*TreeNode, level int) {
	for _, node := range nodes {
		fmt.Printf("%s%s\n", getIndent(level), node.TreeNodeData)
		printTree(node.Children, level+1)
	}
}

// getIndent 返回一个用于缩进的字符串
func getIndent(level int) string {
	return fmt.Sprintf("%s", fmt.Sprintf("%d", level*2))
}
