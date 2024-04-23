package viewport

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/util/blob"
)

func NewCommonData(gfxEngine *graphics.Engine) *CommonData {
	return &CommonData{
		gfxEngine: gfxEngine,
	}
}

type CommonData struct {
	gfxEngine *graphics.Engine

	redMaterial        *graphics.Material
	darkRedMaterial    *graphics.Material
	greenMaterial      *graphics.Material
	darkGreenMaterial  *graphics.Material
	blueMaterial       *graphics.Material
	darkBlueMaterial   *graphics.Material
	grayMaterial       *graphics.Material
	yellowMaterial     *graphics.Material
	darkYellowMaterial *graphics.Material

	gridGeometry *graphics.MeshGeometry
	gridMeshDef  *graphics.MeshDefinition

	nodeMeshGeometry *graphics.MeshGeometry
	nodeMeshDef      *graphics.MeshDefinition

	cameraMeshGeometry *graphics.MeshGeometry
	cameraMeshDef      *graphics.MeshDefinition

	ambientLightMeshGeometry *graphics.MeshGeometry
	ambientLightMeshDef      *graphics.MeshDefinition

	pointLightMeshGeometry *graphics.MeshGeometry
	pointLightMeshDef      *graphics.MeshDefinition

	spotLightMeshGeometry *graphics.MeshGeometry
	spotLightMeshDef      *graphics.MeshDefinition

	directionalLightMeshGeometry *graphics.MeshGeometry
	directionalLightMeshDef      *graphics.MeshDefinition
}

func (d *CommonData) Create() {
	d.createMaterials()
	d.createGridMesh()
	d.createNodeMesh()
	d.createCameraMesh()
	d.createAmbientLightMesh()
	d.createPointLightMesh()
	d.createSpotLightMesh()
	d.createDirectionalLightMesh()
}

