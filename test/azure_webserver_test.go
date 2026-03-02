package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "11fd1bcf-1efc-4704-9a6a-7f83a6b77dae"

func TestAzureWebServer(t *testing.T) {
	testFolder := "."

	test_structure.RunTestStage(t, "setup", func() {
		terraformOptions := newTerraformOptions()

		// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
		terraform.InitAndApply(t, terraformOptions)
		test_structure.SaveTerraformOptions(t, testFolder, terraformOptions)
	})

	test_structure.RunTestStage(t, "validate", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, testFolder)

		t.Run("TestAzureLinuxVMCreation", func(t *testing.T) {
			vmName := terraform.Output(t, terraformOptions, "vm_name")
			resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

			// Confirm VM exists
			assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
		})

		t.Run("TestAzureWebServerNicAttached", func(t *testing.T) {
			vmName := terraform.Output(t, terraformOptions, "vm_name")
			nicName := terraform.Output(t, terraformOptions, "nic_name")
			resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

			// Confirm NIC exists and is connected to the VM
			assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID))
			assert.Contains(t, azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID), nicName)
		})

		t.Run("TestAzureWebServerUbuntuVersion", func(t *testing.T) {
			vmName := terraform.Output(t, terraformOptions, "vm_name")
			resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

			// Confirm the VM is using the expected Ubuntu image
			vmImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
			assert.Equal(t, "Canonical", vmImage.Publisher)
			assert.Equal(t, "0001-com-ubuntu-server-jammy", vmImage.Offer)
			assert.Equal(t, "22_04-lts-gen2", vmImage.SKU)
		})
	})

	test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, testFolder)
		terraform.Destroy(t, terraformOptions)
	})
}

func newTerraformOptions() *terraform.Options {
	return &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "040811108",
			"region":      "canadacentral",
		},
	}
}
