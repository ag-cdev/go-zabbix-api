# Go Zabbix API

Note: This library is adjusted for use with terraform-provider-zabbix and supports Zabbix 7.0 LTS.

[![GoDoc](https://godoc.org/github.com/tpretz/go-zabbix-api?status.svg)](https://godoc.org/github.com/tpretz/go-zabbix-api) [![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE) [![Build Status](https://travis-ci.org/tpretz/go-zabbix-api.svg?branch=master)](https://travis-ci.org/tpretz/go-zabbix-api)

This Go package provides access to the Zabbix API.

Tested on Zabbix 3.2, 3.4, 4.0, 4.2, 4.4, 5.0, 6.0 and 7.0.

This package supports multiple Zabbix resources: trigger, application, host group, host, item, template, LLD rules, graphs, macros, and proxies.

## Install

```bash
go get github.com/tpretz/go-zabbix-api
```

## Getting Started

```go
package main

import (
	"fmt"

	"github.com/tpretz/go-zabbix-api"
)

func main() {
	user := "MyZabbixUsername"
	pass := "MyZabbixPassword"
	api := zabbix.NewAPI("http://localhost/api_jsonrpc.php")
	api.Login(user, pass)

	res, err := api.Version()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected to Zabbix API v%s\n", res)
}
```

## Supported Zabbix Versions

- Zabbix 3.2+
- Zabbix 4.0+
- Zabbix 5.0+
- Zabbix 6.0+
- Zabbix 7.0+ (LTS)

### Zabbix 7.0 Changes

When using Zabbix 7.0+, be aware of these API changes:

- `data_type` and `delta` fields are no longer valid in items (automatically removed)
- SNMP item types: Use `SNMPAgent` (20) instead of deprecated SNMPv1Agent (1) and SNMPv2Agent (4)
- Templates use `TemplateGroups` instead of `Groups` in Zabbix 7.2+
- Trigger expressions require format: `last(/host/key)>N`
- Application API is deprecated (use Tags instead)
- `hostid` parameter not allowed in item.update

## API Methods

### Hosts

- `HostsGet(params)` - Get hosts
- `HostGetByID(id)` - Get host by ID
- `HostGetByHost(host)` - Get host by hostname
- `HostsCreate(hosts)` - Create hosts
- `HostsUpdate(hosts)` - Update hosts
- `HostsDelete(hosts)` - Delete hosts

### Host Groups

- `HostGroupsGet(params)` - Get host groups
- `HostGroupGetByID(id)` - Get host group by ID
- `HostGroupsCreate(groups)` - Create host groups
- `HostGroupsUpdate(groups)` - Update host groups
- `HostGroupsDelete(groups)` - Delete host groups

### Template Groups (Zabbix 7.0+)

- `TemplateGroupsGet(params)` - Get template groups
- `TemplateGroupsCreate(groups)` - Create template groups
- `TemplateGroupsDelete(groups)` - Delete template groups

### Templates

- `TemplatesGet(params)` - Get templates
- `TemplateGetByID(id)` - Get template by ID
- `TemplatesCreate(templates)` - Create templates
- `TemplatesUpdate(templates)` - Update templates
- `TemplatesDelete(templates)` - Delete templates
- `TemplatesMassAdd(templates, hostGroups)` - Link templates to host groups

### Items

- `ItemsGet(params)` - Get items
- `ItemGetByID(id)` - Get item by ID
- `ItemsCreate(items)` - Create items
- `ItemsUpdate(items)` - Update items
- `ItemsDelete(items)` - Delete items
- `ProtoItemsGet/Create/Update/Delete` - Item prototypes (LLD)

### LLD Rules

- `LLDsGet(params)` - Get LLD rules
- `LLDGetByID(id)` - Get LLD rule by ID
- `LLDsCreate(items)` - Create LLD rules
- `LLDsUpdate(items)` - Update LLD rules
- `LLDsDelete(items)` - Delete LLD rules

### Triggers

- `TriggersGet(params)` - Get triggers
- `TriggerGetByID(id)` - Get trigger by ID
- `TriggersCreate(triggers)` - Create triggers
- `TriggersUpdate(triggers)` - Update triggers
- `TriggersDelete(triggers)` - Delete triggers
- `ProtoTriggersGet/Create/Update/Delete` - Trigger prototypes

### Graphs

- `GraphsGet(params)` - Get graphs
- `GraphGetByID(id)` - Get graph by ID
- `GraphsCreate(graphs)` - Create graphs
- `GraphsUpdate(graphs)` - Update graphs
- `GraphsDelete(graphs)` - Delete graphs

### Applications

- `ApplicationsGet(params)` - Get applications
- `ApplicationGetByID(id)` - Get application by ID
- `ApplicationsCreate(apps)` - Create applications
- `ApplicationsDelete(apps)` - Delete applications

Note: Application API is deprecated in Zabbix 5.4+. Use Tags instead.

### Macros

- `MacrosGet(params)` - Get macros
- `MacroGetByID(id)` - Get macro by ID
- `MacrosCreate(macros)` - Create macros
- `MacrosUpdate(macros)` - Update macros
- `MacrosDelete(macros)` - Delete macros

### Proxies

- `ProxiesGet(params)` - Get proxies

### Utility Methods

- `Version()` - Get Zabbix API version
- `IsZabbix7OrGreater()` - Check if connected to Zabbix 7+
- `IsZabbix6OrGreater()` - Check if connected to Zabbix 6+

## Tests

### Considerations

You should run tests before using this package.
Zabbix API doesn't match documentation in some details, which change in patch releases.

Tests are not expected to be destructive, but you are advised to run them against a non-production instance or at least make a backup.

### Run Tests

```bash
export TEST_ZABBIX_URL=http://localhost:8080/zabbix/api_jsonrpc.php
export TEST_ZABBIX_USER=Admin
export TEST_ZABBIX_PASSWORD=zabbix
export TEST_ZABBIX_VERBOSE=1
go test -v
```

`TEST_ZABBIX_URL` may contain HTTP basic auth: `http://username:password@host/api_jsonrpc.php`. Also, in some setups URL should be like `http://host/zabbix/api_jsonrpc.php`.

### Docker Test Environment

A `docker-compose.yml` is provided for testing:

```bash
docker-compose up -d
# Wait for Zabbix to initialize
export TEST_ZABBIX_URL=http://localhost:8080/api_jsonrpc.php
export TEST_ZABBIX_USER=Admin
export TEST_ZABBIX_PASSWORD=zabbix
go test -v
```

## References

Documentation is available on [godoc.org](https://godoc.org/github.com/tpretz/go-zabbix-api).

License: Simplified BSD License (see [LICENSE](LICENSE)).
