package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/Moonlington/FloSelfbot/commands"

	"github.com/google/go-github/github"
	"github.com/inconshreveable/go-update"
)

func getNewestVersionLink() (url string, err error) {
	client := github.NewClient(nil)
	latest, _, err := client.Repositories.GetLatestRelease("Moonlington", "FloSelfbot")
	latesttn := *latest.TagName
	var ext string
	var SOS string
	switch runtime.GOOS {
	case "windows":
		SOS = "Windows"
		ext = "zip"
	case "darwin":
		SOS = "OSX"
		ext = "zip"
	case "linux":
		SOS = "Linux"
		ext = "tgz"
	default:
		err = errors.New("Not a known OS")
		return
	}
	url = fmt.Sprintf("https://github.com/Moonlington/FloSelfbot/releases/download/v%s/Floselfbot_%s_%s.%s", latesttn[1:], SOS, latesttn[1:], ext)
	return
}

func getUpdateFile(url string) (out io.Reader, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if runtime.GOOS == "windows" || runtime.GOOS == "OSX" {
		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		r := bytes.NewReader(b)
		var zr *zip.Reader
		zr, err = zip.NewReader(r, int64(r.Len()))
		if err != nil {
			return
		}
		if len(zr.File) != 1 {
			err = errors.New("Something went wrong with getting the package")
			return
		}
		out, err = zr.File[0].Open()
		return
	} else if runtime.GOOS == "linux" {
		var gzr *gzip.Reader
		gzr, err = gzip.NewReader(resp.Body)
		if err != nil {
			return
		}
		tr := tar.NewReader(gzr)
		tr.Next()
		out = tr
		return
	}
	return
}

func doUpdate() error {
	// request the new file
	url, err := getNewestVersionLink()
	if err != nil {
		return err
	}
	r, err := getUpdateFile(url)
	if err != nil {
		return err
	}
	err = update.Apply(r, update.Options{})
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			fmt.Printf("Failed to rollback from bad update: %v", rerr)
		}
	}
	return err
}

// UpdateCommand struct handles Update Command
type UpdateCommand struct {
	commands.Command
}

func (u *UpdateCommand) message(ctx *commands.Context) {
	ctx.QuickSendEm("Updating...")
	err := doUpdate()
	logerror(err)
	ctx.QuickSendEm("Done!")
}

func (u *UpdateCommand) description() string { return "Updates" }
func (u *UpdateCommand) usage() string       { return "" }
func (u *UpdateCommand) detailed() string    { return "Updates" }
func (u *UpdateCommand) subcommands() map[string]commands.Command {
	return make(map[string]commands.Command)
}
