package integration

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
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
		ac.t.Logf("verifying results of subtest: %d ...", i)

		if _, err := os.Stat(ac.getApprovedTar(i)); os.IsNotExist(err) {
			ac.t.Logf("approved tar isn't exists. Validate .received results and run test with --aprove flag to mark results as approved")
			ac.t.FailNow()
		}

		if err := ac.unpackApprovedDir(i); err != nil {
			ac.t.Logf("failed to unpack approved tar: %s", err)
			ac.t.FailNow()
		}

		if err := verifyDirTree(ac.t, ac.getApprovedDir(i), ac.getReceivedDir(i)); err != nil {
			ac.t.Logf("subtest %d failed: %s", i, err)
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

func verifyDirTree(t *testing.T, approvedDirPath string, receivedDirPath string) error {
	fmt.Printf("%s %s\n", approvedDirPath, receivedDirPath)
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
		if strings.Contains(p, ".git/objects") {
			return nil
		}
		if d.IsDir() {
			// skip .git/objects
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
