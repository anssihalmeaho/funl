package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	fs, _ := ioutil.ReadDir("./stdfun/")
	out, err := os.Create("./funl/stdfunfiles.go")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error in generating std .fun files (%v)", err))
		return
	}

	out.Write([]byte("package funl \n\nfunc init() {\n"))
	for _, f := range fs {
		if before, ok := strings.CutSuffix(f.Name(), ".fnl"); ok {
			keyStr := `"` + before + `"`
			out.Write([]byte("\n\tstdfunMap[" + keyStr + "] = `"))
			f, err := os.Open("./stdfun/" + f.Name())
			if err != nil {
				fmt.Println(fmt.Sprintf("Error in reading std .fun files (%v)", err))
				f.Close()
				break
			}
			io.Copy(out, f)
			out.Write([]byte("`\n"))
			f.Close()
		}
	}
	out.Write([]byte("}\n"))
	out.Close()
}
