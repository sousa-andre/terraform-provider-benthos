---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "benthos_stream Resource - benthos"
subcategory: ""
description: |-
  
---

# benthos_stream (Resource)



## Example Usage

```terraform
resource "benthos_stream" "this" {
  id = "stream4"
  config = jsonencode({
    "input" : {
      "file" : {
        "paths" : ["/tmp/input.json"]
      }
    },
    "output" : {
      "file" : {
        "path" : "/tmp/output.json"
      }
    },
    "pipeline" : {
      "processors" : [
        {
          "mapping" : "root = content().uppercase()"
        }
      ]
    }
  })
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config` (String)

### Read-Only

- `id` (String) The ID of this resource.