func (d *CommonData) Delete() {
	// NOTE: Using defer to ensure deletion but also reverse execution order.
	defer d.deleteMaterials()
	defer d.deleteGridMesh()
	defer d.deleteNodeMesh()
	defer d.deleteCameraMesh()
	defer d.deleteAmbientLightMesh()
	defer d.deletePointLightMesh()
	defer d.deleteSpotLightMesh()
	defer d.deleteDirectionalLightMesh()
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

func (d *CommonData) AmbientLightMeshDefinition() *graphics.MeshDefinition {
	return d.ambientLightMeshDef
}

func (d *CommonData) PointLightMeshDefinition() *graphics.MeshDefinition {
	return d.pointLightMeshDef
}

func (d *CommonData) SpotLightMeshDefinition() *graphics.MeshDefinition {
	return d.spotLightMeshDef
}

func (d *CommonData) DirectionalLightMeshDefinition() *graphics.MeshDefinition {
	return d.directionalLightMeshDef
}

func (d *CommonData) createMaterials() {
	// TODO: Use the same shading for all materials that follow
	// and just adjust the material data for each.

	colorShader := d.gfxEngine.CreateForwardShader(graphics.ShaderInfo{
		SourceCode: `
			uniforms {
				color vec4,
			}

			func #fragment() {
				#color = color
			}	
		`,
	})

	d.redMaterial = d.gfxEngine.CreateMaterial(graphics.MaterialInfo{
		Name: "ColorRed",
		ForwardPasses: []graphics.ForwardRenderPassInfo{
			{
				Shader: colorShader,
			},
		},
	})
	d.redMaterial.SetProperty("color", sprec.NewVec4(1.0, 0.0, 0.0, 1.0))

	d.darkRedMaterial = d.gfxEngine.CreateMaterial(graphics.MaterialInfo{
		Name: "ColorDarkRed",
		ForwardPasses: []graphics.ForwardRenderPassInfo{
			{
				Shader: colorShader,
			},
		},
	})
	d.darkRedMaterial.SetProperty("color", sprec.NewVec4(0.3, 0.0, 0.0, 1.0))

	d.greenMaterial = d.gfxEngine.CreateMaterial(graphics.MaterialInfo{
		Name: "ColorGreen",
		ForwardPasses: []graphics.ForwardRenderPassInfo{
			{
				Shader: colorShader,
			},
		},
	})
	d.greenMaterial.SetProperty("color", sprec.NewVec4(0.0, 1.0, 0.0, 1.0))

	d.darkGreenMaterial = d.gfxEngine.CreateMaterial(graphics.MaterialInfo{
		Name: "ColorDarkGreen",
		ForwardPasses: []graphics.ForwardRenderPassInfo{
			{
				Shader: colorShader,
			},
		},
	})
	d.darkGreenMaterial.SetProperty("color", sprec.NewVec4(0.0, 0.3, 0.0, 1.0))

	d.blueMaterial = d.gfxEngine.CreateMaterial(graphics.MaterialInfo{
		Name: "ColorBlue",
		ForwardPasses: []graphics.ForwardRenderPassInfo{
			{
				Shader: colorShader,
			},
		},
	})
	d.blueMaterial.SetProperty("color", sprec.NewVec4(0.0, 0.0, 1.0, 1.0))

	d.darkBlueMaterial = d.gfxEngine.CreateMaterial(graphics.MaterialInfo{
		Name: "ColorDarkBlue",
		ForwardPasses: []graphics.ForwardRenderPassInfo{
			{
				Shader: colorShader,
			},
		},
	})
	d.darkBlueMaterial.SetProperty("color", sprec.NewVec4(0.0, 0.0, 0.3, 1.0))

	d.grayMaterial = d.gfxEngine.CreateMaterial(graphics.MaterialInfo{
		Name: "ColorGray",
		ForwardPasses: []graphics.ForwardRenderPassInfo{
			{
				Shader: colorShader,
			},
		},
	})
	d.grayMaterial.SetProperty("color", sprec.NewVec4(0.3, 0.3, 0.3, 1.0))

	d.yellowMaterial = d.gfxEngine.CreateMaterial(graphics.MaterialInfo{
		Name: "ColorYellow",
		ForwardPasses: []graphics.ForwardRenderPassInfo{
			{
				Shader: colorShader,
			},
		},
	})
	d.yellowMaterial.SetProperty("color", sprec.NewVec4(1.0, 1.0, 0.0, 1.0))

	d.darkYellowMaterial = d.gfxEngine.CreateMaterial(graphics.MaterialInfo{
		Name: "ColorDarkYellow",
		ForwardPasses: []graphics.ForwardRenderPassInfo{
			{
				Shader: colorShader,
			},
		},
	})
	d.darkYellowMaterial.SetProperty("color", sprec.NewVec4(0.3, 0.3, 0.0, 1.0))
}

func (d *CommonData) deleteMaterials() {
	// Nothing to do currently
}

func (d *CommonData) createGridMesh() {
	const (
		gridSize   = 100.0
		gridOffset = 2.0
	)

	meshBuilder := graphics.NewShapeBuilder()

	// Positive X axis
	meshBuilder.Wireframe(d.redMaterial).
		Line(sprec.ZeroVec3(), sprec.NewVec3(gridSize, 0.0, 0.0))

	// Negative X axis
	meshBuilder.Wireframe(d.darkRedMaterial).
		Line(sprec.ZeroVec3(), sprec.NewVec3(-gridSize, 0.0, 0.0))

	// Positive Z axis
	meshBuilder.Wireframe(d.greenMaterial).
		Line(sprec.ZeroVec3(), sprec.NewVec3(0.0, 0.0, gridSize))

	meshBuilder.Wireframe(d.darkGreenMaterial).
		Line(sprec.ZeroVec3(), sprec.NewVec3(0.0, 0.0, -gridSize))

	// Grid
	lines := meshBuilder.Wireframe(d.grayMaterial)
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

	d.gridGeometry = d.gfxEngine.CreateMeshGeometry(meshBuilder.BuildGeometryInfo())
	d.gridMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildMeshDefinitionInfo(d.gridGeometry))
}

func (d *CommonData) deleteGridMesh() {
	defer d.gridGeometry.Delete()
	defer d.gridMeshDef.Delete()
}

