package workspace

import (
	"fmt"
	"io/ioutil"
	"path"
)

func (description Description) GetServiceList() ([]string, error) {
	schemasDir := path.Join(description.BasePath, description.GetSchemasRootRelPath())
	files, err := ioutil.ReadDir(schemasDir)
	if err != nil {
		return nil, fmt.Errorf("can't collect service list: %w", err)
	}

	res := make([]string, 0, len(files))
	for _, f := range files {
		res = append(res, f.Name())
	}

	return res, nil
}
