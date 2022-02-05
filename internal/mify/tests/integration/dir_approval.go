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

	"github.com/otiai10/copy"
	"github.com/pmezard/go-difflib/difflib"
)

type approvalContext struct {
	subtestSeqNo int
	t            *testing.T
}

func NewApprovalContext(t *testing.T) approvalContext {
	return approvalContext{
		subtestSeqNo: 0,
		t:            t,
	}
}

func (ac approvalContext) getReceivedDir(subtestSeqNo int) string {
	return path.Join(getDataPath(ac.t), fmt.Sprintf("%s.%d.received", ac.t.Name(), subtestSeqNo))
}

func (ac approvalContext) getApprovedDir(subtestSeqNo int) string {
	return path.Join(getDataPath(ac.t), fmt.Sprintf("%s.%d.approved", ac.t.Name(), subtestSeqNo))
}

func (ac *approvalContext) NewSubtest() {
	path := ac.getReceivedDir(ac.subtestSeqNo)
	if err := os.RemoveAll(path); err != nil {
		ac.t.Fatalf("can't prepare received dir: %s", err)
	}

	if err := os.MkdirAll(path, fs.ModePerm); err != nil {
		ac.t.Fatalf("can't prepare received dir: %s", err)
	}
}

func (ac *approvalContext) EndSubtest(actualPath string) {
	receivedPath := ac.getReceivedDir(ac.subtestSeqNo)
	opts := copy.Options{
		Skip: func(src string) (bool, error) {
			return strings.Contains(src, "/.mify"), nil
		},
	}
	if err := copy.Copy(actualPath, receivedPath, opts); err != nil {
		ac.t.Fatalf("can't copy working dir to received: %s", err)
	}
	ac.subtestSeqNo++
}

func (ac *approvalContext) Verify() {
	success := true
	for i := 0; i < ac.subtestSeqNo; i++ {
		if err := verifyDirTree(ac.t, ac.getApprovedDir(i), ac.getReceivedDir(i)); err != nil {
			ac.t.Logf("subtest %d failed: %s", i, err)
			success = false
		}
	}
	if !success {
		ac.t.FailNow()
	}
}

func getDataPath(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("can't get working dir: %s", err)
	}

	return path.Join(wd, "data", t.Name())
}

func verifyDirTree(t *testing.T, approvedDirPath string, receivedDirPath string) error {
	if _, err := os.Stat(approvedDirPath); os.IsNotExist(err) {
		return fmt.Errorf("approved dir data wasn't found. Rename .received to .approved to make directory approved")
	}

	approvedDirTree, err := buildDirTree(approvedDirPath)
	if err != nil {
		return fmt.Errorf("can't get approved directory tree: %w", err)
	}

	receivedDirTree, err := buildDirTree(receivedDirPath)
	if err != nil {
		return fmt.Errorf("can't get received directory tree: %w", err)
	}

	diff := difflib.UnifiedDiff{
		A: difflib.SplitLines(approvedDirTree),
		B: difflib.SplitLines(receivedDirTree),
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	if len(text) > 0 {
		return fmt.Errorf("dir tree differs from approved one:\n%s", text)
	}

	return nil
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
