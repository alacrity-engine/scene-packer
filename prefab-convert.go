package main

import (
	codec "github.com/alacrity-engine/resource-codec"
)

func PrefabMetaToData(prefabMeta *PrefabMeta) *codec.PrefabData {
	return &codec.PrefabData{
		Name:          prefabMeta.Name,
		TransformRoot: TransformMetaToData(prefabMeta.TransformRoot),
	}
}

func TransformMetaToData(trMeta *TransformMeta) *codec.TransformData {
	trData := &codec.TransformData{
		Position: trMeta.Position,
		Angle:    trMeta.Angle,
		Scale:    trMeta.Scale,
		Gmob:     GameObjectMetaToData(trMeta.Gmob),
	}

	if len(trData.Children) > 0 {
		children := make([]*codec.TransformData, 0, len(trMeta.Children))

		for _, child := range trMeta.Children {
			childData := TransformMetaToData(child)
			children = append(children, childData)
		}

		trData.Children = children
	}

	return trData
}

func GameObjectMetaToData(gmobMeta *GameObjectMeta) *codec.GameObjectData {
	gmobData := &codec.GameObjectData{
		Name:       gmobMeta.Name,
		ZUpdate:    gmobMeta.ZUpdate,
		Sprite:     SpriteMetaToData(gmobMeta.Sprite),
		Components: make([]*codec.ComponentData, 0, len(gmobMeta.Components)),
		Draw:       gmobMeta.Draw,
	}

	for _, comp := range gmobMeta.Components {
		compData := ComponentMetaToData(comp)
		gmobData.Components = append(gmobData.Components, compData)
	}

	return gmobData
}

func ComponentMetaToData(compMeta *ComponentMeta) *codec.ComponentData {
	fieldData := make(map[string]interface{}, len(compMeta.Data))

	for fieldName, fieldValue := range compMeta.Data {
		switch fieldVal := fieldValue.(type) {
		case GameObjectPointerMeta:
			fieldData[fieldName] = codec.GameObjectPointerData{
				Name: fieldVal.Name,
			}

		case ComponentPointerMeta:
			fieldData[fieldName] = codec.ComponentPointerData{
				GmobName: fieldVal.GmobName,
				CompType: fieldVal.CompType,
			}

		case ResourcePointerMeta:
			fieldData[fieldName] = codec.ResourcePointerData{
				ResourceType: fieldVal.ResourceType,
				ResourceID:   fieldVal.ResourceID,
			}

		default:
			fieldData[fieldName] = fieldValue
		}
	}

	return &codec.ComponentData{
		TypeName: compMeta.TypeName,
		Active:   compMeta.Active,
		Data:     fieldData,
	}
}

func SpriteMetaToData(spriteMeta *SpriteMeta) *codec.SpriteData {
	return &codec.SpriteData{
		ColorMask:       spriteMeta.ColorMask,
		TargetArea:      spriteMeta.TargetArea,
		ZDraw:           spriteMeta.ZDraw,
		VertexDrawMode:  uint32(spriteMeta.VertexDrawMode),
		TextureDrawMode: uint32(spriteMeta.TextureDrawMode),
		ColorDrawMode:   uint32(spriteMeta.ColorDrawMode),
		ShaderProgramID: spriteMeta.ShaderProgramID,
		TextureID:       spriteMeta.TextureID,
		CanvasID:        spriteMeta.CanvasID,
		BatchID:         spriteMeta.BatchID,
	}
}
