# GoHaveAzureStorage
GoHaveAzureStorage is a Azure Storage Library for Go. As The project is brand new, it is currently in a an alpha state, where only Azure Table Storage is supported. Take a look at the road map, for the planned development on the library.

# Installation
```
go get github.com/ChristianHenrikReich/gohaveazurestorage
```

# Get Started
Until documentation is written, look at tablestorageproxy_test.go for usage.

Short example:
```Go
goHaveAzureStorage := NewWithDebug(Account, Key, false) // Setting last 3rd value true, will enable http req/res dumping
tableStorageProxy := goHaveAzureStorage.NewTableStorageProxy()

httpStatusCode := tableStorageProxy.CreateTable(table)
```

Or without debug:

```Go
goHaveAzureStorage := New(Account, Key)
tableStorageProxy := goHaveAzureStorage.NewTableStorageProxy()

httpStatusCode := tableStorageProxy.CreateTable(table)
```

# Road map
* Table Storage Proxy
  - Implement preflight command
  - Implement support for shared keys
* Get started on Blob Storage Proxy
* Get started on Queue Storage Proxy
* Get started on File Storage Proxy
