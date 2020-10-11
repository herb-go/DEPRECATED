package member

import (
	"net/http"

	"github.com/herb-go/herbmodules/protecter"
)

//Authorizer comonets to Authorize  http request.
//Should be created by Service.Authorize
type Authorizer struct {
	Service      *Service
	PolicyLoader protecter.PolicyLoader
}

//Authorize Authorize http request.
func (a *Authorizer) Authorize(r *http.Request) (bool, error) {
	uid, err := a.Service.IdentifyRequest(r)
	if err != nil {
		return false, err
	}
	if uid == "" {
		return false, nil
	}
	var members = a.Service.GetMembersFromRequest(r)
	if a.Service.StatusProvider != nil {
		_, err = members.LoadStatus(uid)
		if err != nil {
			return false, err
		}
		if !IsAvaliable(members.StatusStore.Get(uid)) {
			return false, nil
		}
	}

	if a.Service.RoleProvider == nil {
		return true, nil
	}
	if a.PolicyLoader == nil {
		return true, nil
	}
	rolesstore, err := members.LoadRoles(uid)
	if err != nil {
		return false, err
	}
	if rolesstore == nil {
		return false, err
	}
	roles := rolesstore.Get(uid)
	if roles == nil {
		return false, nil
	}
	return protecter.Authorize(r, protecter.RoleRolesLoader(roles), a.PolicyLoader)
}
