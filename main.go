package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"

	"cdk.tf/go/stack/generated/hashicorp/azurerm/provider"
	"cdk.tf/go/stack/generated/hashicorp/azurerm/resourcegroup"
	"cdk.tf/go/stack/generated/naming"
)

func NewNamingRenderingFailureStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	subscriptionId := cdktf.NewTerraformVariable(stack, jsii.String("subscriptionId"),
		&cdktf.TerraformVariableConfig{Type: jsii.String("string"),
			Description: jsii.String("Subscription ID"),
		})

	provider.NewAzurermProvider(stack, jsii.String("provider"), &provider.AzurermProviderConfig{
		Features:       &provider.AzurermProviderFeatures{},
		SubscriptionId: subscriptionId.ToString(),
	})

	n := naming.NewNaming(stack, jsii.String("resourceNaming"), &naming.NamingConfig{
		Suffix: jsii.Strings("suffixexample"),
	})

	// the following yields invalid cdk.tf.json output. cdk.tf.json contains the following ...
	// 		"location": "Canada Central",
	// 		"name": "${module.resourceNaming.resource_group}"
	// ... which is invalid because the name is not a string literal, we need to access the map key as such
	// 		"name": "${module.resourceNaming.resource_group.name[\"name\"]}"
	resourcegroup.NewResourceGroup(stack, jsii.String("resource_group"), &resourcegroup.ResourceGroupConfig{
		Name:     n.ResourceGroupOutput(), // ResourceGroupOutput returns a string instead of a Map object.
		Location: jsii.String("Canada Central"),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	NewNamingRenderingFailureStack(app, "naming-output-failure-example")
	app.Synth()
}
