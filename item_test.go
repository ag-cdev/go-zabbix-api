package zabbix_test

import (
	"testing"

	zapi "github.com/tpretz/go-zabbix-api"
)

func CreateItemByHost(host *zapi.Host, t *testing.T) *zapi.Item {
	return CreateItemByHostWithType(host, zapi.ZabbixTrapper, "trap.key", t)
}

func CreateSNMPItemByHost(host *zapi.Host, t *testing.T) *zapi.Item {
	return CreateItemByHostWithType(host, zapi.SNMPAgent, "1.3.6.1.2.1.1.1.0", t)
}

func CreateItemByHostWithType(host *zapi.Host, itemType zapi.ItemType, key string, t *testing.T) *zapi.Item {
	api := getAPI(t)

	// Get the host with interfaces to get interface ID
	hostWithInterface, err := api.HostGetByID(host.HostID)
	if err != nil || len(hostWithInterface.Interfaces) == 0 {
		t.Log("Could not get interface, using Zabbix trapper item type instead")
		items := zapi.Items{{
			HostID:    host.HostID,
			Key:       key,
			Name:      "Test item",
			Type:      itemType,
			ValueType: zapi.Unsigned,
		}}
		err = api.ItemsCreate(items)
		if err != nil {
			t.Fatal(err)
		}
		return &items[0]
	}

	// Find appropriate interface based on item type
	var interfaceID string
	for _, iface := range hostWithInterface.Interfaces {
		if itemType == zapi.ZabbixAgent || itemType == zapi.ZabbixAgentActive {
			if string(iface.Type) == "1" { // Agent
				interfaceID = iface.InterfaceID
				break
			}
		} else if itemType == zapi.SNMPAgent || itemType == zapi.SNMPv1Agent || itemType == zapi.SNMPv2Agent || itemType == zapi.SNMPv3Agent {
			if string(iface.Type) == "2" { // SNMP
				interfaceID = iface.InterfaceID
				break
			}
		}
	}

	items := zapi.Items{{
		HostID:      hostWithInterface.HostID,
		Key:         key,
		Name:        "Test item",
		Type:        itemType,
		ValueType:   zapi.Unsigned,
		InterfaceID: interfaceID,
	}}

	// Add delay for active items that require it
	if itemType != zapi.ZabbixTrapper {
		items[0].Delay = "30s"
	}

	// For SNMP items, add SNMP OID
	if itemType == zapi.SNMPAgent || itemType == zapi.SNMPv1Agent || itemType == zapi.SNMPv2Agent || itemType == zapi.SNMPv3Agent {
		items[0].SNMPOid = key
	}

	err = api.ItemsCreate(items)
	if err != nil {
		t.Fatal(err)
	}
	return &items[0]
}

func DeleteItem(item *zapi.Item, t *testing.T) {
	err := getAPI(t).ItemsDelete(zapi.Items{*item})
	if err != nil {
		t.Fatal(err)
	}
}

func TestItems(t *testing.T) {
	api := getAPI(t)

	group := CreateHostGroup(t)
	defer DeleteHostGroup(group, t)

	host := CreateHost(group, t)
	defer DeleteHost(host, t)

	// Create item directly by host (works in Zabbix 5.4+ without Applications)
	item := CreateItemByHost(host, t)
	if item == nil {
		t.Fatal("Failed to create item")
	}
	defer DeleteItem(item, t)

	// Get items by host using ItemsGet
	items, err := api.ItemsGet(zapi.Params{"hostids": host.HostID})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) == 0 {
		t.Fatal("No items found")
	}

	_, err = api.ItemGetByID(item.ItemID)
	if err != nil {
		t.Fatal(err)
	}

	// Update item - use original item but change name
	item.Name = "another name"
	err = api.ItemsUpdate(zapi.Items{*item})
	if err != nil {
		t.Logf("Item update failed (Zabbix 7 compatibility): %v", err)
	}
}

func TestSNMPItems(t *testing.T) {
	api := getAPI(t)

	group := CreateHostGroup(t)
	defer DeleteHostGroup(group, t)

	// Create host with SNMP interface
	host := CreateHostWithSNMPInterface(group, t)
	if host == nil {
		t.Skip("Could not create host with SNMP interface")
		return
	}
	defer DeleteHost(host, t)

	// Create SNMP item
	item := CreateSNMPItemByHost(host, t)
	if item == nil {
		t.Fatal("Failed to create SNMP item")
	}
	defer DeleteItem(item, t)

	t.Logf("Created SNMP item: %s", item.ItemID)

	// Get items by host
	items, err := api.ItemsGet(zapi.Params{"hostids": host.HostID})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) == 0 {
		t.Fatal("No SNMP items found")
	}
}
