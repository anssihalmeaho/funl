package std

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/anssihalmeaho/funl/funl"
)

const (
	writeMode  = 1
	readMode   = 2
	rwMode     = writeMode + readMode
	appendMode = 4
)

func initSTDFiles(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdfiles"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdFilesFuncs := []stdFuncInfo{
		{
			Name:   "create",
			Getter: getStdFilesCreate,
		},
		{
			Name:   "open",
			Getter: getStdFilesOpen,
		},
		{
			Name:   "write",
			Getter: getStdFilesWrite,
		},
		{
			Name:   "write-at",
			Getter: getStdFilesWriteAt,
		},
		{
			Name:   "writeln",
			Getter: getStdFilesWriteLn,
		},
		{
			Name:   "read",
			Getter: getStdFilesRead,
		},
		{
			Name:   "read-at",
			Getter: getStdFilesReadAt,
		},
		{
			Name:   "read-all",
			Getter: getStdFilesReadAll,
		},
		{
			Name:   "readlines",
			Getter: getStdFilesReadLines,
		},
		{
			Name:   "seek",
			Getter: getStdFilesSeek,
		},
		{
			Name:   "remove",
			Getter: getStdFilesRemove,
		},
		{
			Name:   "rename",
			Getter: getStdFilesRename,
		},
		{
			Name:   "close",
			Getter: getStdFilesClose,
		},
		{
			Name:   "read-dir",
			Getter: getStdFilesReadDir,
		},
		{
			Name:   "finfo-map",
			Getter: getStdFilesFInfoMap,
		},
		{
			Name:   "cwd",
			Getter: getStdFilesCwd,
		},
		{
			Name:   "chdir",
			Getter: getStdFilesChDir,
		},
		{
			Name:   "mkdir",
			Getter: getStdFilesMkDir,
		},
	}

	item := &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.IntValue, Data: writeMode}}
	err = topFrame.Syms.Add("w", item)
	if err != nil {
		return
	}
	item = &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.IntValue, Data: readMode}}
	err = topFrame.Syms.Add("r", item)
	if err != nil {
		return
	}
	item = &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.IntValue, Data: appendMode}}
	err = topFrame.Syms.Add("a", item)
	if err != nil {
		return
	}

	err = setSTDFunctions(topFrame, stdModuleName, stdFilesFuncs, interpreter)
	return
}

type OpaqueFile struct {
	name   string
	path   string
	handle *os.File
}

func (file *OpaqueFile) TypeName() string {
	return "file"
}

func (file *OpaqueFile) Str() string {
	return fmt.Sprintf("file(%#v)", *file)
}

func (file *OpaqueFile) Equals(with funl.OpaqueAPI) bool {
	_, ok := with.(*OpaqueFile)
	if !ok {
		return false
	}
	return false // ==> ?????????????????????
}

type OpaqueFileInfo struct {
	info os.FileInfo
}

func (fi *OpaqueFileInfo) TypeName() string {
	return "fileinfo"
}

func (fi *OpaqueFileInfo) Str() string {
	return fmt.Sprintf("fileinfo(%#v)", *fi)
}

func (fi *OpaqueFileInfo) Equals(with funl.OpaqueAPI) bool {
	_, ok := with.(*OpaqueFileInfo)
	if !ok {
		return false
	}
	return false // ==> ?????????????????????
}

func getStdFilesMkDir(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs just one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		dirName := arguments[0].Data.(string)
		err := os.Mkdir(dirName, 0777)
		var errorText string
		if err != nil {
			errorText = fmt.Sprintf("%s: %v", name, err)
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: errorText}
		return
	}
}

func getStdFilesChDir(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs just one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		dirName := arguments[0].Data.(string)
		err := os.Chdir(dirName)
		var errorText string
		if err != nil {
			errorText = fmt.Sprintf("%s: %v", name, err)
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: errorText}
		return
	}
}

func getStdFilesCwd(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		dirName, err := os.Getwd()
		if err != nil {
			funl.RunTimeError2(frame, "%s: error in getting current working directory", name)
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: dirName}
		return
	}
}

