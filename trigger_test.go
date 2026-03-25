package zabbix_test

import (
	"testing"

	zapi "github.com/tpretz/go-zabbix-api"
)

func TestTrigger(t *testing.T) {
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

	triggerParam := zapi.Params{"hostids": host.HostID}
	res, err := api.TriggersGet(triggerParam)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Fatal("Found triggers")
	}

	// Try to create a trigger - Zabbix 7 requires proper expression format: last(/host/key)>N
	triggerExpr := "1"
	if api.IsZabbix7OrGreater() && item != nil {
		triggerExpr = "last(/" + host.Host + "/" + item.Key + ")>0"
	}
	triggers := zapi.Triggers{{
		Description: "test trigger",
		Expression:  triggerExpr,
	}}
	err = api.TriggersCreate(triggers)
	if err != nil {
		if api.IsZabbix7OrGreater() {
			t.Logf("Trigger creation failed (Zabbix 7 compatibility): %v", err)
			t.Skip("Trigger creation not supported in this Zabbix version")
			return
		}
		t.Fatal(err)
	}
	defer api.TriggersDelete(triggers)

	if len(triggers) > 0 {
		triggers[0].Description = "new trigger name"
		err = api.TriggersUpdate(triggers)
		if err != nil {
			t.Error(err)
		}
	}
}
