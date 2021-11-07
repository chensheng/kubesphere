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
	"bytes"
	"container/list"
	"context"
	"html/template"
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
	kubeclient "k8s.io/client-go/kubernetes"
	devopsv1alpha3 "kubesphere.io/kubesphere/pkg/apis/devops/v1alpha3"
	devopsClient "kubesphere.io/kubesphere/pkg/simple/client/devops"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"k8s.io/client-go/tools/clientcmd"
	clusterinformer "kubesphere.io/kubesphere/pkg/client/informers/externalversions/cluster/v1alpha1"
)

var log = logf.Log.WithName("devopsapp_controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new DevOpsApp Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, devopsClient devopsClient.Interface, clusterInformer clusterinformer.ClusterInformer) error {
	return add(mgr, newReconciler(mgr, devopsClient, clusterInformer))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, devopsClient devopsClient.Interface, clusterInformer clusterinformer.ClusterInformer) reconcile.Reconciler {
	return &ReconcileDevOpsApp{Client: mgr.GetClient(), scheme: mgr.GetScheme(), devopsClient: devopsClient, clusterInformer: clusterInformer, clusterClients: make(map[string]*kubeclient.Clientset)}
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

type TemplateData struct {
	Name        string
	Parameters  map[string]string
	Environment devopsv1alpha3.Environment
	DevOpsApp   devopsv1alpha3.DevOpsApp
}

// ReconcileDevOpsApp reconciles a DevOpsApp object
type ReconcileDevOpsApp struct {
	client.Client
	scheme          *runtime.Scheme
	devopsClient    devopsClient.Interface
	clusterInformer clusterinformer.ClusterInformer
	clusterClients  map[string]*kubeclient.Clientset
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
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if !devopsapp.Spec.AutoUpdate {
		return reconcile.Result{}, nil
	}

	devopsappCopy, err := resolveReferenceSecret(r, rootCtx, *devopsapp)
	if err != nil {
		return reconcile.Result{}, err
	}

	workspace := devopsappCopy.ObjectMeta.Labels["kubesphere.io/workspace"]
	envs := devopsappCopy.Spec.Environments
	if envs == nil || len(envs) == 0 {
		return reconcile.Result{}, nil
	}

	devopsappName := devopsappCopy.ObjectMeta.GenerateName
	if devopsappName == "" {
		devopsappName = devopsappCopy.ObjectMeta.Name
	}

	crdChanged := false
	devopsproject, err := checkDevOpsProject(r, rootCtx, devopsappCopy)
	if err != nil {
		return reconcile.Result{}, err
	}
	if devopsproject.Name != devopsapp.Spec.DevOpsProject {
		devopsapp.Spec.DevOpsProject = devopsproject.Name
		crdChanged = true
	}
	if err := checkCommonCredentials(r, rootCtx, devopsproject, devopsappCopy); err != nil {
		return reconcile.Result{}, err
	}

	for index, env := range envs {
		var clusterClient, err = getClusterClient(r, env.Cluster)
		envNamespace, err := checkEnvNamespace(rootCtx, clusterClient, workspace, devopsappName, *env)
		if err != nil {
			return reconcile.Result{}, err
		}

		if envNamespace.Name != env.Namespace {
			devopsapp.Spec.Environments[index].Namespace = envNamespace.Name
			crdChanged = true
		}

		if err = checkEnvSecrets(clusterClient, rootCtx, devopsappCopy, *env, envNamespace.Name); err != nil {
			return reconcile.Result{}, err
		}

		if err = checkEnvCredentials(r, rootCtx, devopsproject, devopsappCopy, *env); err != nil {
			return reconcile.Result{}, err
		}

		pipeline, err := checkEnvPipeline(r, rootCtx, devopsproject, devopsappCopy, *env)
		if err != nil {
			return reconcile.Result{}, err
		}
		if env.Pipeline != nil && env.Pipeline.Name != pipeline.Name {
			devopsapp.Spec.Environments[index].Pipeline.Name = pipeline.Name
			crdChanged = true
		}
	}

	if crdChanged {
		if err = r.Update(rootCtx, devopsapp); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func resolveReferenceSecret(r *ReconcileDevOpsApp, ctx context.Context, originalDevopsapp devopsv1alpha3.DevOpsApp) (devopsappCopy *devopsv1alpha3.DevOpsApp, err error) {
	devopsapp := originalDevopsapp.DeepCopy()
	git := devopsapp.Spec.Git
	configCenter := devopsapp.Spec.ConfigCenter
	registry := devopsapp.Spec.Registry
	environments := devopsapp.Spec.Environments

	references := list.New()
	if git != nil && git.Reference != "" {
		references.PushFront(git.Reference)
	}
	if configCenter != nil && configCenter.Reference != "" {
		references.PushFront(configCenter.Reference)
	}
	if registry != nil && registry.Reference != "" {
		references.PushFront(registry.Reference)
	}
	if len(environments) > 0 {
		for _, env := range environments {
			if env.ConfigCenter != nil && env.ConfigCenter.Reference != "" {
				references.PushFront(env.ConfigCenter.Reference)
			}
			if env.Registry != nil && env.Registry.Reference != "" {
				references.PushFront(env.Registry.Reference)
			}
			if len(env.Credentials) > 0 {
				for _, credential := range env.Credentials {
					if credential.Reference != "" {
						references.PushFront(credential.Reference)
					}
				}
			}
		}
	}

	secretMap := make(map[string]*devopsv1alpha3.DevOpsSecret)
	for ref := references.Front(); ref != nil; ref = ref.Next() {
		secretName := ref.Value.(string)
		if _, ok := secretMap[secretName]; ok {
			continue
		}
		secret := &devopsv1alpha3.DevOpsSecret{}
		if err := r.Get(ctx, types.NamespacedName{Name: secretName}, secret); err == nil {
			secretMap[secretName] = secret
		}
	}

	if git != nil && git.Reference != "" {
		if secret, ok := secretMap[git.Reference]; ok {
			if git.Username == "" {
				git.Username = secret.Spec.Username
			}
			if git.Password == "" {
				git.Password = secret.Spec.Password
			}
		}
	}
	if configCenter != nil && configCenter.Reference != "" {
		if secret, ok := secretMap[configCenter.Reference]; ok {
			if configCenter.Username == "" {
				configCenter.Username = secret.Spec.Username
			}
			if configCenter.Password == "" {
				configCenter.Password = secret.Spec.Password
			}
		}
	}
	if registry != nil && registry.Reference != "" {
		if secret, ok := secretMap[registry.Reference]; ok {
			if registry.Username == "" {
				registry.Username = secret.Spec.Username
			}
			if registry.Password == "" {
				registry.Password = secret.Spec.Password
			}
		}
	}

	if len(environments) > 0 {
		for _, env := range environments {
			if env.Registry == nil {
				env.Registry = registry
			} else if env.Registry.Reference != "" {
				if secret, ok := secretMap[env.Registry.Reference]; ok {
					if env.Registry.Username == "" {
						env.Registry.Username = secret.Spec.Username
					}
					if env.Registry.Password == "" {
						env.Registry.Password = secret.Spec.Password
					}
				}
			}

			if env.ConfigCenter == nil {
				env.ConfigCenter = configCenter
			} else if env.ConfigCenter.Reference != "" {
				if secret, ok := secretMap[env.ConfigCenter.Reference]; ok {
					if env.ConfigCenter.Username == "" {
						env.ConfigCenter.Username = secret.Spec.Username
					}
					if env.ConfigCenter.Password == "" {
						env.ConfigCenter.Password = secret.Spec.Password
					}
				}
			}

			if len(env.Credentials) > 0 {
				for _, credential := range env.Credentials {
					if credential.Reference == "" {
						continue
					}

					if secret, ok := secretMap[credential.Reference]; ok {
						if credential.Username == "" {
							credential.Username = secret.Spec.Username
						}
						if credential.Password == "" {
							credential.Password = secret.Spec.Password
						}
						if credential.PrivateKey == "" {
							credential.PrivateKey = secret.Spec.PrivateKey
						}
						if credential.Content == "" {
							credential.Content = secret.Spec.Content
						}
						if credential.Secret == "" {
							credential.Secret = secret.Spec.Secret
						}
					}
				}
			}
		}
	}

	return devopsapp, nil
}

func getClusterClient(r *ReconcileDevOpsApp, clusterName string) (clusterClient *kubeclient.Clientset, err error) {
	if clusterClient, ok := r.clusterClients[clusterName]; ok {
		if _, err := clusterClient.Discovery().ServerVersion(); err == nil {
			return clusterClient, nil
		}
	}

	clusterLister := r.clusterInformer.Lister()
	cluster, err := clusterLister.Get(clusterName)
	if err != nil {
		return &kubeclient.Clientset{}, err
	}

	clientConfig, err := clientcmd.NewClientConfigFromBytes(cluster.Spec.Connection.KubeConfig)
	if err != nil {
		return &kubeclient.Clientset{}, err
	}

	clusterConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return &kubeclient.Clientset{}, err
	}

	newClusterClient, err := kubeclient.NewForConfig(clusterConfig)
	if err != nil {
		return &kubeclient.Clientset{}, err
	}

	r.clusterClients[clusterName] = newClusterClient
	return newClusterClient, nil
}

func checkEnvNamespace(ctx context.Context, clusterClient *kubeclient.Clientset, workspace string, devopsappName string, env devopsv1alpha3.Environment) (namespace *corev1.Namespace, err error) {
	envNamespaceName := env.Namespace
	if envNamespaceName == "" {
		envNamespaceName = env.Name + "-" + devopsappName
	}

	envNamespace, err := clusterClient.CoreV1().Namespaces().Get(ctx, envNamespaceName, metav1.GetOptions{})
	if err == nil && envNamespace.ObjectMeta.Labels["kubesphere.io/workspace"] == workspace {
		return envNamespace, nil
	}

	if err != nil {
		envNamespace = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: envNamespaceName,
				Labels: map[string]string{
					"kubesphere.io/workspace": workspace,
				},
			},
		}
	} else {
		envNamespace = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: envNamespaceName,
				Labels: map[string]string{
					"kubesphere.io/workspace": workspace,
				},
			},
		}
	}

	if envNamespace, err = clusterClient.CoreV1().Namespaces().Create(ctx, envNamespace, metav1.CreateOptions{}); err != nil {
		return nil, err
	}

	return envNamespace, nil
}

