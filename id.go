package main

import "strconv"

const (
	sceneBucketPrefix = "__scene"
	listPrefix        = "__list"
)

func sceneBucketNameFromID(sceneID string) string {
	return sceneBucketPrefix + "." + sceneID
}

func listItemFromID(bucketName, itemID string, index int) string {
	return listPrefix + "." + bucketName + "." + itemID + "." + strconv.Itoa(index)
}
