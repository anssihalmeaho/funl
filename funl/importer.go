package funl

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ModuleImporter interface {
	FindModule(importModName string, extensionName string) (targetPath string, content []byte, err error)
}

type pack struct {
	modContents map[string][]byte
}

func newPack() *pack {
	return &pack{
		modContents: make(map[string][]byte),
	}
}

type packageImporter struct {
	mods *pack
}

func (importer *packageImporter) FindModule(importFileName string, extensionName string) (targetPath string, content []byte, err error) {
	var modFound bool
	content, modFound = importer.mods.modContents[importFileName]
	if !modFound {
		err = fmt.Errorf("Module not found: %s", importFileName)
		return
	}
	targetPath = importFileName
	return
}

func GetModsFromTar(tarContent []byte) (map[string][]byte, error) {
	result := map[string][]byte{}

	buf := bytes.NewBuffer(tarContent)

	// Open and iterate through the files in the archive.
	tr := tar.NewReader(buf)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return result, err
		}
		if !hdr.FileInfo().IsDir() {

			modContent, err := io.ReadAll(tr)
			if err != nil {
				return result, err
			}
			_, file := filepath.Split(hdr.Name)
			parts := strings.Split(file, ".")
			if l := len(parts); l == 2 {
				switch parts[1] {
				case "fnl":
					result[parts[0]] = modContent
				case "fpack":
					subm, err := GetModsFromTar(modContent)
					if err != nil {
						return result, err
					}
					for k, v := range subm {
						result[k] = v
					}
				}
			}
		}
	}

	return result, nil
}

func findFromPackageFiles(path, importModName string) (string, []byte, error) {
	fileExtensionName := "fpack"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("\nError in accessing import path: %v \n", err)
		return "", []byte{}, err
	}
	for _, file := range files {
		if file.IsDir() {
			targetPath, content, err := findFromPackageFiles(path+file.Name()+"/", importModName)
			if err == nil {
				return targetPath, content, nil
			}
		} else {
			filenameParts := strings.Split(file.Name(), ".")
			if len(filenameParts) == 2 && filenameParts[1] == fileExtensionName {
				data, err := os.ReadFile(path + file.Name())
				if err != nil {
					continue
				}
				mods, err := GetModsFromTar(data)
				if err != nil {
					continue
				}
				if content, isModFound := mods[importModName]; isModFound {
					return path + file.Name(), content, nil
				}
			}
		}
	}
	return "", []byte{}, fmt.Errorf("Module not found in packages (%s)", importModName)
}

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

type fileImporter struct{}

func (importer *fileImporter) FindModule(importFileName string, extensionName string) (targetPath string, content []byte, err error) {
	importFilePath := os.Getenv("FUNLPATH")

	currentWorkDir, oserr := os.Getwd()
	if oserr != nil {
		err = oserr
		return
	}
	currentWorkDir += "/"

	// lets firs search from current working dir and its subdirectories
	targetPath, fileFound := findSourceFile(convertPathToCommonFormat(currentWorkDir), importFileName, extensionName)

	// then try directory from env.var. and its subdirectories
	if !fileFound {
		if importFilePath != "" {
			targetPath, fileFound = findSourceFile(convertPathToCommonFormat(importFilePath), importFileName, extensionName)
		}
	}

	// lets try to find some packages (.fpack)
	if !fileFound {
		targetPath, content, err = findFromPackageFiles(convertPathToCommonFormat(currentWorkDir), importFileName)
		if err != nil {
			// then try directory from env.var. and its subdirectories
			targetPath, content, err = findFromPackageFiles(convertPathToCommonFormat(importFilePath), importFileName)
		}
		if err == nil {
			return
		}
	}

	if !fileFound {
		err = fmt.Errorf("Module not found: %s", importFileName)
		return
	}

	content, err = ioutil.ReadFile(targetPath)
	if err != nil {
		err = fmt.Errorf("Source file reading failed: %v", err)
		return
	}

	return
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

func readExtModuleFromFile(sid SymID, importPath string, interpreter *Interpreter) (topFrame *Frame, found bool, err error) {
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
	topFrame, err = SetupExtModule(targetPath, interpreter)
	if err == nil {
		found = true
	}

	return
}

func readModuleFromFile(inProcCall bool, sid SymID, importPath string, interpreter *Interpreter) (topFrame *Frame, found bool, err error) {
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

	targetPath, content, err := interpreter.Importer.FindModule(importFileName, fileExtensionName)

	topFrame, err = commonAddFunModToNamespace(inProcCall, targetPath, importModName, content, interpreter)
	if err != nil {
		err = fmt.Errorf("Importing source failed: %v", err)
		return
	}
	found = true
	return
}

// AddNStoCache is for std usage
func AddNStoCache(inProcCall bool, importModName string, nspace *NSpace, interpreter *Interpreter) *Frame {
	// first create top frame for namespace and put to nsDir
	topFrame := newTopFrameForNS(nspace, interpreter)
	nsSid := SymIDMap.Add(importModName)

	// then put imports to namespace
	AddImportsToNamespaceSub(nspace, topFrame, interpreter)

	topFrame.inProcCall = inProcCall // NOTE. this was added later as otherwise proc calls failed at main level

	// then evaluate and assign symbols of namespace
	interpreter.NsDir.FillFromAstNSpaceAndStore(topFrame, nsSid, nspace)

	return topFrame
}

func commonAddFunModToNamespace(inProcCall bool, targetPath, importModName string, content []byte, interpreter *Interpreter) (topFrame *Frame, err error) {
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

	topFrame = AddNStoCache(inProcCall, importModName, nspace, interpreter)
	return
}

func AddFunModToNamespace(importModName string, content []byte, interpreter *Interpreter) (err error) {
	_, err = commonAddFunModToNamespace(true, importModName, importModName, content, interpreter)
	if err != nil {
		err = fmt.Errorf("Importing source failed: %v", err)
	}
	return
}
