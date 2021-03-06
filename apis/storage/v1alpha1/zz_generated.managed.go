/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by angryjet. DO NOT EDIT.

package v1alpha1

import (
	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// GetBindingPhase of this S3BucketPolicy.
func (mg *S3BucketPolicy) GetBindingPhase() runtimev1alpha1.BindingPhase {
	return mg.Status.GetBindingPhase()
}

// GetClaimReference of this S3BucketPolicy.
func (mg *S3BucketPolicy) GetClaimReference() *corev1.ObjectReference {
	return mg.Spec.ClaimReference
}

// GetClassReference of this S3BucketPolicy.
func (mg *S3BucketPolicy) GetClassReference() *corev1.ObjectReference {
	return mg.Spec.ClassReference
}

// GetCondition of this S3BucketPolicy.
func (mg *S3BucketPolicy) GetCondition(ct runtimev1alpha1.ConditionType) runtimev1alpha1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this S3BucketPolicy.
func (mg *S3BucketPolicy) GetDeletionPolicy() runtimev1alpha1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetProviderConfigReference of this S3BucketPolicy.
func (mg *S3BucketPolicy) GetProviderConfigReference() *runtimev1alpha1.Reference {
	return mg.Spec.ProviderConfigReference
}

/*
GetProviderReference of this S3BucketPolicy.
Deprecated: Use GetProviderConfigReference.
*/
func (mg *S3BucketPolicy) GetProviderReference() *runtimev1alpha1.Reference {
	return mg.Spec.ProviderReference
}

// GetReclaimPolicy of this S3BucketPolicy.
func (mg *S3BucketPolicy) GetReclaimPolicy() runtimev1alpha1.ReclaimPolicy {
	return mg.Spec.ReclaimPolicy
}

// GetWriteConnectionSecretToReference of this S3BucketPolicy.
func (mg *S3BucketPolicy) GetWriteConnectionSecretToReference() *runtimev1alpha1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetBindingPhase of this S3BucketPolicy.
func (mg *S3BucketPolicy) SetBindingPhase(p runtimev1alpha1.BindingPhase) {
	mg.Status.SetBindingPhase(p)
}

// SetClaimReference of this S3BucketPolicy.
func (mg *S3BucketPolicy) SetClaimReference(r *corev1.ObjectReference) {
	mg.Spec.ClaimReference = r
}

// SetClassReference of this S3BucketPolicy.
func (mg *S3BucketPolicy) SetClassReference(r *corev1.ObjectReference) {
	mg.Spec.ClassReference = r
}

// SetConditions of this S3BucketPolicy.
func (mg *S3BucketPolicy) SetConditions(c ...runtimev1alpha1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this S3BucketPolicy.
func (mg *S3BucketPolicy) SetDeletionPolicy(r runtimev1alpha1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetProviderConfigReference of this S3BucketPolicy.
func (mg *S3BucketPolicy) SetProviderConfigReference(r *runtimev1alpha1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

/*
SetProviderReference of this S3BucketPolicy.
Deprecated: Use SetProviderConfigReference.
*/
func (mg *S3BucketPolicy) SetProviderReference(r *runtimev1alpha1.Reference) {
	mg.Spec.ProviderReference = r
}

// SetReclaimPolicy of this S3BucketPolicy.
func (mg *S3BucketPolicy) SetReclaimPolicy(r runtimev1alpha1.ReclaimPolicy) {
	mg.Spec.ReclaimPolicy = r
}

// SetWriteConnectionSecretToReference of this S3BucketPolicy.
func (mg *S3BucketPolicy) SetWriteConnectionSecretToReference(r *runtimev1alpha1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}
