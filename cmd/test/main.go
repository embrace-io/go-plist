package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"howett.net/plist"
)

type PBXProj struct {
	ArchiveVersion int                  `plist:"archiveVersion"`
	ObjectVersion  int                  `plist:"objectVersion"`
	Objects        map[string]PBXObject `plist:"objects"`
	RootObject     string               `plist:"rootObject"`
}

// Generic object
type PBXObject struct {
	ISA string `plist:"isa"`
	PBXShellScriptBuildPhase
	PBXContainerItemProxy
}

func (o PBXObject) ToPBXShellScriptBuildPhase(id string) PBXShellScriptBuildPhase {
	phase := o.PBXShellScriptBuildPhase
	phase.ISA = o.ISA
	phase.ID = id
	return phase
}

type Commented struct {

}

type PBXShellScriptBuildPhase struct {
	ID          string `plist:"-"`
	ISA         string `plist:"isa"`
	ShellScript string `plist:"shellScript,omitempty"`
	ShellPath   string `plist:"shellPath,omitempty"`
	Name        string `plist:"name,omitempty"`
}

type PBXContainerItemProxy struct {
	ID         string `plist:"-"`
	ISA        string `plist:"isa"`
	RemoteInfo string `plist:"remoteInfo,omitempty"`
}

func main() {
	raw, err := ioutil.ReadFile("/Users/koheihisakuni/dev/embrace/go/tool/embrace/project.pbxproj")
	if err != nil {
		log.Fatal(err)
		return
	}
	buf := bytes.NewReader(raw)

	var proj PBXProj
	decoder := plist.NewDecoder(buf)
	err = decoder.Decode(&proj)
	if err != nil {
		log.Fatal(err)
		return
	}

	//shellScripts := proj.shellScriptBuildPhases()
	//fmt.Printf(">>> %+v\n", proj.Comments)

	updated, err := plist.MarshalIndent(proj, plist.OpenStepFormat, "  ")
	//writeBuf := &bytes.Buffer{}
	//encoder := plist.NewEncoder(writeBuf)
	//err = encoder.Encode(shellScripts)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("RUNNING TEST %v\n", string(updated))
}
