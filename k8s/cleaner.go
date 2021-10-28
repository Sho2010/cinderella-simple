package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/Sho2010/cinderella-simple/config"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var ManagedLabelMap = map[string]string{
	"cinderella":                   "user",
	"app.kubernetes.io/managed-by": "cinderella",
}

type Cleaner struct {
	client kubernetes.Interface
}

func NewCleaner(client kubernetes.Interface) (Cleaner, error) {
	return Cleaner{
		client: client,
	}, nil
}

func (cleaner *Cleaner) Start(ctx context.Context) error {

	// TODO: configに直接の依存関係を持たない
	// TODO: gen child context
	config := config.GetConfig()
	t := time.NewTicker(config.CleanupTickEvery)
	defer t.Stop()

	for {
		select {
		case now := <-t.C:
			cleaner.cleanupResources(ctx, now)
		case <-ctx.Done():
			fmt.Println("Stop cleaner")
			return ctx.Err()
		}
	}
}

func (c *Cleaner) getListOptions() metav1.ListOptions {
	labelSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{"app.kubernetes.io/managed-by": "cinderella"},
	}

	return metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(labelSelector),
	}
}

func (c *Cleaner) cleanupResources(ctx context.Context, now time.Time) {
	//TODO: Implement me
	fmt.Println(now.Format(time.RFC3339))
	c.listManagedResources(ctx)
}

func (c *Cleaner) listManagedResources(ctx context.Context) (*rbacv1.RoleList, error) {
	list, err := c.client.RbacV1().Roles(metav1.NamespaceAll).List(ctx, c.getListOptions())
	if err != nil {
		panic(err)
	}

	//debug code: audit eventをFireさせる
	RaiseCleanupEvent("Raise audit event test")
	for _, v := range list.Items {
		fmt.Printf("Name: %s\n", v.Name)
	}

	return list, nil
}
