package fileop

import "path"

func PathJoin(root string, followUp ...string) string {
	pathSlice := []string{"/"}
	pathSlice = append(pathSlice, followUp...)

	followDir := path.Join(pathSlice...)

	return path.Join(root, followDir)
}
