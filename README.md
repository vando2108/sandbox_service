# Project layout
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
