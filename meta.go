package main

import (
	"github.com/alacrity-engine/core/math/geometry"
)

type PrefabMeta struct {
	Name          string
	TransformRoot *TransformMeta
}

type TransformMeta struct {
	Position geometry.Vec
	Angle    float64
	Scale    geometry.Vec
	Gmob     *GameObjectMeta
	Children []*TransformMeta
}

type GameObjectMeta struct {
	Name       string
	ZUpdate    float64
	Components []*ComponentMeta
	Sprite     *SpriteMeta
	Draw       bool
}

type ComponentMeta struct {
	Path     string
	TypeName string
	Active   bool
	Data     map[string]interface{}
}

type SpriteMeta struct {
	ColorMask       []float32
	TargetArea      geometry.Rect
	ZDraw           float32
	VertexDrawMode  uint32
	TextureDrawMode uint32
	ColorDrawMode   uint32
	ShaderProgramID string
	TextureID       string
	CanvasID        string
	BatchID         string
}

type GameObjectPointerMeta struct {
	Name string
}

type ComponentPointerMeta struct {
	GmobName string
	CompType string
}

type ResourcePointerMeta struct {
	ResourceType string
	ResourceID   string
}

type BatchMeta struct {
	Name      string
	CanvasID  string
	TextureID string
	ZMin      float32
	ZMax      float32
}

type CanvasMeta struct {
	Name  string
	DrawZ int
}
