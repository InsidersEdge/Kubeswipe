package services

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	v1 "kubefit.com/kubeswipe/api/v1"
	errorsUtil "kubefit.com/kubeswipe/pkg/utils/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Service struct {
	Name      string
	Namespace string
}

func getUnusedServicesInNamespace(ctx context.Context, c client.Client, namespace string, operation string) ([]Service, error) {
	var errors []error
	endpointsList := corev1.EndpointsList{}
	logger := log.FromContext(ctx)
	if err := c.List(ctx, &endpointsList, &client.ListOptions{Namespace: namespace}); err != nil {
		return nil, err
	}

	var unusedServices []Service
	for _, endpoints := range endpointsList.Items {
		if len(endpoints.Subsets) == 0 {
			logger.Info("unused service found in namespace: " + namespace + " with name: " + endpoints.Name + " and namespace: " + endpoints.Namespace)
			if operation == string(v1.CleanUp) {
				service := corev1.Service{}
				err := c.Get(ctx, types.NamespacedName{Name: endpoints.Name, Namespace: endpoints.Namespace}, &service)
				if err != nil {
					if apierrors.IsNotFound(err) {
						logger.Info("service " + service.Name + " not found")
					} else {
						errors = append(errors, err)
					}
				}

				err = c.Delete(ctx, &service)
				if err != nil {
					errors = append(errors, err)
				}
				continue
			}
			unusedServices = append(unusedServices, Service{
				Name:      endpoints.Name,
				Namespace: endpoints.Namespace,
			})
		}
	}

	if len(errors) > 0 {
		errorsUtil.AggregateErrors(errors)
	}

	return unusedServices, nil
}

func GetAllUnusedServices(ctx context.Context, c client.Client) ([]Service, error) {
	var errors []error
	namespaces := &corev1.NamespaceList{}
	if err := c.List(context.TODO(), namespaces); err != nil {
		return nil, err
	}

	var unusedServices []Service
	for _, ns := range namespaces.Items {
		nsServices, err := getUnusedServicesInNamespace(ctx, c, ns.Name, string(v1.Serve))
		if err != nil {
			if len(errors) > 0 {
				errorsUtil.AggregateErrors(errors)
			}
		}
		unusedServices = append(unusedServices, nsServices...)
	}

	if len(errors) > 0 {
		errorsUtil.AggregateErrors(errors)
	}

	return unusedServices, nil
}

func HandleAllUnusedServices(ctx context.Context, c client.Client, cleaner v1.ResourceCleaner) error {
	var errors []error
	namespaces := &corev1.NamespaceList{}
	if err := c.List(context.TODO(), namespaces); err != nil {
		return err
	}

	var unusedServices []Service
	for _, ns := range namespaces.Items {
		nsServices, err := getUnusedServicesInNamespace(ctx, c, ns.Name, string(cleaner.Spec.Operation))
		if err != nil {
			errors = append(errors, err)
		}
		unusedServices = append(unusedServices, nsServices...)
	}

	if len(errors) > 0 {
		errorsUtil.AggregateErrors(errors)
	}

	return nil
}

func DeleteUnunsedServices(ctx context.Context, c client.Client, services []Service) error {
	logger := log.FromContext(ctx)
	var errors []error
	for _, svc := range services {
		service := corev1.Service{}
		err := c.Get(ctx, types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, &service)
		if err != nil {
			if apierrors.IsNotFound(err) {
				logger.Info("service " + service.Name + " not found")
			} else {
				errors = append(errors, err)
			}
		}

		err = c.Delete(ctx, &service)
		if err != nil {
			errors = append(errors, err)
		}
		if len(errors) > 0 {
			errorsUtil.AggregateErrors(errors)
		}
		logger.Info("succesfully cleaned " + service.Name + " service")

	}
	return nil
}
