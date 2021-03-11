package grf

import "github.com/pkg/errors"

type EntryTree struct {
	Root *EntryTreeNode
}

func (t *EntryTree) Traverse(n *EntryTreeNode, f func(*EntryTreeNode)) {
	if n == nil {
		return
	}

	t.Traverse(n.Left, f)
	f(n)
	t.Traverse(n.Right, f)
}

func (t *EntryTree) Insert(value string, data []Entry) error {
	if t.Root == nil {
		t.Root = &EntryTreeNode{Value: value, Data: data}
		return nil
	}

	return t.Root.Insert(value, data)
}

func (t *EntryTree) Find(s string) ([]Entry, bool) {
	if t.Root == nil {
		return nil, false
	}

	return t.Root.Find(s)
}

type EntryTreeNode struct {
	Value string
	Data  []Entry
	Left  *EntryTreeNode
	Right *EntryTreeNode
}

func (n *EntryTreeNode) Insert(value string, data []Entry) error {
	if n == nil {
		return errors.New("could not insert value: nil tree")
	}

	switch {
	case value == n.Value:
		return nil
	case value < n.Value:
		{
			if n.Left == nil {
				n.Left = &EntryTreeNode{
					Value: value,
					Data:  data,
				}

				return nil
			}

			return n.Left.Insert(value, data)
		}
	case value > n.Value:
		{
			if n.Right == nil {
				n.Right = &EntryTreeNode{
					Value: value,
					Data:  data,
				}

				return nil
			}

			return n.Right.Insert(value, data)
		}
	}

	return nil
}

func (n *EntryTreeNode) Find(s string) ([]Entry, bool) {
	if n == nil {
		return nil, false
	}

	switch {
	case s == n.Value:
		return n.Data, true
	case s < n.Value:
		return n.Left.Find(s)
	default:
		return n.Right.Find(s)
	}
}
