package main

import (
	codec "github.com/alacrity-engine/resource-codec"
	bolt "go.etcd.io/bbolt"
)

const (
	animationsKey     = "animations"
	audiosKey         = "audios"
	texturesKey       = "textures"
	fontsKey          = "fonts"
	picturesKey       = "pictures"
	shadersKey        = "shaders"
	shaderProgramsKey = "shader_programs"
)

func preloadItemKeyFromResourceType(resourceType string) string {
	return "preload" + "." + resourceType
}

func preload(
	resourceFile *bolt.DB,
	handleError func(err error),
) func(sceneID, resourceType string, ids []string) {
	return func(sceneID, resourceType string, ids []string) {
		datas := make([][]byte, 0, len(ids))

		for _, id := range ids {
			datas = append(datas, []byte(id))
		}

		sceneBucketName := sceneBucketNameFromID(sceneID)
		preloadItemKey := preloadItemKeyFromResourceType(resourceType)
		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists(
				[]byte(sceneBucketName))

			if err != nil {
				return err
			}

			preloadLHData := buck.Get([]byte(preloadItemKey))
			var lh *codec.ListHeader

			if preloadLHData != nil {
				lh, err = codec.ListHeaderFromBytes(preloadLHData)

				if err != nil {
					return err
				}
			} else {
				lh = &codec.ListHeader{}
			}

			for i, data := range datas {
				key := listItemFromID(sceneBucketName,
					preloadItemKey, int(lh.Count)+i)
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

			err = buck.Put([]byte(preloadItemKey), lhData)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}
