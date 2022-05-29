package model

import (
	"github.com/mokiat/lacking-studio/internal/observer"
)

var (
	ChangeCubeTextureEditor                        = observer.NewChange("cube_texture_editor")
	ChangeCubeTextureEditorPropertiesVisible       = observer.ExtChange(ChangeCubeTextureEditor, "properties_visible")
	ChangeCubeTextureEditorAssetAccordionExpanded  = observer.ExtChange(ChangeCubeTextureEditor, "asset_accordion_expanded")
	ChangeCubeTextureEditorConfigAccordionExpanded = observer.ExtChange(ChangeCubeTextureEditor, "config_accordion_expanded")
)

func NewCubeTextureEditor() *CubeTextureEditor {
	return &CubeTextureEditor{
		Target:              observer.NewTarget(),
		properties:          NewCubeTextureEditorProperties(),
		isPropertiesVisible: true,
	}
}

type CubeTextureEditor struct {
	observer.Target
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
		Target:                    observer.NewTarget(),
		isAssetAccordionExpanded:  false,
		isConfigAccordionExpanded: true,
	}
}

type CubeTextureEditorProperties struct {
	observer.Target
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
