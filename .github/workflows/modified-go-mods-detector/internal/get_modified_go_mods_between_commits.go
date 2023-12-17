package internal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetModifiedGoModsBetweenCommits(gitRefBefore, gitRefAfter string) ([]GoMod, error) {
	modifiedGoMods, err := getModifiedGoMods(gitRefBefore, gitRefAfter)
	if err != nil {
		return nil, err
	}

	indirectlyModifiedGoMods, err := getIndirectlyModifiedGoMods(modifiedGoMods)
	if err != nil {
		return nil, err
	}

	mgm := append(fromMap(modifiedGoMods), fromMap(indirectlyModifiedGoMods)...)
	result := deduplicateGoMods(mgm)

	return result, nil
}

func getModifiedGoMods(gitRefBefore, gitRefAfter string) (map[string]GoMod, error) {
	diffedFiles, err := gitDiff(context.Background(), gitRefBefore, gitRefAfter)
	if err != nil {
		return nil, err
	}

	gmMap := make(map[string]GoMod)
	for _, file := range diffedFiles {
		projectName, projectDir, isLib := parseFilePath(file)

		if projectName == "" || projectDir == "" {
			continue
		}

		isAGoMod, err := isGoMod(projectDir)
		if err != nil {
			return nil, err
		}

		if !isAGoMod {
			continue
		}

		goModName, err := getGoModName(projectDir)
		if err != nil {
			return nil, err
		}

		gmMap[goModName] = GoMod{
			ProjectName: projectName,
			ProjectDir:  projectDir,
			GoModName:   goModName,
			IsLib:       isLib,
		}
	}

	return gmMap, nil
}

func parseFilePath(file string) (string, string, bool) {
	parts := strings.Split(file, "/")

	if len(parts) == 0 {
		return "", "", false
	}

	if strings.HasPrefix(file, "lib/") {
		if len(parts) < 2 {
			return "", "", false
		}
		name := "lib/" + parts[1]
		return name, name, true
	}

	return parts[0], parts[0], false
}

func getIndirectlyModifiedGoMods(diffedGoModules map[string]GoMod) (map[string]GoMod, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	defer func(dir string) {
		_ = os.Chdir(dir)
	}(wd)

	appGoMods, err := findAppGoMods()
	if err != nil {
		return nil, err
	}
	gmMap := make(map[string]GoMod)
	for _, module := range appGoMods {
		containsModifiedLib, err := goModContainsAModifiedLibAsADependency(module, diffedGoModules)
		if err != nil {
			return nil, err
		}
		if !containsModifiedLib {
			continue
		}

		gmMap[module.GoModName] = module
	}

	return gmMap, nil
}

func goModContainsAModifiedLibAsADependency(module GoMod, modifiedGoMods map[string]GoMod) (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	defer func(dir string) {
		_ = os.Chdir(dir)
	}(wd)

	err = os.Chdir(module.ProjectDir)
	if err != nil {
		return false, err
	}

	deps, err := getGoModDeps(context.Background())
	if err != nil {
		return false, err
	}
	for _, dep := range deps {
		for _, dm := range modifiedGoMods {
			if strings.HasPrefix(dep, dm.GoModName) && dm.IsLib {
				return true, nil
			}
		}
	}

	return false, nil
}

func getGoModDeps(ctx context.Context) ([]string, error) {
	cmdArgs := []string{"list", "-f", "'{{.Deps}}'", "./..."}

	cmd := exec.CommandContext(ctx, "go", cmdArgs...)

	var errOut bytes.Buffer
	cmd.Stderr = &errOut

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get go list: %v", err)
	}

	errs := errOut.String()
	if errs != "" {
		return nil, fmt.Errorf("failed to get go list: %v", errs)
	}

	s := out.String()

	deps := strings.Split(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					s, "'[", "",
				), "]'", "",
			), "\n", " ",
		), " ",
	)

	// deduplicate deps
	dm := make(map[string]struct{}, len(deps))
	for _, dep := range deps {
		dm[dep] = struct{}{}
	}
	result := make([]string, 0, len(dm))
	for dep := range dm {
		result = append(result, dep)
	}
	return result, nil
}

type GoListOutput struct {
	Deps []string `json:"Deps"`
}

func deduplicateGoMods(mgm []GoMod) []GoMod {
	mgmMap := make(map[string]GoMod)
	for _, mod := range mgm {
		mgmMap[mod.GoModName] = mod
	}
	result := fromMap(mgmMap)
	return result
}

func fromMap(mgmMap map[string]GoMod) []GoMod {
	result := make([]GoMod, 0, len(mgmMap))
	for _, mod := range mgmMap {
		result = append(result, mod)
	}
	return result
}
