package awsclient

// TODO: Retry upon API failure

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	kubeclientpkg "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	awsCredsSecretIDKey     = "aws_access_key_id"
	awsCredsSecretAccessKey = "aws_secret_access_key"
)

// NewAwsClientInput input for new aws client
type NewAwsClientInput struct {
	AwsCredsSecretIDKey     string
	AwsCredsSecretAccessKey string
	AwsToken                string
	AwsRegion               string
	SecretName              string
	NameSpace               string
}

// Client wraps for AWS SDK (for easier testing)
type Client interface {
	/*
	 * EC2-related Functions
	 */
	// EC2 - to create the security group for the new admin api
	// we can get the instance IDs from Node objects.
	AuthorizeSecurityGroupIngress(*ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
	// for rh-api.
	CreateSecurityGroup(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error)
	// for removing a formerly approved CIDR block from the rh-api. security group
	DeleteSecurityGroup(*ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error)
	// to determine if we need to create the rh-api. security group
	DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	// for removing a formerly approved CIDR block from the rh-api. security group
	RevokeSecurityGroupIngress(*ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error)
	// DescribeSubnets to find subnet for master nodes for incoming elb
	DescribeSubnets(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	// CreateTags to apply tags to EC2 resources
	CreateTags(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error)

	DescribeInstanceStatus(*ec2.DescribeInstanceStatusInput) (*ec2.DescribeInstanceStatusOutput, error)

	// Helper extensions
	// ec2
	SubnetNameToSubnetIDLookup([]string) ([]string, error)
	SubnetIDToVPCLookup([]string) ([]string, error)
	ApplyTagsToResources([]string, map[string]string) error
	setLoadBalancerSecurityGroup(string, *ec2.SecurityGroup) error
	findSecurityGroupByName(string) (*ec2.SecurityGroup, error)
	findSecurityGroupByID(string) (*ec2.SecurityGroup, error)
	createSecurityGroup(string, string, map[string]string) (*ec2.SecurityGroup, error)
	removeIngressRulesFromSecurityGroup(*ec2.SecurityGroup, []*ec2.IpPermission) error
	addIngressRulesToSecurityGroup(*ec2.SecurityGroup, []*ec2.IpPermission) error
	EnsureCIDRAccess(string, string, string, []string, map[string]string) error
}

type AwsClient struct {
	ec2Client ec2iface.EC2API
	//	route53Client route53iface.Route53API
	//	elbClient     elbiface.ELBAPI
	//	elbv2Client   elbv2iface.ELBV2API
}

func NewClient(accessID, accessSecret, token, region string) (*AwsClient, error) {
	awsConfig := &aws.Config{Region: aws.String(region)}
	if token == "" {
		os.Setenv("AWS_ACCESS_KEY_ID", accessID)
		os.Setenv("AWS_SECRET_ACCESS_KEY", accessSecret)
	} else {
		awsConfig.Credentials = credentials.NewStaticCredentials(accessID, accessSecret, token)
	}
	s, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	return &AwsClient{
		ec2Client: ec2.New(s),
		//	elbClient:     elb.New(s),
		//	elbv2Client:   elbv2.New(s),
		//	route53Client: route53.New(s),
	}, nil
}

// GetAWSClient generates an awsclient
// function must include region
// Pass in token if sessions requires a token
// if it includes a secretName and nameSpace it will create credentials from that secret data
// If it includes awsCredsSecretIDKey and awsCredsSecretAccessKey it will build credentials from those
func GetAWSClient(kubeClient kubeclientpkg.Client, input NewAwsClientInput) (*AwsClient, error) {

	// error if region is not included
	if input.AwsRegion == "" {
		return nil, fmt.Errorf("getAWSClient:NoRegion: %v", input.AwsRegion)
	}

	if input.SecretName != "" && input.NameSpace != "" {
		secret := &corev1.Secret{}
		err := kubeClient.Get(context.TODO(),
			types.NamespacedName{
				Name:      input.SecretName,
				Namespace: input.NameSpace,
			},
			secret)
		if err != nil {
			return nil, err
		}
		accessKeyID, ok := secret.Data[awsCredsSecretIDKey]
		if !ok {
			return nil, fmt.Errorf("AWS credentials secret %v did not contain key %v",
				input.SecretName, awsCredsSecretIDKey)
		}
		secretAccessKey, ok := secret.Data[awsCredsSecretAccessKey]
		if !ok {
			return nil, fmt.Errorf("AWS credentials secret %v did not contain key %v",
				input.SecretName, awsCredsSecretAccessKey)
		}

		AwsClient, err := NewClient(string(accessKeyID), string(secretAccessKey), input.AwsToken, input.AwsRegion)
		if err != nil {
			return nil, err
		}
		return AwsClient, nil
	}

	if input.AwsCredsSecretIDKey == "" && input.AwsCredsSecretAccessKey != "" {
		return nil, fmt.Errorf("getAWSClient: NoAwsCredentials or Secret %v", input)
	}

	AwsClient, err := NewClient(input.AwsCredsSecretIDKey, input.AwsCredsSecretAccessKey, input.AwsToken, input.AwsRegion)
	if err != nil {
		return nil, err
	}
	return AwsClient, nil
}

func (c *AwsClient) AuthorizeSecurityGroupIngress(i *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	return c.ec2Client.AuthorizeSecurityGroupIngress(i)
}
func (c *AwsClient) CreateSecurityGroup(i *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
	return c.ec2Client.CreateSecurityGroup(i)
}
func (c *AwsClient) DescribeInstanceStatus(i *ec2.DescribeInstanceStatusInput) (*ec2.DescribeInstanceStatusOutput, error) {
	return c.ec2Client.DescribeInstanceStatus(i)
}
func (c *AwsClient) DeleteSecurityGroup(i *ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
	return c.ec2Client.DeleteSecurityGroup(i)
}
func (c *AwsClient) DescribeSecurityGroups(i *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return c.ec2Client.DescribeSecurityGroups(i)
}
func (c *AwsClient) RevokeSecurityGroupIngress(i *ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error) {
	return c.ec2Client.RevokeSecurityGroupIngress(i)
}
func (c *AwsClient) DescribeSubnets(i *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return c.ec2Client.DescribeSubnets(i)
}
func (c *AwsClient) CreateTags(i *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return c.ec2Client.CreateTags(i)
}
