package model

import (
	"github.com/mokiat/lacking-studio/internal/observer"
	"github.com/mokiat/lacking/game/asset"
)

var (
	ChangeTwoDTextureEditor                        = observer.NewChange("twod_texture_editor")
	ChangeTwoDTextureEditorPropertiesVisible       = observer.ExtChange(ChangeTwoDTextureEditor, "properties_visible")
	ChangeTwoDTextureEditorAssetAccordionExpanded  = observer.ExtChange(ChangeTwoDTextureEditor, "asset_accordion_expanded")
	ChangeTwoDTextureEditorConfigAccordionExpanded = observer.ExtChange(ChangeTwoDTextureEditor, "config_accordion_expanded")
)

func NewTwoDTextureEditor() *TwoDTextureEditor {
	return &TwoDTextureEditor{
		Target:              observer.NewTarget(),
		properties:          NewTwoDTextureEditorProperties(),
		isPropertiesVisible: true,
	}
}

type TwoDTextureEditor struct {
	observer.Target
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
		Target:                    observer.NewTarget(),
		isAssetAccordionExpanded:  false,
		isConfigAccordionExpanded: true,
	}
}

type TwoDTextureEditorProperties struct {
	observer.Target
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

type AssetAccordion struct{}

type TwoDTextureConfigAccordion struct{}

type ITwoDTextureEditor interface {
	Editor

	Target() observer.Target

	IsPropertiesVisible() bool
	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)
	IsConfigAccordionExpanded() bool
	SetConfigAccordionExpanded(expanded bool)

	Wrapping() asset.WrapMode
	Filtering() asset.FilterMode
	DataFormat() asset.TexelFormat

	// ChangeName(name string)
	ChangeContent(path string)
	ChangeWrapping(wrap asset.WrapMode)
	ChangeFiltering(filter asset.FilterMode)
	ChangeDataFormat(format asset.TexelFormat)

	Visualization() Visualization
}
