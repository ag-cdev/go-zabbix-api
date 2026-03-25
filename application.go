package zabbix

import "fmt"

// Application represent Zabbix application object
// Note: Application API was removed in Zabbix 5.4+ and replaced with Tags
// https://www.zabbix.com/documentation/3.2/manual/api/reference/application/object
type Application struct {
	ApplicationID string `json:"applicationid,omitempty"`
	HostID        string `json:"hostid"`
	Name          string `json:"name"`
	TemplateID    string `json:"templateid,omitempty"`
}

// Applications is an array of Application
type Applications []Application

// ErrApplicationAPIDeprecated is returned when trying to use Application API on Zabbix 5.4+
type ErrApplicationAPIDeprecated struct{}

func (e *ErrApplicationAPIDeprecated) Error() string {
	return "Application API is deprecated in Zabbix 5.4+ and replaced with Tags"
}

// ApplicationsGet Wrapper for application.get
// https://www.zabbix.com/documentation/3.2/manual/api/reference/application/get
func (api *API) ApplicationsGet(params Params) (res Applications, err error) {
	if api.Config.Version >= 50400 {
		err = &ErrApplicationAPIDeprecated{}
		return
	}

	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("application.get", params, &res)
	if err != nil {
		return
	}

	return
}

// ApplicationGetByID Gets application by Id only if there is exactly 1 matching application.
func (api *API) ApplicationGetByID(id string) (res *Application, err error) {
	if api.Config.Version >= 50400 {
		err = &ErrApplicationAPIDeprecated{}
		return
	}

	apps, err := api.ApplicationsGet(Params{"applicationids": id})
	if err != nil {
		return
	}

	if len(apps) == 1 {
		res = &apps[0]
	} else {
		e := ExpectedOneResult(len(apps))
		err = &e
	}
	return
}

// ApplicationGetByHostIDAndName Gets application by host Id and name only if there is exactly 1 matching application.
func (api *API) ApplicationGetByHostIDAndName(hostID, name string) (res *Application, err error) {
	if api.Config.Version >= 50400 {
		err = &ErrApplicationAPIDeprecated{}
		return
	}

	apps, err := api.ApplicationsGet(Params{"hostids": hostID, "filter": map[string]string{"name": name}})
	if err != nil {
		return
	}

	if len(apps) == 1 {
		res = &apps[0]
	} else {
		e := ExpectedOneResult(len(apps))
		err = &e
	}
	return
}

// ApplicationsCreate Wrapper for application.create
// https://www.zabbix.com/documentation/3.2/manual/api/reference/application/create
func (api *API) ApplicationsCreate(apps Applications) (err error) {
	if api.Config.Version >= 50400 {
		err = &ErrApplicationAPIDeprecated{}
		return
	}

	response, err := api.CallWithError("application.create", apps)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	applicationids := result["applicationids"].([]interface{})
	for i, id := range applicationids {
		apps[i].ApplicationID = id.(string)
	}
	return
}

// ApplicationsDelete Wrapper for application.delete:
// Cleans ApplicationID in all apps elements if call succeed.
// https://www.zabbix.com/documentation/2.2/manual/appendix/api/application/delete
func (api *API) ApplicationsDelete(apps Applications) (err error) {
	if api.Config.Version >= 50400 {
		err = &ErrApplicationAPIDeprecated{}
		return
	}

	ids := make([]string, len(apps))
	for i, app := range apps {
		ids[i] = app.ApplicationID
	}

	err = api.ApplicationsDeleteByIds(ids)
	if err == nil {
		for i := range apps {
			apps[i].ApplicationID = ""
		}
	}
	return
}

// ApplicationsDeleteByIds Wrapper for application.delete
// https://www.zabbix.com/documentation/2.2/manual/appendix/api/application/delete
func (api *API) ApplicationsDeleteByIds(ids []string) (err error) {
	if len(ids) == 0 {
		return
	}

	if api.Config.Version >= 50400 {
		err = &ErrApplicationAPIDeprecated{}
		return
	}

	response, err := api.CallWithError("application.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	applicationids := result["applicationids"].([]interface{})
	if len(ids) != len(applicationids) {
		err = &ExpectedMore{len(ids), len(applicationids)}
	}
	return
}

// ApplicationGetByHostID is a compatibility function for Zabbix 5.4+
// It retrieves items filtered by host and simulates application-like behavior using tags
func (api *API) ApplicationGetByHostID(hostID string) (res Applications, err error) {
	if api.Config.Version >= 50400 {
		err = fmt.Errorf("use ItemsGet with tags filter instead (Application API removed in 5.4+)")
		return
	}
	return api.ApplicationsGet(Params{"hostids": hostID})
}
