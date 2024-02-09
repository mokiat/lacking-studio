package editor

type NodeKind string

const (
	NodeKindNode       NodeKind = "node"
	NodeKindPointLight NodeKind = "point-light"
)

type NodeList []*Node

type Node struct {
	children NodeList
}
