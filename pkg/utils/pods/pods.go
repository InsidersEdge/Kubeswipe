package pods

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	v1 "kubefit.com/kubeswipe/api/v1"
	errorsUtil "kubefit.com/kubeswipe/pkg/utils/errors"
	filesUtil "kubefit.com/kubeswipe/pkg/utils/files"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	updateInterval           = 6 * time.Hour
	deletionThreshold        = 20
	annotationKey            = "last_cpu_usage_time"
	cpuAnnotationKey         = "cpu_usage"
	updateCountAnnotationKey = "update_count"
)

func DeleteAllPendingAndFailedPods(ctx context.Context, c client.Client, cleaner v1.ResourceCleaner) error {
	fmt.Println("pending and failed pods")
	pods := &corev1.PodList{}
	if err := c.List(ctx, pods); err != nil {
		return err
	}

	var errors []error
	for _, pod := range pods.Items {
		fmt.Println("pod status phase", pod.Status.Phase)
		switch pod.Status.Phase {
		case corev1.PodFailed, corev1.PodSucceeded: // Add PodSucceeded case since we don't want to keep successful pods
			if cleaner.Spec.Resources.Backup {
				err := filesUtil.CreateFile(pod, "pods", pod.Name, cleaner)
				if err != nil {
					errors = append(errors, err)
				}
			}

			err := c.Delete(ctx, &pod)
			if err != nil {
				errors = append(errors, err)
			}
		case corev1.PodPending:
			continue // Skip pending pods

		}

		for _, status := range pod.Status.ContainerStatuses {
			if !status.Ready {
				if cleaner.Spec.Resources.Backup {
					err := filesUtil.CreateFile(pod, "pods", pod.Name, cleaner)
					if err != nil {
						errors = append(errors, err)
					}
				}
				err := c.Delete(ctx, &pod)
				if err != nil {
					errors = append(errors, err)
				}
				fmt.Println("deleted the pod", pod.Name)
				continue
			}
		}

	}

	if len(errors) > 0 {
		errorsUtil.AggregateErrors(errors)
	}

	return nil
}

func DeleteAllUnusedPods(ctx context.Context, c client.Client, cleaner v1.ResourceCleaner) error {
	namespaces := &corev1.NamespaceList{}
	var errors []error
	config := config.GetConfigOrDie()
	mc, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	err = c.List(ctx, namespaces)
	if err != nil {
		return err
	}

	for _, ns := range namespaces.Items {
		podMetrics, err := mc.MetricsV1beta1().PodMetricses(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Println("Error fetching the metrics:", err)
			return err
		}

		for _, po := range podMetrics.Items {
			podContainers := po.Containers

			totalCpu := 0
			for _, container := range podContainers {
				cpuQuantity, _ := container.Usage.Cpu().AsInt64()
				totalCpu += int(cpuQuantity)
			}

			if totalCpu == 0 {

				pod := &corev1.Pod{}
				err := c.Get(ctx, client.ObjectKey{Name: po.Name, Namespace: ns.Name}, pod)
				if err != nil {
					fmt.Printf("Error getting pod %s: %v\n", po.Name, err)
					continue
				}

				annotations := pod.GetAnnotations()
				if annotations == nil {
					annotations = make(map[string]string)
				}

				lastUpdateTime, err := time.Parse(time.RFC3339, annotations[annotationKey])
				if err != nil {
					// If annotation doesn't exist or has invalid time format, create a new one
					annotations[annotationKey] = time.Now().Format(time.RFC3339)
					annotations[cpuAnnotationKey] = fmt.Sprintf("%d", totalCpu)
					annotations[updateCountAnnotationKey] = "1"
					pod.SetAnnotations(annotations)
					err = c.Update(ctx, pod)
					if err != nil {
						return nil
					}
				}

				updateCountStr, ok := annotations[updateCountAnnotationKey]
				if !ok {
					updateCountStr = "0"
				}

				updateCount, err := strconv.Atoi(updateCountStr)
				if err != nil {
					fmt.Printf("Error converting update count to integer: %v\n", err)
					continue
				}

				// Check if update interval has passed
				if time.Since(lastUpdateTime) >= updateInterval {
					lastCPUUsage, err := strconv.Atoi(annotations[cpuAnnotationKey])
					if err != nil {
						fmt.Printf("Error converting last CPU usage to integer: %v\n", err)
						continue
					}

					cpuDifference := int64(math.Abs(float64(int64(totalCpu) - int64(lastCPUUsage))))

					// If CPU usage is the same or increased, update the annotation
					if cpuDifference < 3 {
						annotations[annotationKey] = time.Now().Format(time.RFC3339)
						annotations[cpuAnnotationKey] = fmt.Sprintf("%d", int64(totalCpu))
						annotations[updateCountAnnotationKey] = strconv.Itoa(updateCount + 1)
						pod.SetAnnotations(annotations)
						err = c.Update(ctx, pod)
						if err != nil {
							return nil
						}
					} else {
						// means this pod has some variations in cpu usuage and is not right candidate to be deleted
						delete(annotations, annotationKey)
						delete(annotations, cpuAnnotationKey)
						delete(annotations, updateCountAnnotationKey)
						pod.SetAnnotations(annotations)
						err = c.Update(ctx, pod)
						if err != nil {
							return nil
						}
					}

					// If update count exceeds deletion threshold, delete the pod
					if updateCount >= deletionThreshold {
						if cleaner.Spec.Resources.Backup {
							err := filesUtil.CreateFile(pod, "pods", pod.Name, cleaner)
							if err != nil {
								errors = append(errors, err)
							}
						}
						err = c.Delete(ctx, pod)
						if err != nil {
							return err
						}
					}
				}

			}

		}
	}

	if len(errors) > 0 {
		errorsUtil.AggregateErrors(errors)
	}

	return nil
}
