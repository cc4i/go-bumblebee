package k8s

import (
	"flag"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
	//
	// Uncomment to load all auth plugins
	//_ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

type K8sContext struct {
	Clientset *kubernetes.Clientset
}

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type MyNamespace struct {
	Name    string    `json:"name"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
	Age     int64     `json:"age"` //second
	Labels  []Label   `json:"labels"`
}

type MyDeploymentPods struct {}
type MyServicePods struct {}

type K8sClient interface {
	GetNamespaces() ([]MyNamespace, error)
	GetAllDeploymentsPods(namespace string) ([]MyServicePods, error)
	GetAllServicesPods(namespace string) ([]MyServicePods, error)
	GetAll() string
	GetCRD() string
}

func (cs *K8sContext) GetAllServicesPods(namespace string) ([]MyServicePods, error) {
	servicePods := []MyServicePods{}


	return servicePods, nil
}

func (cs *K8sContext) GetNamespaces() ([]MyNamespace, error) {
	namespaces := []MyNamespace{}
	nl, err := cs.Clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return namespaces, errors.Wrapf(err, "Failed to retrieve namespace from K8s cluster: %s", err.Error())
	}
	for _, n := range nl.Items {
		labels := []Label{}
		for k, v := range n.GetObjectMeta().GetLabels() {
			label := Label{
				Key:   k,
				Value: v,
			}
			labels = append(labels, label)

		}
		namespaces = append(namespaces, MyNamespace{
			Name:    n.Name,
			Status:  string(n.Status.Phase),
			Created: n.GetObjectMeta().GetCreationTimestamp().Local(),
			Age:     time.Since(n.GetObjectMeta().GetCreationTimestamp().Local()).Milliseconds(),
			Labels:  labels,
		})
	}
	return namespaces, nil

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getConfig() *rest.Config {
	inCluster := os.Getenv("IN_CLUSTER_CONFIG")

	if inCluster == "true" {
		inc, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		return inc
	} else {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		exc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		return exc
	}

}
