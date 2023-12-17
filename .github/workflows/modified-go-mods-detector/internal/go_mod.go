package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type GoMod struct {
	ProjectName string `json:"projectName"`
	ProjectDir  string `json:"projectDir"`
	GoModName   string `json:"goModName"`
	IsLib       bool   `json:"isLib"`
}

func findAppGoMods() ([]GoMod, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	defer func(dir string) {
		_ = os.Chdir(dir)
	}(wd)

	appGoMods, err := findGoModsInChildrenDirectories(wd)
	if err != nil {
		return nil, err
	}

	return appGoMods, nil
}

func findGoModsInChildrenDirectories(dir string) ([]GoMod, error) {
	childrenDirectories, err := exploreChildrenDirectories(dir)
	if err != nil {
		return nil, err
	}

	var goMods []GoMod
	for _, childDir := range childrenDirectories {
		isAGoMod, err := isGoMod(childDir)
		if err != nil {
			return nil, err
		}

		if !isAGoMod {
			continue
		}

		goModName, err := getGoModName(childDir)
		if err != nil {
			return nil, err
		}

		goMods = append(goMods, GoMod{
			ProjectName: filepath.Base(childDir),
			ProjectDir:  filepath.Base(childDir),
			GoModName:   goModName,
			IsLib:       false,
		})
	}
	return goMods, nil
}

func isGoMod(dir string) (bool, error) {
	files, err := getDirFiles(dir)
	if err != nil {
		return false, err
	}
	if files == nil {
		return false, nil
	}

	isGoMod := false
	for _, file := range files {
		if file.Name() == "go.mod" {
			isGoMod = true
			break
		}
	}

	return isGoMod, nil
}

func getDirFiles(dir string) ([]os.FileInfo, error) {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		log.Info(fmt.Errorf("error while getting file info for %s: %w", dir, err))
		return nil, nil
	}

	if !fileInfo.IsDir() {
		return nil, nil
	}

	dirF, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	files, err := dirF.Readdir(-1)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func getGoModName(dir string) (string, error) {
	gomod, err := os.Open(filepath.Join(dir, "go.mod"))
	if err != nil {
		return "", err
	}
	defer gomod.Close()

	scanner := bufio.NewScanner(gomod)
	if scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module") {
			fields := strings.Fields(line)
			if len(fields) == 2 {
				moduleName := fields[1]
				return moduleName, nil
			}
		}
	}

	return "", fmt.Errorf("invalid go.mod file")
}

func exploreChildrenDirectories(dirPath string) ([]string, error) {
	var subDirs []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != dirPath {
			subDirs = append(subDirs, path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return subDirs, nil
}