func checkDevOpsProject(r *ReconcileDevOpsApp, ctx context.Context, devopsapp *devopsv1alpha3.DevOpsApp) (devopsproject *devopsv1alpha3.DevOpsProject, err error) {
	workspace := devopsapp.Labels["kubesphere.io/workspace"]
	devopsprojectName := devopsapp.Spec.DevOpsProject

	if devopsprojectName != "" {
		devopsproject := &devopsv1alpha3.DevOpsProject{}
		if err := r.Get(ctx, types.NamespacedName{Name: devopsprojectName}, devopsproject); err == nil {
			return devopsproject, nil
		} else {
			devopsproject.Name = devopsprojectName
			devopsproject.GenerateName = devopsapp.Name
			devopsproject.ObjectMeta.Labels["kubesphere.io/workspace"] = workspace
			err = r.Create(ctx, devopsproject)
			return devopsproject, err
		}
	}

	devopsprojectList := &devopsv1alpha3.DevOpsProjectList{}
	opts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{"kubesphere.io/workspace": workspace}),
	}
	if err := r.List(ctx, devopsprojectList, opts); err == nil {
		for _, project := range devopsprojectList.Items {
			if devopsapp.Name == project.GenerateName {
				return &project, nil
			}
		}
	}

	devopsproject = &devopsv1alpha3.DevOpsProject{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: devopsapp.Name,
			Labels: map[string]string{
				"kubesphere.io/workspace": workspace,
			},
		},
	}
	err = r.Create(ctx, devopsproject)
	return devopsproject, err
}

