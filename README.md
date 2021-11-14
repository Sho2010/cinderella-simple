# cinderella-simple

## Slack Integration

| Command | Description |
| --- | --- |
| SLACK_BOT_TOKEN | Require for API |
| SLACK_APP_TOKEN | Require for socket mode |


## Manifest Template

`.yaml.tmpl`

e.g. `default-role.yaml.tmpl`


```example.yaml.tmpl
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cinderella-tempolary-role
  labels:
    app.kubernetes.io/managed-by: cinderella
  annotations:
    cinderella/claim.user: "{{.User}}"
rules:
- apiGroups:
  - ""
  resources:
  - namespaces
  - nodes
  - events
  verbs: ["get", "list", "watch"]
- apiGroups:
  - ""
  resources:
  - configmaps
  - endpoints
  - pods
  - pods/log
  - pods/exec
  - replicationcontrollers
  - services
  verbs: ["create", "update", "get", "list", "watch", "patch"]

- apiGroups: ["extensions", "apps"]
  resources:
  - deployments
  - replicasets
  - statefulsets
  verbs: ["create", "update", "get", "list", "watch", "patch"]

- apiGroups: ["batch"]
  resources:
  - jobs
  - cronjobs
  verbs: ["create", "update", "delete", "get", "list", "watch", "patch"]

- apiGroups: ["autoscaling"]
  resources: ["horizontalpodautoscalers"]
  verbs: ["create", "update", "get", "list", "watch", "patch"]

- apiGroups:
  - rbac.authorization.k8s.io
  resources:
  - rolebindings
  - roles
  verbs: ["get", "list", "watch"]
```

