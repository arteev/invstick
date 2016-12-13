package helpers

import (
	"math"
	"os"
	"path/filepath"
)

func MkSlice(args ...interface{}) []interface{} {
	return args
}

func Mul(a, b int) int {
	return a * b
}

func Add(a, b int) int {
	return a + b
}

func Calcpages(onpage, count int) int {
	if onpage == 0 {
		return 0
	}
	i := count / onpage
	if math.Mod(float64(count), float64(onpage)) > 0 {
		i++
	}
	return i
}

func MkSliceRange(from, count int) (result []int) {
	for i := from; i < from+count; i++ {
		result = append(result, i)
	}
	return
}

func GetStickerTemplates(path, ext string) []string {
	files := make([]string, 0)
	filepath.Walk(filepath.Join(path), func(path string, info os.FileInfo, err error) error {
		if (!info.IsDir()) && filepath.Ext(info.Name()) == ext {
			files = append(files, info.Name())
		}
		return nil
	})
	return files
}
