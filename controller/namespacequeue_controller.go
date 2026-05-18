package controller

import (
	context "context"
	"strings"

	v1alpha1 "github.com/pushtisonawala/namespace-queue-poc/api/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type NamespaceQueueReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const namespaceQueueFinalizer = "scheduling.volcano.sh/namespacequeue-finalizer"

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceQueueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.NamespaceQueue{}).
		Complete(r)
}

func (r *NamespaceQueueReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling NamespaceQueue",
		"namespace", req.Namespace,
		"name", req.Name,
	)

	var nsq v1alpha1.NamespaceQueue

	if err := r.Get(ctx, req.NamespacedName, &nsq); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("NamespaceQueue not found, ignoring")
			return ctrl.Result{}, nil
		}

		logger.Error(err, "Failed to get NamespaceQueue")
		return ctrl.Result{}, err
	}

	if nsq.Spec.ParentQueue == "" {
		logger.Info("ParentQueue field is empty")

		nsq.Status.State = "Unknown"
		_ = r.Status().Update(ctx, &nsq)

		return ctrl.Result{}, nil
	}

	// Check if parent cluster Queue exists (simulate with annotation for PoC)
	parentQueueName := nsq.Spec.ParentQueue
	parentFound := false

	// In a real implementation, query the cluster for the parent Queue CR
	if strings.ToLower(parentQueueName) == "default" {
		parentFound = true
	}

	if parentFound {
		logger.Info("Parent Queue found", "parentQueue", parentQueueName)
		nsq.Status.State = "Open"
	} else {
		logger.Info("Parent Queue not found", "parentQueue", parentQueueName)
		nsq.Status.State = "Unknown"
	}

	_ = r.Status().Update(ctx, &nsq)

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(&nsq, namespaceQueueFinalizer) {
		controllerutil.AddFinalizer(&nsq, namespaceQueueFinalizer)

		if err := r.Update(ctx, &nsq); err != nil {
			logger.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}

		logger.Info("Finalizer added")
	}

	return ctrl.Result{}, nil
}
