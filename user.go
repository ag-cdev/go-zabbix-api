package zabbix

type UserGroupID struct {
	UserGroupID string `json:"usrgrpid,omitempty"`
}

type UserGroupPermission struct {
	ID         string `json:"id,omitempty"`
	Permission int    `json:"permission,omitempty"`
}

type UserGroup struct {
	UserGroupID string                `json:"usrgrpid,omitempty"`
	Name        string                `json:"name,omitempty"`
	DebugMode   int                   `json:"debug_mode,string,omitempty"`
	GUIAccess   int                   `json:"gui_access,string,omitempty"`
	Status      int                   `json:"status,string,omitempty"`
	Permissions []UserGroupPermission `json:"host_permission,omitempty"`
}

type UserGroups []UserGroup

type User struct {
	UserID   string        `json:"userid,omitempty"`
	Username string        `json:"username"`
	Password string        `json:"passwd,omitempty"`
	RoleID   string        `json:"roleid,omitempty"`
	Name     string        `json:"name,omitempty"`
	Surname  string        `json:"surname,omitempty"`
	Groups   []UserGroupID `json:"usrgrps,omitempty"`
}

type Users []User

func (api *API) UsersGet(params Params) (res Users, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("user.get", params, &res)
	return
}

func (api *API) UsersCreate(users Users) (err error) {
	response, err := api.CallWithError("user.create", users)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	userids := result["userids"].([]interface{})
	for i, id := range userids {
		users[i].UserID = id.(string)
	}
	return
}

func (api *API) UsersUpdate(users Users) (err error) {
	_, err = api.CallWithError("user.update", users)
	return
}

func (api *API) UsersDeleteByIds(ids []string) (err error) {
	_, err = api.CallWithError("user.delete", ids)
	return
}

func (api *API) UserGroupsGet(params Params) (res UserGroups, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("usergroup.get", params, &res)
	return
}

func (api *API) UserGroupsCreate(usergroups UserGroups) (err error) {
	response, err := api.CallWithError("usergroup.create", usergroups)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	usrgrpids := result["usrgrpids"].([]interface{})
	for i, id := range usrgrpids {
		usergroups[i].UserGroupID = id.(string)
	}
	return
}

func (api *API) UserGroupsUpdate(usergroups UserGroups) (err error) {
	_, err = api.CallWithError("usergroup.update", usergroups)
	return
}

func (api *API) UserGroupsDeleteByIds(ids []string) (err error) {
	_, err = api.CallWithError("usergroup.delete", ids)
	return
}
