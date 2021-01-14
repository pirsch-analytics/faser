package favicon

import (
	"fmt"
	"github.com/emvi/logbuch"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func selectFilenameForSize(filename string, size int) (string, int) {
	if size <= 0 {
		return filename, 0
	}

	ext := strings.ToLower(path.Ext(filename))

	if !scalableType(ext) {
		return filename, 0
	}

	size = getValidSize(size)
	return getFilenameForSize(filename, size), size
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

func scale(hostname, filename string, size int) error {
	srcPath := filepath.Join(filesDir, hostname, filename)
	targetPath := filepath.Join(filesDir, hostname, getFilenameForSize(filename, size))
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
	}

	return nil
}
