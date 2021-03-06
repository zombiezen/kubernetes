/*
Copyright 2017 The Kubernetes Authors.

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

package azure

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/disk"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/golang/glog"

	"k8s.io/client-go/util/flowcontrol"
)

// VirtualMachinesClient defines needed functions for azure compute.VirtualMachinesClient
type VirtualMachinesClient interface {
	CreateOrUpdate(resourceGroupName string, VMName string, parameters compute.VirtualMachine, cancel <-chan struct{}) (<-chan compute.VirtualMachine, <-chan error)
	Get(resourceGroupName string, VMName string, expand compute.InstanceViewTypes) (result compute.VirtualMachine, err error)
	List(resourceGroupName string) (result compute.VirtualMachineListResult, err error)
	ListNextResults(lastResults compute.VirtualMachineListResult) (result compute.VirtualMachineListResult, err error)
}

// InterfacesClient defines needed functions for azure network.InterfacesClient
type InterfacesClient interface {
	CreateOrUpdate(resourceGroupName string, networkInterfaceName string, parameters network.Interface, cancel <-chan struct{}) (<-chan network.Interface, <-chan error)
	Get(resourceGroupName string, networkInterfaceName string, expand string) (result network.Interface, err error)
	GetVirtualMachineScaleSetNetworkInterface(resourceGroupName string, virtualMachineScaleSetName string, virtualmachineIndex string, networkInterfaceName string, expand string) (result network.Interface, err error)
}

// LoadBalancersClient defines needed functions for azure network.LoadBalancersClient
type LoadBalancersClient interface {
	CreateOrUpdate(resourceGroupName string, loadBalancerName string, parameters network.LoadBalancer, cancel <-chan struct{}) (<-chan network.LoadBalancer, <-chan error)
	Delete(resourceGroupName string, loadBalancerName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error)
	Get(resourceGroupName string, loadBalancerName string, expand string) (result network.LoadBalancer, err error)
	List(resourceGroupName string) (result network.LoadBalancerListResult, err error)
	ListNextResults(lastResult network.LoadBalancerListResult) (result network.LoadBalancerListResult, err error)
}

// PublicIPAddressesClient defines needed functions for azure network.PublicIPAddressesClient
type PublicIPAddressesClient interface {
	CreateOrUpdate(resourceGroupName string, publicIPAddressName string, parameters network.PublicIPAddress, cancel <-chan struct{}) (<-chan network.PublicIPAddress, <-chan error)
	Delete(resourceGroupName string, publicIPAddressName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error)
	Get(resourceGroupName string, publicIPAddressName string, expand string) (result network.PublicIPAddress, err error)
	List(resourceGroupName string) (result network.PublicIPAddressListResult, err error)
	ListNextResults(lastResults network.PublicIPAddressListResult) (result network.PublicIPAddressListResult, err error)
}

// SubnetsClient defines needed functions for azure network.SubnetsClient
type SubnetsClient interface {
	CreateOrUpdate(resourceGroupName string, virtualNetworkName string, subnetName string, subnetParameters network.Subnet, cancel <-chan struct{}) (<-chan network.Subnet, <-chan error)
	Delete(resourceGroupName string, virtualNetworkName string, subnetName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error)
	Get(resourceGroupName string, virtualNetworkName string, subnetName string, expand string) (result network.Subnet, err error)
	List(resourceGroupName string, virtualNetworkName string) (result network.SubnetListResult, err error)
}

// SecurityGroupsClient defines needed functions for azure network.SecurityGroupsClient
type SecurityGroupsClient interface {
	CreateOrUpdate(resourceGroupName string, networkSecurityGroupName string, parameters network.SecurityGroup, cancel <-chan struct{}) (<-chan network.SecurityGroup, <-chan error)
	Delete(resourceGroupName string, networkSecurityGroupName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error)
	Get(resourceGroupName string, networkSecurityGroupName string, expand string) (result network.SecurityGroup, err error)
	List(resourceGroupName string) (result network.SecurityGroupListResult, err error)
}

// VirtualMachineScaleSetsClient defines needed functions for azure compute.VirtualMachineScaleSetsClient
type VirtualMachineScaleSetsClient interface {
	CreateOrUpdate(resourceGroupName string, VMScaleSetName string, parameters compute.VirtualMachineScaleSet, cancel <-chan struct{}) (<-chan compute.VirtualMachineScaleSet, <-chan error)
	Get(resourceGroupName string, VMScaleSetName string) (result compute.VirtualMachineScaleSet, err error)
	List(resourceGroupName string) (result compute.VirtualMachineScaleSetListResult, err error)
	ListNextResults(lastResults compute.VirtualMachineScaleSetListResult) (result compute.VirtualMachineScaleSetListResult, err error)
	UpdateInstances(resourceGroupName string, VMScaleSetName string, VMInstanceIDs compute.VirtualMachineScaleSetVMInstanceRequiredIDs, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error)
}

// VirtualMachineScaleSetVMsClient defines needed functions for azure compute.VirtualMachineScaleSetVMsClient
type VirtualMachineScaleSetVMsClient interface {
	Get(resourceGroupName string, VMScaleSetName string, instanceID string) (result compute.VirtualMachineScaleSetVM, err error)
	GetInstanceView(resourceGroupName string, VMScaleSetName string, instanceID string) (result compute.VirtualMachineScaleSetVMInstanceView, err error)
	List(resourceGroupName string, virtualMachineScaleSetName string, filter string, selectParameter string, expand string) (result compute.VirtualMachineScaleSetVMListResult, err error)
	ListNextResults(lastResults compute.VirtualMachineScaleSetVMListResult) (result compute.VirtualMachineScaleSetVMListResult, err error)
}

// RoutesClient defines needed functions for azure network.RoutesClient
type RoutesClient interface {
	CreateOrUpdate(resourceGroupName string, routeTableName string, routeName string, routeParameters network.Route, cancel <-chan struct{}) (<-chan network.Route, <-chan error)
	Delete(resourceGroupName string, routeTableName string, routeName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error)
}

// RouteTablesClient defines needed functions for azure network.RouteTablesClient
type RouteTablesClient interface {
	CreateOrUpdate(resourceGroupName string, routeTableName string, parameters network.RouteTable, cancel <-chan struct{}) (<-chan network.RouteTable, <-chan error)
	Get(resourceGroupName string, routeTableName string, expand string) (result network.RouteTable, err error)
}

// StorageAccountClient defines needed functions for azure storage.AccountsClient
type StorageAccountClient interface {
	Create(resourceGroupName string, accountName string, parameters storage.AccountCreateParameters, cancel <-chan struct{}) (<-chan storage.Account, <-chan error)
	Delete(resourceGroupName string, accountName string) (result autorest.Response, err error)
	ListKeys(resourceGroupName string, accountName string) (result storage.AccountListKeysResult, err error)
	ListByResourceGroup(resourceGroupName string) (result storage.AccountListResult, err error)
	GetProperties(resourceGroupName string, accountName string) (result storage.Account, err error)
}

// DisksClient defines needed functions for azure disk.DisksClient
type DisksClient interface {
	CreateOrUpdate(resourceGroupName string, diskName string, diskParameter disk.Model, cancel <-chan struct{}) (<-chan disk.Model, <-chan error)
	Delete(resourceGroupName string, diskName string, cancel <-chan struct{}) (<-chan disk.OperationStatusResponse, <-chan error)
	Get(resourceGroupName string, diskName string) (result disk.Model, err error)
}

// azClientConfig contains all essential information to create an Azure client.
type azClientConfig struct {
	subscriptionID          string
	resourceManagerEndpoint string
	servicePrincipalToken   *adal.ServicePrincipalToken
	rateLimiter             flowcontrol.RateLimiter
}

// azVirtualMachinesClient implements VirtualMachinesClient.
type azVirtualMachinesClient struct {
	client      compute.VirtualMachinesClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzVirtualMachinesClient(config *azClientConfig) *azVirtualMachinesClient {
	virtualMachinesClient := compute.NewVirtualMachinesClient(config.subscriptionID)
	virtualMachinesClient.BaseURI = config.resourceManagerEndpoint
	virtualMachinesClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	virtualMachinesClient.PollingDelay = 5 * time.Second
	configureUserAgent(&virtualMachinesClient.Client)

	return &azVirtualMachinesClient{
		rateLimiter: config.rateLimiter,
		client:      virtualMachinesClient,
	}
}

func (az *azVirtualMachinesClient) CreateOrUpdate(resourceGroupName string, VMName string, parameters compute.VirtualMachine, cancel <-chan struct{}) (<-chan compute.VirtualMachine, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachinesClient.CreateOrUpdate(%q, %q): start", resourceGroupName, VMName)
	defer func() {
		glog.V(10).Infof("azVirtualMachinesClient.CreateOrUpdate(%q, %q): end", resourceGroupName, VMName)
	}()

	return az.client.CreateOrUpdate(resourceGroupName, VMName, parameters, cancel)
}

func (az *azVirtualMachinesClient) Get(resourceGroupName string, VMName string, expand compute.InstanceViewTypes) (result compute.VirtualMachine, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachinesClient.Get(%q, %q): start", resourceGroupName, VMName)
	defer func() {
		glog.V(10).Infof("azVirtualMachinesClient.Get(%q, %q): end", resourceGroupName, VMName)
	}()

	return az.client.Get(resourceGroupName, VMName, expand)
}

func (az *azVirtualMachinesClient) List(resourceGroupName string) (result compute.VirtualMachineListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachinesClient.List(%q): start", resourceGroupName)
	defer func() {
		glog.V(10).Infof("azVirtualMachinesClient.List(%q): end", resourceGroupName)
	}()

	return az.client.List(resourceGroupName)
}

func (az *azVirtualMachinesClient) ListNextResults(lastResults compute.VirtualMachineListResult) (result compute.VirtualMachineListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachinesClient.ListNextResults(%q): start", lastResults)
	defer func() {
		glog.V(10).Infof("azVirtualMachinesClient.ListNextResults(%q): end", lastResults)
	}()

	return az.client.ListNextResults(lastResults)
}

// azInterfacesClient implements InterfacesClient.
type azInterfacesClient struct {
	client      network.InterfacesClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzInterfacesClient(config *azClientConfig) *azInterfacesClient {
	interfacesClient := network.NewInterfacesClient(config.subscriptionID)
	interfacesClient.BaseURI = config.resourceManagerEndpoint
	interfacesClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	interfacesClient.PollingDelay = 5 * time.Second
	configureUserAgent(&interfacesClient.Client)

	return &azInterfacesClient{
		rateLimiter: config.rateLimiter,
		client:      interfacesClient,
	}
}

func (az *azInterfacesClient) CreateOrUpdate(resourceGroupName string, networkInterfaceName string, parameters network.Interface, cancel <-chan struct{}) (<-chan network.Interface, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azInterfacesClient.CreateOrUpdate(%q,%q): start", resourceGroupName, networkInterfaceName)
	defer func() {
		glog.V(10).Infof("azInterfacesClient.CreateOrUpdate(%q,%q): end", resourceGroupName, networkInterfaceName)
	}()

	return az.client.CreateOrUpdate(resourceGroupName, networkInterfaceName, parameters, cancel)
}

func (az *azInterfacesClient) Get(resourceGroupName string, networkInterfaceName string, expand string) (result network.Interface, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azInterfacesClient.Get(%q,%q): start", resourceGroupName, networkInterfaceName)
	defer func() {
		glog.V(10).Infof("azInterfacesClient.Get(%q,%q): end", resourceGroupName, networkInterfaceName)
	}()

	return az.client.Get(resourceGroupName, networkInterfaceName, expand)
}

func (az *azInterfacesClient) GetVirtualMachineScaleSetNetworkInterface(resourceGroupName string, virtualMachineScaleSetName string, virtualmachineIndex string, networkInterfaceName string, expand string) (result network.Interface, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azInterfacesClient.GetVirtualMachineScaleSetNetworkInterface(%q,%q,%q,%q): start", resourceGroupName, virtualMachineScaleSetName, virtualmachineIndex, networkInterfaceName)
	defer func() {
		glog.V(10).Infof("azInterfacesClient.GetVirtualMachineScaleSetNetworkInterface(%q,%q,%q,%q): end", resourceGroupName, virtualMachineScaleSetName, virtualmachineIndex, networkInterfaceName)
	}()

	return az.client.GetVirtualMachineScaleSetNetworkInterface(resourceGroupName, virtualMachineScaleSetName, virtualmachineIndex, networkInterfaceName, expand)
}

// azLoadBalancersClient implements LoadBalancersClient.
type azLoadBalancersClient struct {
	client      network.LoadBalancersClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzLoadBalancersClient(config *azClientConfig) *azLoadBalancersClient {
	loadBalancerClient := network.NewLoadBalancersClient(config.subscriptionID)
	loadBalancerClient.BaseURI = config.resourceManagerEndpoint
	loadBalancerClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	loadBalancerClient.PollingDelay = 5 * time.Second
	configureUserAgent(&loadBalancerClient.Client)

	return &azLoadBalancersClient{
		rateLimiter: config.rateLimiter,
		client:      loadBalancerClient,
	}
}

func (az *azLoadBalancersClient) CreateOrUpdate(resourceGroupName string, loadBalancerName string, parameters network.LoadBalancer, cancel <-chan struct{}) (<-chan network.LoadBalancer, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azLoadBalancersClient.CreateOrUpdate(%q,%q): start", resourceGroupName, loadBalancerName)
	defer func() {
		glog.V(10).Infof("azLoadBalancersClient.CreateOrUpdate(%q,%q): end", resourceGroupName, loadBalancerName)
	}()

	return az.client.CreateOrUpdate(resourceGroupName, loadBalancerName, parameters, cancel)
}

func (az *azLoadBalancersClient) Delete(resourceGroupName string, loadBalancerName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azLoadBalancersClient.Delete(%q,%q): start", resourceGroupName, loadBalancerName)
	defer func() {
		glog.V(10).Infof("azLoadBalancersClient.Delete(%q,%q): end", resourceGroupName, loadBalancerName)
	}()

	return az.client.Delete(resourceGroupName, loadBalancerName, cancel)
}

func (az *azLoadBalancersClient) Get(resourceGroupName string, loadBalancerName string, expand string) (result network.LoadBalancer, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azLoadBalancersClient.Get(%q,%q): start", resourceGroupName, loadBalancerName)
	defer func() {
		glog.V(10).Infof("azLoadBalancersClient.Get(%q,%q): end", resourceGroupName, loadBalancerName)
	}()

	return az.client.Get(resourceGroupName, loadBalancerName, expand)
}

func (az *azLoadBalancersClient) List(resourceGroupName string) (result network.LoadBalancerListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azLoadBalancersClient.List(%q): start", resourceGroupName)
	defer func() {
		glog.V(10).Infof("azLoadBalancersClient.List(%q): end", resourceGroupName)
	}()

	return az.client.List(resourceGroupName)
}

func (az *azLoadBalancersClient) ListNextResults(lastResult network.LoadBalancerListResult) (result network.LoadBalancerListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azLoadBalancersClient.ListNextResults(%q): start", lastResult)
	defer func() {
		glog.V(10).Infof("azLoadBalancersClient.ListNextResults(%q): end", lastResult)
	}()

	return az.client.ListNextResults(lastResult)
}

// azPublicIPAddressesClient implements PublicIPAddressesClient.
type azPublicIPAddressesClient struct {
	client      network.PublicIPAddressesClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzPublicIPAddressesClient(config *azClientConfig) *azPublicIPAddressesClient {
	publicIPAddressClient := network.NewPublicIPAddressesClient(config.subscriptionID)
	publicIPAddressClient.BaseURI = config.resourceManagerEndpoint
	publicIPAddressClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	publicIPAddressClient.PollingDelay = 5 * time.Second
	configureUserAgent(&publicIPAddressClient.Client)

	return &azPublicIPAddressesClient{
		rateLimiter: config.rateLimiter,
		client:      publicIPAddressClient,
	}
}

func (az *azPublicIPAddressesClient) CreateOrUpdate(resourceGroupName string, publicIPAddressName string, parameters network.PublicIPAddress, cancel <-chan struct{}) (<-chan network.PublicIPAddress, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azPublicIPAddressesClient.CreateOrUpdate(%q,%q): start", resourceGroupName, publicIPAddressName)
	defer func() {
		glog.V(10).Infof("azPublicIPAddressesClient.CreateOrUpdate(%q,%q): end", resourceGroupName, publicIPAddressName)
	}()

	return az.client.CreateOrUpdate(resourceGroupName, publicIPAddressName, parameters, cancel)
}

func (az *azPublicIPAddressesClient) Delete(resourceGroupName string, publicIPAddressName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azPublicIPAddressesClient.Delete(%q,%q): start", resourceGroupName, publicIPAddressName)
	defer func() {
		glog.V(10).Infof("azPublicIPAddressesClient.Delete(%q,%q): end", resourceGroupName, publicIPAddressName)
	}()

	return az.client.Delete(resourceGroupName, publicIPAddressName, cancel)
}

func (az *azPublicIPAddressesClient) Get(resourceGroupName string, publicIPAddressName string, expand string) (result network.PublicIPAddress, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azPublicIPAddressesClient.Get(%q,%q): start", resourceGroupName, publicIPAddressName)
	defer func() {
		glog.V(10).Infof("azPublicIPAddressesClient.Get(%q,%q): end", resourceGroupName, publicIPAddressName)
	}()

	return az.client.Get(resourceGroupName, publicIPAddressName, expand)
}

func (az *azPublicIPAddressesClient) List(resourceGroupName string) (result network.PublicIPAddressListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azPublicIPAddressesClient.List(%q): start", resourceGroupName)
	defer func() {
		glog.V(10).Infof("azPublicIPAddressesClient.List(%q): end", resourceGroupName)
	}()

	return az.client.List(resourceGroupName)
}

func (az *azPublicIPAddressesClient) ListNextResults(lastResults network.PublicIPAddressListResult) (result network.PublicIPAddressListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azPublicIPAddressesClient.ListNextResults(%q): start", lastResults)
	defer func() {
		glog.V(10).Infof("azPublicIPAddressesClient.ListNextResults(%q): end", lastResults)
	}()

	return az.client.ListNextResults(lastResults)
}

// azSubnetsClient implements SubnetsClient.
type azSubnetsClient struct {
	client      network.SubnetsClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzSubnetsClient(config *azClientConfig) *azSubnetsClient {
	subnetsClient := network.NewSubnetsClient(config.subscriptionID)
	subnetsClient.BaseURI = config.resourceManagerEndpoint
	subnetsClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	subnetsClient.PollingDelay = 5 * time.Second
	configureUserAgent(&subnetsClient.Client)

	return &azSubnetsClient{
		client:      subnetsClient,
		rateLimiter: config.rateLimiter,
	}
}

func (az *azSubnetsClient) CreateOrUpdate(resourceGroupName string, virtualNetworkName string, subnetName string, subnetParameters network.Subnet, cancel <-chan struct{}) (<-chan network.Subnet, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azSubnetsClient.CreateOrUpdate(%q,%q,%q): start", resourceGroupName, virtualNetworkName, subnetName)
	defer func() {
		glog.V(10).Infof("azSubnetsClient.CreateOrUpdate(%q,%q,%q): end", resourceGroupName, virtualNetworkName, subnetName)
	}()

	return az.client.CreateOrUpdate(resourceGroupName, virtualNetworkName, subnetName, subnetParameters, cancel)
}

func (az *azSubnetsClient) Delete(resourceGroupName string, virtualNetworkName string, subnetName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azSubnetsClient.Delete(%q,%q,%q): start", resourceGroupName, virtualNetworkName, subnetName)
	defer func() {
		glog.V(10).Infof("azSubnetsClient.Delete(%q,%q,%q): end", resourceGroupName, virtualNetworkName, subnetName)
	}()

	return az.client.Delete(resourceGroupName, virtualNetworkName, subnetName, cancel)
}

func (az *azSubnetsClient) Get(resourceGroupName string, virtualNetworkName string, subnetName string, expand string) (result network.Subnet, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azSubnetsClient.Get(%q,%q,%q): start", resourceGroupName, virtualNetworkName, subnetName)
	defer func() {
		glog.V(10).Infof("azSubnetsClient.Get(%q,%q,%q): end", resourceGroupName, virtualNetworkName, subnetName)
	}()

	return az.client.Get(resourceGroupName, virtualNetworkName, subnetName, expand)
}

func (az *azSubnetsClient) List(resourceGroupName string, virtualNetworkName string) (result network.SubnetListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azSubnetsClient.List(%q,%q): start", resourceGroupName, virtualNetworkName)
	defer func() {
		glog.V(10).Infof("azSubnetsClient.List(%q,%q): end", resourceGroupName, virtualNetworkName)
	}()

	return az.client.List(resourceGroupName, virtualNetworkName)
}

// azSecurityGroupsClient implements SecurityGroupsClient.
type azSecurityGroupsClient struct {
	client      network.SecurityGroupsClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzSecurityGroupsClient(config *azClientConfig) *azSecurityGroupsClient {
	securityGroupsClient := network.NewSecurityGroupsClient(config.subscriptionID)
	securityGroupsClient.BaseURI = config.resourceManagerEndpoint
	securityGroupsClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	securityGroupsClient.PollingDelay = 5 * time.Second
	configureUserAgent(&securityGroupsClient.Client)

	return &azSecurityGroupsClient{
		rateLimiter: config.rateLimiter,
		client:      securityGroupsClient,
	}
}

func (az *azSecurityGroupsClient) CreateOrUpdate(resourceGroupName string, networkSecurityGroupName string, parameters network.SecurityGroup, cancel <-chan struct{}) (<-chan network.SecurityGroup, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azSecurityGroupsClient.CreateOrUpdate(%q,%q): start", resourceGroupName, networkSecurityGroupName)
	defer func() {
		glog.V(10).Infof("azSecurityGroupsClient.CreateOrUpdate(%q,%q): end", resourceGroupName, networkSecurityGroupName)
	}()

	return az.client.CreateOrUpdate(resourceGroupName, networkSecurityGroupName, parameters, cancel)
}

func (az *azSecurityGroupsClient) Delete(resourceGroupName string, networkSecurityGroupName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azSecurityGroupsClient.Delete(%q,%q): start", resourceGroupName, networkSecurityGroupName)
	defer func() {
		glog.V(10).Infof("azSecurityGroupsClient.Delete(%q,%q): end", resourceGroupName, networkSecurityGroupName)
	}()

	return az.client.Delete(resourceGroupName, networkSecurityGroupName, cancel)
}

func (az *azSecurityGroupsClient) Get(resourceGroupName string, networkSecurityGroupName string, expand string) (result network.SecurityGroup, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azSecurityGroupsClient.Get(%q,%q): start", resourceGroupName, networkSecurityGroupName)
	defer func() {
		glog.V(10).Infof("azSecurityGroupsClient.Get(%q,%q): end", resourceGroupName, networkSecurityGroupName)
	}()

	return az.client.Get(resourceGroupName, networkSecurityGroupName, expand)
}

func (az *azSecurityGroupsClient) List(resourceGroupName string) (result network.SecurityGroupListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azSecurityGroupsClient.List(%q): start", resourceGroupName)
	defer func() {
		glog.V(10).Infof("azSecurityGroupsClient.List(%q): end", resourceGroupName)
	}()

	return az.client.List(resourceGroupName)
}

// azVirtualMachineScaleSetsClient implements VirtualMachineScaleSetsClient.
type azVirtualMachineScaleSetsClient struct {
	client      compute.VirtualMachineScaleSetsClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzVirtualMachineScaleSetsClient(config *azClientConfig) *azVirtualMachineScaleSetsClient {
	virtualMachineScaleSetsClient := compute.NewVirtualMachineScaleSetsClient(config.subscriptionID)
	virtualMachineScaleSetsClient.BaseURI = config.resourceManagerEndpoint
	virtualMachineScaleSetsClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	virtualMachineScaleSetsClient.PollingDelay = 5 * time.Second
	configureUserAgent(&virtualMachineScaleSetsClient.Client)

	return &azVirtualMachineScaleSetsClient{
		client:      virtualMachineScaleSetsClient,
		rateLimiter: config.rateLimiter,
	}
}

func (az *azVirtualMachineScaleSetsClient) CreateOrUpdate(resourceGroupName string, VMScaleSetName string, parameters compute.VirtualMachineScaleSet, cancel <-chan struct{}) (<-chan compute.VirtualMachineScaleSet, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachineScaleSetsClient.CreateOrUpdate(%q,%q): start", resourceGroupName, VMScaleSetName)
	defer func() {
		glog.V(10).Infof("azVirtualMachineScaleSetsClient.CreateOrUpdate(%q,%q): end", resourceGroupName, VMScaleSetName)
	}()

	return az.client.CreateOrUpdate(resourceGroupName, VMScaleSetName, parameters, cancel)
}

func (az *azVirtualMachineScaleSetsClient) Get(resourceGroupName string, VMScaleSetName string) (result compute.VirtualMachineScaleSet, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachineScaleSetsClient.Get(%q,%q): start", resourceGroupName, VMScaleSetName)
	defer func() {
		glog.V(10).Infof("azVirtualMachineScaleSetsClient.Get(%q,%q): end", resourceGroupName, VMScaleSetName)
	}()

	return az.client.Get(resourceGroupName, VMScaleSetName)
}

func (az *azVirtualMachineScaleSetsClient) List(resourceGroupName string) (result compute.VirtualMachineScaleSetListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachineScaleSetsClient.List(%q,%q): start", resourceGroupName)
	defer func() {
		glog.V(10).Infof("azVirtualMachineScaleSetsClient.List(%q,%q): end", resourceGroupName)
	}()

	return az.client.List(resourceGroupName)
}

func (az *azVirtualMachineScaleSetsClient) ListNextResults(lastResults compute.VirtualMachineScaleSetListResult) (result compute.VirtualMachineScaleSetListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachineScaleSetsClient.ListNextResults(%q): start", lastResults)
	defer func() {
		glog.V(10).Infof("azVirtualMachineScaleSetsClient.ListNextResults(%q): end", lastResults)
	}()

	return az.client.ListNextResults(lastResults)
}

func (az *azVirtualMachineScaleSetsClient) UpdateInstances(resourceGroupName string, VMScaleSetName string, VMInstanceIDs compute.VirtualMachineScaleSetVMInstanceRequiredIDs, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachineScaleSetsClient.UpdateInstances(%q,%q,%q): start", resourceGroupName, VMScaleSetName, VMInstanceIDs)
	defer func() {
		glog.V(10).Infof("azVirtualMachineScaleSetsClient.UpdateInstances(%q,%q,%q): end", resourceGroupName, VMScaleSetName, VMInstanceIDs)
	}()

	return az.client.UpdateInstances(resourceGroupName, VMScaleSetName, VMInstanceIDs, cancel)
}

// azVirtualMachineScaleSetVMsClient implements VirtualMachineScaleSetVMsClient.
type azVirtualMachineScaleSetVMsClient struct {
	client      compute.VirtualMachineScaleSetVMsClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzVirtualMachineScaleSetVMsClient(config *azClientConfig) *azVirtualMachineScaleSetVMsClient {
	virtualMachineScaleSetVMsClient := compute.NewVirtualMachineScaleSetVMsClient(config.subscriptionID)
	virtualMachineScaleSetVMsClient.BaseURI = config.resourceManagerEndpoint
	virtualMachineScaleSetVMsClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	virtualMachineScaleSetVMsClient.PollingDelay = 5 * time.Second
	configureUserAgent(&virtualMachineScaleSetVMsClient.Client)

	return &azVirtualMachineScaleSetVMsClient{
		client:      virtualMachineScaleSetVMsClient,
		rateLimiter: config.rateLimiter,
	}
}

func (az *azVirtualMachineScaleSetVMsClient) Get(resourceGroupName string, VMScaleSetName string, instanceID string) (result compute.VirtualMachineScaleSetVM, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachineScaleSetVMsClient.Get(%q,%q,%q): start", resourceGroupName, VMScaleSetName, instanceID)
	defer func() {
		glog.V(10).Infof("azVirtualMachineScaleSetVMsClient.Get(%q,%q,%q): end", resourceGroupName, VMScaleSetName, instanceID)
	}()

	return az.client.Get(resourceGroupName, VMScaleSetName, instanceID)
}

func (az *azVirtualMachineScaleSetVMsClient) GetInstanceView(resourceGroupName string, VMScaleSetName string, instanceID string) (result compute.VirtualMachineScaleSetVMInstanceView, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachineScaleSetVMsClient.GetInstanceView(%q,%q,%q): start", resourceGroupName, VMScaleSetName, instanceID)
	defer func() {
		glog.V(10).Infof("azVirtualMachineScaleSetVMsClient.GetInstanceView(%q,%q,%q): end", resourceGroupName, VMScaleSetName, instanceID)
	}()

	return az.client.GetInstanceView(resourceGroupName, VMScaleSetName, instanceID)
}

func (az *azVirtualMachineScaleSetVMsClient) List(resourceGroupName string, virtualMachineScaleSetName string, filter string, selectParameter string, expand string) (result compute.VirtualMachineScaleSetVMListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachineScaleSetVMsClient.List(%q,%q,%q): start", resourceGroupName, virtualMachineScaleSetName, filter)
	defer func() {
		glog.V(10).Infof("azVirtualMachineScaleSetVMsClient.List(%q,%q,%q): end", resourceGroupName, virtualMachineScaleSetName, filter)
	}()

	return az.client.List(resourceGroupName, virtualMachineScaleSetName, filter, selectParameter, expand)
}

func (az *azVirtualMachineScaleSetVMsClient) ListNextResults(lastResults compute.VirtualMachineScaleSetVMListResult) (result compute.VirtualMachineScaleSetVMListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azVirtualMachineScaleSetVMsClient.ListNextResults(%q,%q,%q): start", lastResults)
	defer func() {
		glog.V(10).Infof("azVirtualMachineScaleSetVMsClient.ListNextResults(%q,%q,%q): end", lastResults)
	}()

	return az.client.ListNextResults(lastResults)
}

// azRoutesClient implements RoutesClient.
type azRoutesClient struct {
	client      network.RoutesClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzRoutesClient(config *azClientConfig) *azRoutesClient {
	routesClient := network.NewRoutesClient(config.subscriptionID)
	routesClient.BaseURI = config.resourceManagerEndpoint
	routesClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	routesClient.PollingDelay = 5 * time.Second
	configureUserAgent(&routesClient.Client)

	return &azRoutesClient{
		client:      routesClient,
		rateLimiter: config.rateLimiter,
	}
}

func (az *azRoutesClient) CreateOrUpdate(resourceGroupName string, routeTableName string, routeName string, routeParameters network.Route, cancel <-chan struct{}) (<-chan network.Route, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azRoutesClient.CreateOrUpdate(%q,%q,%q): start", resourceGroupName, routeTableName, routeName)
	defer func() {
		glog.V(10).Infof("azRoutesClient.CreateOrUpdate(%q,%q,%q): end", resourceGroupName, routeTableName, routeName)
	}()

	return az.client.CreateOrUpdate(resourceGroupName, routeTableName, routeName, routeParameters, cancel)
}

func (az *azRoutesClient) Delete(resourceGroupName string, routeTableName string, routeName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azRoutesClient.Delete(%q,%q,%q): start", resourceGroupName, routeTableName, routeName)
	defer func() {
		glog.V(10).Infof("azRoutesClient.Delete(%q,%q,%q): end", resourceGroupName, routeTableName, routeName)
	}()

	return az.client.Delete(resourceGroupName, routeTableName, routeName, cancel)
}

// azRouteTablesClient implements RouteTablesClient.
type azRouteTablesClient struct {
	client      network.RouteTablesClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzRouteTablesClient(config *azClientConfig) *azRouteTablesClient {
	routeTablesClient := network.NewRouteTablesClient(config.subscriptionID)
	routeTablesClient.BaseURI = config.resourceManagerEndpoint
	routeTablesClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	routeTablesClient.PollingDelay = 5 * time.Second
	configureUserAgent(&routeTablesClient.Client)

	return &azRouteTablesClient{
		client:      routeTablesClient,
		rateLimiter: config.rateLimiter,
	}
}

func (az *azRouteTablesClient) CreateOrUpdate(resourceGroupName string, routeTableName string, parameters network.RouteTable, cancel <-chan struct{}) (<-chan network.RouteTable, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azRouteTablesClient.CreateOrUpdate(%q,%q): start", resourceGroupName, routeTableName)
	defer func() {
		glog.V(10).Infof("azRouteTablesClient.CreateOrUpdate(%q,%q): end", resourceGroupName, routeTableName)
	}()

	return az.client.CreateOrUpdate(resourceGroupName, routeTableName, parameters, cancel)
}

func (az *azRouteTablesClient) Get(resourceGroupName string, routeTableName string, expand string) (result network.RouteTable, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azRouteTablesClient.Get(%q,%q): start", resourceGroupName, routeTableName)
	defer func() {
		glog.V(10).Infof("azRouteTablesClient.Get(%q,%q): end", resourceGroupName, routeTableName)
	}()

	return az.client.Get(resourceGroupName, routeTableName, expand)
}

// azStorageAccountClient implements StorageAccountClient.
type azStorageAccountClient struct {
	client      storage.AccountsClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzStorageAccountClient(config *azClientConfig) *azStorageAccountClient {
	storageAccountClient := storage.NewAccountsClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	storageAccountClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	storageAccountClient.PollingDelay = 5 * time.Second
	configureUserAgent(&storageAccountClient.Client)

	return &azStorageAccountClient{
		client:      storageAccountClient,
		rateLimiter: config.rateLimiter,
	}
}

func (az *azStorageAccountClient) Create(resourceGroupName string, accountName string, parameters storage.AccountCreateParameters, cancel <-chan struct{}) (<-chan storage.Account, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azStorageAccountClient.Create(%q,%q): start", resourceGroupName, accountName)
	defer func() {
		glog.V(10).Infof("azStorageAccountClient.Create(%q,%q): end", resourceGroupName, accountName)
	}()

	return az.client.Create(resourceGroupName, accountName, parameters, cancel)
}

func (az *azStorageAccountClient) Delete(resourceGroupName string, accountName string) (result autorest.Response, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azStorageAccountClient.Delete(%q,%q): start", resourceGroupName, accountName)
	defer func() {
		glog.V(10).Infof("azStorageAccountClient.Delete(%q,%q): end", resourceGroupName, accountName)
	}()

	return az.client.Delete(resourceGroupName, accountName)
}

func (az *azStorageAccountClient) ListKeys(resourceGroupName string, accountName string) (result storage.AccountListKeysResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azStorageAccountClient.ListKeys(%q,%q): start", resourceGroupName, accountName)
	defer func() {
		glog.V(10).Infof("azStorageAccountClient.ListKeys(%q,%q): end", resourceGroupName, accountName)
	}()

	return az.client.ListKeys(resourceGroupName, accountName)
}

func (az *azStorageAccountClient) ListByResourceGroup(resourceGroupName string) (result storage.AccountListResult, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azStorageAccountClient.ListByResourceGroup(%q): start", resourceGroupName)
	defer func() {
		glog.V(10).Infof("azStorageAccountClient.ListByResourceGroup(%q): end", resourceGroupName)
	}()

	return az.client.ListByResourceGroup(resourceGroupName)
}

func (az *azStorageAccountClient) GetProperties(resourceGroupName string, accountName string) (result storage.Account, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azStorageAccountClient.GetProperties(%q,%q): start", resourceGroupName, accountName)
	defer func() {
		glog.V(10).Infof("azStorageAccountClient.GetProperties(%q,%q): end", resourceGroupName, accountName)
	}()

	return az.client.GetProperties(resourceGroupName, accountName)
}

// azDisksClient implements DisksClient.
type azDisksClient struct {
	client      disk.DisksClient
	rateLimiter flowcontrol.RateLimiter
}

func newAzDisksClient(config *azClientConfig) *azDisksClient {
	disksClient := disk.NewDisksClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	disksClient.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	disksClient.PollingDelay = 5 * time.Second
	configureUserAgent(&disksClient.Client)

	return &azDisksClient{
		client:      disksClient,
		rateLimiter: config.rateLimiter,
	}
}

func (az *azDisksClient) CreateOrUpdate(resourceGroupName string, diskName string, diskParameter disk.Model, cancel <-chan struct{}) (<-chan disk.Model, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azDisksClient.CreateOrUpdate(%q,%q): start", resourceGroupName, diskName)
	defer func() {
		glog.V(10).Infof("azDisksClient.CreateOrUpdate(%q,%q): end", resourceGroupName, diskName)
	}()

	return az.client.CreateOrUpdate(resourceGroupName, diskName, diskParameter, cancel)
}

func (az *azDisksClient) Delete(resourceGroupName string, diskName string, cancel <-chan struct{}) (<-chan disk.OperationStatusResponse, <-chan error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azDisksClient.Delete(%q,%q): start", resourceGroupName, diskName)
	defer func() {
		glog.V(10).Infof("azDisksClient.Delete(%q,%q): end", resourceGroupName, diskName)
	}()

	return az.client.Delete(resourceGroupName, diskName, cancel)
}

func (az *azDisksClient) Get(resourceGroupName string, diskName string) (result disk.Model, err error) {
	az.rateLimiter.Accept()
	glog.V(10).Infof("azDisksClient.Get(%q,%q): start", resourceGroupName, diskName)
	defer func() {
		glog.V(10).Infof("azDisksClient.Get(%q,%q): end", resourceGroupName, diskName)
	}()

	return az.client.Get(resourceGroupName, diskName)
}
