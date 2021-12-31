//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlfredDobradi/ledgerlog/internal/cli"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

const pkgBase string = "github.com/AlfredDobradi/ledgerlog/cmd"

var (
	ldflags         = "-s -w -X main.commitHash=$COMMIT_HASH -X main.buildTime=$BUILD_TIME -X main.tag=$VERSION_TAG"
	targetOS        = []string{"linux", "darwin"}
	checksumFormats = []string{"sha256", "md5"}
	cleanupTargets  = []string{"./target"}
)

// Daemon builds the daemon binary
func (Build) Daemon() error {
	return build("daemon")
}

// Client builds the client binary
func (Build) Client() error {
	return build("client")
}

func build(pkg string) error {
	pkgPath := fmt.Sprintf("%s/%s", pkgBase, pkg)

	for _, os := range targetOS {
		env := flagEnv()
		env["GOOS"] = os
		output := fmt.Sprintf("./target/%s/%s/%s", os, pkg, pkg)
		fmt.Printf("Building package '%s' in '%s'...", pkgPath, output)
		if err := sh.RunWith(env, "go", "build", "-o", output, "-ldflags", ldflags, pkgPath); err != nil {
			cli.Failure()
			return fmt.Errorf("Failed building package for %s: %w", os, err)
		}
		cli.Success()

		fmt.Printf("Generating sha256 and md5 checksum files for %s target...", os)
		if err := generateCheckSumFiles(pkg, os); err != nil {
			cli.Failure()
			return fmt.Errorf("Failed writing checksum files for %s: %w", os, err)
		}
		cli.Success()
	}
	return nil
}

func Clean() error {
	for _, target := range cleanupTargets {
		fmt.Printf("Removing %s...", target)
		if err := os.RemoveAll(target); err != nil {
			cli.Failure()
			return fmt.Errorf("Failed removing %s: %w", target, err)
		}
		cli.Success()
	}
	return nil
}

func flagEnv() map[string]string {
	return map[string]string{
		"COMMIT_HASH": getCommitHash(),
		"BUILD_TIME":  time.Now().Format("2006-01-02T15:04:05Z0700"),
		"VERSION_TAG": getVersionTag(),
		"CGO_ENABLED": "0",
	}
}

func getCommitHash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

func getVersionTag() string {
	versiontag, err := sh.Output("cat", ".release_version")
	if err != nil {
		versiontag, _ = sh.Output("git", "describe", "--tags", "--abbrev=0")
	}

	return versiontag
}

func generateCheckSumFiles(pkg, target string) error {

	targetpath := fmt.Sprintf("target/%s/%s/%s", target, pkg, pkg)
	for _, checksumFormat := range checksumFormats {
		out, err := sh.Output(fmt.Sprintf("%ssum", checksumFormat), targetpath)
		if err != nil {
			return err
		}

		f, openerr := os.OpenFile(fmt.Sprintf("%s.%s", targetpath, checksumFormat), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
		defer f.Close()
		if openerr != nil {
			return openerr
		}

		if _, err := f.Write([]byte(strings.Split(out, " ")[0])); err != nil {
			return err
		}
	}

	return nil
}
