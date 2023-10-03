package main

import (
	"github.com/alacrity-engine/core/definitions"
	codec "github.com/alacrity-engine/resource-codec"
	bolt "go.etcd.io/bbolt"
)

const (
	canvasesKey = "canvases"
	batchesKey  = "batches"
)

// list makes Lua perceive the data
// as a slice instead of a map.
func list(data []interface{}) []interface{} {
	return data
}

func createCanvases(
	resourceFile *bolt.DB, handleError func(err error),
) func(canvasDefinitions []definitions.CanvasDefinition) {
	return func(canvasDefinitions []definitions.CanvasDefinition) {
		datas := make([][]byte, len(canvasDefinitions))

		for _, canvasDefinition := range canvasDefinitions {
			canvasData := &codec.CanvasData{
				Name:       canvasDefinition.Name,
				DrawZ:      canvasDefinition.DrawZ,
				Projection: canvasDefinition.Projection,
			}

			data, err := canvasData.ToBytes()
			handleError(err)
			datas = append(datas, data)
		}

		data, err := codec.SerializeBlobs(datas)
		handleError(err)

		err = resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists(
				[]byte(sceneBucketName))

			if err != nil {
				return err
			}

			err = buck.Put([]byte(canvasesKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}

func createBatches(
	resourceFile *bolt.DB, handleError func(err error),
) func(batchDefinitions []definitions.BatchDefinition) {
	return func(batchDefinitions []definitions.BatchDefinition) {
		datas := make([][]byte, len(batchDefinitions))

		for _, batchDefinition := range batchDefinitions {
			batchData := codec.BatchData{
				Name:      batchDefinition.Name,
				CanvasID:  batchDefinition.CanvasID,
				TextureID: batchDefinition.TextureID,
				ZMin:      batchDefinition.ZMin,
				ZMax:      batchDefinition.ZMax,
			}

			data, err := batchData.ToBytes()
			handleError(err)
			datas = append(datas, data)
		}

		data, err := codec.SerializeBlobs(datas)
		handleError(err)

		err = resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists(
				[]byte(sceneBucketName))

			if err != nil {
				return err
			}

			err = buck.Put([]byte(batchesKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}
