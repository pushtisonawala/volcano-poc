package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	GroupName     = "scheduling.volcano.sh"
	GroupVersion  = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&NamespaceQueue{},
		&NamespaceQueueList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
