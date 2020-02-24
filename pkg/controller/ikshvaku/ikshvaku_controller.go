package ikshvaku

import (
	"context"
	"sort"
	"time"

	cachev1alpha1 "github.com/ikshvaku/ikshvaku/pkg/apis/cache/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
        "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_ikshvaku")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Ikshvaku Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileIkshvaku{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ikshvaku-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Ikshvaku
	err = c.Watch(&source.Kind{Type: &cachev1alpha1.Ikshvaku{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Ikshvaku
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1alpha1.Ikshvaku{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileIkshvaku implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileIkshvaku{}

// ReconcileIkshvaku reconciles a Ikshvaku object
type ReconcileIkshvaku struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Ikshvaku object and makes changes based on the state read
// and what is in the Ikshvaku.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// an Ikshvaku deployment for each Ikshvaku CR
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileIkshvaku) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Ikshvaku")

	// Fetch the Ikshvaku instance
	ikshvaku := &cachev1alpha1.Ikshvaku{}
	err := r.client.Get(context.TODO(), request.NamespacedName, ikshvaku)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Ikshvaku resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Ikshvaku")
		return reconcile.Result{}, err
	}

	// Check if the deployment exists, if not create one
        found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ikshvaku.Name, Namespace: ikshvaku.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		//Define a new deployment
		dep := r.deploymentForIkshvaku(ikshvaku)
		reqLogger.Info("Creating a new deployent", "Deployment.namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			//Deployment creation failed
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		//Deployment created successfully
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get deployment")
		return reconcile.Result{}, nil
	}

	//Ensure the deployment size is same as the spec
	size := ikshvaku.Spec.ClusterSize
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return reconcile.Result{}, err
		}
	}

	
        podList := &corev1.PodList{}
        listOpts := []client.ListOption{
                client.InNamespace(ikshvaku.Namespace),
                client.MatchingLabels(labelsForIkshvaku(ikshvaku.Name)),
        }

        if err = r.client.List(context.TODO(), podList, listOpts...); err != nil {
		reqLogger.Error(err, "Failed to list pods", "Ikshvaku.Namespace", ikshvaku.Namespace, "Ikshvaku.Name", ikshvaku.Name)
		return reconcile.Result{}, err
	}

	// There's nothing more to do if the cluster size is lesser than or equal to 2
        if size <= 2 {
		return reconcile.Result{Requeue: true}, nil
	}
	running := runningPods(podList.Items)

	if len(running) != int(size) {
		return reconcile.Result{Requeue: true}, nil
	}

	sortPods(running)

	index := 0
	numMasters := 2 - (int(size) % 2)
	ikshvaku.Status.Masters = []string{}
	ikshvaku.Status.Slaves = []string{}
	for _, pod := range running {
		if index == 0 || (index == 1 && numMasters == 2) {
			pod.ObjectMeta.Labels["serverType"] = "master"
			ikshvaku.Status.Masters = append(ikshvaku.Status.Masters, pod.Name)
		} else if index < (int(size) - (numMasters % 2))/2 + 1 {
			pod.ObjectMeta.Labels["serverType"] = "cat1slave"
			ikshvaku.Status.Slaves = append(ikshvaku.Status.Slaves, pod.Name)
		} else {
			pod.ObjectMeta.Labels["serverType"] = "cat2slave"
			ikshvaku.Status.Slaves = append(ikshvaku.Status.Slaves, pod.Name)
		}
		r.client.Update(context.TODO(), &pod)
		index = index + 1
	}

	r.client.Status().Update(context.TODO(), ikshvaku)	

	//Spec has been updated, return and requeue
        return reconcile.Result{Requeue: true}, nil
}

func sortPods(pods []corev1.Pod) {
	sort.Slice(pods, func(i, j int) bool {
		return pods[i].CreationTimestamp.Before(&pods[j].CreationTimestamp)
	})
}

func convert(str string) time.Time {
        layout := "2006-01-02T15:04:05Z"
        t, err := time.Parse(layout, str)

        if err != nil { }
        return t
}

func runningPods(pods []corev1.Pod) (ret []corev1.Pod) {
	for _, pod := range pods {
		if pod.Status.Phase == "Running" {
			ret = append(ret, pod)
		}
	}
	return
}	


//deploymentForIkshvaku returns a ikshvaku Deployment object
func (r *ReconcileIkshvaku) deploymentForIkshvaku(m *cachev1alpha1.Ikshvaku) *appsv1.Deployment {
	ls := labelsForIkshvaku(m.Name)
	replicas := m.Spec.ClusterSize

	dep := &appsv1.Deployment {
		ObjectMeta: metav1.ObjectMeta{
			Name: m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:ls,
				},
				Spec: corev1.PodSpec {
					Containers: []corev1.Container{{
						Image: "nginx",
						Name: "nginx",
					}},
				},
			},
		},
	}
	// Set Ikshvaku instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

func labelsForIkshvaku(name string) map[string]string {
	return map[string]string{"app":"ikshvaku", "ikshvaku_cr": name}
}
