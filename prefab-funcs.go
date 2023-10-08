package main

import (
	"fmt"

	codec "github.com/alacrity-engine/resource-codec"
	bolt "go.etcd.io/bbolt"
)

const (
	gmobsKey          = "gmobs"
	prefabsBucketName = "prefabs"
)

// TODO: make createGameObjects method
// callable many many times so users
// could fill a scene partially.

func createGameObjects(
	resourceFile *bolt.DB, handleError func(err error),
) func(sceneID string, prefabs []*PrefabMeta) {
	return func(sceneID string, prefabs []*PrefabMeta) {
		datas := make([][]byte, 0, len(prefabs))

		for _, prefab := range prefabs {
			prefabData := PrefabMetaToData(prefab)
			data, err := prefabData.ToBytes()
			handleError(err)

			datas = append(datas, data)
		}

		data, err := codec.SerializeBlobs(datas)
		handleError(err)

		err = resourceFile.Update(func(tx *bolt.Tx) error {
			buck := tx.Bucket([]byte(sceneBucketNameFromID(sceneID)))

			if buck == nil {
				return fmt.Errorf(
					"bucket '%s' not found", sceneBucketNameFromID(sceneID))
			}

			err := buck.Put([]byte(gmobsKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}

func storePrefabs(
	resourceFile *bolt.DB, handleError func(err error),
) func(prefabs []*PrefabMeta) {
	return func(prefabs []*PrefabMeta) {
		datas := make([]*codec.PrefabData, 0, len(prefabs))

		for _, prefab := range prefabs {
			prefabData := PrefabMetaToData(prefab)
			datas = append(datas, prefabData)
		}

		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists([]byte(prefabsBucketName))

			if err != nil {
				return err
			}

			if buck == nil {
				return fmt.Errorf(
					"bucket '%s' not found", prefabsBucketName)
			}

			for _, prefabData := range datas {
				data, err := prefabData.ToBytes()

				if err != nil {
					return err
				}

				err = buck.Put([]byte(prefabData.Name), data)

				if err != nil {
					return err
				}
			}

			return nil
		})
		handleError(err)
	}
}
