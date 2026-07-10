package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "135c0fd2-233e-4da4-bc1d-eac95325b187"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "deha0036",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variables
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists and is connected to the VM
	actualNics := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
	assert.Contains(t, actualNics, nicName)

	// Confirm the VM is running the correct Ubuntu version
	actualImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
	assert.Equal(t, "Canonical", actualImage.Publisher)
	assert.Equal(t, "0001-com-ubuntu-server-jammy", actualImage.Offer)
	assert.Equal(t, "22_04-lts-gen2", actualImage.SKU)
}