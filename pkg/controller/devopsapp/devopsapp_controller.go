/*
Copyright 2020 The KubeSphere Authors.

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

package devopsapp

import (
	"context"
	"strings"

	"encoding/json"

	"encoding/base64"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	devopsv1alpha3 "kubesphere.io/kubesphere/pkg/apis/devops/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new DevOpsApp Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDevOpsApp{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("devopsapp-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to DevOpsApp
	err = c.Watch(&source.Kind{Type: &devopsv1alpha3.DevOpsApp{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by DevOpsApp - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &devopsv1alpha3.DevOpsApp{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileDevOpsApp{}

type DockerConfigJson struct {
	Auths DockerConfigMap `json:"auths"`
}

type DockerConfigMap map[string]DockerConfigEntry

type DockerConfigEntry struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     string `json:"auth,omitempty"`
}

// ReconcileDevOpsApp reconciles a DevOpsApp object
type ReconcileDevOpsApp struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a DevOpsApp object and makes changes based on the state read
// and what is in the DevOpsApp.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=devops.kubesphere.io,resources=devopsapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=devops.kubesphere.io,resources=devopsapps/status,verbs=get;update;patch
func (r *ReconcileDevOpsApp) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the DevOpsApp instance
	rootCtx := context.Background()
	devopsapp := &devopsv1alpha3.DevOpsApp{}
	err := r.Get(rootCtx, request.NamespacedName, devopsapp)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	envs := devopsapp.Spec.Environments
	if envs == nil || len(envs) == 0 {
		return reconcile.Result{}, nil
	}

	devopsappName := devopsapp.ObjectMeta.GenerateName
	if devopsappName == "" {
		devopsappName = devopsapp.ObjectMeta.Name
	}

	for _, env := range envs {
		projectName := env.Name + "-" + devopsappName
		project := &corev1.Namespace{}
		if err = r.Get(rootCtx, types.NamespacedName{Name: projectName}, project); err != nil {
			return reconcile.Result{}, nil
		} else {
			harborSecret := &corev1.Secret{}
			if err = r.Get(rootCtx, types.NamespacedName{Name: "harbor", Namespace: projectName}, harborSecret); err != nil {
				registryUrl := devopsapp.Spec.Registry.Url

				rawAuth := devopsapp.Spec.Registry.Username + ":" + devopsapp.Spec.Registry.Password
				encodedAuth := make([]byte, base64.StdEncoding.EncodedLen(len(rawAuth)))
				base64.StdEncoding.Encode(encodedAuth, []byte(rawAuth))
				auth := string(encodedAuth)
				if strings.HasSuffix(auth, "==") {
					auth = auth[0 : len(auth)-2]
				}

				dockerConfigJson := &DockerConfigJson{Auths: DockerConfigMap{}}
				dockerConfigJson.Auths[registryUrl] = DockerConfigEntry{
					Username: devopsapp.Spec.Registry.Username,
					Password: devopsapp.Spec.Registry.Password,
					Email:    devopsapp.Spec.Registry.Email,
					Auth:     auth,
				}
				secretData, _ := json.Marshal(dockerConfigJson)

				harborSecret.Name = "harbor"
				harborSecret.Namespace = projectName
				harborSecret.Type = corev1.SecretTypeDockerConfigJson
				harborSecret.StringData = map[string]string{corev1.DockerConfigJsonKey: string(secretData)}
				if err = r.Create(rootCtx, harborSecret); err != nil {
					return reconcile.Result{}, err
				}
			} else {

			}
		}
	}

	return reconcile.Result{}, nil
}
