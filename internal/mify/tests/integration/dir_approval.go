package integration

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"log"
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
	ignoreFunc   func(path string) bool
}

func NewApprovalContext(t *testing.T) approvalContext {
	return approvalContext{
		subtestSeqNo: 0,
		t:            t,
	}
}

func (ac approvalContext) getReceivedDir(subtestSeqNo int) string {
	return path.Join(getResultsPath(ac.t), fmt.Sprintf("%s.%d.received", ac.t.Name(), subtestSeqNo))
}

func (ac approvalContext) getApprovedDir(subtestSeqNo int) string {
	return path.Join(getResultsPath(ac.t), fmt.Sprintf("%s.%d.approved", ac.t.Name(), subtestSeqNo))
}

func (ac approvalContext) getApprovedTar(subtestSeqNo int) string {
	return path.Join(getDataPath(ac.t), fmt.Sprintf("%s.%d.approved.tar", ac.t.Name(), subtestSeqNo))
}

// Pack approved dir to tar
func (ac approvalContext) packApprovedDir(subtestSeqNo int) error {
	approvedTarPath := ac.getApprovedTar(subtestSeqNo)
	if err := os.RemoveAll(approvedTarPath); err != nil {
		return err
	}

	if err := os.MkdirAll(path.Dir(approvedTarPath), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(approvedTarPath)
	if err != nil {
		return err
	}

	tw := tar.NewWriter(f)

	approvedDir := ac.getApprovedDir(subtestSeqNo)
	err = filepath.WalkDir(approvedDir, func(p string, d fs.DirEntry, _ error) error {
		relPath := strings.TrimPrefix(p, approvedDir)
		if d.IsDir() {
			hdr := &tar.Header{
				Typeflag: tar.TypeDir,
				Name:     relPath,
				Mode:     0600,
			}

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		} else {
			f, err := os.ReadFile(p)
			if err != nil {
				return err
			}

			hdr := &tar.Header{
				Name: relPath,
				Mode: 0600,
				Size: int64(len(f)),
			}

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			if _, err := tw.Write(f); err != nil {
				return err
			}
		}

		return nil
	})

	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		return err
	}

	return nil
}

// Unpack approved tar
func (ac approvalContext) unpackApprovedDir(subtestSeqNo int) error {
	approvedDirPath := ac.getApprovedDir(subtestSeqNo)
	if err := os.RemoveAll(approvedDirPath); err != nil {
		return err
	}

	f, err := os.Open(ac.getApprovedTar(subtestSeqNo))
	if err != nil {
		return err
	}

	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		p := path.Join(approvedDirPath, hdr.Name)
		if hdr.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(p, fs.ModePerm); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(path.Dir(p), fs.ModePerm); err != nil {
				return err
			}

			targetFile, err := os.Create(p)
			if err != nil {
				return err
			}

			if _, err := io.Copy(targetFile, tr); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ac *approvalContext) NewSubtest() {
	ac.t.Logf("start new subtest: %d", ac.subtestSeqNo)

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
			return strings.Contains(src, "/.mify") || strings.Contains(src, "/venv"), nil
		},
	}
	if err := copy.Copy(actualPath, receivedPath, opts); err != nil {
		ac.t.Fatalf("can't copy working dir to received: %s", err)
	}
	ac.subtestSeqNo++
}

func (ac *approvalContext) SetIgnoreFunc(ignoreFunc func(path string) bool) {
	ac.ignoreFunc = ignoreFunc
}

func (ac *approvalContext) Verify() {
	success := true
	for i := 0; i < ac.subtestSeqNo; i++ {
		ac.t.Logf("verifying results of subtest: %d ...", i)

		if _, err := os.Stat(ac.getApprovedTar(i)); os.IsNotExist(err) {
			ac.t.Logf("approved tar isn't exists. Validate .received results and run test with --aprove flag to mark results as approved")
			ac.t.FailNow()
		}

		if err := ac.unpackApprovedDir(i); err != nil {
			ac.t.Logf("failed to unpack approved tar: %s", err)
			ac.t.FailNow()
		}

		diff, err := diffDir(ac, ac.getApprovedDir(i), ac.getReceivedDir(i))
		if err != nil {
			ac.t.Logf("failed to calc diff: %s", err)
			ac.t.FailNow()
		}

		if len(diff) > 0 {
			ac.t.Logf("subtest %d failed. Unapproved changes were found: %s", i, diff)
			success = false
			continue
		}

		ac.t.Logf("subtest %d successed", i)
	}
	if !success {
		ac.t.FailNow()
	}
}

func (ac *approvalContext) Approve() {
	for i := 0; i < ac.subtestSeqNo; i++ {
		approvedPath := ac.getApprovedDir(i)
		if err := os.RemoveAll(approvedPath); err != nil {
			ac.t.Logf("approve %d failed: %s", i, err)
			ac.t.FailNow()
		}

		if err := copy.Copy(ac.getReceivedDir(i), approvedPath); err != nil {
			ac.t.Logf("approve %d failed: %s", i, err)
			ac.t.FailNow()
		}

		if err := ac.packApprovedDir(i); err != nil {
			ac.t.Logf("approve pack %d failed: %s", i, err)
			ac.t.FailNow()
		}
	}
}

func getDataPath(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("can't get working dir: %s", err)
	}

	return path.Join(wd, "data", t.Name())
}

func getResultsPath(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("can't get results dir: %s", err)
	}

	return path.Join(wd, "results", t.Name())
}

func shouldIgnorePath(ac *approvalContext, path string, dir fs.DirEntry) (bool, error) {
	// skip .git
	if strings.Contains(path, ".git") {
		return true, nil
	}
	if strings.Contains(path, "venv/") {
		return true, nil
	}

	if ac.ignoreFunc != nil && ac.ignoreFunc(path) {
		return true, nil
	}

	return dir.IsDir(), nil
}

func diffReceivedAndApproved(p string, aprovedDirPath string, receivedDirPath string) (string, error) {
	relPath := ""
	if strings.HasPrefix(p, aprovedDirPath) {
		relPath = strings.TrimPrefix(p, aprovedDirPath)
	} else {
		relPath = strings.TrimPrefix(p, receivedDirPath)
	}

	receivedFilePath := path.Join(receivedDirPath, relPath)
	receivedContent, err := os.ReadFile(receivedFilePath)
	if err != nil {
		return "", err
	}

	approvedFilePath := path.Join(aprovedDirPath, relPath)
	approvedContent, err := os.ReadFile(approvedFilePath)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(receivedContent)),
		B:        difflib.SplitLines(string(approvedContent)),
		FromFile: receivedFilePath,
		ToFile:   approvedFilePath,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)

	return text, nil
}

func diffDir(ac *approvalContext, aprovedDirPath string, receivedDirPath string) (string, error) {
	checkedFilePaths := make(map[string]struct{})
	res := ""

	handlePath := func(p string, d fs.DirEntry, _ error) error {
		if _, ok := checkedFilePaths[p]; ok {
			return nil
		}

		ignore, err := shouldIgnorePath(ac, p, d)
		if err != nil {
			return err
		}
		if ignore {
			return nil
		}

		diff, err := diffReceivedAndApproved(p, aprovedDirPath, receivedDirPath)
		if err != nil {
			return err
		}

		if len(diff) > 0 {
			res += diff
			res += "\n"
		}

		checkedFilePaths[p] = struct{}{}

		return nil
	}

	err := filepath.WalkDir(aprovedDirPath, handlePath)
	if err != nil {
		return "", err
	}
	err = filepath.WalkDir(receivedDirPath, handlePath)
	if err != nil {
		return "", err
	}

	return res, nil
}
