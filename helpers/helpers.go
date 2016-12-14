package helpers

import (
	"math"
	"os"
	"path/filepath"
)

//MkSlice returns array from arguments. Helpers for text/template
func MkSlice(args ...interface{}) []interface{} {
	return args
}

//Mul returns multiplication a*b . Helpers for text/template
func Mul(a, b int) int {
	return a * b
}

//Add returns addition a+b. Helpers for text/template
func Add(a, b int) int {
	return a + b
}

//Calcpages - calculates the number of pages.
//count-amount of elements.onpage- number of elements on the page. Helpers for text/template
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

//MkSliceRange returns an array of integers "from" to "from+count" . Helpers for text/template
func MkSliceRange(from, count int) (result []int) {
	for i := from; i < from+count; i++ {
		result = append(result, i)
	}
	return
}

//GetStickerTemplates returns files templates into "path" with extention "ext"
func GetStickerTemplates(path, ext string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(filepath.Join(path), func(path string, info os.FileInfo, err error) error {
		if (!info.IsDir()) && filepath.Ext(info.Name()) == ext {
			files = append(files, info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err

	}
	return files, nil
}
