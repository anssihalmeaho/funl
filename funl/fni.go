package funl

import (
	"fmt"
	"plugin"
)

type FNIApi interface {
	RegExtProc(ExtProcType, string) error
}

type FNIHandler struct {
	topFrame *Frame
}

func (fni *FNIHandler) RegExtProc(extProc ExtProcType, extProcName string) (err error) {
	/*
		nsSid := SymIDMap.Add(extProcName)
		nsDir.Put(nsSid, fni.topFrame)

		epVal := Value{Kind: ExtProcValue, Data: extProc}
		item := &Item{Type: ValueItem, Data: epVal}
		err = fni.topFrame.Syms.Add(extProcName, item)
	*/
	return
}

type ExtSetupHandler func(FNIApi) error

func SetupExtModule(targetPath string) (topFrame *Frame, err error) {
	var plug *plugin.Plugin
	plug, err = plugin.Open(targetPath)
	if err != nil {
		err = fmt.Errorf("Plugin file could not be read: %s: %v", targetPath, err)
		return
	}
	v, err := plug.Lookup("Setup")
	if err != nil {
		err = fmt.Errorf("Setup function not found in plugin: %s: %v", targetPath, err)
		return
	}
	setupHandler := v.(func(FNIApi) error)

	topFrame = &Frame{
		Syms:     NewSymt(),
		OtherNS:  make(map[SymID]ImportInfo),
		Imported: make(map[SymID]*Frame),
	}
	napi := &FNIHandler{topFrame: topFrame}
	err = setupHandler(napi)
	return
}
