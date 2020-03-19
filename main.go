package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/anssihalmeaho/funl"
	"github.com/anssihalmeaho/std"
)

//go:generate go run ./stdfun/stdfun_generator.go

const doProfiling = false
const doStackMeas = false
const doMemMeas = false

func initFunSourceSTD() (err error) {
	type funModInfo struct {
		name    string
		content string
	}
	var funmods = []funModInfo{
		{
			name:    "stdfu",
			content: stdfu,
		},
		{
			name:    "stdset",
			content: stdset,
		},
		{
			name:    "stddbc",
			content: stddbc,
		},
		{
			name:    "stdfilu", //note. this needs stdfu
			content: stdfilu,
		},
	}
	for _, fm := range funmods {
		err = funl.AddFunModToNamespace(fm.name, []byte(fm.content))
		if err != nil {
			return
		}
	}
	return
}

func main() {
	if doProfiling {
		f, err := os.Create("fup.prof")
		if err != nil {
			fmt.Println("Error in profile setup: " + err.Error())
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Runtime error: ", r)
		}
	}()

	var err error
	var fargs string
	flag.StringVar(&fargs, "args", "", fmt.Sprintf("arguments for %s", os.Args[0]))
	var name string
	flag.StringVar(&name, "name", "main", "function or procedure to be evaluated (main is default)")
	var modName string
	flag.StringVar(&modName, "mod", "main", "module (namespace) to be evaluated (main is default)")
	replPtr := flag.Bool("repl", false, "starts REPL")
	silentPtr := flag.Bool("silent", false, "does not print result of evaluation when returning from program (not silent is default)")
	noPrintPtr := flag.Bool("noprint", false, "prevents printing from functions (by print-operator)")
	doRTEPrintPtr := flag.Bool("rteprint", false, "enables printing RTE location and scope")
	var evalStr string
	flag.StringVar(&evalStr, "eval", "", "evaluate expression")
	flag.Parse()

	if *noPrintPtr {
		funl.PrintingDisabledInFunctions = true
	}
	if *doRTEPrintPtr {
		funl.PrintingRTElocationAndScopeEnabled = true
	}

	var parsedArgs []*funl.Item
	if fargs != "" {
		parsedArgs, err = funl.GetArgs(fargs)
		if err != nil {
			fmt.Println("Error in parsing arguments: ", err)
			return
		}
	}

	var skipSrcFile bool
	if modName != "main" {
		skipSrcFile = true
	}

	var content []byte
	var srcFileName string
	if !skipSrcFile {
		if *replPtr {
			srcFileName = "repl.fun"
		}

		if evalStr != "" {
			content = []byte(fmt.Sprintf("ns main main = proc() %s end endns", evalStr))
			name = "main"
		} else {
			others := flag.Args()
			if srcFileName == "" {
				if len(others) != 1 {
					fmt.Println("Source file not given correctly")
					return
				}
				srcFileName = others[0]
			}
			if srcFileName == "repl.fun" {
				content = []byte(repl)
			} else {
				content, err = ioutil.ReadFile(srcFileName)
				if err != nil {
					fmt.Println(fmt.Sprintf("Source file reading failed: %v", err))
					return
				}
			}
		}
	} else {
		srcFileName = modName
		var argStr string
		if fargs != "" {
			argStr = "," + fargs
		}
		content = []byte(fmt.Sprintf("ns main import %s main = proc() call(%s.%s%s) end endns", modName, modName, name, argStr))
		name = "main"
	}

	var retValue funl.Value
	retValue, err = funl.FunlMainWithArgs(string(content), parsedArgs, name, srcFileName, std.InitSTD, initFunSourceSTD)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %v", err))
		return
	}
	if !*replPtr {
		if !*silentPtr {
			fmt.Println(fmt.Sprintf("%#v", retValue))
		}
	}

	if doStackMeas {
		var mems runtime.MemStats
		runtime.ReadMemStats(&mems)
		fmt.Println("in use   : ", mems.StackInuse)
		fmt.Println("stack sys: ", mems.StackSys)
	}
	if doMemMeas {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("\nAlloc = %v MB", transformBytesToMegaBytes(m.Alloc))
		fmt.Printf("\tTotalAlloc = %v MB", transformBytesToMegaBytes(m.TotalAlloc))
		fmt.Printf("\tSys = %v MB", transformBytesToMegaBytes(m.Sys))
		fmt.Printf("\tNumGC = %v\n", m.NumGC)
	}
}

func transformBytesToMegaBytes(b uint64) uint64 {
	return b / 1024 / 1024
}
