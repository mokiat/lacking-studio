package viewport

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/graphics/shading"
)

func NewCommonData(gfxEngine *graphics.Engine) *CommonData {
	return &CommonData{
		gfxEngine: gfxEngine,
	}
}

type CommonData struct {
	gfxEngine *graphics.Engine

	redMaterialDef       *graphics.MaterialDefinition
	darkRedMaterialDef   *graphics.MaterialDefinition
	greenMaterialDef     *graphics.MaterialDefinition
	darkGreenMaterialDef *graphics.MaterialDefinition
	blueMaterialDef      *graphics.MaterialDefinition
	darkBlueMaterialDef  *graphics.MaterialDefinition
	grayMaterialDef      *graphics.MaterialDefinition
	yellowMaterialDef    *graphics.MaterialDefinition

	gridMeshDef   *graphics.MeshDefinition
	cameraMeshDef *graphics.MeshDefinition
}

func (d *CommonData) Create() {
	d.createMaterials()
	d.createGridMesh()
	d.createCameraMesh()
}

func (d *CommonData) Delete() {
	// NOTE: Using defer to ensure deletion but also reverse execution order.
	defer d.deleteMaterials()
	defer d.deleteGridMesh()
	defer d.deleteCameraMesh()
}

func (d *CommonData) GridMeshDefinition() *graphics.MeshDefinition {
	return d.gridMeshDef
}

func (d *CommonData) CameraMeshDefinition() *graphics.MeshDefinition {
	return d.cameraMeshDef
}

func (d *CommonData) createMaterials() {
	// TODO: Use the same shading for all materials that follow
	// and just adjust the material data for each.

	redShading := d.gfxEngine.CreateShading(graphics.ShadingInfo{
		ForwardFunc: func(palette *shading.ForwardPalette) {
			palette.OutputColor(palette.ConstVec4(1.0, 0.0, 0.0, 1.0))
		},
	})
	d.redMaterialDef = d.gfxEngine.CreateMaterialDefinition(graphics.MaterialDefinitionInfo{
		Shading: redShading,
	})

	darkRedShading := d.gfxEngine.CreateShading(graphics.ShadingInfo{
		ForwardFunc: func(palette *shading.ForwardPalette) {
			palette.OutputColor(palette.ConstVec4(0.1, 0.0, 0.0, 1.0))
		},
	})
	d.darkRedMaterialDef = d.gfxEngine.CreateMaterialDefinition(graphics.MaterialDefinitionInfo{
		Shading: darkRedShading,
	})

	greenShading := d.gfxEngine.CreateShading(graphics.ShadingInfo{
		ForwardFunc: func(palette *shading.ForwardPalette) {
			palette.OutputColor(palette.ConstVec4(0.0, 1.0, 0.0, 1.0))
		},
	})
	d.greenMaterialDef = d.gfxEngine.CreateMaterialDefinition(graphics.MaterialDefinitionInfo{
		Shading: greenShading,
	})

	darkGreenShading := d.gfxEngine.CreateShading(graphics.ShadingInfo{
		ForwardFunc: func(palette *shading.ForwardPalette) {
			palette.OutputColor(palette.ConstVec4(0.0, 0.1, 0.0, 1.0))
		},
	})
	d.darkGreenMaterialDef = d.gfxEngine.CreateMaterialDefinition(graphics.MaterialDefinitionInfo{
		Shading: darkGreenShading,
	})

	blueShading := d.gfxEngine.CreateShading(graphics.ShadingInfo{
		ForwardFunc: func(palette *shading.ForwardPalette) {
			palette.OutputColor(palette.ConstVec4(0.0, 0.0, 1.0, 1.0))
		},
	})
	d.blueMaterialDef = d.gfxEngine.CreateMaterialDefinition(graphics.MaterialDefinitionInfo{
		Shading: blueShading,
	})

	darkBlueShading := d.gfxEngine.CreateShading(graphics.ShadingInfo{
		ForwardFunc: func(palette *shading.ForwardPalette) {
			palette.OutputColor(palette.ConstVec4(0.0, 0.0, 0.1, 1.0))
		},
	})
	d.darkBlueMaterialDef = d.gfxEngine.CreateMaterialDefinition(graphics.MaterialDefinitionInfo{
		Shading: darkBlueShading,
	})

	grayShading := d.gfxEngine.CreateShading(graphics.ShadingInfo{
		ForwardFunc: func(palette *shading.ForwardPalette) {
			palette.OutputColor(palette.ConstVec4(0.5, 0.5, 0.5, 1.0))
		},
	})
	d.grayMaterialDef = d.gfxEngine.CreateMaterialDefinition(graphics.MaterialDefinitionInfo{
		Shading: grayShading,
	})

	yellowShading := d.gfxEngine.CreateShading(graphics.ShadingInfo{
		ForwardFunc: func(palette *shading.ForwardPalette) {
			palette.OutputColor(palette.ConstVec4(1.0, 1.0, 0.0, 1.0))
		},
	})
	d.yellowMaterialDef = d.gfxEngine.CreateMaterialDefinition(graphics.MaterialDefinitionInfo{
		Shading: yellowShading,
	})
}

func (d *CommonData) deleteMaterials() {
	// Nothing to do currently
}

