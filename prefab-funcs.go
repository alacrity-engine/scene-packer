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

		sceneBucketName := sceneBucketNameFromID(sceneID)
		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists(
				[]byte(sceneBucketName))

			if err != nil {
				return err
			}

			gmobsLHData := buck.Get([]byte(gmobsKey))
			var lh *codec.ListHeader

			if gmobsLHData != nil {
				lh, err = codec.ListHeaderFromBytes(gmobsLHData)

				if err != nil {
					return err
				}
			} else {
				lh = &codec.ListHeader{}
			}

			for i, data := range datas {
				key := listItemFromID(sceneBucketName,
					gmobsKey, int(lh.Count)+i)
				err = buck.Put([]byte(key), data)

				if err != nil {
					return err
				}
			}

			lh.Count += int32(len(datas))
			lhData, err := lh.ToBytes()

			if err != nil {
				return err
			}

			err = buck.Put([]byte(gmobsKey), lhData)

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
