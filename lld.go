package zabbix

import "encoding/json"

type (
	LLDEvalType     string
	LLDOperatorType string
)

const (
	LLDAndOr    LLDEvalType     = "0"
	LLDAnd      LLDEvalType     = "1"
	LLDOr       LLDEvalType     = "2"
	LLDCustom   LLDEvalType     = "3"
	LLDMatch    LLDOperatorType = "8"
	LLDNotMatch LLDOperatorType = "9"
)

type LLDRuleFilterCondition struct {
	Macro     string          `json:"macro"`
	Value     string          `json:"value"`
	FormulaID string          `json:"formulaid,omitempty"`
	Operator  LLDOperatorType `json:"operator,omitempty"`
}

type LLDRuleFilterConditions []LLDRuleFilterCondition

type LLDRuleFilter struct {
	Conditions  LLDRuleFilterConditions `json:"conditions"`
	EvalType    LLDEvalType             `json:"evaltype"`
	EvalFormula string                  `json:"eval_formula,omitempty"`
	Formula     string                  `json:"formula"`
}

type LLDMacroPath struct {
	Macro string `json:"lld_macro"`
	Path  string `json:"path"`
}

type LLDMacroPaths []LLDMacroPath

// Item represent Zabbix lld object
// https://www.zabbix.com/documentation/3.2/manual/api/reference/item/object
type LLDRule struct {
	ItemID       string   `json:"itemid,omitempty"`
	Delay        string   `json:"delay,omitempty"`
	HostID       string   `json:"hostid,omitempty"`
	InterfaceID  string   `json:"interfaceid,omitempty"`
	Key          string   `json:"key_"`
	Name         string   `json:"name"`
	Type         ItemType `json:"type,string"`
	AuthType     string   `json:"authtype,omitempty"`
	DelayFlex    string   `json:"delay_flex,omitempty"`
	Description  string   `json:"description,omitempty"`
	Error        string   `json:"error,omitempty"`
	IpmiSensor   string   `json:"ipmi_sensor,omitempty"`
	LifeTime     string   `json:"lifetime,omitempty"`
	Params       string   `json:"params,omitempty"`
	PrivateKey   string   `json:"privatekey,omitempty"`
	PublicKey    string   `json:"publickey,omitempty"`
	Status       string   `json:"status,omitempty"`
	TrapperHosts string   `json:"trapper_hosts,omitempty"`
	MasterItemID string   `json:"master_itemid,omitempty"`

	// ssh / telnet
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Port     string `json:"port,omitempty"`

	// HTTP Agent Fields
	Url             string          `json:"url,omitempty"`
	RequestMethod   string          `json:"request_method,omitempty"`
	AllowTraps      string          `json:"allow_traps,omitempty"`
	PostType        string          `json:"post_type,omitempty"`
	RetrieveMode    string          `json:"retrieve_mode,omitempty"`
	Posts           string          `json:"posts,omitempty"`
	StatusCodes     string          `json:"status_codes,omitempty"`
	Timeout         string          `json:"timeout,omitempty"`
	VerifyHost      string          `json:"verify_host,omitempty"`
	VerifyPeer      string          `json:"verify_peer,omitempty"`
	Headers         HttpHeaders     `json:"-"`
	HeadersArray    KeyValueArray   `json:"-"`
	RawHeaders      json.RawMessage `json:"headers,omitempty"`
	QueryFields     KeyValueArray   `json:"query_fields,omitempty"`
	Proxy           string          `json:"http_proxy,omitempty"`
	FollowRedirects string          `json:"follow_redirects,omitempty"`

	zbxVersion int `json:"-"`

	// SNMP Fields
	SNMPOid              string `json:"snmp_oid,omitempty"`
	SNMPCommunity        string `json:"snmp_community,omitempty"`
	SNMPv3AuthPassphrase string `json:"snmpv3_authpassphrase,omitempty"`
	SNMPv3AuthProtocol   string `json:"snmpv3_authprotocol,omitempty"`
	SNMPv3ContextName    string `json:"snmpv3_contextname,omitempty"`
	SNMPv3PrivPasshrase  string `json:"snmpv3_privpassphrase,omitempty"`
	SNMPv3PrivProtocol   string `json:"snmpv3_privprotocol,omitempty"`
	SNMPv3SecurityLevel  string `json:"snmpv3_securitylevel,omitempty"`
	SNMPv3SecurityName   string `json:"snmpv3_securityname,omitempty"`

	Preprocessors Preprocessors `json:"preprocessing,omitempty"`
	Filter        LLDRuleFilter `json:"filter"`
	MacroPaths    LLDMacroPaths `json:"lld_macro_paths,omitempty"`
}