func checkCommonCredentials(r *ReconcileDevOpsApp, ctx context.Context, devopsproject *devopsv1alpha3.DevOpsProject, devopsapp *devopsv1alpha3.DevOpsApp) error {
	if err := doCheckBasicAuthCredential(r, ctx, devopsproject.Name, "git", devopsapp.Spec.Git.Username, devopsapp.Spec.Git.Password); err != nil {
		return err
	}

	if err := doCheckBasicAuthCredential(r, ctx, devopsproject.Name, "harbor", devopsapp.Spec.Registry.Username, devopsapp.Spec.Registry.Password); err != nil {
		return err
	}

	return nil
}

func checkEnvSecrets(clusterClient *kubeclient.Clientset, ctx context.Context, devopsapp *devopsv1alpha3.DevOpsApp, env devopsv1alpha3.Environment, namespace string) error {
	registry := env.Registry
	secretName := "harbor"
	rawAuth := registry.Username + ":" + registry.Password
	encodedAuth := make([]byte, base64.StdEncoding.EncodedLen(len(rawAuth)))
	base64.StdEncoding.Encode(encodedAuth, []byte(rawAuth))
	auth := string(encodedAuth)
	if strings.HasSuffix(auth, "==") {
		auth = auth[0 : len(auth)-2]
	}

	var secretInstance, err = clusterClient.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		dockerConfigJson := &DockerConfigJson{Auths: DockerConfigMap{}}
		dockerConfigJson.Auths[registry.Url] = DockerConfigEntry{
			Username: registry.Username,
			Password: registry.Password,
			Email:    registry.Email,
			Auth:     auth,
		}
		secretData, _ := json.Marshal(dockerConfigJson)

		secretInstance = &corev1.Secret{}
		secretInstance.Name = secretName
		secretInstance.Namespace = namespace
		secretInstance.Type = corev1.SecretTypeDockerConfigJson
		secretInstance.StringData = map[string]string{corev1.DockerConfigJsonKey: string(secretData)}
		_, err := clusterClient.CoreV1().Secrets(namespace).Create(ctx, secretInstance, metav1.CreateOptions{})
		return err
	}

	secretData := secretInstance.Data[corev1.DockerConfigJsonKey]
	dockerConfigJson := &DockerConfigJson{Auths: DockerConfigMap{}}
	json.Unmarshal(secretData, dockerConfigJson)
	dockerConfigEntry := dockerConfigJson.Auths[registry.Url]
	if dockerConfigEntry.Username == registry.Username && dockerConfigEntry.Password == registry.Password {
		return nil
	}

	dockerConfigJson.Auths = DockerConfigMap{}
	dockerConfigJson.Auths[registry.Url] = DockerConfigEntry{
		Username: registry.Username,
		Password: registry.Password,
		Email:    registry.Email,
		Auth:     auth,
	}
	secretDataBytes, _ := json.Marshal(dockerConfigJson)
	secretInstance.StringData = map[string]string{corev1.DockerConfigJsonKey: string(secretDataBytes)}
	_, err = clusterClient.CoreV1().Secrets(namespace).Update(ctx, secretInstance, metav1.UpdateOptions{})
	return err
}

