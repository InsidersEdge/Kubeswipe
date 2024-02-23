package services

import (
	"context"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Service struct {
	Name      string
	Namespace string
}

func getUnusedServicesInNamespace(ctx context.Context, client client.Client, namespace string) ([]Service, error) {
	endpointsList := v1.EndpointsList{}
	logger := log.FromContext(ctx)
	if err := client.List(context.TODO(), &endpointsList); err != nil {
		return nil, err
	}

	var unusedServices []Service
	for _, endpoints := range endpointsList.Items {
		if len(endpoints.Subsets) == 0 {
			logger.Info("unused service found in namespace: " + namespace + " with name: " + endpoints.Name + " and namespace: " + endpoints.Namespace)
			unusedServices = append(unusedServices, Service{
				Name:      endpoints.Name,
				Namespace: endpoints.Namespace,
			})
		}
	}

	return unusedServices, nil
}

func GetAllUnusedServices(ctx context.Context, client client.Client) ([]Service, error) {
	namespaces := &v1.NamespaceList{}
	if err := client.List(context.TODO(), namespaces); err != nil {
		return nil, err
	}

	var unusedServices []Service
	for _, ns := range namespaces.Items {
		nsServices, err := getUnusedServicesInNamespace(ctx, client, ns.Name)
		if err != nil {
			return nil, err
		}
		unusedServices = append(unusedServices, nsServices...)
	}

	return unusedServices, nil
}

func DeleteUnunsedServices(ctx context.Context, client client.Client, services []Service) error {
	logger := log.FromContext(ctx)
	for _, svc := range services {
		service := v1.Service{}
		err := client.Get(ctx, types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, &service)
		if err != nil {
			if apierrors.IsNotFound(err) {
				logger.Info("service " + service.Name + " not found")
			} else {
				logger.Error(err, "error getting service")
				return err
			}
		}

		err = client.Delete(ctx, &service)
		if err != nil {
			logger.Error(err, "error deleting service")
			return err
		}
		logger.Info("succesfully cleaned " + service.Name + " service")

	}
	return nil
}
