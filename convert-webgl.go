package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const SRCDIR = `WebGL/`
const DSTDIR = `WebGLLocal/`

func main() {
	err := CleanDstDir()
	if err != nil {
		log.Fatalf(`%v`, err)
		return
	}
	// Traverse all files in source directory recursively
	err = filepath.Walk(SRCDIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Path %s access failed: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Generate pathes
		srcFilePath := strings.ReplaceAll(path, `\`, `/`)
		assetPath := strings.TrimPrefix(srcFilePath, SRCDIR)
		dstFilePath := DSTDIR + assetPath

		// Make sure the destination path exists
		// Use 0755 so created directories are traversable
		err = os.MkdirAll(filepath.Dir(dstFilePath), 0755)
		if err != nil {
			return err
		}
		if !isFileNeedConvert(assetPath) {
			return CopyFile(srcFilePath)
		} else {
			return ConvertFile(srcFilePath)
		}
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", SRCDIR, err)
	}
}

// Clean destination directory
func CleanDstDir() error {
	_, err := os.Stat(DSTDIR)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	err = os.RemoveAll(DSTDIR)
	if err != nil {
		return err
	}
	return nil
}

// Check if specific file needs convert
func isFileNeedConvert(assetPath string) bool {
	if !strings.HasPrefix(assetPath, `Build/`) &&
		!strings.HasPrefix(assetPath, `StreamingAssets/`) {
		return false
	}
	if strings.HasSuffix(assetPath, `.js`) {
		return false
	}
	return true
}

// Copy file from source path to the corresponding path in Destination Directory
func CopyFile(srcFilePath string) error {
	assetPath := strings.TrimPrefix(srcFilePath, SRCDIR)
	dstFilePath := DSTDIR + assetPath

	// Open the source file
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy file content
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	if assetPath == `TemplateData/fetch.js` {
		_, err := dstFile.WriteString("\nlocal = true;")
		if err != nil {
			return err
		}
	}

	// Make sure file content sync
	return dstFile.Sync()
}

// Read source file content, convert to base64 encode and write into corresponding .js file
func ConvertFile(srcFilePath string) error {
	assetPath := strings.TrimPrefix(srcFilePath, SRCDIR)
	dstFilePath := DSTDIR + assetPath

	// Read source file content
	srcFileData, err := os.ReadFile(srcFilePath)
	if err != nil {
		return err
	}

	// Create the destination file
	dstFile, err := os.Create(dstFilePath + `.js`)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Convert to base64 encode and write into corresponding .js file
	varName := filepath.Base(srcFilePath)
	varName = strings.ReplaceAll(varName, `.`, `_`)
	varName = strings.ReplaceAll(varName, `-`, `_`)
	covertStr := base64.StdEncoding.EncodeToString(srcFileData)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(`const %s_bin = window.atob('%s');`, varName, covertStr))
	builder.WriteString(fmt.Sprintf(`const %s_len = %s_bin.length;`, varName, varName))
	builder.WriteString(fmt.Sprintf(`const %s = new Uint8Array(%s_len);`, varName, varName))
	builder.WriteString(fmt.Sprintf(`for (let i = 0; i < %s_len; i++){%s[i] = %s_bin.charCodeAt(i);}`, varName, varName, varName))
	builder.WriteString(fmt.Sprintf(`window._uinty_asset_%s=%s;`, varName, varName))
	_, err = dstFile.WriteString(builder.String())
	if err != nil {
		return err
	}

	// Make sure file content sync
	return dstFile.Sync()
}
