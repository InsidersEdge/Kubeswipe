package pods

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	errorsUtil "kubefit.com/kubeswipe/pkg/utils/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeleteAllPendingAndFailedPods(ctx context.Context, c client.Client) error {
	pods := &corev1.PodList{}
	if err := c.List(ctx, pods); err != nil {
		return err
	}

	var errors []error
	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case corev1.PodFailed, corev1.PodSucceeded: // Add PodSucceeded case since we don't want to keep successful pods
			err := c.Delete(ctx, &pod)
			if err != nil {
				errors = append(errors, err)
			}
		case corev1.PodPending:
			continue // Skip pending pods
		default:
			continue // Ignore other phases
		}
	}

	if len(errors) > 0 {
		errorsUtil.AggregateErrors(errors)
	}

	return nil
}
