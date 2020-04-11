package k8s

import (
	"bytes"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
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

type K8sClient interface {
	GetNamespaces() string
	GetAllServicesPods(namespace string) string
	GetAllDeploymentsPods(namespace string) string
	GetAll() string
	GetCRD() string
}

func (cs *K8sContext)GetNamespaces() string {
	nl, err := cs.Clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err!=nil {
		log.Error(err)
		return err.Error()
	}
	var buf bytes.Buffer
	buf.WriteString("Name\t\tStatus\t\tCreated\t\tAge\t\tLabels\n")
	for _, n := range nl.Items {
		buf.WriteString(n.Name+"\t\t")
		buf.WriteString(string(n.Status.Phase)+"\t\t")
		buf.WriteString(n.GetObjectMeta().GetCreationTimestamp().String()+"\t\t")

		h := fmt.Sprintf("%.2fh\t\t", time.Since(n.GetObjectMeta().GetCreationTimestamp().Local()).Hours())
		buf.WriteString(h)
		for k,v := range n.GetObjectMeta().GetLabels() {
			buf.WriteString(k+"="+v+";")
		}

		buf.WriteString("\n")
	}
	return buf.String()
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

	return nil
}