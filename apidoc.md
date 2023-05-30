## Flow
### Get flow list

GET `/api/flow`

query:

```yaml
offset: number
limit: number
```

response:

```json
{
  "list": [
    {
      "id": 5,
      "name": "demo",
      "description": "First demo",
      "graph": {
        "nodes": [
          {
            "id": "echo_1",
            "width": 100,
            "height": 50,
            "position": {
              "x": 0,
              "y": 0
            },
            "type": "Echo",
            "data": {
              "name": {
                "en": "Echo"
              },
              "icon": "",
              "description": null,
              "source": {
                "type": "builtin",
                "cmd_type": "",
                "git_url": "",
                "go_script": {
                  "script": ""
                }
              },
              "input_anchors": null,
              "input_params": [
                {
                  "id": "",
                  "name": {
                    "en": "Message"
                  },
                  "key": "message",
                  "type": "string"
                }
              ],
              "output_anchors": null,
              "inputs": {
                "message": "Hello"
              }
            },
            "position_absolute": {
              "x": 0,
              "y": 0
            }
          }
        ],
        "output_node_id": "echo_1"
      },
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

### Create flow

body:
```json
{
    "name": "demo",
    "description": "First demo",
    "graph": {
        "nodes": [
            {
                "id":"echo_1",
                "type": "Echo",
                "width":100,
                "height": 50,
                "postition":{
                    "x": 40,
                    "y": 40
                },
                "data":{
                    "name": {
                        "en":"Echo"
                    },
                    "source":{
                        "type": "builtin"
                    },
                    "input_params":[
                        {
                            "name":{
                                "en":"Message"
                            },
                            "key": "message",
                            "type":"string"
                        }
                    ],
                    "inputs": {
                        "message": "Hello"
                    }
                }
            }
        ],
        "output_node_id": "echo_1"
    }
}
```

### Delete flow

query:
```yaml
id: number
```


### Run flow

body:
```json
{
  "id": 1
}
```

## Component

### Get component by keys
TODO 