package fileop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathJoin(t *testing.T) {
	t.Run("should join path", func(t *testing.T) {
		dir := PathJoin("root", "/a", "/b", "c")
		assert.Equal(t, "root/a/b/c", dir)
	})

	t.Run("should join path when have ..", func(t *testing.T) {
		dir := PathJoin("root", "a", "/b", "../c")
		assert.Equal(t, "root/a/c", dir)
	})

	t.Run("should not go before root", func(t *testing.T) {
		dir := PathJoin("root", "../../a")
		assert.Equal(t, "root/a", dir)
	})

	t.Run("should not go before root and use last element file name as the final file name", func(t *testing.T) {
		dir := PathJoin("root", "../../a", "../b")
		assert.Equal(t, "root/b", dir)
	})
}
