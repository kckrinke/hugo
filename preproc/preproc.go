// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package preproc

import (
	"fmt"
	jww "github.com/spf13/jwalterweatherman"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime/debug"
)

func Ddie(f string, argv ...interface{}) {
	fmt.Printf(fmt.Sprintf("[preproc-die] %v\n", f), argv...)
	debug.PrintStack()
	os.Exit(1)
}

func Dd(f string, argv ...interface{}) {
	jww.FEEDBACK.Printf(fmt.Sprintf("[preproc] %v", f), argv...)
	fmt.Printf(fmt.Sprintf("[preproc] %v\n", f), argv...)
}

func PreProcess(b []byte, path string) []byte {
	if ppcmd := os.Getenv("HUGO_PREPROC_CMD"); ppcmd != "" {
		// check if ppcmd exists
		if _, err := os.Stat(ppcmd); os.IsNotExist(err) {
			jww.ERROR.Printf("HUGO_PREPROC_CMD file not found: %v", ppcmd)
			return b
		}
		// construct and start cmd
		c := exec.Command(ppcmd, path)
		w, _ := c.StdinPipe()
		r, _ := c.StdoutPipe()
		if err := c.Start(); err != nil {
			jww.FATAL.Fatalln("HUGO_PREPROC_CMD: %v", err)
		}
		// write b to stdin, close stdin
		w.Write(b)
		w.Close()
		// read all of stdout as nb, return nb
		nb, err := ioutil.ReadAll(r)
		if err != nil {
			jww.ERROR.Printf("HUGO_PREPROC_CMD: %v", err)
			return b
		}
		Dd("os.Exec %v \"%v\"", ppcmd, path)
		if err = c.Wait(); err != nil {
			jww.ERROR.Printf("HUGO_PREPROC_CMD: %v", err)
			return b
		}
		return nb
	}
	return b
}
