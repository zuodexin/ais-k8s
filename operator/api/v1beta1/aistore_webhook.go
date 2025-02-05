// Package contains declaration of AIS Kubernetes Custom Resource Definitions
/*
 * Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
 */

package v1beta1

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var aistorelog = logf.Log.WithName("aistore-resource")

func (r *AIStore) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:path=/validate-ais-nvidia-com-v1beta1-aistore,mutating=false,failurePolicy=fail,sideEffects=None,groups=ais.nvidia.com,resources=aistores,verbs=create;update,versions=v1beta1,name=vaistore.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &AIStore{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *AIStore) ValidateCreate() (admission.Warnings, error) {
	aistorelog.Info("validate create", "name", r.Name)

	if r.Spec.Size <= 0 {
		return nil, errInvalidClusterSize(r.Spec.Size)
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AIStore) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	aistorelog.Info("validate update", "name", r.Name)

	return nil, r.vup(old)
}

func (r *AIStore) vup(old runtime.Object) error {
	if r.Spec.Size <= 0 {
		return errInvalidClusterSize(r.Spec.Size)
	}

	prev, ok := old.(*AIStore)
	if !ok {
		return nil
	}

	// TODO: better validation, maybe using AIS IterFields?
	if !reflect.DeepEqual(r.Spec.ProxySpec, prev.Spec.ProxySpec) {
		return errCannotUpdateSpec("proxySpec")
	}

	if !reflect.DeepEqual(r.Spec.TargetSpec, prev.Spec.TargetSpec) {
		// TODO: For now, just log error if target specs are updated. Eventually, implement
		// logic that compares target specs accurately.
		err := errCannotUpdateSpec("targetSpec")
		aistorelog.Error(err, fmt.Sprintf("%v != %v", r.Spec.TargetSpec, prev.Spec.TargetSpec))
	}

	if !reflect.DeepEqual(r.Spec.DisablePodAntiAffinity, prev.Spec.DisablePodAntiAffinity) {
		return errCannotUpdateSpec("disablePodAntiAffinity")
	}

	if r.Spec.EnableExternalLB != prev.Spec.EnableExternalLB {
		return errCannotUpdateSpec("enableExternalLB")
	}

	if r.Spec.HostpathPrefix != prev.Spec.HostpathPrefix {
		return errCannotUpdateSpec("hostpathPrefix")
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AIStore) ValidateDelete() (admission.Warnings, error) {
	aistorelog.Info("validate delete", "name", r.Name)
	return nil, nil
}

// errors
func errInvalidClusterSize(size int32) error {
	return fmt.Errorf("invalid cluster size %d, should be at least 1", size)
}

func errCannotUpdateSpec(specName string) error {
	return fmt.Errorf("cannot update spec %q for an existing cluster", specName)
}