func (i LLDRule) MarshalJSON() ([]byte, error) {
	headers := i.HeadersArray
	if i.Headers != nil && len(i.Headers) > 0 {
		for k, v := range i.Headers {
			headers = append(headers, KeyValue{Name: k, Value: v})
		}
	}

	// For Zabbix 7+, omit formulaid from filter conditions and hostid for LLD rules
	if i.zbxVersion >= 70000 {
		type AliasNoFormulaid LLDRule
		aux := struct {
			AliasNoFormulaid
			Headers KeyValueArray `json:"headers,omitempty"`
			Filter  struct {
				Conditions []struct {
					Macro    string          `json:"macro"`
					Value    string          `json:"value"`
					Operator LLDOperatorType `json:"operator,omitempty"`
				} `json:"conditions"`
				EvalType LLDEvalType `json:"evaltype"`
				Formula  string      `json:"formula"`
			} `json:"filter"`
		}{
			AliasNoFormulaid: AliasNoFormulaid{
				HostID:               i.HostID,
				ItemID:               i.ItemID,
				Delay:                i.Delay,
				InterfaceID:          i.InterfaceID,
				Key:                  i.Key,
				Name:                 i.Name,
				Type:                 i.Type,
				AuthType:             i.AuthType,
				DelayFlex:            i.DelayFlex,
				Description:          i.Description,
				Error:                i.Error,
				IpmiSensor:           i.IpmiSensor,
				LifeTime:             i.LifeTime,
				Params:               i.Params,
				PrivateKey:           i.PrivateKey,
				PublicKey:            i.PublicKey,
				Status:               i.Status,
				TrapperHosts:         i.TrapperHosts,
				MasterItemID:         i.MasterItemID,
				Username:             i.Username,
				Password:             i.Password,
				Port:                 i.Port,
				Url:                  i.Url,
				RequestMethod:        i.RequestMethod,
				AllowTraps:           i.AllowTraps,
				PostType:             i.PostType,
				RetrieveMode:         i.RetrieveMode,
				Posts:                i.Posts,
				StatusCodes:          i.StatusCodes,
				Timeout:              i.Timeout,
				VerifyHost:           i.VerifyHost,
				VerifyPeer:           i.VerifyPeer,
				QueryFields:          i.QueryFields,
				Proxy:                i.Proxy,
				FollowRedirects:      i.FollowRedirects,
				SNMPOid:              i.SNMPOid,
				SNMPCommunity:        i.SNMPCommunity,
				SNMPv3AuthPassphrase: i.SNMPv3AuthPassphrase,
				SNMPv3AuthProtocol:   i.SNMPv3AuthProtocol,
				SNMPv3ContextName:    i.SNMPv3ContextName,
				SNMPv3PrivPasshrase:  i.SNMPv3PrivPasshrase,
				SNMPv3PrivProtocol:   i.SNMPv3PrivProtocol,
				SNMPv3SecurityLevel:  i.SNMPv3SecurityLevel,
				SNMPv3SecurityName:   i.SNMPv3SecurityName,
				Preprocessors:        i.Preprocessors,
				MacroPaths:           i.MacroPaths,
			},
			Headers: headers,
		}

		// Copy filter conditions without formulaid
		aux.Filter.EvalType = i.Filter.EvalType
		aux.Filter.Formula = i.Filter.Formula
		aux.Filter.Conditions = []struct {
			Macro    string          `json:"macro"`
			Value    string          `json:"value"`
			Operator LLDOperatorType `json:"operator,omitempty"`
		}{}
		for _, c := range i.Filter.Conditions {
			aux.Filter.Conditions = append(aux.Filter.Conditions, struct {
				Macro    string          `json:"macro"`
				Value    string          `json:"value"`
				Operator LLDOperatorType `json:"operator,omitempty"`
			}{
				Macro:    c.Macro,
				Value:    c.Value,
				Operator: c.Operator,
			})
		}

		return json.Marshal(aux)
	}

	type Alias LLDRule
	aux := struct {
		Alias
		Headers KeyValueArray `json:"headers,omitempty"`
	}{
		Alias:   Alias(i),
		Headers: headers,
	}

	return json.Marshal(aux)
}

