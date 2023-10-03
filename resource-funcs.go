package main

import (
	codec "github.com/alacrity-engine/resource-codec"
	bolt "go.etcd.io/bbolt"
)

const (
	preloadAnimationsKey     = "preload_animations"
	preloadAudiosKey         = "preload_audios"
	preloadTexturesKey       = "preload_textures"
	preloadFontsKey          = "preload_fonts"
	preloadPicturesKey       = "preload_pictures"
	preloadShadersKey        = "preload_shaders"
	preloadShaderProgramsKey = "preload_shader_programs"
)

func preloadAnimations(resourceFile *bolt.DB, handleError func(err error)) func(names []string) {
	return func(names []string) {
		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists([]byte(sceneBucketName))

			if err != nil {
				return err
			}

			data, err := codec.EncodeTag(names)

			if err != nil {
				return err
			}

			err = buck.Put([]byte(preloadAnimationsKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}

func preloadAudios(resourceFile *bolt.DB, handleError func(err error)) func(names []string) {
	return func(names []string) {
		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists([]byte(sceneBucketName))

			if err != nil {
				return err
			}

			data, err := codec.EncodeTag(names)

			if err != nil {
				return err
			}

			err = buck.Put([]byte(preloadAudiosKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}

func preloadTextures(resourceFile *bolt.DB, handleError func(err error)) func(names []string) {
	return func(names []string) {
		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists([]byte(sceneBucketName))

			if err != nil {
				return err
			}

			data, err := codec.EncodeTag(names)

			if err != nil {
				return err
			}

			err = buck.Put([]byte(preloadTexturesKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}

func preloadFonts(resourceFile *bolt.DB, handleError func(err error)) func(names []string) {
	return func(names []string) {
		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists([]byte(sceneBucketName))

			if err != nil {
				return err
			}

			data, err := codec.EncodeTag(names)

			if err != nil {
				return err
			}

			err = buck.Put([]byte(preloadFontsKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}

func preloadPictures(resourceFile *bolt.DB, handleError func(err error)) func(names []string) {
	return func(names []string) {
		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists([]byte(sceneBucketName))

			if err != nil {
				return err
			}

			data, err := codec.EncodeTag(names)

			if err != nil {
				return err
			}

			err = buck.Put([]byte(preloadPicturesKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}

func preloadShaders(resourceFile *bolt.DB, handleError func(err error)) func(names []string) {
	return func(names []string) {
		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists([]byte(sceneBucketName))

			if err != nil {
				return err
			}

			data, err := codec.EncodeTag(names)

			if err != nil {
				return err
			}

			err = buck.Put([]byte(preloadShadersKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}

func preloadShaderPrograms(resourceFile *bolt.DB, handleError func(err error)) func(names []string) {
	return func(names []string) {
		err := resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists([]byte(sceneBucketName))

			if err != nil {
				return err
			}

			data, err := codec.EncodeTag(names)

			if err != nil {
				return err
			}

			err = buck.Put([]byte(preloadShaderProgramsKey), data)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}
