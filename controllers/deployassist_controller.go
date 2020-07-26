/*


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

package controllers

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	deployassistv1alpha1 "deployassist/api/v1alpha1"
)

var (
	mgr manager.Manager
)

// DeployassistReconciler reconciles a Deployassist object
type DeployassistReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	Mgr         manager.Manager
	Controllers []CrController
}

type CrController struct {
	Controller controller.Controller
	Stop       chan struct{}
	Name       string
}

// +kubebuilder:rbac:groups=deployassist.apiextensions.k8s.io,resources=deployassists,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=deployassist.apiextensions.k8s.io,resources=deployassists/status,verbs=get;update;patch

func (r *DeployassistReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("deployassist", req.NamespacedName)
	// your logic her

	var deployAssist deployassistv1alpha1.Deployassist
	if err := r.Get(ctx, req.NamespacedName, &deployAssist); err != nil {
		log.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	myFinalizerName := "stopRemoveControllers"
	// examine DeletionTimestamp to determine if object is under deletion
	if deployAssist.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(deployAssist.ObjectMeta.Finalizers, myFinalizerName) {
			deployAssist.ObjectMeta.Finalizers = append(deployAssist.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), &deployAssist); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(deployAssist.ObjectMeta.Finalizers, myFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteExternalResources(&deployAssist); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			deployAssist.ObjectMeta.Finalizers = removeString(deployAssist.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), &deployAssist); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}
	fmt.Println("hi")
	fmt.Println(deployAssist)
	fmt.Println("REQUEST NAMESPACE")
	fmt.Println(r.Controllers)
	if len(r.Controllers) > 0 {
		r.Controllers = r.Controllers[:len(r.Controllers)-1]
	}
	fmt.Println(r.Controllers)
	r.Controllers = append(r.Controllers, NewUnmanagedController(r.Mgr, deployAssist.Name))

	return ctrl.Result{}, nil
}

func (r *DeployassistReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Mgr = mgr
	return ctrl.NewControllerManagedBy(mgr).
		For(&deployassistv1alpha1.Deployassist{}).Complete(r)
}

func (r *DeployassistReconciler) deleteExternalResources(deployAssist *deployassistv1alpha1.Deployassist) error {

	for i, controller := range r.Controllers {
		if controller.Name == deployAssist.Name {
			close(controller.Stop)
			r.Controllers = append(r.Controllers[:i], r.Controllers[i+1:]...)
		}

	}
	return nil
}

// This example creates a new controller named "pod-controller" to watch Pods
// and call a no-op reconciler. The controller is not added to the provided
// manager, and must thus be started and stopped by the caller.
func NewUnmanagedController(mgr manager.Manager, name string) CrController {
	// mgr is a manager.Manager

	// Configure creates a new controller but does not add it to the supplied
	// manager.
	fmt.Println("hfds")

	c, err := controller.NewUnmanaged("pod-controller", mgr, controller.Options{
		Reconciler: reconcile.Func(func(_ reconcile.Request) (reconcile.Result, error) {
			return reconcile.Result{}, nil
		}),
	})
	if err != nil {
		log.Error(err, "unable to create pod-controller")
		os.Exit(1)
	}

	if err := c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{}); err != nil {
		log.Error(err, "unable to watch pods")
		os.Exit(1)
	}

	// Create a stop channel for our controller. The controller will stop when
	// this channel is closed.
	stop := make(chan struct{})

	// Start our controller in a goroutine so that we do not block.
	go func() {
		// Block until our controller manager is elected leader. We presume our
		// entire process will terminate if we lose leadership, so we don't need
		// to handle that.

		// Elected not available? maybe in newer version of
		<-mgr.Elected()

		// Start our controller. This will block until the stop channel is
		// closed, or the controller returns an error.
		if err := c.Start(stop); err != nil {
			log.Error(err, "cannot run experiment controller")
		}
	}()

	return CrController{
		Controller: c,
		Stop:       stop,
		Name:       name,
	}
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
