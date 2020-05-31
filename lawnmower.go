/*
Copyright 2016 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file.")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file.")
	}

	var exceptapigroups *string
	exceptapigroups = flag.String("except", "", "(optional) except apiGroups sepalated by comma. EX) rbac.authorization.k8s.io,networking.k8s.io")

	var crname *string
	crname = flag.String("clusterRoleName", "restricted-cluster-admin", "(optional) resource name of cluster role.")

	var crfname *string
	crfname = flag.String("outputFileName", "restricted-cluster-admin-cluster-role.yaml", "(optional) name of output cluster role file.")

	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	apigrouplist, err := clientset.Discovery().ServerGroups()
	if err != nil {
		panic(err.Error())
	}

	// api group name list
	var apignl []string
	for _, apigroup := range apigrouplist.Groups {
		apignl = append(apignl, apigroup.Name)
	}

	for _, e := range strings.Split(*exceptapigroups, ",") {
		apignl = Remove(apignl, e)
	}

	// Set ClusterRole template
	cr := rbacv1.ClusterRole{}
	cr.APIVersion = "rbac.authorization.k8s.io/v1"
	cr.Kind = "ClusterRole"
	cr.Name = *crname

	// Add Resource rule
	policyRule := rbacv1.PolicyRule{}
	policyRule.APIGroups = apignl
	policyRule.Verbs = []string{"*"}
	policyRule.Resources = []string{"*"}
	cr.Rules = append(cr.Rules, policyRule)

	// Add NonResourceURLs
	policyRule = rbacv1.PolicyRule{}
	policyRule.NonResourceURLs = []string{"*"}
	policyRule.Verbs = []string{"*"}
	cr.Rules = append(cr.Rules, policyRule)

	// ClusterRole is not parsed to yaml directly because type of ClusterRole has some tags for arranging manifest.
	crJSON, err := json.Marshal(cr)
	if err != nil {
		fmt.Errorf("Could not marshal: %v", err)
	}

	crYaml, err := yaml.JSONToYAML(crJSON)
	if err != nil {
		fmt.Errorf("Could not convert from json to yaml: %v", err)
	}

	err = ioutil.WriteFile(*crfname, []byte(crYaml), 0644)
	if err != nil {
		fmt.Errorf("Could not write result: %v", err)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func Remove(s []string, str string) []string {
	var res = []string{}
	for _, v := range s {
		if v == str {
			continue
		}
		res = append(res, v)
	}
	return res
}