func checkEnvCredentials(r *ReconcileDevOpsApp, ctx context.Context, devopsproject *devopsv1alpha3.DevOpsProject, devopsapp *devopsv1alpha3.DevOpsApp, env devopsv1alpha3.Environment) error {
	if len(env.Credentials) == 0 {
		return nil
	}

	for _, credential := range env.Credentials {
		var err error
		if credential.Type == "basic-auth" {
			err = doCheckBasicAuthCredential(r, ctx, devopsproject.Name, credential.Name, credential.Username, credential.Password)
		} else if credential.Type == "kubeconfig" {
			err = doCheckKubeConfigCredential(r, ctx, devopsproject.Name, credential.Name, credential.Content)
		} else if credential.Type == "secret-text" {
			err = doCheckSecretTextCredential(r, ctx, devopsproject.Name, credential.Name, credential.Secret)
		} else if credential.Type == "ssh-auth" {
			err = doCheckSshAuthCredential(r, ctx, devopsproject.Name, credential.Name, credential.Username, credential.Password, credential.PrivateKey)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func doCheckBasicAuthCredential(r *ReconcileDevOpsApp, ctx context.Context, namespace string, name string, username string, password string) error {
	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret); err != nil {
		secret = &corev1.Secret{
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

	if string(secret.Data["username"]) == username && string(secret.Data["password"]) == password {
		return nil
	}

	secret.Data = map[string][]byte{
		"username": []byte(username),
		"password": []byte(password),
	}
	return r.Update(ctx, secret)
}

func doCheckKubeConfigCredential(r *ReconcileDevOpsApp, ctx context.Context, namespace string, name string, content string) error {
	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret); err != nil {
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

	if string(secret.Data["content"]) == content {
		return nil
	}

	secret.Data = map[string][]byte{
		"content": []byte(content),
	}
	return r.Update(ctx, secret)
}

func doCheckSecretTextCredential(r *ReconcileDevOpsApp, ctx context.Context, namespace string, name string, secretText string) error {
	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret); err != nil {
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

	if string(secret.Data["secret"]) == secretText {
		return nil
	}

	secret.Data = map[string][]byte{
		"secret": []byte(secretText),
	}
	return r.Update(ctx, secret)
}

func doCheckSshAuthCredential(r *ReconcileDevOpsApp, ctx context.Context, namespace string, name string, username string, password string, privateKey string) error {
	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret); err != nil {
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

	if string(secret.Data["username"]) == username && string(secret.Data["passphrase"]) == password && string(secret.Data["private_key"]) == privateKey {
		return nil
	}

	secret.Data = map[string][]byte{
		"username":    []byte(username),
		"passphrase":  []byte(password),
		"private_key": []byte(privateKey),
	}
	return r.Update(ctx, secret)
}

func checkEnvPipeline(r *ReconcileDevOpsApp, ctx context.Context, devopsproject *devopsv1alpha3.DevOpsProject, devopsapp *devopsv1alpha3.DevOpsApp, env devopsv1alpha3.Environment) (pipeline *devopsv1alpha3.Pipeline, err error) {
	if env.Pipeline == nil || env.Pipeline.TemplateName == "" {
		return &devopsv1alpha3.Pipeline{}, nil
	}

	pipelineName := env.Pipeline.Name
	if pipelineName == "" {
		pipelineName = "pipeline-" + env.Name
	}

	pipelineTemplate := &devopsv1alpha3.PipelineTemplate{}
	if err := r.Get(ctx, types.NamespacedName{Name: env.Pipeline.TemplateName}, pipelineTemplate); err != nil {
		return &devopsv1alpha3.Pipeline{}, err
	}

	pipeline = &devopsv1alpha3.Pipeline{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: devopsproject.Name, Name: pipelineName}, pipeline); err != nil {
		noScmPipeline, err := resolveNoScmPipeline(pipelineName, devopsproject, devopsapp, env, pipelineTemplate)
		if err != nil {
			return &devopsv1alpha3.Pipeline{}, err
		}

		pipeline.Name = pipelineName
		pipeline.Namespace = devopsproject.Name
		pipeline.Spec = devopsv1alpha3.PipelineSpec{
			Type:     "pipeline",
			Pipeline: noScmPipeline,
		}
		err = r.Create(ctx, pipeline)
		return pipeline, err
	}

	newNoScmPipeline, err := resolveNoScmPipeline(pipelineName, devopsproject, devopsapp, env, pipelineTemplate)
	if err != nil {
		return &devopsv1alpha3.Pipeline{}, err
	}

	existingNoScmPipeline := pipeline.Spec.Pipeline
	if !isPipelineChanged(existingNoScmPipeline, newNoScmPipeline) {
		return pipeline, nil
	}

	pipeline.Spec.Pipeline = newNoScmPipeline
	err = r.Update(ctx, pipeline)
	return pipeline, err
}

func isPipelineChanged(existingNoScmPipeline *devopsv1alpha3.NoScmPipeline, newNoScmPipeline *devopsv1alpha3.NoScmPipeline) bool {
	if existingNoScmPipeline.DisableConcurrent != newNoScmPipeline.DisableConcurrent {
		return true
	}

	if (existingNoScmPipeline.Discarder == nil && newNoScmPipeline.Discarder != nil) || (existingNoScmPipeline.Discarder != nil && newNoScmPipeline.Discarder == nil) {
		return true
	}

	if existingNoScmPipeline.Discarder.DaysToKeep != newNoScmPipeline.Discarder.DaysToKeep || existingNoScmPipeline.Discarder.NumToKeep != newNoScmPipeline.Discarder.NumToKeep {
		return true
	}

	if existingNoScmPipeline.Jenkinsfile != newNoScmPipeline.Jenkinsfile {
		return true
	}

	if (existingNoScmPipeline.RemoteTrigger == nil && newNoScmPipeline.RemoteTrigger != nil) || (existingNoScmPipeline.RemoteTrigger != nil && newNoScmPipeline.RemoteTrigger == nil) {
		return true
	}

	if (existingNoScmPipeline.RemoteTrigger != nil && newNoScmPipeline.RemoteTrigger != nil) && (existingNoScmPipeline.RemoteTrigger.Token != newNoScmPipeline.RemoteTrigger.Token) {
		return true
	}

	if (existingNoScmPipeline.TimerTrigger == nil && newNoScmPipeline.TimerTrigger != nil) || (existingNoScmPipeline.TimerTrigger != nil && newNoScmPipeline.TimerTrigger == nil) {
		return true
	}

	if (existingNoScmPipeline.TimerTrigger != nil && newNoScmPipeline.TimerTrigger != nil) && (existingNoScmPipeline.TimerTrigger.Cron != newNoScmPipeline.TimerTrigger.Cron || existingNoScmPipeline.TimerTrigger.Interval != newNoScmPipeline.TimerTrigger.Interval) {
		return true
	}

	if len(existingNoScmPipeline.Parameters) != len(newNoScmPipeline.Parameters) {
		return true
	}

	paramMap := make(map[string]devopsv1alpha3.Parameter)
	for _, param := range existingNoScmPipeline.Parameters {
		paramMap[param.Name] = param
	}
	for _, param := range newNoScmPipeline.Parameters {
		existingParam, ok := paramMap[param.Name]
		if !ok || existingParam.Type != param.Type || existingParam.DefaultValue != param.DefaultValue || existingParam.Description != param.Description {
			return true
		}
	}

	return false
}

func resolveNoScmPipeline(pipelineName string, devopsproject *devopsv1alpha3.DevOpsProject, devopsapp *devopsv1alpha3.DevOpsApp, env devopsv1alpha3.Environment, pipelineTemplate *devopsv1alpha3.PipelineTemplate) (noScmPipeline *devopsv1alpha3.NoScmPipeline, err error) {
	goTplName := devopsproject.Name + "_" + pipelineName
	goTemplate, err := template.New(goTplName).Parse(pipelineTemplate.Spec.Jenkinsfile)
	if err != nil {
		return &devopsv1alpha3.NoScmPipeline{}, err
	}

	parameters := make(map[string]string)
	if len(pipelineTemplate.Spec.JenkinsfileParameters) > 0 {
		for _, param := range pipelineTemplate.Spec.JenkinsfileParameters {
			parameters[param.Name] = param.DefaultValue
			if len(env.Pipeline.Parameters) > 0 && env.Pipeline.Parameters[param.Name] != "" {
				parameters[param.Name] = env.Pipeline.Parameters[param.Name]
			}
		}
	}

	templateData := TemplateData{
		DevOpsApp:   *devopsapp,
		Environment: env,
		Parameters:  parameters,
	}
	buffer := bytes.NewBuffer(nil)
	if err = goTemplate.Execute(buffer, templateData); err != nil {
		return &devopsv1alpha3.NoScmPipeline{}, err
	}

	noScmPipeline = &devopsv1alpha3.NoScmPipeline{
		Name:              pipelineName,
		Discarder:         pipelineTemplate.Spec.Discarder,
		Parameters:        pipelineTemplate.Spec.BuildParameters,
		DisableConcurrent: pipelineTemplate.Spec.DisableConcurrent,
		TimerTrigger:      pipelineTemplate.Spec.TimerTrigger,
		RemoteTrigger:     pipelineTemplate.Spec.RemoteTrigger,
		Jenkinsfile:       buffer.String(),
	}
	return noScmPipeline, nil
}
