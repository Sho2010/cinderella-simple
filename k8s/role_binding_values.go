package k8s

type RoleBindingValues struct {
	Subject                 string
	ServiceAccountName      string
	ServiceAccountNamespace string
	RoleBindingName         string
	RoleName                string
}

