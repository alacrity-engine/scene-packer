package main

import (
	codec "github.com/alacrity-engine/resource-codec"
	bolt "go.etcd.io/bbolt"
)

const (
	canvasesKey        = "canvases"
	batchesKey         = "batches"
	canvasesBucketName = "canvases"
)

// list makes Lua perceive the data
// as a slice instead of a map.
func list(data []interface{}) []interface{} {
	return data
}

func colorMask(data [16]float32) []float32 {
	return data[:]
}

func colorRepeat4(data [4]float32) []float32 {
	mask := make([]float32, 16)

	for i := 0; i < 4; i++ {
		mask[i] = data[i]
		mask[i+4] = data[i]
		mask[i+8] = data[i]
		mask[i+12] = data[i]
	}

	return mask
}

func createCanvases(
	resourceFile *bolt.DB, handleError func(err error),
) func(canvasMetas []CanvasMeta) {
	return func(canvasMetas []CanvasMeta) {
		datas := make([]*codec.CanvasData, len(canvasMetas))

		for _, canvasDefinition := range canvasMetas {
			canvasData := &codec.CanvasData{
				Name:  canvasDefinition.Name,
				DrawZ: canvasDefinition.DrawZ,
			}
			datas = append(datas, canvasData)
		}

		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists(
				[]byte(canvasesBucketName))

			if err != nil {
				return err
			}

			for _, canvasData := range datas {
				data, err := canvasData.ToBytes()

				if err != nil {
					return err
				}

				err = buck.Put([]byte(canvasData.Name), data)

				if err != nil {
					return err
				}
			}

			return nil
		})
		handleError(err)
	}
}

func createBatches(
	resourceFile *bolt.DB, handleError func(err error),
) func(sceneID string, batchMetas []BatchMeta) {
	return func(sceneID string, batchMetas []BatchMeta) {
		datas := make([][]byte, len(batchMetas))

		for _, batchDefinition := range batchMetas {
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
				[]byte(sceneBucketNameFromID(sceneID)))

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
