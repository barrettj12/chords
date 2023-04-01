package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/barrettj12/chords/src/dblayer"
)

// diff compares the local DB to remote.
//
//	chords diff
func diff(st state, args []string) {
	// Create temp dir to hold files
	tempdir, err := os.MkdirTemp("", "chords-diff")
	check(err)
	defer os.RemoveAll(tempdir)
	fmt.Printf("writing files to %s\n", tempdir)
	zipFile := filepath.Join(tempdir, "data.zip")

	// sftp connection to remote VM
	// TODO: this takes a long time. Might be quicker to zip everything up
	// on the VM, and copy the zip file over.
	sshCmd := exec.Command("fly", "ssh", "sftp", "shell")
	// sshCmd.Stdout = os.Stdout
	// sshCmd.Stderr = os.Stderr
	stdin, err := sshCmd.StdinPipe()
	check(err)

	check(sshCmd.Start())
	io.WriteString(stdin, fmt.Sprintf("get /data %s\n", zipFile))
	io.WriteString(stdin, "\x04") // exit signal
	fmt.Println("waiting for ssh command to exit")
	check(sshCmd.Wait())
	fmt.Println("ssh command exited")

	// Open zip file
	zipReader, err := zip.OpenReader(zipFile)
	check(err)
	defer zipReader.Close()

	// Initialise a map to keep track of which files exist locally/remotely
	fileMap := make(map[string]fileInfo, len(zipReader.File))

	// Put remote files in fileMap
	for _, f := range zipReader.File {
		filename, err := filepath.Rel("/data/", f.Name)
		check(err)
		fileMap[filename] = fileInfo{remoteFile: f}
	}

	// Put local files in fileMap
	filepath.WalkDir(st.dbPath, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() {
			filename, err := filepath.Rel(st.dbPath, path)
			check(err)
			fileMap[filename] = fileInfo{
				existsLocally: true,
				remoteFile:    fileMap[filename].remoteFile,
			}
		}
		return nil
	})

	// Look for differences
	diffFound := false
	reportDiff := func(format string, a ...any) {
		diffFound = true
		fmt.Printf(format, a...)
	}

	for filename, info := range fileMap {
		if info.existsLocally && info.remoteFile == nil {
			reportDiff("file %q exists locally but not remotely\n", filename)
			continue
		}
		if !info.existsLocally && info.remoteFile != nil {
			reportDiff("file %q exists remotely but not locally\n", filename)
			continue
		}

		// Check local file and remote file match
		localFileContents, err := os.ReadFile(filepath.Join(st.dbPath, filename))
		check(err)

		remoteFileReader, err := info.remoteFile.Open()
		check(err)
		remoteFileContents, err := io.ReadAll(remoteFileReader)
		check(err)

		if strings.HasSuffix(filename, ".json") {
			// For JSON files, the exact formatting is not important.
			// Just check that the fields are the same.
			localMeta := &dblayer.SongMeta{}
			err = json.Unmarshal(localFileContents, localMeta)
			check(err)

			remoteMeta := &dblayer.SongMeta{}
			err = json.Unmarshal(remoteFileContents, remoteMeta)
			check(err)

			if *localMeta != *remoteMeta {
				reportDiff(`metadata %q is different between local and remote
local meta: %#v
remote meta: %#v
---------------------------------------
`, filename, localMeta, remoteMeta)
			}
			continue
		}

		if !bytes.Equal(localFileContents, remoteFileContents) {
			reportDiff(`file %q is different between local and remote
--------- local -----------------------
%s
--------- remote ----------------------
%s
---------------------------------------
`, filename, localFileContents, remoteFileContents)
		}
	}

	if !diffFound {
		fmt.Println("no differences found between local and remote :)")
	}
}

type fileInfo struct {
	existsLocally bool
	remoteFile    *zip.File
}
