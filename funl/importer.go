package funl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func findSourceFile(path, importModName, fileExtensionName string) (pathAndFilename string, found bool) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("\nError in accessing import path: %v \n", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			pathAndFilename, found = findSourceFile(path+file.Name()+"/", importModName, fileExtensionName)
			if found {
				return
			}
		}
		filenameParts := strings.Split(file.Name(), ".")
		if len(filenameParts) != 2 {
			continue
		}
		if filenameParts[1] != fileExtensionName {
			continue
		}
		if filenameParts[0] != importModName {
			continue
		}
		return path + file.Name(), true
	}
	return "", false
}

func convertPathToCommonFormat(path string) string {
	pathWithoutVolume := path[len(filepath.VolumeName(path)):]
	result := filepath.ToSlash(pathWithoutVolume)
	// lets also add / in the end if missing
	if l := len(result); l > 0 {
		lastCh := result[len(result)-1:]
		if lastCh != "/" {
			result += "/"
		}
	}
	return result
}

func readExtModuleFromFile(sid SymID, importPath string) (topFrame *Frame, found bool, err error) {
	// TODO: duplicate code...
	importSpecs := make(map[string]string)
	if importPath != "" {
		for _, onepart := range strings.Split(importPath, ";") {
			splits := strings.Split(onepart, ":")
			if len(splits) != 2 {
				continue
			}
			importSpecs[splits[0]] = splits[1]
		}
	}

	fileExtensionName := "so"
	importModName := SymIDMap.AsString(sid)
	importFileName := importModName
	if importModName == "" {
		return
	}
	if specModName, nameFound := importSpecs["file"]; nameFound {
		specFilenameParts := strings.Split(specModName, ".")
		if len(specFilenameParts) == 2 {
			importFileName = specFilenameParts[0]
			fileExtensionName = specFilenameParts[1]
		}
	}
	importFilePath := os.Getenv("FUNLPATH")

	currentWorkDir, oserr := os.Getwd()
	if oserr != nil {
		err = oserr
		return
	}
	currentWorkDir += "/"

	// lets firs search from current working dir and its subdirectories
	targetPath, fileFound := findSourceFile(convertPathToCommonFormat(currentWorkDir), importFileName, fileExtensionName)

	// then try directory from env.var. and its subdirectories
	if !fileFound {
		if importFilePath != "" {
			targetPath, fileFound = findSourceFile(convertPathToCommonFormat(importFilePath), importFileName, fileExtensionName)
		}
	}
	if !fileFound {
		err = fmt.Errorf("Module not found: %s", importModName)
		return
	}

	// open external plugin module
	topFrame, err = SetupExtModule(targetPath)
	if err == nil {
		found = true
	}

	return
}

func readModuleFromFile(inProcCall bool, sid SymID, importPath string) (topFrame *Frame, found bool, err error) {
	importSpecs := make(map[string]string)
	if importPath != "" {
		for _, onepart := range strings.Split(importPath, ";") {
			splits := strings.Split(onepart, ":")
			if len(splits) != 2 {
				continue
			}
			importSpecs[splits[0]] = splits[1]
		}
	}

	fileExtensionName := "fnl"
	importModName := SymIDMap.AsString(sid)
	importFileName := importModName
	if importModName == "" {
		return
	}
	if specModName, nameFound := importSpecs["file"]; nameFound {
		specFilenameParts := strings.Split(specModName, ".")
		if len(specFilenameParts) == 2 {
			importFileName = specFilenameParts[0]
			fileExtensionName = specFilenameParts[1]
		}
	}
	importFilePath := os.Getenv("FUNLPATH")

	currentWorkDir, oserr := os.Getwd()
	if oserr != nil {
		err = oserr
		return
	}
	currentWorkDir += "/"

	// lets firs search from current working dir and its subdirectories
	targetPath, fileFound := findSourceFile(convertPathToCommonFormat(currentWorkDir), importFileName, fileExtensionName)

	// then try directory from env.var. and its subdirectories
	if !fileFound {
		if importFilePath != "" {
			targetPath, fileFound = findSourceFile(convertPathToCommonFormat(importFilePath), importFileName, fileExtensionName)
		}
	}
	if !fileFound {
		err = fmt.Errorf("Module not found: %s", importModName)
		return
	}

	content, err := ioutil.ReadFile(targetPath)
	if err != nil {
		err = fmt.Errorf("Source file reading failed: %v", err)
		return
	}

	topFrame, err = commonAddFunModToNamespace(inProcCall, targetPath, importModName, content)
	if err != nil {
		err = fmt.Errorf("Importing source failed: %v", err)
		return
	}
	found = true
	return
}

// AddNStoCache is for std usage
func AddNStoCache(inProcCall bool, importModName string, nspace *NSpace) *Frame {
	// first create top frame for namespace and put to nsDir
	topFrame := newTopFrameForNS(nspace)
	nsSid := SymIDMap.Add(importModName)

	// then put imports to namespace
	AddImportsToNamespaceSub(nspace, topFrame)

	topFrame.inProcCall = inProcCall // NOTE. this was added later as otherwise proc calls failed at main level

	// then evaluate and assign symbols of namespace
	nsDir.FillFromAstNSpaceAndStore(topFrame, nsSid, nspace)

	return topFrame
}

func commonAddFunModToNamespace(inProcCall bool, targetPath, importModName string, content []byte) (topFrame *Frame, err error) {
	parser := NewParser(NewDefaultOperators(), &targetPath)
	var nsName string
	var nspace *NSpace
	nsName, nspace, err = parser.Parse(string(content))
	if err != nil {
		err = fmt.Errorf("Parse error: %v", err)
		return
	}
	if nsName != importModName {
		err = fmt.Errorf("Mismatch in module name: %s vs %s", nsName, importModName)
		return
	}

	topFrame = AddNStoCache(inProcCall, importModName, nspace)
	return
}

func AddFunModToNamespace(importModName string, content []byte) (err error) {
	_, err = commonAddFunModToNamespace(true, importModName, importModName, content)
	if err != nil {
		err = fmt.Errorf("Importing source failed: %v", err)
	}
	return
}
