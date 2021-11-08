package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/Sho2010/cinderella-simple/claim"
	"github.com/Sho2010/cinderella-simple/config"
	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
)

// NOTE: https://github.com/kubernetes/client-go/issues/193

type ResourceCreater struct {
	client                  kubernetes.Interface
	serviceAccountNamespace string
	claim                   claim.Claim
}

var _requireLabels = map[string]string{
	"app.kubernetes.io/managed-by": "cinderella",
}

var (
	defaultRoleManifest        = "default-role"
	defaultRoleBindingManifest = "default-role-binding"
)

func NewResourceCreater(client kubernetes.Interface, serviceAccountNamespace string, claim claim.Claim) (*ResourceCreater, error) {

	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	if claim == nil {
		return nil, fmt.Errorf("claim is nil")
	}

	return &ResourceCreater{
		client:                  client,
		serviceAccountNamespace: serviceAccountNamespace,
		claim:                   claim,
	}, nil
}

func (rc *ResourceCreater) Create() error {

	// TODO: each
	// for _, ns := range rc.claim.GetNamespaces() {
	// 	role.Namespace = ns
	// 	rc.client.RbacV1().Roles(ns)
	// 	_, err := roleClient.Create(context.TODO(), role, metav1.CreateOptions{})
	// 	if err != nil {
	// 		return nil, fmt.Errorf("do %w", err)
	// 	}
	// }

	ns := rc.claim.GetNamespaces()[0]
	role, err := rc.createRole(defaultRoleManifest, ns)
	if err != nil {
		return err
	}

	rb, err := rc.createRoleBinding(
		defaultRoleBindingManifest,
		defaultRoleManifest,
		ns)
	if err != nil {
		return err
	}

	sa, err := rc.createServiceAccount()
	if err != nil {
		return err
	}
	fmt.Println(role)
	fmt.Println(rb)
	fmt.Println(sa)

	RaiseResourceCreateEvent(fmt.Sprintf("Create success %s/%s/%s", role.GetName(), rb.GetName(), sa.GetName()))
	return nil
}

func (rc *ResourceCreater) createServiceAccount() (metav1.Object, error) {
	//TODO: err
	saName, _ := rc.claim.GetServiceAccountName()

	sa := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      saName,
			Namespace: rc.serviceAccountNamespace,
		},
	}
	saClient := rc.client.CoreV1().ServiceAccounts(rc.serviceAccountNamespace)

	if sa.Labels == nil {
		sa.Labels = make(map[string]string)
	}
	if err := rc.mergeLabels(sa.Labels); err != nil {
		panic(err)
	}

	if sa.Annotations == nil {
		sa.Annotations = make(map[string]string)
	}
	if err := rc.mergeAnnotations(sa.Annotations); err != nil {
		panic(err)
	}

	ret, err := saClient.Create(context.TODO(), &sa, metav1.CreateOptions{})
	if err != nil {
		// FIXME: Applyにしたかったところだけどメッチャクチャめんどくさそうなのでとりあえず握りつぶす
		if errors.IsAlreadyExists(err) {
			return ret, nil
		}
		return nil, err
	}
	return ret, nil
}

