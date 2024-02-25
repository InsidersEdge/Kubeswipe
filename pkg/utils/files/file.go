package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	v1 "kubefit.com/kubeswipe/api/v1"
	errorsUtil "kubefit.com/kubeswipe/pkg/utils/errors"
)

func CreateFile(o interface{}, name string, dirName string, cleaner v1.ResourceCleaner) error {
	var errors []error

	podYAML, err := yaml.Marshal(o)
	if err != nil {
		errors = append(errors, err)
	}

	// Define the directory path
	swipeDir := v1.SwipeDIR + "/" + dirName
	if cleaner.Spec.Resources.BackupDir != "" {
		swipeDir = cleaner.Spec.Resources.BackupDir + "/" + dirName
	}
	path := filepath.Join(swipeDir, fmt.Sprintf("%s.yaml", name))

	// Check if the directory exists, create it if not
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		errors = append(errors, err)
	}

	// Write the YAML contents to a file
	if err := ioutil.WriteFile(path, podYAML, 0o644); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		errorsUtil.AggregateErrors(errors)
	}

	return nil
}
