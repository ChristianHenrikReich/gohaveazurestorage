# GoHaveAzureStorage
GoHaveAzureStorage is a Azure Storage Library for Go. As The project is brand new, it is currently in a an alpha state, where only Azure Table Storage is supported. Take a look at the road map, for the planned development on the library.

# Installation
```
go get github.com/ChristianHenrikReich/gohaveazurestorage
```

# Get Started
More detailed information about usage is being written in the [Wiki](https://github.com/ChristianHenrikReich/gohaveazurestorage/wiki/1.-Home), and tablestorageproxy_test.go is also a good source for information.

```Go

package main

import (
	"fmt"
	"gohaveazurestorage"
)

//Either the primary or the secondary key found on Azure Portal
var key = "PrimaryOrSecondaryKey"

//Storage account name
var account = "Account"

func main() {
	// Create an instance of the lib
	goHaveAzureStorage := gohaveazurestorage.New(account, key)
  
  // Example with debug, if last parameter is true then request/response sessions dumps is written to console
  //goHaveAzureStorage := gohaveazurestorage.NewWithDebug(account, key, true)

	// From the lib instace, we can create multiple client instances
	tableStorage := goHaveAzureStorage.NewTableStorage()

	//Creating a table
	httpStatusCode := tableStorage.CreateTable("Table")
	if httpStatusCode != 201 {
		fmt.Println("Create table error")
	}
}

```

# Road map
* Table Storage
  - Implement preflight command
  - Implement support for shared keys
* Get started on Blob Storage
* Get started on Queue Storage
* Get started on File Storage