//TODO: arg適当すぎるのでリファクタ
func (rc *ResourceCreater) createRoleBinding(bindingName string, roleName string, namespace string) (metav1.Object, error) {
	scheme := runtime.NewScheme()
	codecFactory := serializer.NewCodecFactory(scheme)
	deserializer := codecFactory.UniversalDeserializer()

	f := config.SearchManifest(bindingName)
	if len(f) == 0 {
		return nil, fmt.Errorf("v1/RoleBinding %s manifest file not found", bindingName)
	}

	// yaml, err := os.ReadFile(f)
	// if err != nil {
	// 	return nil, fmt.Errorf("do %w", err)
	// }

	saName, _ := rc.claim.GetServiceAccountName()
	roleName, err := claim.NormalizeDNS1123(roleName + "-" + rc.claim.GetSubject())
	bindName, err := claim.NormalizeDNS1123(bindingName + "-" + rc.claim.GetSubject())
	values := RoleBindingValues{
		ServiceAccountName:      saName,
		ServiceAccountNamespace: rc.serviceAccountNamespace,
		RoleBindingName:         bindName,
		RoleName:                roleName,
	}

	data, err := templateExecute(f, values)
	if err != nil {
		return nil, fmt.Errorf("Role manifest template execute fail: %w", err)
	}

	obj, _, err := deserializer.Decode(data, nil, &rbacv1.RoleBinding{})
	if err != nil {
		return nil, err
	}
	rb := obj.(*rbacv1.RoleBinding)
	rbClient := rc.client.RbacV1().RoleBindings(namespace)
	ret, err := rbClient.Create(context.TODO(), rb, metav1.CreateOptions{})

	if err != nil {
		fmt.Printf("%#v", err)
		// FIXME: Applyにしたかったところだけどメッチャクチャめんどくさそうなのでとりあえず握りつぶす
		if errors.IsAlreadyExists(err) {
			return ret, nil
		}
		panic(err)
		// return nil, err
	}

	return ret, nil
}

func (rc *ResourceCreater) createRole(roleName string, namespace string) (metav1.Object, error) {

	scheme := runtime.NewScheme()
	codecFactory := serializer.NewCodecFactory(scheme)
	deserializer := codecFactory.UniversalDeserializer()

	f := config.SearchManifest(roleName)
	if len(f) == 0 {
		return nil, fmt.Errorf("v1/Role %s manifest file not found", roleName)
	}

	str, err := claim.NormalizeDNS1123(roleName + "-" + rc.claim.GetSubject())
	values := struct{ RoleName string }{
		RoleName: str,
	}

	data, err := templateExecute(f, values)
	if err != nil {
		return nil, fmt.Errorf("Role manifest template execute fail: %w", err)
	}
	o, _, err := deserializer.Decode(data, nil, &rbacv1.Role{})
	if err != nil {
		return nil, err
	}
	role := o.(*rbacv1.Role)
	roleClient := rc.client.RbacV1().Roles(namespace)

	if role.Labels == nil {
		role.Labels = make(map[string]string)
	}
	if err := rc.mergeLabels(role.Labels); err != nil {
		panic(err)
	}

	if role.Annotations == nil {
		role.Annotations = make(map[string]string)
	}
	if err := rc.mergeAnnotations(role.Annotations); err != nil {
		panic(err)
	}

	ret, err := roleClient.Create(context.TODO(), role, metav1.CreateOptions{})

	if err != nil {
		// FIXME: Applyにしたかったところだけどメッチャクチャめんどくさそうなのでとりあえず握りつぶす
		// Resourceが既に存在する場合は握りつぶす
		if errors.IsAlreadyExists(err) {
			return ret, nil
		}
		return nil, err
	}
	return ret, nil
}

func (rc *ResourceCreater) mergeLabels(dist map[string]string) error {
	if err := mergo.Map(&dist, rc.claim.GetLabels()); err != nil {
		return err
	}

	// 必須ラベルがないと動作しないので追加
	if err := mergo.Map(&dist, _requireLabels); err != nil {
		return err
	}
	return nil
}

func (rc *ResourceCreater) mergeAnnotations(dist map[string]string) error {
	if err := mergo.Map(&dist, rc.claim.GetAnnotations()); err != nil {
		return err
	}
	return nil
}

func templateExecute(tmplFile string, values interface{}) ([]byte, error) {
	r, w := io.Pipe()
	defer r.Close()

	go func() error {
		defer w.Close()
		tmpl, err := template.ParseFiles(tmplFile)
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, values)
		if err != nil {
			panic(err)
		}
		return nil
	}()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.Bytes(), nil
}
