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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	devopsv1alpha3 "kubesphere.io/kubesphere/pkg/apis/devops/v1alpha3"
	devopsClient "kubesphere.io/kubesphere/pkg/simple/client/devops"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("devopsapp_controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new DevOpsApp Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, devopsClient devopsClient.Interface) error {
	return add(mgr, newReconciler(mgr, devopsClient))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, devopsClient devopsClient.Interface) reconcile.Reconciler {
	return &ReconcileDevOpsApp{Client: mgr.GetClient(), scheme: mgr.GetScheme(), devopsClient: devopsClient}
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
	scheme       *runtime.Scheme
	devopsClient devopsClient.Interface
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
			continue
		}

		if err = checkHarborSecret(r, rootCtx, devopsapp.Spec.Registry, env, projectName); err != nil {
			return reconcile.Result{}, err
		}

		if err = checkDevOpsCredentials(r, rootCtx, devopsapp, env); err != nil {
			return reconcile.Result{}, err
		}

	}

	return reconcile.Result{}, nil
}

func checkHarborSecret(r *ReconcileDevOpsApp, ctx context.Context, registry *devopsv1alpha3.Registry, env devopsv1alpha3.Environment, namespace string) error {
	envRegistry := &devopsv1alpha3.Registry{
		Url:      registry.Url,
		Username: registry.Username,
		Password: registry.Password,
		Email:    registry.Email,
	}
	if env.Registry != nil {
		envRegistry.Url = env.Registry.Url
		envRegistry.Username = env.Registry.Username
		envRegistry.Password = env.Registry.Password
		envRegistry.Email = env.Registry.Email
	}

	secretName := "harbor"
	secretInstance := &corev1.Secret{}
	rawAuth := envRegistry.Username + ":" + envRegistry.Password
	encodedAuth := make([]byte, base64.StdEncoding.EncodedLen(len(rawAuth)))
	base64.StdEncoding.Encode(encodedAuth, []byte(rawAuth))
	auth := string(encodedAuth)
	if strings.HasSuffix(auth, "==") {
		auth = auth[0 : len(auth)-2]
	}
	if err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: namespace}, secretInstance); err != nil {
		dockerConfigJson := &DockerConfigJson{Auths: DockerConfigMap{}}
		dockerConfigJson.Auths[envRegistry.Url] = DockerConfigEntry{
			Username: envRegistry.Username,
			Password: envRegistry.Password,
			Email:    envRegistry.Email,
			Auth:     auth,
		}
		secretData, _ := json.Marshal(dockerConfigJson)

		secretInstance.Name = secretName
		secretInstance.Namespace = namespace
		secretInstance.Type = corev1.SecretTypeDockerConfigJson
		secretInstance.StringData = map[string]string{corev1.DockerConfigJsonKey: string(secretData)}
		if err = r.Create(ctx, secretInstance); err != nil {
			return err
		} else {
			return nil
		}
	}

	secretData := secretInstance.Data[corev1.DockerConfigJsonKey]
	dockerConfigJson := &DockerConfigJson{Auths: DockerConfigMap{}}
	json.Unmarshal(secretData, dockerConfigJson)
	dockerConfigEntry := dockerConfigJson.Auths[envRegistry.Url]
	if dockerConfigEntry.Username == envRegistry.Username && dockerConfigEntry.Password == envRegistry.Password {
		return nil
	}

	dockerConfigJson.Auths = DockerConfigMap{}
	dockerConfigJson.Auths[envRegistry.Url] = DockerConfigEntry{
		Username: envRegistry.Username,
		Password: envRegistry.Password,
		Email:    envRegistry.Email,
		Auth:     auth,
	}
	secretDataBytes, _ := json.Marshal(dockerConfigJson)
	secretInstance.StringData = map[string]string{corev1.DockerConfigJsonKey: string(secretDataBytes)}
	if err := r.Update(ctx, secretInstance); err != nil {
		return err
	}
	return nil
}

