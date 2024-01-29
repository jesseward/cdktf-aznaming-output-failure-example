## Summary

This demonstrates an issue with the creation of the `cdktf` bindings for the `"Azure/naming/azurerm"` module (https://github.com/Azure/terraform-azurerm-naming). The root of the issue is that the language specific functions generated for this module yields a token that references a key->value map type, and not the actual output (string) value.

For example, calling `ResourceGroupOutput()` yields `${module.resourceNaming.resource_group}`, where the `resource_group` output value is actually a map. When running tf plan, we're greeted with the error `module.resourceNaming.resource_group is object with 8 attributes`

The module reference in question can be found at  [main.tf](https://github.com/Azure/terraform-azurerm-naming/blob/8a1c8616d4cd05423e53c3260a016919ce0df33d/main.tf#L1869-L1878) and [output.tf](https://github.com/Azure/terraform-azurerm-naming/blob/8a1c8616d4cd05423e53c3260a016919ce0df33d/outputs.tf#L924-L927)

Generated cdktf.json is

```json
        "name": "${module.resourceNaming.resource_group}"
```

Needed cdktf.json is

```json
        "name": "${module.resourceNaming.resource_group[\"name\"]}"
```


## To reproduce

```sh
cdktf get # fetch the Go TF bindings
go run main.go # build, compile and run the module. This performs the Synth()..
cd cdktf.out/stacks/naming-output-failure-example # the cdk.tf.json is placed within
terraform init # use a local statefile for this test
terraform plan # yield the error
```

## View `tf plan` results

```sh
$ terraform plan
...
Plan: 2 to add, 0 to change, 0 to destroy.
╷
│ Error: Incorrect attribute value type
│
│   on cdk.tf.json line 44, in resource.azurerm_resource_group.resource_group:
│   44:         "name": "${module.resourceNaming.resource_group}"
│     ├────────────────
│     │ module.resourceNaming.resource_group is object with 8 attributes
│
│ Inappropriate value for attribute "name": string required.
```

## Example of the Synth'd `cdk.tf.json`

The following block of `json` is computed during `Synth()` execution. The `ResourceGroupOutput()` call yields `"name": "${module.resourceNaming.resource_group}"` which is a reference to the object (with 8 keys) instead of direct access to the map or `name` key.

Ideally `ResourceGroupOutput()` returns `*map[string]interface{}` that would allow us to fetch our desired key eg `name`.

```json
  "resource": {
    "azurerm_resource_group": {
      "resource_group": {
        "//": {
          "metadata": {
            "path": "naming-output-failure-example/resource_group",
            "uniqueId": "resource_group"
          }
        },
        "location": "Canada Central",
        "name": "${module.resourceNaming.resource_group}"
      }
    }
  },
  ```