package utils

// MaxSize defines the max supported file size
var MaxSize int64

func init() {
	MaxSize = 1024 * 1024 * 1024
}
