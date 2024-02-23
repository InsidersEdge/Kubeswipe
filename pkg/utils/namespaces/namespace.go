package namespaces

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ForceDeleteTerminatingNamespaces(ctx context.Context, c client.Client) error {
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
				fmt.Printf("Error executing kubectl patch: %v\n", err)
			} else {
				fmt.Printf("Namespace %s patched successfully\n", ns.Name)
			}
			if err := c.Delete(ctx, &ns); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting namespace %s: %v\n", ns.Name, err)
			} else {
				fmt.Printf("Namespace %s deleted successfully\n", ns.Name)
			}
		}
	}
	return nil
}
