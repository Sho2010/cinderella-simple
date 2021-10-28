package k8s

import (
	"context"
	"fmt"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
)

// NOTE: https://github.com/kubernetes/client-go/issues/193

type ResourceCreater struct {
	client kubernetes.Interface
}

// Role File
// RoleBinding file
// ServiceAccount file

func (rc *ResourceCreater) Create() error {
	return nil
}

func (rc *ResourceCreater) createServiceAccount() {
}

func (rc *ResourceCreater) createRoleBinding() {
}

func (rc *ResourceCreater) createRole() {

	scheme := runtime.NewScheme()
	codecFactory := serializer.NewCodecFactory(scheme)
	deserializer := codecFactory.UniversalDeserializer()

	yaml, err := ioutil.ReadFile("k8s/templates/role.yaml.tmpl")
	if err != nil {
		panic(err)
	}

	o, _, err := deserializer.Decode(yaml, nil, &rbacv1.Role{})
	if err != nil {
		panic(err)
	}
	role := o.(*rbacv1.Role)

	namespace := "default"
	roleClient := rc.client.RbacV1().Roles(namespace)

	var opts metav1.CreateOptions = metav1.CreateOptions{}
	ret, err := roleClient.Create(context.TODO(), role, opts)
	if err != nil {
		//TODO: error handling
		panic(err)
	}

	fmt.Println(ret)
}

// func CreateRole(client kubernetes.Interface) {
// 	scheme := runtime.NewScheme()
// 	codecFactory := serializer.NewCodecFactory(scheme)
// 	deserializer := codecFactory.UniversalDeserializer()
//
// 	yaml, err := ioutil.ReadFile("k8s/templates/role.yaml.tmpl")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	object, _, err := deserializer.Decode(yaml, nil, &rbacv1.Role{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	role := object.(*rbacv1.Role)
//
// 	namespace := "default"
// 	roleClient := client.RbacV1().Roles(namespace)
//
// 	var opts metav1.CreateOptions = metav1.CreateOptions{}
// 	result, err := roleClient.Create(context.TODO(), role, opts)
// 	if err != nil {
// 		//TODO: error handling
// 		panic(err)
// 	}
//
// }
