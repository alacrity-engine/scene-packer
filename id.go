package main

func sceneBucketNameFromID(sceneID string) string {
	return sceneBucketPrefix + "." + sceneID
}