// Items is an array of Item
type LLDRules []LLDRule

func (api *API) lldsHeadersUnmarshal(item LLDRules) {
	for i := 0; i < len(item); i++ {
		h := item[i]

		item[i].Headers = HttpHeaders{}

		if len(h.RawHeaders) == 0 {
			continue
		}

		asStr := string(h.RawHeaders)
		if asStr == "[]" {
			continue
		}

		// Try to unmarshal as map first (Zabbix < 7 format)
		out := HttpHeaders{}
		err := json.Unmarshal(h.RawHeaders, &out)
		if err != nil {
			// Try array format (Zabbix 7+)
			var arr KeyValueArray
			if err2 := json.Unmarshal(h.RawHeaders, &arr); err2 == nil {
				out = HttpHeaders{}
				for _, kv := range arr {
					out[kv.Name] = kv.Value
				}
			} else {
				api.printf("got error during unmarshal %s", err)
				panic(err)
			}
		}
		item[i].Headers = out
	}
}

func prepLLDs(item LLDRules) {
	for i := 0; i < len(item); i++ {
		h := item[i]

		if h.Headers == nil {
			continue
		}
		asB, _ := json.Marshal(h.Headers)
		item[i].RawHeaders = json.RawMessage(asB)
	}
}

// ItemsGet Wrapper for item.get
// https://www.zabbix.com/documentation/3.2/manual/api/reference/item/get
func (api *API) LLDsGet(params Params) (res LLDRules, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("discoveryrule.get", params, &res)
	api.lldsHeadersUnmarshal(res)
	return
}

// ItemGetByID Gets item by Id only if there is exactly 1 matching host.
func (api *API) LLDGetByID(id string) (res *LLDRule, err error) {
	items, err := api.LLDsGet(Params{"itemids": id})
	if err != nil {
		return
	}

	if len(items) != 1 {
		e := ExpectedOneResult(len(items))
		err = &e
		return
	}
	res = &items[0]
	return
}

// ItemsCreate Wrapper for item.create
// https://www.zabbix.com/documentation/3.2/manual/api/reference/item/create
func (api *API) LLDsCreate(items LLDRules) (err error) {
	version := api.Config.Version
	for i := range items {
		items[i].zbxVersion = version
	}
	prepLLDs(items)
	response, err := api.CallWithError("discoveryrule.create", items)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	itemids := result["itemids"].([]interface{})
	for i, id := range itemids {
		items[i].ItemID = id.(string)
	}
	return
}

// ItemsUpdate Wrapper for item.update
// https://www.zabbix.com/documentation/3.2/manual/api/reference/item/update
func (api *API) LLDsUpdate(items LLDRules) (err error) {
	version := api.Config.Version
	for i := range items {
		items[i].zbxVersion = version
	}
	prepLLDs(items)
	_, err = api.CallWithError("discoveryrule.update", items)
	return
}

// ItemsDelete Wrapper for item.delete
// Cleans ItemId in all items elements if call succeed.
// https://www.zabbix.com/documentation/3.2/manual/api/reference/item/delete
func (api *API) LLDsDelete(items LLDRules) (err error) {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ItemID
	}

	err = api.LLDDeleteByIds(ids)
	if err == nil {
		for i := range items {
			items[i].ItemID = ""
		}
	}
	return
}

// ItemsDeleteByIds Wrapper for item.delete
// https://www.zabbix.com/documentation/3.2/manual/api/reference/item/delete
func (api *API) LLDDeleteByIds(ids []string) (err error) {
	deleteIds, err := api.LLDDeleteIDs(ids)
	if err != nil {
		return
	}
	l := len(deleteIds)
	if len(ids) != l {
		err = &ExpectedMore{len(ids), l}
	}
	return
}

// ItemsDeleteIDs Wrapper for item.delete
// Delete the item and return the id of the deleted item
func (api *API) LLDDeleteIDs(ids []string) (itemids []interface{}, err error) {
	response, err := api.CallWithError("discoveryrule.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	itemids1, ok := result["ruleids"].([]interface{})
	if !ok {
		itemids2 := result["ruleids"].(map[string]interface{})
		for _, id := range itemids2 {
			itemids = append(itemids, id)
		}
	} else {
		itemids = itemids1
	}
	return
}
