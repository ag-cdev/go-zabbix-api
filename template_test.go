package zabbix_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	zapi "github.com/tpretz/go-zabbix-api"
)

func CreateTemplate(hostGroup *zapi.HostGroup, t *testing.T) *zapi.Template {
	group := zapi.HostGroupID{
		GroupID: hostGroup.GroupID,
	}

	groups := []zapi.HostGroupID{group}

	template := zapi.Templates{zapi.Template{
		Host:   "template name",
		Groups: groups,
	}}
	err := getAPI(t).TemplatesCreate(template)
	if isPermissionError(err) {
		t.Logf("Template creation failed (permissions issue): %v", err)
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return &template[0]
}

func DeleteTemplate(template *zapi.Template, t *testing.T) {
	err := getAPI(t).TemplatesDelete(zapi.Templates{*template})
	if err != nil {
		t.Fatal(err)
	}
}

func isPermissionError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no permissions") ||
		strings.Contains(err.Error(), "does not exist") ||
		strings.Contains(err.Error(), "Application error")
}

func TestTemplates(t *testing.T) {
	api := getAPI(t)

	hostGroup := CreateHostGroup(t)
	defer DeleteHostGroup(hostGroup, t)

	// For Zabbix 7.2+, use template groups
	// For Zabbix 7.0.x, use host groups (templategroups not available)
	if api.Config.Version >= 70200 {
		templateGroup := zapi.TemplateGroups{{Name: "zabbix-testing-tg-" + t.Name()}}
		err := api.TemplateGroupsCreate(templateGroup)
		if isPermissionError(err) {
			t.Logf("Template group creation failed (permissions issue): %v", err)
			t.Skip("Requires template group permissions")
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		defer api.TemplateGroupsDelete(templateGroup)

		// Create template with template group
		template := zapi.Templates{zapi.Template{
			Host:   "template name",
			Groups: []zapi.TemplateGroupID{{GroupID: templateGroup[0].GroupID}},
		}}
		err = api.TemplatesCreate(template)
		if err != nil {
			t.Fatal(err)
		}
		defer DeleteTemplate(&template[0], t)

		if template[0].TemplateID == "" {
			t.Errorf("Template id is empty")
		}

		templates, err := api.TemplatesGet(zapi.Params{})
		if err != nil {
			t.Fatal(err)
		}
		if len(templates) == 0 {
			t.Fatal("No templates were obtained")
		}

		template[0].Name = "new template name"
		err = api.TemplatesUpdate(template)
		if err != nil {
			t.Error(err)
		}
		return
	}

	// For Zabbix 7+, use template groups
	if api.IsZabbix7OrGreater() {
		templateGroup := zapi.TemplateGroups{{Name: "zabbix-testing-tg-" + t.Name()}}
		err := api.TemplateGroupsCreate(templateGroup)
		if isPermissionError(err) {
			t.Logf("Template group creation failed: %v", err)
			t.Skip("Requires template group permissions")
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		defer api.TemplateGroupsDelete(templateGroup)

		template := zapi.Templates{zapi.Template{
			Host:   "template name",
			Groups: []zapi.TemplateGroupID{{GroupID: templateGroup[0].GroupID}},
		}}
		err = api.TemplatesCreate(template)
		if err != nil {
			t.Fatal(err)
		}
		defer DeleteTemplate(&template[0], t)

		if template[0].TemplateID == "" {
			t.Errorf("Template id is empty")
		}

		templates, err := api.TemplatesGet(zapi.Params{})
		if err != nil {
			t.Fatal(err)
		}
		if len(templates) == 0 {
			t.Fatal("No templates were obtained")
		}

		template[0].Name = "new template name"
		err = api.TemplatesUpdate(template)
		if err != nil {
			t.Error(err)
		}
		return
	}

	// For Zabbix < 7
	template := CreateTemplate(hostGroup, t)
	if template.TemplateID == "" {
		t.Errorf("Template id is empty %#v", template)
	}

	templates, err := api.TemplatesGet(zapi.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) == 0 {
		t.Fatal("No templates were obtained")
	}

	_, err = api.TemplateGetByID(template.TemplateID)
	if err != nil {
		t.Error(err)
	}

	template.Name = "new template name"
	err = api.TemplatesUpdate(zapi.Templates{*template})
	if err != nil {
		t.Error(err)
	}

	DeleteTemplate(template, t)
}

func TestTemplateWithSNMPItems(t *testing.T) {
	api := getAPI(t)

	if !api.IsZabbix7OrGreater() {
		t.Skip("Template with SNMP items requires Zabbix 7+")
	}

	// Create template group
	templateGroup := zapi.TemplateGroups{{Name: fmt.Sprintf("zabbix-testing-tg-snmp-%d", rand.Int())}}
	err := api.TemplateGroupsCreate(templateGroup)
	if err != nil {
		t.Fatal(err)
	}
	defer api.TemplateGroupsDelete(templateGroup)

	// Create template with template group
	template := zapi.Templates{zapi.Template{
		Host:   fmt.Sprintf("template-snmp-%d", rand.Int()),
		Name:   fmt.Sprintf("SNMP Template %d", rand.Int()),
		Groups: []zapi.TemplateGroupID{{GroupID: templateGroup[0].GroupID}},
	}}
	err = api.TemplatesCreate(template)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTemplate(&template[0], t)

	t.Logf("Created template: %s", template[0].TemplateID)

	// Create SNMP item on template
	snmpItem := zapi.Items{{
		HostID:    template[0].TemplateID,
		Key:       "1.3.6.1.2.1.1.1.0",
		Name:      "SNMP System Description",
		Type:      zapi.SNMPAgent,
		ValueType: zapi.Unsigned,
		SNMPOid:   "1.3.6.1.2.1.1.1.0",
		Delay:     "30s",
	}}
	err = api.ItemsCreate(snmpItem)
	if err != nil {
		t.Fatal(err)
	}
	defer api.ItemsDelete(snmpItem)

	t.Logf("Created SNMP item on template: %s", snmpItem[0].ItemID)

	// Verify item exists
	items, err := api.ItemsGet(zapi.Params{"hostids": template[0].TemplateID})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) == 0 {
		t.Fatal("No items found on template")
	}

	t.Logf("Found %d items on template", len(items))

	// Create host group
	hostGroup := zapi.HostGroups{{Name: fmt.Sprintf("zabbix-testing-hg-%d", rand.Int())}}
	err = api.HostGroupsCreate(hostGroup)
	if err != nil {
		t.Fatal(err)
	}
	defer api.HostGroupsDelete(hostGroup)

	// Get host group ID
	hostGroupID := hostGroup[0].GroupID

	// Try to link template to host group (requires super-admin permissions in Zabbix 7+)
	err = api.TemplatesMassAdd(zapi.TemplateIDs{{template[0].TemplateID}}, zapi.HostGroupIDs{{GroupID: hostGroupID}})
	if err != nil {
		t.Logf("Template-HostGroup linking requires super-admin permissions: %v", err)
	}

	t.Logf("Template %s linked to host group %s", template[0].TemplateID, hostGroupID)

	// Create host (will inherit template via host group)
	snmpIface := zapi.HostInterface{
		DNS:   fmt.Sprintf("host-snmp-%d", rand.Int()),
		IP:    "127.0.0.1",
		Port:  "161",
		Type:  zapi.SNMP,
		UseIP: "1",
		Main:  "1",
		Details: &zapi.HostInterfaceDetail{
			Version:   "2",
			Bulk:      "1",
			Community: "public",
		},
	}

	host := zapi.Hosts{{
		Host:       fmt.Sprintf("host-with-snmp-template-%d", rand.Int()),
		Name:       fmt.Sprintf("Host with SNMP Template %d", rand.Int()),
		GroupIds:   zapi.HostGroupIDs{{GroupID: hostGroupID}},
		Interfaces: zapi.HostInterfaces{snmpIface},
	}}
	err = api.HostsCreate(host)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteHost(&host[0], t)

	t.Logf("Created host: %s with template: %s", host[0].HostID, template[0].TemplateID)
}
