package model

import (
	"github.com/mokiat/lacking/ui/mvc"
)

var (
	ChangeTwoDTextureEditor                        = mvc.NewChange("twod_texture_editor")
	ChangeTwoDTextureEditorPropertiesVisible       = mvc.SubChange(ChangeTwoDTextureEditor, "properties_visible")
	ChangeTwoDTextureEditorAssetAccordionExpanded  = mvc.SubChange(ChangeTwoDTextureEditor, "asset_accordion_expanded")
	ChangeTwoDTextureEditorConfigAccordionExpanded = mvc.SubChange(ChangeTwoDTextureEditor, "config_accordion_expanded")
)

func NewTwoDTextureEditor() *TwoDTextureEditor {
	return &TwoDTextureEditor{
		Observable:          mvc.NewObservable(),
		properties:          NewTwoDTextureEditorProperties(),
		isPropertiesVisible: true,
	}
}

type TwoDTextureEditor struct {
	mvc.Observable
	properties          *TwoDTextureEditorProperties
	isPropertiesVisible bool
}

func (e *TwoDTextureEditor) Properties() *TwoDTextureEditorProperties {
	return e.properties
}

func (e *TwoDTextureEditor) IsPropertiesVisible() bool {
	return e.isPropertiesVisible
}

func (e *TwoDTextureEditor) SetPropertiesVisible(visible bool) {
	e.isPropertiesVisible = visible
	e.SignalChange(ChangeTwoDTextureEditorPropertiesVisible)
}

func NewTwoDTextureEditorProperties() *TwoDTextureEditorProperties {
	return &TwoDTextureEditorProperties{
		Observable:                mvc.NewObservable(),
		isAssetAccordionExpanded:  false,
		isConfigAccordionExpanded: true,
	}
}

type TwoDTextureEditorProperties struct {
	mvc.Observable
	isAssetAccordionExpanded  bool
	isConfigAccordionExpanded bool
}

func (p *TwoDTextureEditorProperties) IsAssetAccordionExpanded() bool {
	return p.isAssetAccordionExpanded
}

func (p *TwoDTextureEditorProperties) SetAssetAccordionExpanded(expanded bool) {
	p.isAssetAccordionExpanded = expanded
	p.SignalChange(ChangeTwoDTextureEditorAssetAccordionExpanded)
}

func (p *TwoDTextureEditorProperties) IsConfigAccordionExpanded() bool {
	return p.isConfigAccordionExpanded
}

func (p *TwoDTextureEditorProperties) SetConfigAccordionExpanded(expanded bool) {
	p.isConfigAccordionExpanded = expanded
	p.SignalChange(ChangeTwoDTextureEditorConfigAccordionExpanded)
}
