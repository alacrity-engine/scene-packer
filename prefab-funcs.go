package main

import (
	"fmt"

	codec "github.com/alacrity-engine/resource-codec"
	bolt "go.etcd.io/bbolt"
)

const (
	gmobsKey   = "gmobs"
	prefabsKey = "prefabs"
)

func createGameObjects(
	resourceFile *bolt.DB, handleError func(err error),
) func(prefabs []*PrefabMeta) {
	return func(prefabs []*PrefabMeta) {
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
			buck := tx.Bucket([]byte(sceneBucketName))

			if buck == nil {
				return fmt.Errorf(
					"bucket '%s' not found", sceneBucketName)
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
			buck := tx.Bucket([]byte(sceneBucketName))

			if buck == nil {
				return fmt.Errorf(
					"bucket '%s' not found", sceneBucketName)
			}

			err := buck.Put([]byte(prefabsKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}
