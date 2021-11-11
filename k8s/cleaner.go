package k8s

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var ManagedLabelMap = map[string]string{
	"cinderella":                   "user",
	"app.kubernetes.io/managed-by": "cinderella",
}

type CleanupError struct {
	msg string
	err error
}

const (
	minTickEvery = time.Second * 10
	// maxValidPeriod = time.Hour
	maxValidPeriod = time.Second * 20
)

var _ error = &CleanupError{}

func (m CleanupError) Error() string {
	return fmt.Sprintf("%s: %s", m.msg, m.err)
}

func (m CleanupError) Unwrap() error {
	return m.err
}

type Cleaner struct {
	client         kubernetes.Interface
	tickEvery      time.Duration
	maxValidPeriod time.Duration
}

func NewCleaner(client kubernetes.Interface, tickEvery time.Duration) (Cleaner, error) {
	return Cleaner{
		client:         client,
		tickEvery:      tickEvery,
		maxValidPeriod: maxValidPeriod,
	}, nil
}

func (c *Cleaner) Start(ctx context.Context) error {

	if c.tickEvery < (time.Second * 10) {
		c.tickEvery = time.Second * 10
	}

	t := time.NewTicker(c.tickEvery)
	defer t.Stop()

	for {
		select {
		case now := <-t.C:
			c.CleanupResources(ctx, now)
		case <-ctx.Done():
			fmt.Println("Stop cleaner")
			return ctx.Err()
		}
	}
}

type deleteError struct {
	err    error
	target metav1.Object
}

type deleteErrors []deleteError

func (d deleteErrors) Add(err error, target metav1.Object) deleteErrors {
	return append(d, deleteError{
		err:    err,
		target: target,
	})
}

// 一応外部から手動で呼び出せるようにしておく
func (c *Cleaner) CleanupResources(ctx context.Context, now time.Time) {

	deletedObjects := []metav1.Object{}
	errs := deleteErrors{}

	roleList, err := c.listRoles(ctx)
	if err != nil {
		panic(err)
	}

	for i, v := range roleList.Items {
		del, err := c.deleteResource(ctx, &v)
		if err != nil {
			errs.Add(err, &roleList.Items[i])
		} else if del {
			deletedObjects = append(deletedObjects, &roleList.Items[i])
		}
	}

	rbList, err := c.listRoleBindings(ctx)
	if err != nil {
		panic(err)
	}

	for i, v := range rbList.Items {
		del, err := c.deleteResource(ctx, &v)
		if err != nil {
			errs.Add(err, &rbList.Items[i])
		} else if del {
			deletedObjects = append(deletedObjects, &rbList.Items[i])
		}
	}

	saList, err := c.listServiceAccounts(ctx)
	if err != nil {
		panic(err)
	}

	for i, v := range saList.Items {
		del, err := c.deleteResource(ctx, &v)
		if err != nil {
			errs.Add(err, &saList.Items[i])
		} else if del {
			deletedObjects = append(deletedObjects, &saList.Items[i])
		}
	}

	publishCleanupEvent(deletedObjects, errs)
}

func (c *Cleaner) listRoles(ctx context.Context) (*rbacv1.RoleList, error) {
	list, err := c.client.RbacV1().Roles(metav1.NamespaceAll).List(ctx, c.getListOptions())
	if err != nil {
		return nil, &CleanupError{
			msg: "Failed to list managed role resources",
			err: err,
		}
	}
	return list, nil
}

func (c *Cleaner) listRoleBindings(ctx context.Context) (*rbacv1.RoleBindingList, error) {
	list, err := c.client.RbacV1().RoleBindings(metav1.NamespaceAll).List(ctx, c.getListOptions())
	if err != nil {
		return nil, &CleanupError{
			msg: "Failed to list managed roleBindings resources",
			err: err,
		}
	}
	return list, nil
}

func (c *Cleaner) listServiceAccounts(ctx context.Context) (*corev1.ServiceAccountList, error) {

	list, err := c.client.CoreV1().ServiceAccounts(_serviceAccountNamespace).List(ctx, c.getListOptions())
	if err != nil {
		return nil, &CleanupError{
			msg: "Failed to list managed roleBindings resources",
			err: err,
		}
	}
	return list, nil
}

func (c *Cleaner) isExpired(createdAt time.Time) bool {
	//TODO: annotationに個別のExpireを記述してそれを参照するようにする
	// expire := getByAnnotation()
	// return expire.Before(time.Now()) || createdAt.Add(c.maxValidPeriod).Before(time.Now())

	// fmt.Println("---")
	// fmt.Printf("createdAt: %s\n", createdAt)
	// fmt.Printf("expiredAt: %s\n", createdAt.Add(c.maxValidPeriod))
	// fmt.Printf("now: %s\n", time.Now())
	// fmt.Println("---")

	return createdAt.Add(c.maxValidPeriod).Before(time.Now())
}

func (c *Cleaner) deleteResource(ctx context.Context, obj metav1.Object) (bool, error) {
	if !c.isExpired(obj.GetCreationTimestamp().Time) {
		fmt.Printf("%s/%s is valid\n", obj.GetNamespace(), obj.GetName())
		return false, nil
	}

	var err error
	switch obj.(type) {
	case *rbacv1.Role:
		role := obj.(*rbacv1.Role)
		client := c.client.RbacV1().Roles(role.Namespace)
		err = client.Delete(ctx, role.Name, metav1.DeleteOptions{})
	case *rbacv1.RoleBinding:
		rb := obj.(*rbacv1.RoleBinding)
		client := c.client.RbacV1().RoleBindings(rb.Namespace)
		err = client.Delete(ctx, rb.Name, metav1.DeleteOptions{})
	case *corev1.ServiceAccount:
		sa := obj.(*corev1.ServiceAccount)
		client := c.client.CoreV1().ServiceAccounts(sa.Namespace)
		err = client.Delete(ctx, sa.Name, metav1.DeleteOptions{})
	default:
		//誰かが手動でlabel書いたケースとかで一応ありえる
		return false, fmt.Errorf("unknown resource type: %w", err)
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Cleaner) getListOptions() metav1.ListOptions {
	labelSelector := &metav1.LabelSelector{
		MatchLabels: _managedResourceLabels,
	}

	return metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(labelSelector),
	}
}

// func logObjectInfo(){
// 		fmt.Printf("Name: %s\n", v.Name)
// 		fmt.Printf("Namespace: %s\n", v.Namespace)
// 		fmt.Printf("CreatedAt: %s\n", v.CreationTimestamp.Format(time.RFC3339))
// 		fmt.Printf("Expired: %t\n", c.isExpired(v.CreationTimestamp.Time))
// }
