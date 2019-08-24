package main

import (
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 1 {
		fmt.Printf("no resource mentioned, example command run: \n./ls create <podname>\n")
		os.Exit(-1)
	}

	podName := argsWithoutProg[1]

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

	ns := "default"
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: ns,
			Labels: map[string]string{
				"app": "mypod",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  podName,
					Image: "bitnami/nginx",
				},
			},
		},
	}

	_, err = cs.CoreV1().Pods(ns).Create(pod)
	if err != nil {
		fmt.Printf("could not create pods in %q namespace: %v\n", ns, err)
	}
	fmt.Println("Created pod", podName, "successfully")
}

func getKubeConfig() string {
	kubeconfigEnv := os.Getenv("KUBECONFIG")
	if kubeconfigEnv != "" {
		return kubeconfigEnv
	}

	return os.ExpandEnv("$HOME/.kube/config")
}
