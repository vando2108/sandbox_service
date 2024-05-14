# Project details
## Project layout
Here is a brief explanation of my project structure, which is based on the guidelines from ```golang-standards/project-layout```.
```
.
├── Makefile                          # run 'make generate' to generate pb
├── cmd                               
│   ├── README.md                     # Read the README.md for more doc
│   ├── sandbox     
│   │   └── main.go                  
│   └── script                      
│       └── tool.go                   
├── configs                           # Configuration files
├── go.mod                       
├── go.sum                          
├── internal                                        # Private application and library code
│   ├── app                                         # Core application logic
│   │   └── sandbox                                 # Sandbox-specific logic
│   │       ├── app.go                              # Main sandbox application code
│   │       ├── handler                               # Handlers for various routes and functionalities
│   │       │   ├── authenticator_handler.go          # Handler for authentication
│   │       │   ├── error.go                          # Error handling logic
│   │       │   └── sandbox_handler.go                # Handler for sandbox-specific routes
│   │       ├── model                               # Data models
│   │       │   ├── environment.go          
│   │       │   ├── environment_status.go     
│   │       │   └── user.go                  
│   │       └── repository                          # Data repository interfaces and implementations
│   │           ├── implement                         # Implementations of repositories
│   │           │   ├── memory                          # In-memory implementations
│   │           │   │   ├── mem_nonce_cache.go   
│   │           │   │   └── mem_user_repository.go 
│   │           │   └── redis                         # Redis implementations
│   │           │       └── redis_nonce_cache.go   
│   │           └── interface                         # Repository interfaces
│   │               ├── environment.go         
│   │               ├── nonce_cache_repository.go
│   │               └── user_repository.go      
│   └── work_queue                                  # Work queue implementation
│       └── work_queue.go             
├── pb                              
│   ├── authenticator.pb.go     
│   ├── authenticator_grpc.pb.go     
│   ├── error.pb.go                  
│   ├── sandbox.pb.go               
│   └── sandbox_grpc.pb.go          
├── proto                             # Protocol buffer definitions
│   ├── authenticator.proto         
│   ├── error.proto               
│   └── sandbox.proto                
├── test                              # Test-related files and directories
└── utils                             # Utility functions and helpers
    └── utils.go                      # Main utility functions
```

## Project note
### Repository pattern
I chose to use the repository pattern to manage the data logic. Initially, I planned to use Redis for caching and PostgreSQL for the database. However, due to time constraints, I decided to adopt the repository pattern, which provides a flexible way to switch between in-memory storage and other databases.

### Data models
```
type User struct {
	ID        int
	Publickey string
}

type Environment struct {
	ID                string
	UserID            string
	HardwareInput     string
	SnapshotURL       string
	Images            []string
	ServiceNames      []string
	NumberOfInstances int32
}

type RequirmentStatus struct {
	ID            string
	RequirementID string
	Status        string
}
```
I want to keep the Environment information and their statuses in two separate tables to facilitate easy updates by workers for each environment's status. This design enables effective monitoring of all environments, allowing for tracking changes such as errors, starts, stops, and enqueue events.

These models should include timestamps and additional information, but for the purposes of this mock project, I've kept the definitions simple.

### Services
#### Response structure
```
message Response {
  bool Success = 1;				// = false if any error occurs when the server tries to execute your request
  ErrorCode ErrorCode = 2;		// Defined in proto/error.proto
  string ErrorMessage = 3;		// Defined in internal/app/sandbox/handler/error.go
}
```

#### Authenticator service
```
syntax = "proto3";

package authenticator;

import "error.proto";

option go_package = "./pb";

service Authenticator {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc NonceConfirm(NonceConfirmRequest) returns (NonceConfirmResponse);
}

message RegisterRequest { string Publickey = 1; }

message RegisterResponse {
  bool Success = 1;
  string HashedNonce = 2;
  ErrorCode ErrorCode = 3;
  string ErrorMessage = 4;
}

message NonceConfirmRequest {
  string Publickey = 1;
  string Nonce = 2;
}

message NonceConfirmResponse {
  bool Success = 1;
  ErrorCode ErrorCode = 2;
  string ErrorMessage = 3;
}

```

1. To start using our service users will need to register by sending their publickey through the ```Register``` method. For the server side, whenever the server receives this request it will start the process by these steps.
   	- Validate the publickey sent by the user.
	- Check if this user exists in our system or not (by checking the database).
	- Generate a random nonce and hash it by using the provided publickey.
	- Keep the draw value of nonce to cache service (publickey - nonce keypair).
2. After receiving the ```hashed_nonce``` return by the server, the user will need to decrypt it and send it to the server for confirmation.
   	- Validate the publickey sent by the user.
   	- Try to get the nonce from the cache service by using their publickey.
   	- In case nonce is matched, the server will insert a new entry to the user table.
   	- Response to the user (success, nonce not exist, miss match, ....)
  
#### Sandbox service
This service has the responsibility of managing all information about the environment.

```
syntax = "proto3";

package sandbox;

import "error.proto";

option go_package = "./pb";

message Requirements {
  string HardwareInput = 1;
  string SnapshotUrl = 2;
  repeated string Image = 3;
  repeated string ServiceName = 4;
  int32 NumberOfInstance = 5;
}

service Sandbox {
  rpc CreateNewEnvironment(CreateNewEnvironmentRequest)
      returns (CreateNewEnvironmentResponse);
}

message CreateNewEnvironmentRequest {
  string Publickey = 1;
  Requirements Requirements = 2;
  string Signature = 6;
}

message CreateNewEnvironmentResponse {
  bool Success = 1;
  string EnvID = 2;
  ErrorCode ErrorCode = 3;
  string ErrorMessage = 4;
}
```

```
type SandboxHandler struct {
	userRepo  *repoInterface.UserRepository
	workQueue *workqueue.WorkQueue
	pb.UnimplementedSandboxServer
}
```
This service handler contains a work_queue.

To create a new environment, users will send a request with their requirements, publickey, and their signature created by using their private key.
Whenever the server receives the ```CreateNewEnvironmentRequest```it will start by
- Validate the request.
- Check if the user has registered.
- Verify the signature.
- Insert new entry to environment table.
- Insert new entry to environment_status table with status = pending.
- Enqueue a new task, which will be handled later by our workers.

There's an issue I've identified while documenting this: the CreateNewEnvironmentRequest doesn't include any randomized elements. If a user loses their request data (which includes a signature), it becomes possible for someone else to repeatedly submit this request, potentially leading to abuse. To address this vulnerability, I think implementing a confirmation method similar to AuthenticatorService would be a prudent measure. This method would help ensure that each request is authenticated and unique, preventing spam and enhancing security.


## Notes
I also want to add some unit tests, but I don't have enough time to find a framework and learn to use it.

I also want to use .env and separate config files to give our server a flexible way to config (e.g. choose the cache service in-mem or redis). Have implemented a redis_nonce_cache and will keep it as a future work.
