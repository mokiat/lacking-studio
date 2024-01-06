package editor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mokiat/gog/ds"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var ModelImport = co.Define(&modelImportComponent{})

type ModelImportData struct {
	Model *pack.Model
}

type ModelImportCallbackData struct {
	OnImport func(model *pack.Model)
}

type modelImportComponent struct {
	co.BaseComponent

	model *pack.Model

	nodesExpanded      bool
	texturesExpanded   bool
	materialsExpanded  bool
	meshesExpanded     bool
	animationsExpanded bool

	selectedNodes     *ds.Set[*pack.Node]
	selectedTextures  *ds.Set[*pack.Image]
	selectedMaterials *ds.Set[*pack.Material]
	selectedMeshes    *ds.Set[*pack.MeshDefinition]

	onImport func(model *pack.Model)
}

func (c *modelImportComponent) OnCreate() {
	c.nodesExpanded = true
	c.texturesExpanded = true
	c.materialsExpanded = true
	c.meshesExpanded = true
	c.animationsExpanded = true

	c.selectedNodes = ds.NewSet[*pack.Node](0)
	c.selectedTextures = ds.NewSet[*pack.Image](0)
	c.selectedMaterials = ds.NewSet[*pack.Material](0)
	c.selectedMeshes = ds.NewSet[*pack.MeshDefinition](0)

	data := co.GetData[ModelImportData](c.Properties())
	c.model = data.Model

	callbackData := co.GetCallbackData[ModelImportCallbackData](c.Properties())
	c.onImport = callbackData.OnImport
	if c.onImport == nil {
		c.onImport = func(model *pack.Model) {}
	}

	c.handleSelectAll()
}

func (c *modelImportComponent) Render() co.Instance {
	return co.New(std.Modal, func() {
		co.WithLayoutData(layout.Data{
			Width:            opt.V(800),
			Height:           opt.V(600),
			HorizontalCenter: opt.V(0),
			VerticalCenter:   opt.V(0),
		})

		co.WithChild("header", co.New(std.Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})

			co.WithChild("select-all", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Text: "Select All",
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleSelectAll,
				})
			}))

			co.WithChild("separator", co.New(std.ToolbarSeparator, nil))

			co.WithChild("deselect-all", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Text: "Deselect All",
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleDeselectAll,
				})
			}))
		}))

		co.WithChild("article", co.New(std.ScrollPane, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ScrollPaneData{
				DisableHorizontal: true,
			})

			co.WithChild("content", co.New(std.Element, func() {
				co.WithLayoutData(layout.Data{
					GrowHorizontally: true,
				})
				co.WithData(std.ElementData{
					Layout: layout.Vertical(layout.VerticalSettings{
						ContentSpacing: 2,
					}),
					Padding: ui.UniformSpacing(2),
				})

				co.WithChild("nodes", co.New(std.Accordion, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.AccordionData{
						Title:    "Nodes",
						Expanded: c.nodesExpanded,
					})
					co.WithCallbackData(std.AccordionCallbackData{
						OnToggle: c.handleNodesToggle,
					})

					c.eachNode(func(index, depth int, node *pack.Node, selected bool) {
						co.WithChild(strconv.Itoa(index), co.New(std.Checkbox, func() {
							co.WithData(std.CheckboxData{
								Checked: selected,
								Label:   fmt.Sprintf("%s%s", strings.Repeat("  ", depth), node.Name),
							})
							co.WithCallbackData(std.CheckboxCallbackData{
								OnToggle: func(active bool) {
									c.handleNodeSelection(node, active)
								},
							})
						}))
					})
				}))

				co.WithChild("textures", co.New(std.Accordion, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.AccordionData{
						Title:    "Textures",
						Expanded: c.texturesExpanded,
					})
					co.WithCallbackData(std.AccordionCallbackData{
						OnToggle: c.handleTexturesToggle,
					})

					c.eachTexture(func(index int, texture *pack.Image, selected bool) {
						co.WithChild(strconv.Itoa(index), co.New(std.Checkbox, func() {
							co.WithData(std.CheckboxData{
								Checked: selected,
								Label:   texture.Name,
							})
							co.WithCallbackData(std.CheckboxCallbackData{
								OnToggle: func(active bool) {
									c.handleTextureSelection(texture, active)
								},
							})
						}))
					})
				}))

				co.WithChild("materials", co.New(std.Accordion, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.AccordionData{
						Title:    "Materials",
						Expanded: c.materialsExpanded,
					})
					co.WithCallbackData(std.AccordionCallbackData{
						OnToggle: c.handleMaterialsToggle,
					})

					c.eachMaterial(func(index int, material *pack.Material, selected bool) {
						co.WithChild(strconv.Itoa(index), co.New(std.Checkbox, func() {
							co.WithData(std.CheckboxData{
								Checked: selected,
								Label:   material.Name,
							})
							co.WithCallbackData(std.CheckboxCallbackData{
								OnToggle: func(active bool) {
									c.handleMaterialSelection(material, active)
								},
							})
						}))
					})
				}))

				co.WithChild("meshes", co.New(std.Accordion, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.AccordionData{
						Title:    "Meshes",
						Expanded: c.meshesExpanded,
					})
					co.WithCallbackData(std.AccordionCallbackData{
						OnToggle: c.handleMeshesToggle,
					})

					c.eachMesh(func(index int, mesh *pack.MeshDefinition, selected bool) {
						co.WithChild(strconv.Itoa(index), co.New(std.Checkbox, func() {
							co.WithData(std.CheckboxData{
								Checked: selected,
								Label:   mesh.Name,
							})
							co.WithCallbackData(std.CheckboxCallbackData{
								OnToggle: func(active bool) {
									c.handleMeshSelection(mesh, active)
								},
							})
						}))
					})
				}))

				co.WithChild("animations", co.New(std.Accordion, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.AccordionData{
						Title:    "Animations",
						Expanded: c.animationsExpanded,
					})
					co.WithCallbackData(std.AccordionCallbackData{
						OnToggle: c.handleAnimationsToggle,
					})
				}))
			}))
		}))

		co.WithChild("footer", co.New(std.Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentBottom,
			})
			co.WithData(std.ToolbarData{
				Positioning: std.ToolbarPositioningBottom,
			})

			co.WithChild("import", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Text: "Import",
				})
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleImport,
				})
			}))

			co.WithChild("cancel", co.New(std.ToolbarButton, func() {
				co.WithData(std.ToolbarButtonData{
					Text: "Cancel",
				})
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
				})
				co.WithCallbackData(std.ToolbarButtonCallbackData{
					OnClick: c.handleCancel,
				})
			}))
		}))
	})
}

