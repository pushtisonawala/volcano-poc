package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	v1alpha1 "github.com/pushtisonawala/namespace-queue-poc/api/v1alpha1"
	"github.com/pushtisonawala/namespace-queue-poc/controller"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	ctrl.SetLogger(logr.Discard())

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	restCfg := config.GetConfigOrDie()
	mgr, err := manager.New(restCfg, manager.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Fatalf("unable to start manager: %v", err)
	}

	reconciler := &controller.NamespaceQueueReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
	// Add SetupWithManager stub if not present
	if err := reconciler.SetupWithManager(mgr); err != nil {
		log.Fatalf("unable to create controller: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("starting manager")
	if err := mgr.Start(ctx); err != nil {
		log.Fatalf("problem running manager: %v", err)
	}
}
