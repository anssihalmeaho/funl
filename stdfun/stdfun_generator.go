package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

//go generate
//go build -o funl.exe

type myWriter struct {
	orig io.Writer
}

func (myw *myWriter) Write(p []byte) (n int, err error) {
	var newp []byte
	for _, bv := range p {
		if bv == 10 {
			newp = append(newp, []byte(" ")...)
		}
		newp = append(newp, bv)
	}
	return myw.orig.Write(newp)
}

func main() {
	fs, _ := ioutil.ReadDir("./stdfun/")
	out, err := os.Create("stdfunfiles.go")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error in generating std .fun files (%v)", err))
		return
	}

	mw := &myWriter{orig: out}

	out.Write([]byte("package main \n\nfunc init() {\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".fun") {
			keyStr := `"` + strings.TrimSuffix(f.Name(), ".fun") + `"`
			out.Write([]byte("\n\tstdfunMap[" + keyStr + "] = `"))
			f, err := os.Open("./stdfun/" + f.Name())
			if err != nil {
				fmt.Println(fmt.Sprintf("Error in reading std .fun files (%v)", err))
				f.Close()
				break
			}
			io.Copy(mw, f)
			out.Write([]byte("`\n"))
			f.Close()
		}
	}
	out.Write([]byte("}\n"))
	out.Close()
}
