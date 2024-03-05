package favicon

import (
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/pirsch-analytics/faser/pkg/server"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var sizes = []int{
	16,
	32,
	64,
	96,
	128,
	196,
}

func scale(hostname, filename string, size int) string {
	size = getValidSize(size)
	dir := server.Config().Cache.Dir
	srcPath := filepath.Join(dir, hostname, filename)
	newFilename := getFilenameForSize(filename, size)
	targetPath := filepath.Join(dir, hostname, newFilename)
	cmd := exec.Command("magick",
		"convert",
		srcPath,
		"-resize",
		fmt.Sprintf("%dx%d>", size, size),
		targetPath)
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		logbuch.Error("Error resizing favicon", logbuch.Fields{
			"err":    err,
			"src":    srcPath,
			"target": targetPath,
			"size":   size,
			"stdout": string(stdout),
		})
		return ""
	}

	return newFilename
}

func getFilenameForSize(filename string, size int) string {
	ext := strings.ToLower(path.Ext(filename))
	filenameWithoutExt := filename[:len(filename)-len(ext)]
	return filenameWithoutExt + "-" + strconv.Itoa(size) + ext
}

func scalableType(ext string) bool {
	return ext == ".ico" || ext == ".png" || ext == ".gif" || ext == ".jpg" || ext == ".jpeg"
}

func getValidSize(size int) int {
	for _, s := range sizes {
		if size <= s {
			size = s
			break
		}
	}

	if size > sizes[len(sizes)-1] {
		size = sizes[len(sizes)-1]
	}

	return size
}
