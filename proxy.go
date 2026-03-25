package zabbix

// Proxy represent Zabbix proxy object
// https://www.zabbix.com/documentation/3.2/manual/api/reference/proxy/object
type Proxy struct {
	ProxyID        string `json:"proxyid,omitempty"`
	Host           string `json:"host"`
	Name           string `json:"name,omitempty"`
	OperatingMode  int    `json:"operating_mode,string,omitempty"`
	Description    string `json:"description,omitempty"`
	TLSConnect     int    `json:"tls_connect,string,omitempty"`
	TLSAccept      int    `json:"tls_accept,string,omitempty"`
	TLSIssuer      string `json:"tls_issuer,omitempty"`
	TLSSubject     string `json:"tls_subject,omitempty"`
	TLSPSKIdentity string `json:"tls_psk_identity,omitempty"`
	TLSPSK         string `json:"tls_psk,omitempty"`
	ProxyAddress   string `json:"proxy_address,omitempty"`
}

// Proxies is an array of Proxy
type Proxies []Proxy

// ProxiesGet Wrapper for proxy.get
// https://www.zabbix.com/documentation/3.2/manual/api/reference/proxy/get
func (api *API) ProxiesGet(params Params) (res Proxies, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("proxy.get", params, &res)
	return
}

// ProxiesCreate Wrapper for proxy.create
func (api *API) ProxiesCreate(proxies Proxies) (err error) {
	response, err := api.CallWithError("proxy.create", proxies)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	proxyids := result["proxyids"].([]interface{})
	for i, id := range proxyids {
		proxies[i].ProxyID = id.(string)
	}
	return
}

// ProxiesUpdate Wrapper for proxy.update
func (api *API) ProxiesUpdate(proxies Proxies) (err error) {
	_, err = api.CallWithError("proxy.update", proxies)
	return
}

// ProxiesDelete Wrapper for proxy.delete
func (api *API) ProxiesDelete(proxies Proxies) (err error) {
	ids := make([]string, len(proxies))
	for i, proxy := range proxies {
		ids[i] = proxy.ProxyID
	}
	_, err = api.CallWithError("proxy.delete", ids)
	return
}

// ProxiesDeleteByIds Wrapper for proxy.delete
func (api *API) ProxiesDeleteByIds(ids []string) (err error) {
	_, err = api.CallWithError("proxy.delete", ids)
	return
}
