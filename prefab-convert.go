package main

import (
	"github.com/alacrity-engine/core/definitions"
	codec "github.com/alacrity-engine/resource-codec"
)

func PrefabDefinitionToData(prefabDef *definitions.Prefab) *codec.PrefabData {
	return &codec.PrefabData{
		Name:          prefabDef.Name,
		TransformRoot: TransformDefinitionToData(prefabDef.TransformRoot),
	}
}

func TransformDefinitionToData(trDef *definitions.TransformDefinition) *codec.TransformData {
	trData := &codec.TransformData{
		Position: trDef.Position,
		Angle:    trDef.Angle,
		Scale:    trDef.Scale,
		Gmob:     GameObjectDefinitionToData(trDef.Gmob),
	}

	if len(trData.Children) > 0 {
		children := make([]*codec.TransformData, 0, len(trDef.Children))

		for _, child := range trDef.Children {
			childData := TransformDefinitionToData(child)
			children = append(children, childData)
		}

		trData.Children = children
	}

	return trData
}

func GameObjectDefinitionToData(gmobDef *definitions.GameObjectDefinition) *codec.GameObjectData {
	gmobData := &codec.GameObjectData{
		Name:       gmobDef.Name,
		ZUpdate:    gmobDef.ZUpdate,
		Sprite:     SpriteDefinitionToData(gmobDef.Sprite),
		Components: make([]*codec.ComponentData, 0, len(gmobDef.Components)),
		Draw:       gmobDef.Draw,
	}

	for _, comp := range gmobDef.Components {
		compData := ComponentDefinitionToData(comp)
		gmobData.Components = append(gmobData.Components, compData)
	}

	return gmobData
}

func ComponentDefinitionToData(compDef *definitions.ComponentDefinition) *codec.ComponentData {
	fieldData := make(map[string]interface{}, len(compDef.Data))

	for fieldName, fieldValue := range compDef.Data {
		switch fieldVal := fieldValue.(type) {
		case definitions.GameObjectPointer:
			fieldData[fieldName] = codec.GameObjectPointerData{
				Name: fieldVal.Name,
			}

		case definitions.ComponentPointer:
			fieldData[fieldName] = codec.ComponentPointerData{
				GmobName: fieldVal.GmobName,
				CompType: fieldVal.CompType,
			}

		case definitions.ResourcePointer:
			fieldData[fieldName] = codec.ResourcePointerData{
				ResourceType: fieldVal.ResourceType,
				ResourceID:   fieldVal.ResourceID,
			}

		default:
			fieldData[fieldName] = fieldValue
		}
	}

	return &codec.ComponentData{
		TypeName: compDef.TypeName,
		Active:   compDef.Active,
		Data:     fieldData,
	}
}

func SpriteDefinitionToData(spriteDef *definitions.SpriteDefinition) *codec.SpriteData {
	colorMask := spriteDef.ColorMask.Data()
	colorMaskData := colorMask[:]

	return &codec.SpriteData{
		ColorMask:       colorMaskData,
		TargetArea:      spriteDef.TargetArea,
		VertexDrawMode:  uint32(spriteDef.VertexDrawMode),
		TextureDrawMode: uint32(spriteDef.TextureDrawMode),
		ColorDrawMode:   uint32(spriteDef.ColorDrawMode),
		ShaderProgramID: spriteDef.ShaderProgramID,
		TextureID:       spriteDef.TextureID,
		CanvasID:        spriteDef.CanvasID,
		BatchID:         spriteDef.BatchID,
	}
}
