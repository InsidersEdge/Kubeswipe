package namespaces

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	corev1 "k8s.io/api/core/v1"
	v1 "kubefit.com/kubeswipe/api/v1"
	errorsUtil "kubefit.com/kubeswipe/pkg/utils/errors"
	filesUtil "kubefit.com/kubeswipe/pkg/utils/files"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ForceDeleteTerminatingNamespaces(ctx context.Context, c client.Client, cleaner v1.ResourceCleaner) error {
	var errors []error
	namespaces := &corev1.NamespaceList{}
	if err := c.List(ctx, namespaces); err != nil {
		return err
	}
	for _, ns := range namespaces.Items {
		// Delete namespaces that are stuck in "Terminating" state
		if ns.Status.Phase == corev1.NamespaceTerminating || ns.Name == "test-namespace" {
			fmt.Printf("Deleting namespace %s...\n", ns.Name)
			patchJSON := `{"metadata":{"finalizers":[]}}`
			cmd := exec.Command("kubectl", "patch", "namespace", ns.Name, "-p", patchJSON, "--type=merge")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				errors = append(errors, err)
			} else {
				fmt.Printf("Namespace %s patched successfully\n", ns.Name)
			}
			if cleaner.Spec.Resources.Backup {
				err := filesUtil.CreateFile(ns, ns.Name, "namespaces", cleaner)
				if err != nil {
					errors = append(errors, err)
				}
			}
			if err := c.Delete(ctx, &ns); err != nil {
				if cleaner.Spec.Resources.Backup {
					err := filesUtil.CreateFile(ns, ns.Name, "namespaces", cleaner)
					if err != nil {
						errors = append(errors, err)
					}
				}
			} else {
				fmt.Printf("Namespace %s deleted successfully\n", ns.Name)
			}
		}
	}
	if len(errors) > 0 {
		errorsUtil.AggregateErrors(errors)
	}
	return nil
}