func (d *CommonData) createNodeMesh() {
	meshBuilder := graphics.NewShapeBuilder()

	meshBuilder.Solid(d.yellowMaterial).
		Cuboid(sprec.ZeroVec3(), sprec.IdentityQuat(), sprec.NewVec3(0.2, 0.2, 0.2))

	meshBuilder.Wireframe(d.darkYellowMaterial).
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

	meshBuilder.Wireframe(d.redMaterial).
		Line(sprec.ZeroVec3(), sprec.NewVec3(0.2, 0.0, 0.0))
	meshBuilder.Wireframe(d.greenMaterial).
		Line(sprec.ZeroVec3(), sprec.NewVec3(0.0, 0.0, 0.2))
	meshBuilder.Wireframe(d.blueMaterial).
		Line(sprec.ZeroVec3(), sprec.NewVec3(0.0, 0.2, 0.0))

	d.nodeMeshGeometry = d.gfxEngine.CreateMeshGeometry(meshBuilder.BuildGeometryInfo())
	d.nodeMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildMeshDefinitionInfo(d.nodeMeshGeometry))
}

func (d *CommonData) deleteNodeMesh() {
	defer d.nodeMeshGeometry.Delete()
	defer d.nodeMeshDef.Delete()
}

func (d *CommonData) createCameraMesh() {
	meshBuilder := graphics.NewShapeBuilder()

	solids := meshBuilder.Solid(d.yellowMaterial)
	solids.Cuboid(
		sprec.ZeroVec3(),
		sprec.IdentityQuat(),
		sprec.NewVec3(0.2, 0.3, -0.5),
	)
	solids.Cylinder(
		sprec.NewVec3(0.0, 0.25, -0.1),
		sprec.RotationQuat(sprec.Degrees(90), sprec.BasisZVec3()),
		0.15, 0.1, 20,
	)
	solids.Cylinder(
		sprec.NewVec3(0.0, 0.23, 0.13),
		sprec.RotationQuat(sprec.Degrees(90), sprec.BasisZVec3()),
		0.1, 0.1, 20,
	)
	solids.Cone(
		sprec.NewVec3(0.0, 0.0, -0.3),
		sprec.RotationQuat(sprec.Degrees(90), sprec.BasisXVec3()),
		0.2, 0.3, 20,
	)

	d.cameraMeshGeometry = d.gfxEngine.CreateMeshGeometry(meshBuilder.BuildGeometryInfo())
	d.cameraMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildMeshDefinitionInfo(d.cameraMeshGeometry))
}

func (d *CommonData) deleteCameraMesh() {
	defer d.cameraMeshGeometry.Delete()
	defer d.cameraMeshDef.Delete()
}