func getStdFilesFInfoMap(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs just one", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		fInfo, ok := arguments[0].Data.(*OpaqueFileInfo)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not fileinfo value", name)
		}
		fileInfo := fInfo.info
		moperands := []*funl.Item{
			&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "name"}},
			&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: fileInfo.Name()}},

			&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "size"}},
			&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.IntValue, Data: int(fileInfo.Size())}},

			&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "mode"}},
			&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%v", fileInfo.Mode())}},

			&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "modtime"}},
			&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%v", fileInfo.ModTime())}},

			&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "is-dir"}},
			&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.BoolValue, Data: fileInfo.IsDir()}},
		}
		retVal = funl.HandleMapOP(frame, moperands)
		return
	}
}

func getStdFilesReadDir(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs just one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		dirName := arguments[0].Data.(string)
		absPath, err := filepath.Abs(dirName)
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: %v", name, err)}
			return
		}
		fileInfos, err := ioutil.ReadDir(absPath)
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: %v", name, err)}
			return
		}
		var moperands []*funl.Item
		for _, fileInfo := range fileInfos {
			fileNameVal := &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: fileInfo.Name()}}
			v := &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueFileInfo{info: fileInfo}}}
			moperands = append(moperands, fileNameVal, v)
		}
		retVal = funl.HandleMapOP(frame, moperands)
		return
	}
}

func getStdFilesSeek(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		file, ok := arguments[0].Data.(*OpaqueFile)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not file value", name)
		}
		if arguments[1].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value as 2nd argument", name)
		}
		offset, ok := arguments[1].Data.(int)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not int value", name)
		}
		if arguments[2].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value as 3rd argument", name)
		}
		whence, ok := arguments[2].Data.(int)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not int value", name)
		}

		newOffset, err := file.handle.Seek(int64(offset), whence)
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: %v", name, err)}
			return
		}
		retVal = funl.Value{Kind: funl.IntValue, Data: int(newOffset)}
		return
	}
}

func getStdFilesReadAt(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		file, ok := arguments[0].Data.(*OpaqueFile)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not file value", name)
		}
		if arguments[1].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value as 2nd argument", name)
		}
		maxcount, ok := arguments[1].Data.(int)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not int value", name)
		}
		if arguments[2].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value as 3rd argument", name)
		}
		offset, ok := arguments[2].Data.(int)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not int value", name)
		}

		rbuf := make([]byte, maxcount+1)
		n, err := file.handle.ReadAt(rbuf, int64(offset))
		var isEOF bool
		switch err {
		case nil:
		case io.EOF:
			isEOF = true
		default:
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: %v", name, err)}
			return
		}
		if !isEOF {
			if n < (maxcount + 1) {
				isEOF = true
			}
		}
		if n == (maxcount + 1) {
			n = maxcount
		}
		retValues := []funl.Value{
			funl.Value{
				Kind: funl.BoolValue,
				Data: isEOF,
			},
			funl.Value{
				Kind: funl.OpaqueValue,
				Data: &OpaqueByteArray{data: rbuf[:n]},
			},
		}
		retVal = funl.MakeListOfValues(frame, retValues)
		return
	}
}

func getStdFilesReadLines(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		file, ok := arguments[0].Data.(*OpaqueFile)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not file value", name)
		}

		var lineValues []funl.Value
		scanner := bufio.NewScanner(file.handle)
		for scanner.Scan() {
			lineValues = append(lineValues, funl.Value{Kind: funl.StringValue, Data: scanner.Text()})
		}
		if err := scanner.Err(); err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: %v", name, err)}
			return
		}
		retVal = funl.MakeListOfValues(frame, lineValues)
		return
	}
}

func getStdFilesReadAll(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		file, ok := arguments[0].Data.(*OpaqueFile)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not file value", name)
		}

		buf := bytes.NewBuffer(nil)
		io.Copy(buf, file.handle)
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueByteArray{data: buf.Bytes()}}
		return
	}
}

