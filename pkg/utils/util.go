package utils

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "kubefit.com/kubeswipe/api/v1"
	errorsUtil "kubefit.com/kubeswipe/pkg/utils/errors"
	"kubefit.com/kubeswipe/pkg/utils/namespaces"
	"kubefit.com/kubeswipe/pkg/utils/pods"
	"kubefit.com/kubeswipe/pkg/utils/services"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func HandleAllUnusedResources(ctx context.Context, client client.Client, cleaner v1.ResourceCleaner) error {
	logger := log.FromContext(ctx)
	if len(cleaner.Spec.Resources.Include) == 0 && len(cleaner.Spec.Resources.Exclude) == 0 {
		err := CleanAllResources(ctx, client, cleaner)
		if err != nil {
			logger.Error(err, "cleaning all resources")
			return err
		}

	} else {
		resourcesMap := make(map[string]bool)
		for _, includedResource := range cleaner.Spec.Resources.Include {
			resourcesMap[includedResource.Name] = true
		}
		for _, excludedResource := range cleaner.Spec.Resources.Exclude {
			resourcesMap[excludedResource.Name] = false
		}

		err := HandleUnusedResourcesInSteps(ctx, client, cleaner, resourcesMap)
		if err != nil {
			logger.Error(err, "handling unused resources")
			return err
		}

	}

	return nil
}

func HandleUnusedResourcesInSteps(ctx context.Context, client client.Client, cleaner v1.ResourceCleaner, resourceMap map[string]bool) error {
	logger := log.FromContext(ctx)
	for resourceName, included := range resourceMap {
		if resourceName == "Namespace" && included {
			err := namespaces.ForceDeleteTerminatingNamespaces(ctx, client,cleaner)
			if err != nil {
				logger.Error(err, "force deleting namespaces")
				return err
			}
		}
		if resourceName == "Service" && included {
			err := services.HandleAllUnusedServices(ctx, client, cleaner)
			if err != nil {
				logger.Error(err, "handling services")
				return err
			}
		}
		if resourceName == "Pod" && included {
			err := pods.DeleteAllPendingAndFailedPods(ctx, client, cleaner)
			if err != nil {
				logger.Error(err, "handling services")
				return err
			}
		}
	}
	return nil
}

func CleanAllResources(ctx context.Context, client client.Client, cleaner v1.ResourceCleaner) error {
	var errors []error

	err := namespaces.ForceDeleteTerminatingNamespaces(ctx, client, cleaner)
	if err != nil {
		errors = append(errors, err)
	}
	err = services.HandleAllUnusedServices(ctx, client, cleaner)
	if err != nil {
		errors = append(errors, err)
	}

	err = pods.DeleteAllPendingAndFailedPods(ctx, client, cleaner)
	if err != nil {
		errors = append(errors, err)
	}

	if cleaner.Spec.SwipePolicy == v1.Moderate {
		err = pods.DeleteAllUnusedPods(ctx, client,cleaner)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		errorsUtil.AggregateErrors(errors)
	}

	return nil
}
