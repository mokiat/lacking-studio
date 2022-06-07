package model

import (
	"github.com/mokiat/lacking/ui/mvc"
)

var (
	ChangeBinaryEditor                       = mvc.NewChange("binary_editor")
	ChangeBinaryEditorAssetAccordionExpanded = mvc.SubChange(ChangeBinaryEditor, "asset_accordion_expanded")
	ChangeBinaryEditorInfoAccordionExpanded  = mvc.SubChange(ChangeBinaryEditor, "info_accordion_expanded")
)

func NewBinaryEditor(editor *Editor) *BinaryEditor {
	return &BinaryEditor{
		Editor:     editor,
		properties: NewBinaryEditorProperties(),
	}
}

type BinaryEditor struct {
	*Editor
	properties *BinaryEditorProperties
}

func (e *BinaryEditor) Properties() *BinaryEditorProperties {
	return e.properties
}

func NewBinaryEditorProperties() *BinaryEditorProperties {
	return &BinaryEditorProperties{
		Observable:               mvc.NewObservable(),
		isAssetAccordionExpanded: false,
		isInfoAccordionExpanded:  true,
	}
}

type BinaryEditorProperties struct {
	mvc.Observable
	isAssetAccordionExpanded bool
	isInfoAccordionExpanded  bool
}

func (p *BinaryEditorProperties) IsAssetAccordionExpanded() bool {
	return p.isAssetAccordionExpanded
}

func (p *BinaryEditorProperties) SetAssetAccordionExpanded(expanded bool) {
	p.isAssetAccordionExpanded = expanded
	p.SignalChange(ChangeBinaryEditorAssetAccordionExpanded)
}

func (p *BinaryEditorProperties) IsInfoAccordionExpanded() bool {
	return p.isInfoAccordionExpanded
}

func (p *BinaryEditorProperties) SetInfoAccordionExpanded(expanded bool) {
	p.isInfoAccordionExpanded = expanded
	p.SignalChange(ChangeBinaryEditorInfoAccordionExpanded)
}
