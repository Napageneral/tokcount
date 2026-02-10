package output

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Napageneral/tokcount/internal/count"
)

type treeNode struct {
	name     string
	path     string
	tokens   int
	children map[string]*treeNode
}

// RenderTree returns an ASCII full directory token breakdown.
func RenderTree(result *count.Result) string {
	root := &treeNode{
		name:     ".",
		path:     ".",
		tokens:   result.TotalTokens,
		children: make(map[string]*treeNode),
	}

	for relPath, tokens := range result.DirectoryTokens {
		if relPath == "." || tokens <= 0 {
			continue
		}
		insertTreeNode(root, relPath, tokens)
	}

	var b strings.Builder
	b.WriteString("Directory tree:\n")
	b.WriteString(fmt.Sprintf(".  %s tokens (100.0%%)\n", formatInt(result.TotalTokens)))

	children := sortedChildren(root)
	for i, child := range children {
		renderTreeNode(&b, child, "", i == len(children)-1, result.TotalTokens)
	}
	return b.String()
}

func insertTreeNode(root *treeNode, relPath string, tokens int) {
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	current := root
	currentPath := ""

	for _, part := range parts {
		if part == "." || part == "" {
			continue
		}
		if currentPath == "" {
			currentPath = part
		} else {
			currentPath = currentPath + "/" + part
		}

		child := current.children[part]
		if child == nil {
			child = &treeNode{
				name:     part,
				path:     currentPath,
				children: make(map[string]*treeNode),
			}
			current.children[part] = child
		}
		current = child
	}

	current.tokens = tokens
}

func renderTreeNode(b *strings.Builder, node *treeNode, prefix string, isLast bool, totalTokens int) {
	branch := "|- "
	childPrefix := prefix + "|  "
	if isLast {
		branch = "\\- "
		childPrefix = prefix + "   "
	}

	percent := 0.0
	if totalTokens > 0 {
		percent = (float64(node.tokens) / float64(totalTokens)) * 100
	}

	b.WriteString(fmt.Sprintf("%s%s%s/ %s tokens (%.1f%%)\n",
		prefix,
		branch,
		node.name,
		formatInt(node.tokens),
		percent,
	))

	children := sortedChildren(node)
	for i, child := range children {
		renderTreeNode(b, child, childPrefix, i == len(children)-1, totalTokens)
	}
}

func sortedChildren(node *treeNode) []*treeNode {
	children := make([]*treeNode, 0, len(node.children))
	for _, child := range node.children {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		if children[i].tokens == children[j].tokens {
			return children[i].path < children[j].path
		}
		return children[i].tokens > children[j].tokens
	})
	return children
}
