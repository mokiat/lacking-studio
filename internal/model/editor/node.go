package editor

import "github.com/mokiat/gomath/dprec"

type NodeKind string

const (
	NodeKindNode       NodeKind = "node"
	NodeKindPointLight NodeKind = "point-light"
)

type NodeList []Node

func (l NodeList) HasNode(node Node) bool {
	for _, n := range l {
		if n == node {
			return true
		}
		if extendable, ok := n.(ExtendableNode); ok {
			if extendable.Children().HasNode(node) {
				return true
			}
		}
	}
	return false
}

type Node interface {
	Kind() NodeKind
	Name() string
	SetName(name string)
	Parent() ExtendableNode
	SetParent(parent ExtendableNode)
	Delete()
}

type PositionableNode interface {
	Node
	Position() dprec.Vec3
	SetPosition(position dprec.Vec3)
}

type RotatableNode interface {
	Node
	Rotation() dprec.Quat
	SetRotation(rotation dprec.Quat)
}

type ScalableNode interface {
	Node
	Scale() dprec.Vec3
	SetScale(scale dprec.Vec3)
}

type ExtendableNode interface {
	Node
	AppendChild(child Node)
	RemoveChild(child Node)
	Children() NodeList
}
