// WIP

package audit

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NOTE
// 多分EventRecorderを使っていい感じにイベントを記録する
// https://github.com/kubernetes/client-go/blob/master/tools/record/event.go#L317-L320

type KubernetesEventHandler struct {
	client kubernetes.Interface
}

func NewKubernetesEventHandler(client kubernetes.Interface) (KubernetesEventHandler, error) {
	return KubernetesEventHandler{
		client: client,
	}, nil
}

func (h *KubernetesEventHandler) Start(event <-chan AuditEvent) {
	fmt.Println("Start handler")

	for e := range event {
		fmt.Println("Audit event received and create k8s event")
		fmt.Println(e.GetMessage())
	}
}

func create() corev1.Event {
	e := corev1.Event{
		TypeMeta:            metav1.TypeMeta{},
		ObjectMeta:          metav1.ObjectMeta{},
		InvolvedObject:      corev1.ObjectReference{},
		Reason:              "",
		Message:             "",
		Source:              corev1.EventSource{},
		FirstTimestamp:      metav1.Time{},
		LastTimestamp:       metav1.Time{},
		Count:               0,
		Type:                "",
		EventTime:           metav1.MicroTime{},
		Series:              &corev1.EventSeries{},
		Action:              "",
		Related:             &corev1.ObjectReference{},
		ReportingController: "",
		ReportingInstance:   "",
	}
	return e
}
