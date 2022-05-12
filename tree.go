package tree

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Node struct {
	name    string
	child   *Node
	sibling *Node
}

func PrintTree(path string) error {
	node, err := Tree(path)
	if err != nil {
		return err
	}

	printTree(node, 0)

	return nil
}

func printTree(node *Node, depth int) {
	fmt.Printf("%s%s\n", strings.Repeat(" ", depth), node.name)

	if node.child != nil {
		printTree(node.child, depth+1)
	}

	if node.sibling != nil {
		printTree(node.sibling, depth)
	}
}

func Tree(path string) (*Node, error) {
	var root Node
	if err := walk(&root, path); err != nil {
		return nil, err
	}

	return &root, nil
}

func walk(parent *Node, path string) error {
	cmd := exec.Command("ls", path)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run 'ls' cmd: %w", err)
	}

	for {
		line, err := out.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read output of 'ls' cmd: %w", err)
		}

		filename := line[:len(line)-1] // remove "\n"
		childPath := filepath.Join(path, filename)

		child := Node{name: filename}
		if parent.child == nil {
			parent.child = &child
		} else {
			sibling := parent.child
			for sibling.sibling != nil {
				sibling = sibling.sibling
			}
			sibling.sibling = &child
		}

		fileInfo, err := os.Lstat(childPath)
		if err != nil {
			return fmt.Errorf("failed to run Lstat: %w", err)
		}

		if !fileInfo.IsDir() {
			continue
		}

		if err := walk(&child, childPath); err != nil {
			return err
		}
	}

	return nil
}
