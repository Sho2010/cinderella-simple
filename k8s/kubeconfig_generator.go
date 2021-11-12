package k8s

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/Sho2010/cinderella-simple/claim"
	"github.com/Sho2010/cinderella-simple/config"
	"github.com/Sho2010/cinderella-simple/encrypt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//go:embed templates/*
var _templates embed.FS

const (
	KubeconfigTemplate = "templates/kubeconfig.tmpl"
)

type KubeconfigGenerator struct {
	Client kubernetes.Interface
}

type KubeconfigValues struct {
	ClusterName string
	CA          string
	User        string
	Token       string
	Namespace   string
	Server      string
}

func (gen *KubeconfigGenerator) Generate(writer io.Writer, name, namespace string) error {
	sa, err := gen.findSA(name, namespace)
	if err != nil {
		return err
	}

	values, err := gen.buildFromSA(sa)
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFS(_templates, KubeconfigTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(writer, values)
	if err != nil {
		return err
	}

	return nil
}

func (gen *KubeconfigGenerator) buildFromSA(sa *v1.ServiceAccount) (KubeconfigValues, error) {
	ctx := context.TODO()

	if len(sa.Secrets) == 0 {
		return KubeconfigValues{}, fmt.Errorf("ServiceAccount referenced secret not found")
	}
	ref := sa.Secrets[0]

	//namespaceをrefから取らないのはSAに付随するsecretのObjectReferenceにはNameしか入ってないから
	secret, err := gen.Client.CoreV1().Secrets(sa.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})

	if err != nil {
		panic(err)
	}

	// NOTE secret.Data のBase64に関して
	// 基本的にdecode状態の[]byte がやってくるっぽいのでそのまま使いたい場合はstring
	// encodeする必要がある場合は自分でencodeしてやる必要がある
	// https://github.com/kubernetes/client-go/issues/198
	values := KubeconfigValues{
		ClusterName: "cluster",
		CA:          base64.StdEncoding.EncodeToString(secret.Data[v1.ServiceAccountRootCAKey]),
		User:        "cinderella",
		Token:       string(secret.Data[v1.ServiceAccountTokenKey]),
		Namespace:   string(secret.Data[v1.ServiceAccountNamespaceKey]),
		Server:      config.GetConfig().KubeServer,
	}
	return values, nil

}

func (gen *KubeconfigGenerator) findSA(name, namespace string) (*v1.ServiceAccount, error) {
	ctx := context.TODO()
	sa, err := gen.Client.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return sa, nil
}

func (gen *KubeconfigGenerator) storeKubeconfig(ctx context.Context, sa *v1.ServiceAccount) error {

	b := false
	b2 := false

	// NOTE: APIVersion, Kindを固定値で書いてるのはTypeMetaは基本的にクリアされるので既存のオブジェクトからは取れない
	// See: https://github.com/kubernetes/client-go/issues/308#issuecomment-378165425
	owner := []metav1.OwnerReference{
		{
			APIVersion:         "v1",
			Kind:               "ServiceAccount",
			Name:               sa.Name,
			UID:                sa.UID,
			Controller:         &b,
			BlockOwnerDeletion: &b2,
		},
	}

	imm := true
	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:         "test-zip",
			GenerateName: "cinderella",
			Namespace:    "glass",
			Labels:       map[string]string{},
			Annotations:  map[string]string{},

			// SAの従属オブジェクトとすることによってSAが消された時に自動で消されるようにする
			//https://kubernetes.io/ja/docs/concepts/workloads/controllers/garbage-collection/
			OwnerReferences: owner,
		},
		Immutable: &imm,
		Data:      map[string][]byte{},
		Type:      v1.SecretTypeOpaque,
	}

	_, err := gen.Client.CoreV1().Secrets("glass").Create(ctx, &secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	return nil
}

func CreateEncryptedFile(claim claim.Claim) (string, error) {
	tmpFile, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	defer tmpFile.Close()

	if err := WriteEncryptedFile(tmpFile, claim); err != nil {
		return "", fmt.Errorf("Create kubeconfig file fail: %w", err)
	}
	return tmpFile.Name(), nil
}

func WriteEncryptedFile(writer io.Writer, claim claim.Claim) error {
	r, w := io.Pipe()
	defer r.Close()

	//TODO: goroutineのエラーを補足する
	go func() {
		defer w.Close()
		client, _ := GetDefaultClient()
		gen := KubeconfigGenerator{
			Client: client,
		}
		sa, err := claim.GetServiceAccountName()
		if err != nil {
			//FIXME
			panic(err)
		}
		gen.Generate(w, sa, _serviceAccountNamespace)
	}()

	var enc encrypt.FileEncrypter

	enc, err := encrypt.CreateEncrypter(claim.GetEncryptType())
	if err != nil {
		return fmt.Errorf("Create encrypter failed: %w", err)
	}

	// TODO: もうちょっとパスワードの受け渡し方法考える
	if zip, ok := enc.(*encrypt.ZipEncrypter); ok {
		zip.Password = claim.GetZipPassword()
	}

	enc.Encrypt(writer, r)
	return nil
}
