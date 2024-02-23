package namespaces

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ForceDeleteTerminatingNamespaces(ctx context.Context, c client.Client) error {
	namespaces := &corev1.NamespaceList{}
	if err := c.List(ctx, namespaces); err != nil {
		return err
	}
	for _, ns := range namespaces.Items {
		// Delete namespaces that are stuck in "Terminating" state
		if ns.Status.Phase == corev1.NamespaceTerminating {
			fmt.Printf("Deleting namespace %s...\n", ns.Name)
			patch := client.RawPatch(types.JSONPatchType, []byte(`[{"op": "remove", "path": "/metadata/finalizers"}]`))
			if err := c.Patch(ctx, &ns, patch); err != nil {
				fmt.Fprintf(os.Stderr, "Error patching namespace %s: %v\n", ns.Name, err)
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
