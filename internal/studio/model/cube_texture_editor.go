package model

import "github.com/mokiat/lacking/ui"

type CubeTextureEditor interface {
	Editor

	IsAssetAccordionExpanded() bool
	SetAssetAccordionExpanded(expanded bool)
	IsSourceAccordionExpanded() bool
	SetSourceAccordionExpanded(expanded bool)
	IsConfigAccordionExpanded() bool
	SetConfigAccordionExpanded(expanded bool)

	SourceFilename() string
	SourcePreview() ui.Image

	ChangeSource(uri string)
	ReloadSource()
}
