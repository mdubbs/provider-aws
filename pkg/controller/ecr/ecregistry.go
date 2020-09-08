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

package ecr

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsecr "github.com/aws/aws-sdk-go-v2/service/ecr"
	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/provider-aws/apis/ecr/v1alpha1"
	awsv1alpha3 "github.com/crossplane/provider-aws/apis/v1alpha3"
	awsclients "github.com/crossplane/provider-aws/pkg/clients"
	"github.com/crossplane/provider-aws/pkg/clients/ecr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	controllerName = "ecregistry.aws.crossplane.io"
	finalizer      = "finalizer." + controllerName
)

const (
	errNotRegistry       = "managed resource is not a Registry custom resource"
	errGetProvider       = "cannot get provider"
	errGetProviderSecret = "cannot get provider secret"
	errECRClient         = "cannot create ECR client"
	errNotFound          = "cannot get repository"
	errCreateFailed      = "cannot create repository"
	errDeleteFailed      = "cannot delete repository"
	errUpdateFailed      = "cannot update repository"
)

var (
	ctx           = context.Background()
	result        = reconcile.Result{}
	resultRequeue = reconcile.Result{Requeue: true}
)

type connector struct {
	kube        client.Client
	newClientFn func(ctx context.Context, credentials []byte, region string, auth awsclients.AuthMethod) (ecr.Client, error)
}
type external struct {
	client ecr.Client
	kube   client.Client
}

// SetupRegistry adds a controller that reconciles Registry
func SetupRegistry(mgr ctrl.Manager, l logging.Logger) error {
	name := managed.ControllerName(v1alpha1.RegistryGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.Registry{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.RegistryGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), newClientFn: ecr.NewClient}),
			managed.WithInitializers(),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Registry)
	if !ok {
		return nil, errors.New(errNotRegistry)
	}

	p := &awsv1alpha3.Provider{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.Spec.ProviderReference.Name}, p); err != nil {
		return nil, errors.Wrap(err, errGetProvider)
	}

	if aws.BoolValue(p.Spec.UseServiceAccount) {
		ecrClient, err := c.newClientFn(ctx, []byte{}, p.Spec.Region, awsclients.UsePodServiceAccount)
		return &external{client: ecrClient, kube: c.kube}, errors.Wrap(err, errECRClient)
	}

	if p.GetCredentialsSecretReference() == nil {
		return nil, errors.New(errGetProviderSecret)
	}

	s := &corev1.Secret{}
	n := types.NamespacedName{Namespace: p.Spec.CredentialsSecretRef.Namespace, Name: p.Spec.CredentialsSecretRef.Name}
	if err := c.kube.Get(ctx, n, s); err != nil {
		return nil, errors.Wrap(err, errGetProviderSecret)
	}

	ecrClient, err := c.newClientFn(ctx, s.Data[p.Spec.CredentialsSecretRef.Key], p.Spec.Region, awsclients.UseProviderSecret)
	return &external{client: ecrClient, kube: c.kube}, errors.Wrap(err, errECRClient)
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Registry)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotRegistry)
	}

	res, err := e.client.DescribeRepositoriesRequest(&awsecr.DescribeRepositoriesInput{
		MaxResults: aws.Int64(1),
		RepositoryNames: []string{
			fmt.Sprintf("%s/%s", cr.Namespace, cr.Name),
		},
	}).Send(ctx)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(resource.Ignore(ecr.IsNotFound, err), errNotFound)
	}

	cr.Status.SetConditions(runtimev1alpha1.Available())
	cr.Status.RegistryID = *res.Repositories[0].RegistryId
	cr.Status.RepositoryName = *res.Repositories[0].RepositoryName
	cr.Status.RepositoryURI = *res.Repositories[0].RepositoryUri

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Registry)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotRegistry)
	}
	cr.SetConditions(runtimev1alpha1.Creating())

	_, err := e.client.CreateRepositoryRequest(&awsecr.CreateRepositoryInput{
		ImageScanningConfiguration: &awsecr.ImageScanningConfiguration{ScanOnPush: aws.Bool(false)},
		ImageTagMutability:         awsecr.ImageTagMutabilityMutable,
		RepositoryName:             aws.String(fmt.Sprintf("%s/%s", cr.Namespace, cr.Name)),
	}).Send(ctx)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateFailed)
	}

	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Registry)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotRegistry)
	}

	// attempt to update image scanning config (scan on push)
	_, err := e.client.PutImageScanningConfigurationRequest(&awsecr.PutImageScanningConfigurationInput{
		ImageScanningConfiguration: &awsecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(cr.Spec.ScanOnPush),
		},
		RegistryId:     aws.String(cr.Status.RegistryID),
		RepositoryName: aws.String(cr.Status.RepositoryName),
	}).Send(ctx)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateFailed)
	}

	// attempt to update image immutability
	updateImmutability := &awsecr.PutImageTagMutabilityInput{
		RegistryId:         aws.String(cr.Status.RegistryID),
		RepositoryName:     aws.String(cr.Status.RepositoryName),
		ImageTagMutability: awsecr.ImageTagMutabilityImmutable,
	}
	if !cr.Spec.ImmutableTags {
		updateImmutability.ImageTagMutability = awsecr.ImageTagMutabilityMutable
	}
	_, err = e.client.PutImageTagMutabilityRequest(updateImmutability).Send(ctx)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateFailed)
	}

	return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateFailed)
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Registry)
	if !ok {
		return errors.New(errNotRegistry)
	}
	if cr.Status.RepositoryURI == "" {
		return nil
	}
	cr.SetConditions(runtimev1alpha1.Deleting())

	_, err := e.client.DeleteRepositoryRequest(&awsecr.DeleteRepositoryInput{
		Force:          aws.Bool(true),
		RegistryId:     aws.String(cr.Status.RegistryID),
		RepositoryName: aws.String(cr.Status.RepositoryName),
	}).Send(ctx)
	return errors.Wrap(resource.Ignore(ecr.IsNotFound, err), errDeleteFailed)
}
