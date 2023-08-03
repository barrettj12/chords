package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// diff compares the local DB to remote.
//
//	chords diff
func diff(st state, args []string) {
	// Create temp dir to hold files
	tempdir, err := os.MkdirTemp("", "chords-diff-")
	check(err)
	defer func() {
		err := os.RemoveAll(tempdir)
		if err == nil {
			fmt.Printf("cleaned up temp dir %s\n", tempdir)
		} else {
			fmt.Printf("error removing temp dir %q: %v\n", tempdir, err)
		}
	}()

	fmt.Printf("writing files to %s\n", tempdir)
	zipFile := filepath.Join(tempdir, "data.tar.gz")

	// Zip remote /data and pull to local machine
	pullZippedData(zipFile)

	// Extract tar archive
	err = exec.Command("tar", "-zxvf", zipFile, "-C", tempdir).Run()
	check(err)

	// Initialise a map to keep track of which files exist locally/remotely
	fileMap := make(map[string]fileInfo, 0)
	remoteDataDir := filepath.Join(tempdir, "data")

	// Put local files in fileMap
	filepath.WalkDir(remoteDataDir, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() {
			filename, err := filepath.Rel(remoteDataDir, path)
			check(err)
			fileMap[filename] = fileInfo{
				existsRemotely: true,
			}
		}
		return nil
	})

	// Put local files in fileMap
	filepath.WalkDir(st.dbPath, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() {
			filename, err := filepath.Rel(st.dbPath, path)
			check(err)
			fileMap[filename] = fileInfo{
				existsLocally:  true,
				existsRemotely: fileMap[filename].existsRemotely,
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
		if info.existsLocally && !info.existsRemotely {
			reportDiff("file %q exists locally but not remotely\n", filename)
			continue
		}
		if !info.existsLocally && info.existsRemotely {
			reportDiff("file %q exists remotely but not locally\n", filename)
			continue
		}

		// Check local file and remote file match
		localFileContents, err := os.ReadFile(filepath.Join(st.dbPath, filename))
		check(err)

		remoteFileContents, err := os.ReadFile(filepath.Join(remoteDataDir, filename))
		check(err)

		if strings.HasSuffix(filename, ".json") {
			// For JSON files, the exact formatting is not important.
			// Just check that the fields are the same.
			var localData, remoteData any
			err = json.Unmarshal(localFileContents, &localData)
			check(err)
			err = json.Unmarshal(remoteFileContents, &remoteData)
			check(err)

			if !reflect.DeepEqual(localData, remoteData) {
				reportDiff(`json data %q is different between local and remote
local data:  %#v
remote data: %#v
---------------------------------------
`, filename, localData, remoteData)
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
	} else {
		// Should return with exit code 1, but this will skip deferred calls
		// e.g. cleanup. So don't use os.Exit.
		// os.Exit(1)
	}
}

type fileInfo struct {
	existsLocally  bool
	existsRemotely bool
}

func pullZippedData(localZipFile string) {
	remoteZipFile := fmt.Sprintf("/tmp/data-%d.tar.gz", time.Now().Unix())

	// SSH in and zip /data for quick file transfer
	sshCmd := newSSHCommand("fly", "ssh", "console")
	sshCmd.exitSignal = "logout\n"
	sshCmd.Execf("tar -zcvf %s /data\n", remoteZipFile)
	check(sshCmd.Exit())

	// Transfer zip file via sftp
	sftpCmd := newSSHCommand("fly", "ssh", "sftp", "get", remoteZipFile, localZipFile)
	sftpCmd.exitSignal = "\x04"
	check(sftpCmd.Exit())

	// Delete temp file on VM
	rmCmd := newSSHCommand("fly", "ssh", "console")
	sshCmd.exitSignal = "logout\n"
	rmCmd.Execf("rm %s", remoteZipFile)
	// For some reason we have to force kill this, else it gets stuck.
	check(rmCmd.cmd.Cancel())
}

type sshCommand struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stderr *bytes.Buffer

	exitSignal string
}

func newSSHCommand(name string, args ...string) *sshCommand {
	cmd := exec.CommandContext(context.Background(), name, args...)
	stdin, err := cmd.StdinPipe()
	check(err)
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	check(cmd.Start())

	return &sshCommand{
		cmd:    cmd,
		stdin:  stdin,
		stderr: stderr,
	}
}

func (c *sshCommand) Execf(format string, v ...any) {
	io.WriteString(c.stdin, fmt.Sprintf(format, v...))
}

func (c *sshCommand) Exit() error {
	io.WriteString(c.stdin, c.exitSignal)
	fmt.Println("waiting for ssh command to exit")

	err := c.cmd.Wait()
	if err != nil {
		return errors.Join(err, fmt.Errorf(c.stderr.String()))
	}

	fmt.Println("ssh command exited successfully")
	return nil
}
