// gomuks - A terminal Matrix client written in Go.
// Copyright (C) 2022 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package filepicker

import (
	"bytes"
	"errors"
	"os/exec"
	"os"
	"strings"
)

var exists = []string{}
var cmdArgs = map[string]string{
	"zenity":	"--file-selection",
	"yad":		"--file",
	"kdialog":	"--getopenfilename",
}

func init() {
	var path string
	var dialogs = desktopDialogHierarchy()
	for _, prog := range dialogs {
		path, _ = exec.LookPath(prog)
		if len(path) > 0 {
			exists = append(exists, path)
		}
	}
}

func IsSupported() bool {
	return len(exists) > 0
}

func desktopDialogHierarchy() ([3]string) {
	var XDG_CURRENT_DESKTOP = strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
	var DESKTOP_SESSION = strings.ToLower(os.Getenv("DESKTOP_SESSION"))

	if XDG_CURRENT_DESKTOP == "kde" || DESKTOP_SESSION == "kde" {
		return [3]string{"kdialog", "yad", "zenity"}
	} else if XDG_CURRENT_DESKTOP == "lxqt" || DESKTOP_SESSION == "lxqt" {
		return [3]string{"kdialog", "yad", "zenity"}
	} else if XDG_CURRENT_DESKTOP == "gnome" || DESKTOP_SESSION == "gnome" {
		return [3]string{"zenity", "yad", "kdialog"}
	} else {
		return [3]string{"yad", "zenity", "kdialog"}
	}
}



func Open() (string, error) {
	var dialog []string = strings.Split(exists[0], "/")
	cmd := exec.Command(exists[0], cmdArgs[dialog[len(dialog)-1]])
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(output.String()), nil
}
