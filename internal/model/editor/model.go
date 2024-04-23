package editor

import (
	"fmt"

	"github.com/mokiat/gog/ds"
	"github.com/mokiat/gomath/dprec"
	registrymodel "github.com/mokiat/lacking-studio/internal/model/registry"
	"github.com/mokiat/lacking-studio/internal/visualization"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewModel(eventBus *mvc.EventBus, asset *registrymodel.Asset, vis *visualization.Fragment) *Model {
	return &Model{
		eventBus: eventBus,
		asset:    asset,
		vis:      vis,

		navigatorPage: NavigatorPageNodes,
		inspectorPage: InspectorPageAsset,
	}
}

type Model struct {
	eventBus *mvc.EventBus
	asset    *registrymodel.Asset
	vis      *visualization.Fragment

	navigatorPage NavigatorPage
	inspectorPage InspectorPage

	nodes NodeList

	textures []*Texture

	selection any
}

func (m *Model) Visualization() *visualization.Fragment {
	return m.vis
}

func (m *Model) ID() string {
	return m.asset.ID()
}

func (m *Model) Name() string {
	return m.asset.Name()
}

func (m *Model) Image() *ui.Image {
	return m.asset.Image()
}

func (m *Model) Asset() *registrymodel.Asset {
	return m.asset
}

func (m *Model) CanSave() bool {
	return false
}

func (m *Model) NavigatorPage() NavigatorPage {
	return m.navigatorPage
}

func (m *Model) SetNavigatorPage(page NavigatorPage) {
	if page != m.navigatorPage {
		m.navigatorPage = page
		m.eventBus.Notify(NavigatorPageChangedEvent{
			Editor: m,
		})
	}
}

func (m *Model) InspectorPage() InspectorPage {
	return m.inspectorPage
}

func (m *Model) SetInspectorPage(page InspectorPage) {
	if page != m.inspectorPage {
		m.inspectorPage = page
		m.eventBus.Notify(InspectorPageChangedEvent{
			Editor: m,
		})
	}
}

func (m *Model) Textures() []*Texture {
	return m.textures
}

func (m *Model) Selection() any {
	return m.selection
}

func (m *Model) SetSelection(selection any) {
	if selection != m.selection {
		m.selection = selection
		m.eventBus.Notify(SelectionChangedEvent{
			Editor: m,
		})
	}
}

func (m *Model) Nodes() NodeList {
	return m.nodes
}

func (m *Model) AddNode(node Node) {
	m.nodes = append(m.nodes, node)
	m.eventBus.Notify(NodesChangedEvent{
		Editor: m,
	})
}

func (m *Model) CreatePointLight(info PointLightInfo) *PointLight {
	lightVis := m.vis.CreatePointLight()
	return newPointLight(m, lightVis, info.Name)
}

func (m *Model) Save(scope co.Scope /* FIXME */) error {
	img := m.vis.TakeSnapshot(ui.NewSize(128, 128))

	if err := m.asset.Resource().SetPreview(img); err != nil {
		return fmt.Errorf("error setting preview: %w", err)
	}

	if currentImg := m.asset.Image(); currentImg != nil {
		currentImg.Destroy()
		m.asset.SetImage(nil)

		// FIXME: This should probably be done by the respective UI component.
		// If needed, there can be a shared image pool that handles this.
		newImg := co.CreateImage(scope, img)
		m.asset.SetImage(newImg)
	}

	fragment := m.buildFragment()
	if err := m.asset.Resource().SaveContent(fragment); err != nil {
		return fmt.Errorf("error saving content: %w", err)
	}

	log.Info("SAVED!!")

	return nil
}

func (m *Model) Undo() {
	log.Info("UNDO!")
}

func (m *Model) Redo() {
	log.Info("REDO!")
}

func (m *Model) buildFragment() asset.Model {
	var result asset.Model

	nodes := ds.NewList[Node](0)
	m.collectNodes(nodes, m.nodes)

	nodeIndex := make(map[Node]int)
	for i, node := range nodes.Items() {
		nodeIndex[node] = i
	}

	result.Nodes = make([]asset.Node, nodes.Size())
	for i, node := range nodes.Items() {
		var parentIndex = asset.UnspecifiedNodeIndex
		if pIndex, ok := nodeIndex[node.Parent()]; ok {
			parentIndex = int32(pIndex)
		}

		translation := dprec.ZeroVec3()
		if positionable, ok := node.(PositionableNode); ok {
			translation = positionable.Position()
		}

		rotation := dprec.IdentityQuat()
		if rotatable, ok := node.(RotatableNode); ok {
			rotation = rotatable.Rotation()
		}

		scale := dprec.NewVec3(1.0, 1.0, 1.0)
		if scalable, ok := node.(ScalableNode); ok {
			scale = scalable.Scale()
		}

		result.Nodes[i] = asset.Node{
			Name:        node.Name(),
			ParentIndex: parentIndex,
			Translation: translation,
			Rotation:    rotation,
			Scale:       scale,
			Mask:        asset.NodeMaskNone, // TODO
		}
	}

	pointLights := ds.NewList[*PointLight](0)
	m.collectPointLights(pointLights, m.nodes)

	result.PointLights = make([]asset.PointLight, pointLights.Size())
	for i, pointLight := range pointLights.Items() {
		result.PointLights[i] = asset.PointLight{
			NodeIndex:    uint32(nodeIndex[pointLight]),
			EmitColor:    dprec.Vec3Prod(pointLight.EmitColor(), pointLight.EmitIntensity()),
			EmitDistance: pointLight.EmitRange(),
		}
	}

	return result
}

func (m *Model) collectNodes(nodes *ds.List[Node], source NodeList) {
	for _, node := range source {
		nodes.Add(node)
		if extendable, ok := node.(ExtendableNode); ok {
			m.collectNodes(nodes, extendable.Children())
		}
	}
}

func (m *Model) collectPointLights(pointLights *ds.List[*PointLight], source NodeList) {
	for _, node := range source {
		if pointLight, ok := node.(*PointLight); ok {
			pointLights.Add(pointLight)
		}
		if extendable, ok := node.(ExtendableNode); ok {
			m.collectPointLights(pointLights, extendable.Children())
		}
	}
}
