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
	nodeMeshDef   *graphics.MeshDefinition
	cameraMeshDef *graphics.MeshDefinition
}

func (d *CommonData) Create() {
	d.createMaterials()
	d.createGridMesh()
	d.createNodeMesh()
	d.createCameraMesh()
}

func (d *CommonData) Delete() {
	// NOTE: Using defer to ensure deletion but also reverse execution order.
	defer d.deleteMaterials()
	defer d.deleteGridMesh()
	defer d.deleteNodeMesh()
	defer d.deleteCameraMesh()
}

func (d *CommonData) GridMeshDefinition() *graphics.MeshDefinition {
	return d.gridMeshDef
}

func (d *CommonData) NodeMeshDefinition() *graphics.MeshDefinition {
	return d.nodeMeshDef
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

	meshBuilder := graphics.NewSimpleMeshBuilder()

	// Positive X axis
	meshBuilder.Wireframe(d.redMaterialDef).
		Line(sprec.ZeroVec3(), sprec.NewVec3(gridSize, 0.0, 0.0))

	// Negative X axis
	meshBuilder.Wireframe(d.darkRedMaterialDef).
		Line(sprec.ZeroVec3(), sprec.NewVec3(-gridSize, 0.0, 0.0))

	// Positive Z axis
	meshBuilder.Wireframe(d.greenMaterialDef).
		Line(sprec.ZeroVec3(), sprec.NewVec3(0.0, 0.0, gridSize))

	meshBuilder.Wireframe(d.darkGreenMaterialDef).
		Line(sprec.ZeroVec3(), sprec.NewVec3(0.0, 0.0, -gridSize))

	// Grid
	lines := meshBuilder.Wireframe(d.grayMaterialDef)
	for i := 1; i <= int(gridSize/gridOffset); i++ {
		// Along X axis
		lines.Line(
			sprec.NewVec3(-gridSize, 0.0, -float32(i)*gridOffset),
			sprec.NewVec3(gridSize, 0.0, -float32(i)*gridOffset),
		)
		lines.Line(
			sprec.NewVec3(-gridSize, 0.0, float32(i)*gridOffset),
			sprec.NewVec3(gridSize, 0.0, float32(i)*gridOffset),
		)
		// Along Z axis
		lines.Line(
			sprec.NewVec3(-float32(i)*gridOffset, 0.0, -gridSize),
			sprec.NewVec3(-float32(i)*gridOffset, 0.0, gridSize),
		)
		lines.Line(
			sprec.NewVec3(float32(i)*gridOffset, 0.0, -gridSize),
			sprec.NewVec3(float32(i)*gridOffset, 0.0, gridSize),
		)
	}

	d.gridMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildInfo())
}

func (d *CommonData) deleteGridMesh() {
	d.gridMeshDef.Delete()
}

