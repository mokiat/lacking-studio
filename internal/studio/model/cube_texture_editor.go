package model

import (
	"github.com/mokiat/lacking/ui/mvc"
)

var (
	ChangeCubeTextureEditor                        = mvc.NewChange("cube_texture_editor")
	ChangeCubeTextureEditorPropertiesVisible       = mvc.SubChange(ChangeCubeTextureEditor, "properties_visible")
	ChangeCubeTextureEditorAssetAccordionExpanded  = mvc.SubChange(ChangeCubeTextureEditor, "asset_accordion_expanded")
	ChangeCubeTextureEditorConfigAccordionExpanded = mvc.SubChange(ChangeCubeTextureEditor, "config_accordion_expanded")
)

func NewCubeTextureEditor() *CubeTextureEditor {
	return &CubeTextureEditor{
		Observable:          mvc.NewObservable(),
		properties:          NewCubeTextureEditorProperties(),
		isPropertiesVisible: true,
	}
}

type CubeTextureEditor struct {
	mvc.Observable
	properties          *CubeTextureEditorProperties
	isPropertiesVisible bool
}

func (e *CubeTextureEditor) Properties() *CubeTextureEditorProperties {
	return e.properties
}

func (e *CubeTextureEditor) IsPropertiesVisible() bool {
	return e.isPropertiesVisible
}

func (e *CubeTextureEditor) SetPropertiesVisible(visible bool) {
	e.isPropertiesVisible = visible
	e.SignalChange(ChangeCubeTextureEditorPropertiesVisible)
}

func NewCubeTextureEditorProperties() *CubeTextureEditorProperties {
	return &CubeTextureEditorProperties{
		Observable:                mvc.NewObservable(),
		isAssetAccordionExpanded:  false,
		isConfigAccordionExpanded: true,
	}
}

type CubeTextureEditorProperties struct {
	mvc.Observable
	isAssetAccordionExpanded  bool
	isConfigAccordionExpanded bool
}

func (p *CubeTextureEditorProperties) IsAssetAccordionExpanded() bool {
	return p.isAssetAccordionExpanded
}

func (p *CubeTextureEditorProperties) SetAssetAccordionExpanded(expanded bool) {
	p.isAssetAccordionExpanded = expanded
	p.SignalChange(ChangeCubeTextureEditorAssetAccordionExpanded)
}

func (p *CubeTextureEditorProperties) IsConfigAccordionExpanded() bool {
	return p.isConfigAccordionExpanded
}

func (p *CubeTextureEditorProperties) SetConfigAccordionExpanded(expanded bool) {
	p.isConfigAccordionExpanded = expanded
	p.SignalChange(ChangeCubeTextureEditorConfigAccordionExpanded)
}
