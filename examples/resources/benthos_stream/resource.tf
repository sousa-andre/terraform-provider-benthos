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