func (d *CommonData) createNodeMesh() {
	meshBuilder := graphics.NewSimpleMeshBuilder()

	meshBuilder.Solid(d.yellowMaterialDef).
		Cuboid(sprec.ZeroVec3(), sprec.IdentityQuat(), sprec.NewVec3(0.2, 0.2, 0.2))

	meshBuilder.Wireframe(d.yellowMaterialDef).
		// front-top-left
		Line(sprec.NewVec3(-0.2, 0.2, 0.2), sprec.NewVec3(-0.1, 0.2, 0.2)).
		Line(sprec.NewVec3(-0.2, 0.2, 0.2), sprec.NewVec3(-0.2, 0.1, 0.2)).
		Line(sprec.NewVec3(-0.2, 0.2, 0.2), sprec.NewVec3(-0.2, 0.2, 0.1)).
		// front-top-right
		Line(sprec.NewVec3(0.2, 0.2, 0.2), sprec.NewVec3(0.1, 0.2, 0.2)).
		Line(sprec.NewVec3(0.2, 0.2, 0.2), sprec.NewVec3(0.2, 0.1, 0.2)).
		Line(sprec.NewVec3(0.2, 0.2, 0.2), sprec.NewVec3(0.2, 0.2, 0.1)).
		// front-bottom-left
		Line(sprec.NewVec3(-0.2, -0.2, 0.2), sprec.NewVec3(-0.1, -0.2, 0.2)).
		Line(sprec.NewVec3(-0.2, -0.2, 0.2), sprec.NewVec3(-0.2, -0.1, 0.2)).
		Line(sprec.NewVec3(-0.2, -0.2, 0.2), sprec.NewVec3(-0.2, -0.2, 0.1)).
		// front-bottom-right
		Line(sprec.NewVec3(0.2, -0.2, 0.2), sprec.NewVec3(0.1, -0.2, 0.2)).
		Line(sprec.NewVec3(0.2, -0.2, 0.2), sprec.NewVec3(0.2, -0.1, 0.2)).
		Line(sprec.NewVec3(0.2, -0.2, 0.2), sprec.NewVec3(0.2, -0.2, 0.1)).
		// back-top-left
		Line(sprec.NewVec3(-0.2, 0.2, -0.2), sprec.NewVec3(-0.1, 0.2, -0.2)).
		Line(sprec.NewVec3(-0.2, 0.2, -0.2), sprec.NewVec3(-0.2, 0.1, -0.2)).
		Line(sprec.NewVec3(-0.2, 0.2, -0.2), sprec.NewVec3(-0.2, 0.2, -0.1)).
		// back-top-right
		Line(sprec.NewVec3(0.2, 0.2, -0.2), sprec.NewVec3(0.1, 0.2, -0.2)).
		Line(sprec.NewVec3(0.2, 0.2, -0.2), sprec.NewVec3(0.2, 0.1, -0.2)).
		Line(sprec.NewVec3(0.2, 0.2, -0.2), sprec.NewVec3(0.2, 0.2, -0.1)).
		// back-bottom-left
		Line(sprec.NewVec3(-0.2, -0.2, -0.2), sprec.NewVec3(-0.1, -0.2, -0.2)).
		Line(sprec.NewVec3(-0.2, -0.2, -0.2), sprec.NewVec3(-0.2, -0.1, -0.2)).
		Line(sprec.NewVec3(-0.2, -0.2, -0.2), sprec.NewVec3(-0.2, -0.2, -0.1)).
		// back-bottom-right
		Line(sprec.NewVec3(0.2, -0.2, -0.2), sprec.NewVec3(0.1, -0.2, -0.2)).
		Line(sprec.NewVec3(0.2, -0.2, -0.2), sprec.NewVec3(0.2, -0.1, -0.2)).
		Line(sprec.NewVec3(0.2, -0.2, -0.2), sprec.NewVec3(0.2, -0.2, -0.1))

	meshBuilder.Wireframe(d.redMaterialDef).
		Line(sprec.ZeroVec3(), sprec.NewVec3(0.2, 0.0, 0.0))
	meshBuilder.Wireframe(d.greenMaterialDef).
		Line(sprec.ZeroVec3(), sprec.NewVec3(0.0, 0.0, 0.2))
	meshBuilder.Wireframe(d.blueMaterialDef).
		Line(sprec.ZeroVec3(), sprec.NewVec3(0.0, 0.2, 0.0))

	d.nodeMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildInfo())
}

func (d *CommonData) deleteNodeMesh() {
	d.nodeMeshDef.Delete()
}

func (d *CommonData) createCameraMesh() {
	meshBuilder := graphics.NewSimpleMeshBuilder()

	meshBuilder.Solid(d.yellowMaterialDef).
		Cuboid(sprec.ZeroVec3(), sprec.IdentityQuat(), sprec.NewVec3(0.2, 0.3, 0.5)).
		Cylinder(sprec.NewVec3(0.0, 0.25, 0.1), sprec.RotationQuat(sprec.Degrees(90), sprec.BasisZVec3()), 0.15, 0.1, 20).
		Cylinder(sprec.NewVec3(0.0, 0.23, -0.13), sprec.RotationQuat(sprec.Degrees(90), sprec.BasisZVec3()), 0.1, 0.1, 20).
		Cone(sprec.NewVec3(0.0, 0.0, 0.3), sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisXVec3()), 0.2, 0.3, 20)

	d.cameraMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildInfo())
}

func (d *CommonData) deleteCameraMesh() {
	d.cameraMeshDef.Delete()
}
