package main

import (
	"fmt"
	"os"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

	var cs *kubernetes.Clientset
	cs.CoreV1().Pods("kube-system").List(meta_v1.ListOptions{})
}