func (c *modelImportComponent) handleSelectAll() {
	c.eachNode(func(_ int, _ int, node *pack.Node, _ bool) {
		c.selectedNodes.Add(node)
	})
	c.eachTexture(func(_ int, texture *pack.Image, _ bool) {
		c.selectedTextures.Add(texture)
	})
	c.eachMaterial(func(_ int, material *pack.Material, _ bool) {
		c.selectedMaterials.Add(material)
	})
	c.eachMesh(func(_ int, mesh *pack.MeshDefinition, _ bool) {
		c.selectedMeshes.Add(mesh)
	})
	c.Invalidate()
}

func (c *modelImportComponent) handleDeselectAll() {
	c.eachNode(func(_ int, _ int, node *pack.Node, _ bool) {
		c.selectedNodes.Remove(node)
	})
	c.eachTexture(func(_ int, texture *pack.Image, _ bool) {
		c.selectedTextures.Remove(texture)
	})
	c.eachMaterial(func(_ int, material *pack.Material, _ bool) {
		c.selectedMaterials.Remove(material)
	})
	c.eachMesh(func(_ int, mesh *pack.MeshDefinition, _ bool) {
		c.selectedMeshes.Remove(mesh)
	})
	c.Invalidate()
}

func (c *modelImportComponent) handleNodesToggle() {
	c.nodesExpanded = !c.nodesExpanded
	c.Invalidate()
}

func (c *modelImportComponent) handleNodeSelection(node *pack.Node, active bool) {
	if active {
		c.selectedNodes.Add(node)
	} else {
		c.selectedNodes.Remove(node)
	}
	c.Invalidate()
}

func (c *modelImportComponent) handleTexturesToggle() {
	c.texturesExpanded = !c.texturesExpanded
	c.Invalidate()
}

func (c *modelImportComponent) handleTextureSelection(texture *pack.Image, active bool) {
	if active {
		c.selectedTextures.Add(texture)
	} else {
		c.selectedTextures.Remove(texture)
	}
	c.Invalidate()
}

func (c *modelImportComponent) handleMaterialsToggle() {
	c.materialsExpanded = !c.materialsExpanded
	c.Invalidate()
}

func (c *modelImportComponent) handleMaterialSelection(material *pack.Material, active bool) {
	if active {
		c.selectedMaterials.Add(material)
	} else {
		c.selectedMaterials.Remove(material)
	}
	c.Invalidate()
}

func (c *modelImportComponent) handleMeshesToggle() {
	c.meshesExpanded = !c.meshesExpanded
	c.Invalidate()
}

func (c *modelImportComponent) handleMeshSelection(mesh *pack.MeshDefinition, active bool) {
	if active {
		c.selectedMeshes.Add(mesh)
	} else {
		c.selectedMeshes.Remove(mesh)
	}
	c.Invalidate()
}

func (c *modelImportComponent) handleAnimationsToggle() {
	c.animationsExpanded = !c.animationsExpanded
	c.Invalidate()
}

func (c *modelImportComponent) handleImport() {
	co.CloseOverlay(c.Scope())
	c.onImport(c.model)
}

func (c *modelImportComponent) handleCancel() {
	co.CloseOverlay(c.Scope())
}

func (c *modelImportComponent) eachNode(cb func(index, depth int, node *pack.Node, selected bool)) {
	var index int
	var visit func(node *pack.Node, depth int)
	visit = func(node *pack.Node, depth int) {
		cb(index, depth, node, c.selectedNodes.Contains(node))
		index++
		for _, child := range node.Children {
			visit(child, depth+1)
		}
	}

	for _, node := range c.model.RootNodes {
		visit(node, 0)
	}
}

func (c *modelImportComponent) eachTexture(cb func(index int, texture *pack.Image, selected bool)) {
	for i, texture := range c.model.Textures {
		cb(i, texture, c.selectedTextures.Contains(texture))
	}
}

func (c *modelImportComponent) eachMaterial(cb func(index int, material *pack.Material, selected bool)) {
	for i, material := range c.model.Materials {
		cb(i, material, c.selectedMaterials.Contains(material))
	}
}

func (c *modelImportComponent) eachMesh(cb func(index int, mesh *pack.MeshDefinition, selected bool)) {
	for i, mesh := range c.model.MeshDefinitions {
		cb(i, mesh, c.selectedMeshes.Contains(mesh))
	}
}
