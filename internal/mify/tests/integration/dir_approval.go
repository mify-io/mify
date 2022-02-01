package integration

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func CreateReceivedDir(t *testing.T) string {
	path := getReceivedDir(t)
	if err := os.RemoveAll(path); err != nil {
		t.Fatalf("can't prepare received dir: %s", err)
	}

	if err := os.MkdirAll(path, fs.ModePerm); err != nil {
		t.Fatalf("can't prepare received dir: %s", err)
	}

	return path
}

func VerifyWithApproved(t *testing.T) {
	verifyDir(t, getApprovedDir(t), getReceivedDir(t))
}

func getDataPath(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("can't get working dir: %s", err)
	}

	return path.Join(wd, "data")
}

func getReceivedDir(t *testing.T) string {
	return path.Join(getDataPath(t), fmt.Sprintf("%s.received", t.Name()))
}

func getApprovedDir(t *testing.T) string {
	return path.Join(getDataPath(t), fmt.Sprintf("%s.approved", t.Name()))
}

func verifyDir(t *testing.T, approvedDirPath string, receivedDirPath string) {
	verifyDirTree(t, approvedDirPath, receivedDirPath)
}

func verifyDirTree(t *testing.T, approvedDirPath string, receivedDirPath string) {
	approvedDirTree, err := buildDirTree(approvedDirPath)
	if err != nil {
		t.Fatalf("can't get approved directory tree: %s", err)
	}

	if _, err := os.Stat(approvedDirPath); os.IsNotExist(err) {
		t.Fatalf("Approved dir data wasn't found. Rename .received to .approved to make directory approved")
	}

	receivedDirTree, err := buildDirTree(receivedDirPath)
	if err != nil {
		t.Fatalf("can't get received directory tree: %s", err)
	}

	diff := difflib.UnifiedDiff{
		A: difflib.SplitLines(approvedDirTree),
		B: difflib.SplitLines(receivedDirTree),
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	if len(text) > 0 {
		t.Fatalf("Dir tree differs from approved one:\n%s", text)
	}
}

func buildDirTree(path string) (string, error) {
	res := ""
	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			files, err := ioutil.ReadDir(p)
			if err != nil {
				return err
			}

			// Ignore empty dirs, because git removes them from approved
			if len(files) == 0 {
				return nil
			}
		}

		res += fmt.Sprintf("%s\n", strings.TrimPrefix(p, path))
		return nil
	})

	if err != nil {
		return "", err
	}

	return res, nil
}