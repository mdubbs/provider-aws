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

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	awsclients "github.com/crossplane/provider-aws/pkg/clients"
)

// Service defines the ECR Client Operations
// type Service interface {
// 	CreateOrUpdateRepository(repo *v1alpha1.Registry) error
// 	DeleteRepository(repo *v1alpha1.Registry) error
// }

// Client defines ECR Client operations
type Client interface {
	CreateRepositoryRequest(input *ecr.CreateRepositoryInput) ecr.CreateRepositoryRequest
	DescribeRepositoriesRequest(input *ecr.DescribeRepositoriesInput) ecr.DescribeRepositoriesRequest
	DeleteRepositoryRequest(input *ecr.DeleteRepositoryInput) ecr.DeleteRepositoryRequest
	PutImageScanningConfigurationRequest(input *ecr.PutImageScanningConfigurationInput) ecr.PutImageScanningConfigurationRequest
	PutImageTagMutabilityRequest(input *ecr.PutImageTagMutabilityInput) ecr.PutImageTagMutabilityRequest
}

// NewClient creates new Queue Client with provided AWS Configurations/Credentials
func NewClient(ctx context.Context, credentials []byte, region string, auth awsclients.AuthMethod) (Client, error) {
	cfg, err := auth(ctx, credentials, awsclients.DefaultSection, region)
	if cfg == nil {
		return nil, err
	}
	return ecr.New(*cfg), err
}

// IsNotFound checks if the error returned by AWS API says that the queue being probed doesn't exist
func IsNotFound(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == ecr.ErrCodeRepositoryNotFoundException {
			return true
		}
	}

	return false
}