func (d *CommonData) createAmbientLightMesh() {
	const (
		coneRadius   = 0.05
		coneHeight   = 0.1
		coneSegments = 12
	)

	meshBuilder := graphics.NewShapeBuilder()
	solids := meshBuilder.Solid(d.yellowMaterial)
	solids.Cone(
		sprec.NewVec3(0.0, -0.2, 0.0),
		sprec.IdentityQuat(),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(0.0, 0.2, 0.0),
		sprec.RotationQuat(sprec.Degrees(180), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(0.2, 0.0, 0.0),
		sprec.RotationQuat(sprec.Degrees(90), sprec.BasisZVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(-0.2, 0.0, 0.0),
		sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisZVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(0.0, 0.0, -0.2),
		sprec.RotationQuat(sprec.Degrees(90), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(0.0, 0.0, 0.2),
		sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)

	d.ambientLightMeshGeometry = d.gfxEngine.CreateMeshGeometry(meshBuilder.BuildGeometryInfo())
	d.ambientLightMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildMeshDefinitionInfo(d.ambientLightMeshGeometry))
}

func (d *CommonData) deleteAmbientLightMesh() {
	defer d.ambientLightMeshGeometry.Delete()
	defer d.ambientLightMeshDef.Delete()
}

func (d *CommonData) createPointLightMesh() {
	const (
		sphereRadius   = 0.1
		sphereSegments = 8
		coneRadius     = 0.05
		coneHeight     = 0.1
		coneSegments   = 12
	)

	meshBuilder := graphics.NewShapeBuilder()
	solids := meshBuilder.Solid(d.yellowMaterial)
	solids.Sphere(
		sprec.ZeroVec3(),
		sphereRadius,
		sphereSegments,
	)
	solids = meshBuilder.Solid(d.darkYellowMaterial)
	solids.Cone(
		sprec.NewVec3(0.0, 0.2, 0.0),
		sprec.IdentityQuat(),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(0.0, -0.2, 0.0),
		sprec.RotationQuat(sprec.Degrees(180), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(-0.2, 0.0, 0.0),
		sprec.RotationQuat(sprec.Degrees(90), sprec.BasisZVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(0.2, 0.0, 0.0),
		sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisZVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(0.0, 0.0, 0.2),
		sprec.RotationQuat(sprec.Degrees(90), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(0.0, 0.0, -0.2),
		sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)

	d.pointLightMeshGeometry = d.gfxEngine.CreateMeshGeometry(meshBuilder.BuildGeometryInfo())
	d.pointLightMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildMeshDefinitionInfo(d.pointLightMeshGeometry))
}

func (d *CommonData) deletePointLightMesh() {
	defer d.pointLightMeshGeometry.Delete()
	defer d.pointLightMeshDef.Delete()
}

func (d *CommonData) createSpotLightMesh() {
	const (
		sphereRadius   = 0.1
		sphereSegments = 8
		coneRadius     = 0.05
		coneHeight     = 0.1
		coneSegments   = 12
	)

	meshBuilder := graphics.NewShapeBuilder()
	solids := meshBuilder.Solid(d.yellowMaterial)
	solids.Sphere(
		sprec.ZeroVec3(),
		sphereRadius,
		sphereSegments,
	)
	solids = meshBuilder.Solid(d.darkYellowMaterial)
	solids.Cone(
		sprec.NewVec3(0.0, 0.0, -0.2),
		sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)

	// TODO: Split these lines as separate meshes
	lines := meshBuilder.Wireframe(d.darkYellowMaterial)
	lines.Circle(
		sprec.NewVec3(0.0, 0.0, -0.4),
		sprec.IdentityQuat(),
		0.1, 20,
	)
	lines.Circle(
		sprec.NewVec3(0.0, 0.0, -0.4),
		sprec.IdentityQuat(),
		0.05, 20,
	)

	d.spotLightMeshGeometry = d.gfxEngine.CreateMeshGeometry(meshBuilder.BuildGeometryInfo())
	d.spotLightMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildMeshDefinitionInfo(d.spotLightMeshGeometry))
}

func (d *CommonData) deleteSpotLightMesh() {
	defer d.spotLightMeshGeometry.Delete()
	defer d.spotLightMeshDef.Delete()
}

func (d *CommonData) createDirectionalLightMesh() {
	const (
		coneRadius   = 0.05
		coneHeight   = 0.1
		coneSegments = 12
	)

	meshBuilder := graphics.NewShapeBuilder()
	solids := meshBuilder.Solid(d.yellowMaterial)
	solids.Cuboid(
		sprec.ZeroVec3(),
		sprec.IdentityQuat(),
		sprec.NewVec3(0.3, 0.3, 0.02),
	)
	solids = meshBuilder.Solid(d.darkYellowMaterial)
	solids.Cone(
		sprec.NewVec3(-0.1, 0.1, -0.15),
		sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(0.1, 0.1, -0.15),
		sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(-0.1, -0.1, -0.15),
		sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)
	solids.Cone(
		sprec.NewVec3(0.1, -0.1, -0.15),
		sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisXVec3()),
		coneRadius, coneHeight, coneSegments,
	)

	d.directionalLightMeshGeometry = d.gfxEngine.CreateMeshGeometry(meshBuilder.BuildGeometryInfo())
	d.directionalLightMeshDef = d.gfxEngine.CreateMeshDefinition(meshBuilder.BuildMeshDefinitionInfo(d.directionalLightMeshGeometry))
}

func (d *CommonData) deleteDirectionalLightMesh() {
	defer d.directionalLightMeshGeometry.Delete()
	defer d.directionalLightMeshDef.Delete()
}

type ColorUniform struct {
	Color sprec.Vec4
}

func (u ColorUniform) Std140Plot(plotter *blob.Plotter) {
	plotter.PlotSPVec4(u.Color)
}

func (u ColorUniform) Std140Size() int {
	return 4 * render.SizeF32
}
