package repo

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"

	"github.com/getcouragenow/sys-share/sys-account/service/go/pkg"
	sharedAuth "github.com/getcouragenow/sys-share/sys-account/service/go/pkg/shared"
	"github.com/getcouragenow/sys/sys-core/service/go/pkg/coredb"
)

func (ad *SysAccountRepo) getAccountAndRole(id, email string) (*pkg.Account, error) {
	queryParams := map[string]interface{}{}
	if id != "" {
		queryParams["id"] = id
	}
	if email != "" {
		queryParams["email"] = email
	}
	acc, err := ad.store.GetAccount(&coredb.QueryParams{Params: queryParams})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot find user account: %v", sharedAuth.Error{Reason: sharedAuth.ErrAccountNotFound})
	}
	daoRoles, err := ad.store.ListRole(&coredb.QueryParams{Params: map[string]interface{}{
		"account_id": acc.ID,
	}})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot find user role: %v", sharedAuth.Error{Reason: sharedAuth.ErrAccountNotFound})
	}
	var pkgRoles []*pkg.UserRoles
	for _, daoRole := range daoRoles {
		pkgRole, err := daoRole.ToPkgRole()
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "cannot find user role: %v", sharedAuth.Error{Reason: sharedAuth.ErrAccountNotFound})
		}
		pkgRoles = append(pkgRoles, pkgRole)
	}
	return acc.ToPkgAccount(pkgRoles)
}

func (ad *SysAccountRepo) listAccountsAndRoles(filter *coredb.QueryParams, orderBy string, limit, cursor int64) ([]*pkg.Account, *int64, error) {
	listAccounts, next, err := ad.store.ListAccount(filter, orderBy, limit, cursor)
	if err != nil {
		return nil, nil, err
	}
	var accounts []*pkg.Account

	for _, acc := range listAccounts {
		daoRoles, err := ad.store.ListRole(&coredb.QueryParams{Params: map[string]interface{}{
			"account_id": acc.ID,
		}})
		if err != nil {
			return nil, nil, status.Errorf(codes.NotFound, "cannot find user roles: %v", sharedAuth.Error{Reason: sharedAuth.ErrAccountNotFound, Err: err})
		}
		var pkgRoles []*pkg.UserRoles
		for _, daoRole := range daoRoles {
			pkgRole, err := daoRole.ToPkgRole()
			if err != nil {
				return nil, nil, status.Errorf(codes.NotFound, "cannot find user roles: %v", sharedAuth.Error{Reason: sharedAuth.ErrAccountNotFound, Err: err})
			}
			pkgRoles = append(pkgRoles, pkgRole)
		}
		account, err := acc.ToPkgAccount(pkgRoles)
		if err != nil {
			return nil, nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, &next, nil
}

func (ad *SysAccountRepo) getCursor(currentCursor string) (int64, error) {
	if currentCursor != "" {
		return strconv.ParseInt(currentCursor, 10, 64)
	} else {
		return 0, nil
	}
}

