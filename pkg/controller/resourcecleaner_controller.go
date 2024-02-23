/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kubeswipev1 "kubefit.com/kubeswipe/api/v1"
	v1 "kubefit.com/kubeswipe/api/v1"
	"kubefit.com/kubeswipe/pkg/utils"
	"kubefit.com/kubeswipe/pkg/utils/services"
)

// ResourceCleanerReconciler reconciles a ResourceCleaner object
type ResourceCleanerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kubeswipe.kubefit.com,resources=resourcecleaners,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubeswipe.kubefit.com,resources=resourcecleaners/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubeswipe.kubefit.com,resources=resourcecleaners/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ResourceCleaner object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *ResourceCleanerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	cleaner := &v1.ResourceCleaner{}

	// if resources empty then just return
	err := r.Client.Get(ctx, req.NamespacedName, cleaner)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("cleaner not found")
		}
		logger.Error(err, "failed to get the cleaner resource")
	}

	err = utils.HandleALLUnusedResources(ctx, r.Client, *cleaner)

	if err != nil {
		logger.Error(err, "error handling unused resources")
	}

	// reconcile after some specified duration based on the schedule
	if cleaner.Spec.Schedule != "" {
		schedule, err := cron.ParseStandard(cleaner.Spec.Schedule)
		if err != nil {
			logger.Info("Can't parse the schedule")
		}

		next := schedule.Next(time.Now())
		duration := time.Until(next)
		fmt.Println("duration is", duration.Seconds())
		return ctrl.Result{RequeueAfter: duration}, nil
	}
	return ctrl.Result{RequeueAfter: time.Second * 60}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceCleanerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	logger := log.FromContext(context.Background())
	mux := http.NewServeMux()
	mux.HandleFunc("/getservice", r.GetServiceHandler)

	// Create a context with cancel function
	_, cancel := context.WithCancel(context.Background())

	// Enable CORS middleware
	handler := enableCORS(mux)

	// Start HTTP server in a Goroutine
	go func() {
		server := &http.Server{Addr: ":5000", Handler: handler}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err, "unable to start HTTP server")
			cancel() // Cancel context on error to stop the server
		}
	}()

	// Listen for SIGINT and SIGTERM signals to gracefully shut down the server
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeswipev1.ResourceCleaner{}).
		Complete(r)
}

// GetServiceHandler handles requests to /getservice
func (r *ResourceCleanerReconciler) GetServiceHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	unusedServices, err := services.GetAllUnusedServices(req.Context(), r.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(unusedServices)
}

// enableCORS is a middleware function to enable CORS
func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, req)
	})
}
