/*******************************************************************************

 Project: Tourney

 Module: build

 Description: build engines

 Author(s): Andrew Backes
 Created: 10/03/2015

*******************************************************************************/

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"fmt"
	
)

type BuildSpec struct {
	Repo string
	Branch string
	BuildFile string
	EngineFile string
}

func (B BuildSpec) exec(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir, _ = filepath.Abs(B.Dir())
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}

func (B BuildSpec) Name() string {
	n := strings.Split(B.Repo, "/")
	return n[len(n)-1]
}

func (B BuildSpec) Folder() string {
	return B.Name() + "-" +B.Branch
}

func (B BuildSpec) Dir() string {
	return filepath.Join(Settings.BuildDirectory, B.Folder())
}

func (B *BuildSpec) GitClone() error {
	// remove old repo:
	fmt.Print("Cloning repo: " + B.Repo + ":" + B.Branch + "... ")
	d, _ := filepath.Abs(B.Dir())
	if err := os.RemoveAll(d); err != nil {
		fmt.Println(err)
	}
	if err := os.MkdirAll(B.Dir(), os.ModePerm); err != nil {
		return err
	}
	err := B.exec( "git", "clone", "-b", B.Branch, B.Repo, "." )
	if err != nil {
		fmt.Println("FAILED.", err)
	} else {
		fmt.Println("Success.")
	}
	return err
}

func (B *BuildSpec) Build() error {
	fmt.Print("Building... ")
	fp, _ := filepath.Abs( filepath.Join(B.Dir(), B.BuildFile) )
	err := B.exec( fp )
	// rename the file that was built:
	B.RenameEngineFile()
	if err != nil {
		fmt.Println("FAILED.", err)
	} else {
		fmt.Println("Success.")
	}
	return err
}

func (B *BuildSpec) FullEngineFile() string {
	if filepath.IsAbs(B.EngineFile) {
		return B.EngineFile
	}
	f, _ := filepath.Abs( filepath.Join(B.Dir(), B.EngineFile ) )
	return f
}

// RenameEngineFile appends the branch to the engine's filename
func (B *BuildSpec) RenameEngineFile() {
	oldpath := B.FullEngineFile()
	d, f := filepath.Split(oldpath)
	newfile := strings.TrimSuffix( f, ".exe" )
	newfile += "-" + B.Branch
	if strings.HasSuffix(f, ".exe") {
		newfile += ".exe"
	}
	newpath := filepath.Join(d, newfile)
	if err := os.Rename(oldpath, newpath); err == nil {
		B.EngineFile = newpath
	}
}