func checkDevOpsCredentials(r *ReconcileDevOpsApp, ctx context.Context, devopsapp *devopsv1alpha3.DevOpsApp, env devopsv1alpha3.Environment) error {
	devopsappName := devopsapp.Name
	workspace := devopsapp.Labels["kubesphere.io/workspace"]

	devOpsProjectList := &devopsv1alpha3.DevOpsProjectList{}
	opts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{"kubesphere.io/workspace": workspace}),
	}
	if err := r.List(ctx, devOpsProjectList, opts); err != nil || len(devOpsProjectList.Items) == 0 {
		return nil
	}

	devOpsProject := devopsv1alpha3.DevOpsProject{}
	for _, project := range devOpsProjectList.Items {
		if devopsappName == project.GenerateName {
			devOpsProject = project
			break
		}
	}
	if devOpsProject.Name == "" {
		return nil
	}

	devopsNamespace := devOpsProject.Name
	if err := doCheckBasicAuthCredential(r, ctx, devopsNamespace, "git", devopsapp.Spec.Git.Username, devopsapp.Spec.Git.Password); err != nil {
		return err
	}
	if err := doCheckBasicAuthCredential(r, ctx, devopsNamespace, "harbor", devopsapp.Spec.Registry.Username, devopsapp.Spec.Registry.Password); err != nil {
		return err
	}

	if len(env.Credentials) == 0 {
		return nil
	}

	for _, credential := range env.Credentials {
		var err error
		if credential.Type == "basic-auth" {
			err = doCheckBasicAuthCredential(r, ctx, devopsNamespace, credential.Name, credential.Username, credential.Password)
		} else if credential.Type == "kubeconfig" {
			err = doCheckKubeConfigCredential(r, ctx, devopsNamespace, credential.Name, credential.Content)
		} else if credential.Type == "secret-text" {
			err = doCheckSecretTextCredential(r, ctx, devopsNamespace, credential.Name, credential.Secret)
		} else if credential.Type == "ssh-auth" {
			err = doCheckSshAuthCredential(r, ctx, devopsNamespace, credential.Name, credential.Username, credential.Password, credential.PrivateKey)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func doCheckBasicAuthCredential(r *ReconcileDevOpsApp, ctx context.Context, namespace string, name string, username string, password string) error {
	if _, err := r.devopsClient.GetCredentialInProject(namespace, name); err != nil {
		secret := &corev1.Secret{
			Type: devopsv1alpha3.SecretTypeBasicAuth,
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: map[string][]byte{
				"username": []byte(username),
				"password": []byte(password),
			},
		}
		return r.Create(ctx, secret)
	}

	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret); err != nil {
		return nil
	}

	if string(secret.Data["username"]) != username || string(secret.Data["password"]) != password {
		secret.Data = map[string][]byte{
			"username": []byte(username),
			"password": []byte(password),
		}
		return r.Update(ctx, secret)
	}
	return nil
}

func doCheckKubeConfigCredential(r *ReconcileDevOpsApp, ctx context.Context, namespace string, name string, content string) error {
	if _, err := r.devopsClient.GetCredentialInProject(namespace, name); err != nil {
		secret := &corev1.Secret{
			Type: devopsv1alpha3.SecretTypeKubeConfig,
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: map[string][]byte{
				"content": []byte(content),
			},
		}
		return r.Create(ctx, secret)
	}

	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret); err != nil {
		return nil
	}

	if string(secret.Data["content"]) != content {
		secret.Data = map[string][]byte{
			"content": []byte(content),
		}
		return r.Update(ctx, secret)
	}
	return nil
}

func doCheckSecretTextCredential(r *ReconcileDevOpsApp, ctx context.Context, namespace string, name string, secretText string) error {
	if _, err := r.devopsClient.GetCredentialInProject(namespace, name); err != nil {
		secret := &corev1.Secret{
			Type: devopsv1alpha3.SecretTypeSecretText,
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: map[string][]byte{
				"secret": []byte(secretText),
			},
		}
		return r.Create(ctx, secret)
	}

	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret); err != nil {
		return nil
	}

	if string(secret.Data["secret"]) != secretText {
		secret.Data = map[string][]byte{
			"secret": []byte(secretText),
		}
		return r.Update(ctx, secret)
	}
	return nil
}

func doCheckSshAuthCredential(r *ReconcileDevOpsApp, ctx context.Context, namespace string, name string, username string, password string, privateKey string) error {
	if _, err := r.devopsClient.GetCredentialInProject(namespace, name); err != nil {
		secret := &corev1.Secret{
			Type: devopsv1alpha3.SecretTypeSSHAuth,
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: map[string][]byte{
				"username":    []byte(username),
				"passphrase":  []byte(password),
				"private_key": []byte(privateKey),
			},
		}
		return r.Create(ctx, secret)
	}

	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret); err != nil {
		return nil
	}

	if string(secret.Data["username"]) != username || string(secret.Data["passphrase"]) != password || string(secret.Data["private_key"]) != privateKey {
		secret.Data = map[string][]byte{
			"username":    []byte(username),
			"passphrase":  []byte(password),
			"private_key": []byte(privateKey),
		}
		return r.Update(ctx, secret)
	}
	return nil
}
