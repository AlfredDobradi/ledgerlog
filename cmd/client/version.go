package main

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime"
	"time"
)

type VersionCmd struct{}

var versionTemplate = `LedgerD:
    Version:    {{ .Version }}
    Go OS:      {{ .GoOS }}
    Go Arch:    {{ .GoArch }}
    Go version: {{ .GoVersion }}
    Built:      {{ .Built }}
    Commit:     {{ .Commit }}`

func (VersionCmd) Run(ctx *Context) error {
	tpl, parseErr := template.New("").Parse(versionTemplate)
	if parseErr != nil {
		return parseErr
	}

	versionInfo := struct{ Version, GoOS, GoArch, GoVersion, Built, Commit string }{
		Version:   tag,
		GoOS:      runtime.GOOS,
		GoArch:    runtime.GOARCH,
		GoVersion: runtime.Version(),
		Built:     buildTime,
		Commit:    commitHash,
	}

	out := bytes.NewBuffer(nil)
	t, err := time.Parse("2006-01-02T15:04:05Z0700", versionInfo.Built)
	if err == nil {
		versionInfo.Built = t.Format(time.ANSIC)
	}

	if err := tpl.Execute(out, versionInfo); err != nil {
		panic(fmt.Errorf("while compiling version template: %s", err))
	}

	fmt.Println(out.String())
	return nil
}