func getStdFilesRead(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		file, ok := arguments[0].Data.(*OpaqueFile)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not file value", name)
		}
		if arguments[1].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value as 2nd argument", name)
		}
		maxcount, ok := arguments[1].Data.(int)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not int value", name)
		}

		rbuf := make([]byte, maxcount+1)
		n, err := file.handle.Read(rbuf)
		var isEOF bool
		switch err {
		case nil:
		case io.EOF:
			isEOF = true
		default:
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: %v", name, err)}
			return
		}
		if !isEOF {
			if n < (maxcount + 1) {
				isEOF = true
			}
		}
		if n == (maxcount + 1) {
			n = maxcount
		}
		retValues := []funl.Value{
			funl.Value{
				Kind: funl.BoolValue,
				Data: isEOF,
			},
			funl.Value{
				Kind: funl.OpaqueValue,
				Data: &OpaqueByteArray{data: rbuf[:n]},
			},
		}
		retVal = funl.MakeListOfValues(frame, retValues)
		return
	}
}

func getStdFilesWriteAt(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		file, ok := arguments[0].Data.(*OpaqueFile)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not file value", name)
		}
		if arguments[1].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		byteArray, ok := arguments[1].Data.(*OpaqueByteArray)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not bytearray value", name)
		}
		if arguments[2].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value as 3rd argument", name)
		}
		offset, ok := arguments[2].Data.(int)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not int value", name)
		}

		n, err := file.handle.WriteAt(byteArray.data, int64(offset))
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: %v", name, err)}
			return
		}
		retVal = funl.Value{Kind: funl.IntValue, Data: n}
		return
	}
}

func getStdFilesWriteLn(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		file, ok := arguments[0].Data.(*OpaqueFile)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not file value", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		str, ok := arguments[1].Data.(string)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not string value", name)
		}

		n, err := file.handle.WriteString(str + "\n")
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: %v", name, err)}
			return
		}
		retVal = funl.Value{Kind: funl.IntValue, Data: n}
		return
	}
}

func getStdFilesWrite(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		file, ok := arguments[0].Data.(*OpaqueFile)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not file value", name)
		}
		if arguments[1].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		byteArray, ok := arguments[1].Data.(*OpaqueByteArray)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not bytearray value", name)
		}

		n, err := file.handle.Write(byteArray.data)
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: %v", name, err)}
			return
		}
		retVal = funl.Value{Kind: funl.IntValue, Data: n}
		return
	}
}

func getStdFilesOpen(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value as file name", name)
		}
		fileName, ok := arguments[0].Data.(string)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not string value", name)
		}
		var fmode int
		if arguments[1].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value as 2nd argument", name)
		}
		fmode, ok = arguments[1].Data.(int)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not int value", name)
		}
		var flags int
		switch fmode & rwMode {
		case writeMode:
			flags |= os.O_WRONLY
		case readMode:
			flags |= os.O_RDONLY
		case rwMode:
			flags |= os.O_RDWR
		default:
			funl.RunTimeError2(frame, "%s: unsupported mode (%d)", name, fmode)
		}
		switch fmode & appendMode {
		case appendMode:
			flags |= os.O_APPEND
		}

		dir, fname := filepath.Split(fileName)
		file := &OpaqueFile{name: fname, path: dir}
		fh, err := os.OpenFile(fileName, flags, 0)
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: open failed: %v", name, err)}
			return
		}
		file.handle = fh
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: file}
		return
	}
}

func getStdFilesClose(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		file, ok := arguments[0].Data.(*OpaqueFile)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not file value", name)
		}
		err := file.handle.Close()
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: %v", name, err)}
			return
		}
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}

func getStdFilesCreate(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value as file name", name)
		}
		fileName, ok := arguments[0].Data.(string)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not string value", name)
		}
		dir, fname := filepath.Split(fileName)
		file := &OpaqueFile{name: fname, path: dir}
		fp, err := os.Create(fileName)
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: creation failed: %v", name, err)}
			return
		}
		file.handle = fp
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: file}
		return
	}
}

func getStdFilesRename(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value as source file name", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value as target file name", name)
		}
		srcFileName, ok := arguments[0].Data.(string)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not string value", name)
		}
		trgFileName, ok := arguments[1].Data.(string)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not string value", name)
		}

		err := os.Rename(srcFileName, trgFileName)
		var text string
		if err != nil {
			text = err.Error()
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: text}
		return
	}
}

func getStdFilesRemove(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value as file name", name)
		}
		fileName, ok := arguments[0].Data.(string)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not string value", name)
		}

		err := os.Remove(fileName)
		var text string
		if err != nil {
			text = err.Error()
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: text}
		return
	}
}
