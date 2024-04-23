package editor

import (
	"math/rand"
	"slices"

	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking-studio/internal/visualization"
	"github.com/mokiat/lacking/ui/mvc"
)

type PointLightInfo struct {
	Name string
}

func newPointLight(editorModel *Model, vis *visualization.PointLight, name string) *PointLight {
	return &PointLight{
		editorModel: editorModel,
		eventBus:    editorModel.eventBus,

		vis: vis,

		name:          name,
		position:      dprec.ZeroVec3(),                                              // FIXME
		emitColor:     dprec.NewVec3(rand.Float64(), rand.Float64(), rand.Float64()), // FIXME
		emitIntensity: 100.0,                                                         // FIXME
		emitRange:     20.0,
	}
}

var _ Node = (*PointLight)(nil)
var _ PositionableNode = (*PointLight)(nil)
var _ ExtendableNode = (*PointLight)(nil)

type PointLight struct {
	editorModel *Model
	eventBus    *mvc.EventBus

	vis *visualization.PointLight

	name          string
	parent        ExtendableNode
	children      NodeList
	position      dprec.Vec3
	emitColor     dprec.Vec3
	emitIntensity float64
	emitRange     float64
}

func (l *PointLight) Kind() NodeKind {
	return NodeKindPointLight
}

func (l *PointLight) Name() string {
	return l.name
}

func (l *PointLight) SetName(name string) {
	if name != l.name {
		l.name = name
		l.eventBus.Notify(PointLightChangedEvent{
			Light: l,
		})
	}
}

func (l *PointLight) Delete() {
	defer l.vis.Delete()
}

func (l *PointLight) Parent() ExtendableNode {
	return l.parent
}

func (l *PointLight) SetParent(parent ExtendableNode) {
	if parent != l.parent {
		l.parent = parent
		l.eventBus.Notify(PointLightChangedEvent{
			Light: l,
		})
	}
}

func (l *PointLight) AppendChild(child Node) {
	l.children = append(l.children, child)
}

func (l *PointLight) RemoveChild(child Node) {
	l.children = slices.DeleteFunc(l.children, func(candidate Node) bool {
		return candidate == child
	})
}

func (l *PointLight) Children() NodeList {
	return l.children
}

func (l *PointLight) Position() dprec.Vec3 {
	return l.position
}

func (l *PointLight) SetPosition(position dprec.Vec3) {
	if position != l.position {
		l.position = position

		// FIXME: This needs to be absolute and not relative!
		l.vis.SetMatrix(dprec.TranslationMat4(position.X, position.Y, position.Z))

		l.eventBus.Notify(PointLightChangedEvent{
			Light: l,
		})
	}
}

func (l *PointLight) EmitColor() dprec.Vec3 {
	return l.emitColor
}

func (l *PointLight) SetEmitColor(color dprec.Vec3) {
	if color != l.emitColor {
		l.emitColor = color
		l.eventBus.Notify(PointLightChangedEvent{
			Light: l,
		})
	}
}

func (l *PointLight) EmitIntensity() float64 {
	return l.emitIntensity
}

func (l *PointLight) SetEmitIntensity(intensity float64) {
	if intensity != l.emitIntensity {
		l.emitIntensity = intensity
		l.eventBus.Notify(PointLightChangedEvent{
			Light: l,
		})
	}
}

func (l *PointLight) EmitRange() float64 {
	return l.emitRange
}

func (l *PointLight) SetEmitRange(emitRange float64) {
	if emitRange != l.emitRange {
		l.emitRange = emitRange
		l.eventBus.Notify(PointLightChangedEvent{
			Light: l,
		})
	}
}

type PointLightChangedEvent struct {
	Light *PointLight
}
