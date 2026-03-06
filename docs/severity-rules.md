# Severity Rules

## OpenAPI

| Change                                  | Severity     |
| --------------------------------------- | ------------ |
| Endpoint / method removed               | breaking     |
| Endpoint / method added                 | non-breaking |
| Parameter removed                       | breaking     |
| Parameter added                         | non-breaking |
| Parameter type changed                  | breaking     |
| Parameter required: optional → required | breaking     |
| Parameter required: required → optional | non-breaking |
| Request body removed                    | breaking     |
| Response code removed                   | breaking     |
| Field removed                           | breaking     |
| Field added                             | non-breaking |
| Field type changed                      | breaking     |
| Field required: optional → required     | breaking     |

## GraphQL

| Change                                              | Severity     |
| --------------------------------------------------- | ------------ |
| Type removed                                        | breaking     |
| Type added                                          | non-breaking |
| Type kind changed (e.g. Object → Interface)         | breaking     |
| Output field removed                                | breaking     |
| Output field added                                  | non-breaking |
| Output field deprecated                             | info         |
| Output field type: non-null → nullable (`T!` → `T`) | breaking     |
| Output field type: nullable → non-null (`T` → `T!`) | non-breaking |
| Argument removed                                    | breaking     |
| Argument added (required)                           | breaking     |
| Argument added (optional)                           | non-breaking |
| Enum value removed                                  | breaking     |
| Enum value added                                    | non-breaking |
| Union member removed                                | breaking     |
| Union member added                                  | non-breaking |
| Input field removed                                 | breaking     |
| Input field added (required)                        | breaking     |
| Input field added (optional)                        | non-breaking |
| Interface removed from type                         | breaking     |
| Interface added to type                             | non-breaking |

## gRPC

| Change                                    | Severity     |
| ----------------------------------------- | ------------ |
| Service removed                           | breaking     |
| Service added                             | non-breaking |
| RPC removed                               | breaking     |
| RPC added                                 | non-breaking |
| RPC request type changed                  | breaking     |
| RPC response type changed                 | breaking     |
| RPC streaming mode changed                | breaking     |
| Message removed                           | breaking     |
| Message added                             | non-breaking |
| Field removed                             | breaking     |
| Field added                               | non-breaking |
| Field type changed                        | breaking     |
| Field number changed                      | breaking     |
| Field label changed (singular ↔ repeated) | breaking     |