func (d *CommonData) createGridMesh() {
	const (
		gridSize   = 100.0
		gridOffset = 2.0
	)

	gridMeshBuilder := graphics.NewMeshBuilder(
		graphics.MeshBuilderWithCoords(),
	)

	// Positive X axis
	vertexOffset := gridMeshBuilder.VertexOffset()
	gridMeshBuilder.Vertex().
		Coord(0.0, 0.0, 0.0)
	gridMeshBuilder.Vertex().
		Coord(gridSize, 0.0, 0.0)
	indexStart, indexEnd := gridMeshBuilder.IndexLine(vertexOffset, vertexOffset+1)
	gridMeshBuilder.Fragment(graphics.PrimitiveLines, d.redMaterialDef, indexStart, indexEnd-indexStart)

	// Negative X axis
	vertexOffset = gridMeshBuilder.VertexOffset()
	gridMeshBuilder.Vertex().
		Coord(0.0, 0.0, 0.0)
	gridMeshBuilder.Vertex().
		Coord(-gridSize, 0.0, 0.0)

	indexStart, indexEnd = gridMeshBuilder.IndexLine(vertexOffset, vertexOffset+1)
	gridMeshBuilder.Fragment(graphics.PrimitiveLines, d.darkRedMaterialDef, indexStart, indexEnd-indexStart)

	// Positive Z axis
	vertexOffset = gridMeshBuilder.VertexOffset()
	gridMeshBuilder.Vertex().
		Coord(0.0, 0.0, 0.0)
	gridMeshBuilder.Vertex().
		Coord(0.0, 0.0, gridSize)
	indexStart, indexEnd = gridMeshBuilder.IndexLine(vertexOffset, vertexOffset+1)
	gridMeshBuilder.Fragment(graphics.PrimitiveLines, d.greenMaterialDef, indexStart, indexEnd-indexStart)

	// Negative Z axis
	vertexOffset = gridMeshBuilder.VertexOffset()
	gridMeshBuilder.Vertex().
		Coord(0.0, 0.0, 0.0)
	gridMeshBuilder.Vertex().
		Coord(0.0, 0.0, -gridSize)
	indexStart, indexEnd = gridMeshBuilder.IndexLine(vertexOffset, vertexOffset+1)
	gridMeshBuilder.Fragment(graphics.PrimitiveLines, d.darkGreenMaterialDef, indexStart, indexEnd-indexStart)

	// Grid
	indexStart = gridMeshBuilder.IndexOffset()
	for i := 1; i <= int(gridSize/gridOffset); i++ {
		// Along X axis
		vertexOffset := gridMeshBuilder.VertexOffset()
		gridMeshBuilder.Vertex().
			Coord(-gridSize, 0.0, -float32(i)*gridOffset)
		gridMeshBuilder.Vertex().
			Coord(gridSize, 0.0, -float32(i)*gridOffset)
		gridMeshBuilder.IndexLine(vertexOffset, vertexOffset+1)

		vertexOffset = gridMeshBuilder.VertexOffset()
		gridMeshBuilder.Vertex().
			Coord(-gridSize, 0.0, float32(i)*gridOffset)
		gridMeshBuilder.Vertex().
			Coord(gridSize, 0.0, float32(i)*gridOffset)
		gridMeshBuilder.IndexLine(vertexOffset, vertexOffset+1)

		// Along Z axis
		vertexOffset = gridMeshBuilder.VertexOffset()
		gridMeshBuilder.Vertex().
			Coord(-float32(i)*gridOffset, 0.0, -gridSize)
		gridMeshBuilder.Vertex().
			Coord(-float32(i)*gridOffset, 0.0, gridSize)
		gridMeshBuilder.IndexLine(vertexOffset, vertexOffset+1)

		vertexOffset = gridMeshBuilder.VertexOffset()
		gridMeshBuilder.Vertex().
			Coord(float32(i)*gridOffset, 0.0, -gridSize)
		gridMeshBuilder.Vertex().
			Coord(float32(i)*gridOffset, 0.0, gridSize)
		gridMeshBuilder.IndexLine(vertexOffset, vertexOffset+1)
	}
	indexEnd = gridMeshBuilder.IndexOffset()
	gridMeshBuilder.Fragment(graphics.PrimitiveLines, d.grayMaterialDef, indexStart, indexEnd-indexStart)

	d.gridMeshDef = d.gfxEngine.CreateMeshDefinition(gridMeshBuilder.BuildInfo())
}

func (d *CommonData) deleteGridMesh() {
	d.gridMeshDef.Delete()
}

func (d *CommonData) createCameraMesh() {
	meshBuilder := graphics.NewSimpleMeshBuilder(d.yellowMaterialDef)

	meshBuilder.Solid().
		Cuboid(sprec.ZeroVec3(), sprec.IdentityQuat(), sprec.NewVec3(0.2, 0.3, 0.5)).
		Cylinder(sprec.NewVec3(0.0, 0.25, 0.1), sprec.RotationQuat(sprec.Degrees(90), sprec.BasisZVec3()), 0.15, 0.1, 20).
		Cylinder(sprec.NewVec3(0.0, 0.23, -0.13), sprec.RotationQuat(sprec.Degrees(90), sprec.BasisZVec3()), 0.1, 0.1, 20).
		Cone(sprec.NewVec3(0.0, 0.0, 0.3), sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisXVec3()), 0.2, 0.3, 20)

	d.cameraMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildInfo())
}

func (d *CommonData) deleteCameraMesh() {
	d.cameraMeshDef.Delete()
}
