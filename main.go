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
	// 		"name": "${module.resourceNaming.resource_group.name}"
	resourcegroup.NewResourceGroup(stack, jsii.String("resource_group"), &resourcegroup.ResourceGroupConfig{
		Name:     n.ResourceGroupOutput(), // ResourceGroupOutput returns a string instead of a Map object.
		Location: jsii.String("Canada Central"),
	})

	// See comment at https://github.com/hashicorp/terraform-cdk/issues/3477#issuecomment-1926338050
	//  Fn_Lookup states: retrieves the value of a single element from a map, given its key. If the given key does not exist,
	//					   the given default value is returned instead.
	// Using the following will generate the correct cdk.tf.json output. For example... Notice the name field is now referenced correctly.
	// "resource_group_workaround": {
	// 	"//": {
	// 	"metadata": {
	// 		"path": "naming-output-failure-example/resource_group_workaround",
	// 		"uniqueId": "resource_group_workaround"
	// 	}
	// 	},
	// 	"location": "Canada Central",
	// 	"name": "${module.resourceNaming.resource_group.name}"
	// }
	resourcegroup.NewResourceGroup(stack, jsii.String("resource_group_workaround"), &resourcegroup.ResourceGroupConfig{
		Name:     cdktf.Token_AsString(cdktf.Fn_Lookup(n.ResourceGroupOutput(), jsii.String("name"), nil), nil),
		Location: jsii.String("Canada Central"),
	})
	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	NewNamingRenderingFailureStack(app, "naming-output-failure-example")
	app.Synth()
}
