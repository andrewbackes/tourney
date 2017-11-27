package instructions

import (
	"fmt"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/util"
	"os"
	"path/filepath"
	"strings"
)

type EngineUploadInstruction struct {
	models.Engine
	filepath string
}

func (e *EngineUploadInstruction) Execute() {
	p := prepare(e)
	if e.filepath != "" {
		upload(p)
	}
	util.PostJSON(util.GetAPIURL()+"/engines", p.Engine)
}

func upload(e *EngineUploadInstruction) {
	if _, err := os.Stat(e.filepath); os.IsNotExist(err) {
		fmt.Println(e.filepath + " does not exist")
		os.Exit(1)
	}
	err := util.PostFile(e.filepath, e.URL)
	if err != nil {
		panic(err)
	}
	fmt.Println("Uploaded", e.filepath)
}

func (e *EngineUploadInstruction) Validate() {
	if e.filepath == "" && e.URL == "" {
		util.Fail("Please specify either a url or local file path for the engine.")
	}
	if strings.HasSuffix(e.filepath, ".zip") && e.Executable == "" {
		util.Fail("Please specify the relative path to the executable within the .zip file.\n	usage: tourney-cli --name myengine --version 1.0 --os linux --executable bin/myengine somepath/myengine.zip")
	}
	if e.Name == "" || e.Version == "" || e.Os == "" {
		util.Fail("'add engine' command requires --name, --version, and --os")
	}
}

func NewEngineUploadIntruction(params map[string]string) *EngineUploadInstruction {
	e := &EngineUploadInstruction{
		Engine: models.Engine{
			Name:       params["name"],
			Version:    params["version"],
			Os:         params["os"],
			Executable: params["executable"],
			URL:        params["url"],
		},
		filepath: params["filepath"],
	}
	return e
}

func prepare(e *EngineUploadInstruction) *EngineUploadInstruction {
	p := *e
	filename := filepath.Base(e.filepath)
	if e.URL == "" {
		p.URL = util.GetAPIURL() + "/engineFiles/" + e.Name + "/" + e.Version + "/" + e.Os + "/" + filename
	}
	if e.Executable == "" {
		p.Executable = filename
	}
	return &p
}
