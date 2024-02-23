package utils

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "kubefit.com/kubeswipe/api/v1"
	"kubefit.com/kubeswipe/pkg/utils/namespaces"
	"kubefit.com/kubeswipe/pkg/utils/services"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func HandleALLUnusedResources(ctx context.Context, client client.Client, cleaner v1.ResourceCleaner) error {
	logger := log.FromContext(ctx)
	if len(cleaner.Spec.Resources.Include) == 0 && len(cleaner.Spec.Resources.Exclude) == 0 {
		err := namespaces.ForceDeleteTerminatingNamespaces(ctx, client)
		if err != nil {
			logger.Error(err, "force deleting namespaces")
			return err
		}
		// if resource there then fetch and monitor them and apply logic
		err = services.HandleALLUnusedServices(ctx, client, cleaner)
		if err != nil {
			logger.Error(err, "handling services")
			return err
		}

	}
	return nil
}
