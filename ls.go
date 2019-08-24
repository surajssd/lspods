package main

import (
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Printf("no resource mentioned, example command run: \n./ls pods\n")
		os.Exit(-1)
	}

	k8sresource := argsWithoutProg[0]
	if k8sresource != "pods" {
		fmt.Printf("error: the server doesn't have a resource type %q\n", k8sresource)
		os.Exit(-1)
	}

	kubeconfig := getKubeConfig()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("can't build config from kubeconfig at %s: %v\n", kubeconfig, err)
		os.Exit(-1)
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("can't get kubernetes client: %v\n", err)
		os.Exit(-1)
	}

	ns := "nginx"
	watcher, err := cs.CoreV1().Pods(ns).Watch(meta_v1.ListOptions{})
	if err != nil {
		fmt.Printf("could not watch pods in %q namespace: %v\n", ns, err)
	}

	fmt.Println("Watching on Pods in the ns:", ns)
	fmt.Println("")
	ch := watcher.ResultChan()
	for event := range ch {
		pod := event.Object.(*corev1.Pod)
		fmt.Printf("%s pod %v\n", event.Type, pod.Name)
	}
}

func getKubeConfig() string {
	kubeconfigEnv := os.Getenv("KUBECONFIG")
	if kubeconfigEnv != "" {
		return kubeconfigEnv
	}

	return os.ExpandEnv("$HOME/.kube/config")
}
