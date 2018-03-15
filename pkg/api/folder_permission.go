package api

import (
	"time"

	"github.com/xformation/sdp/pkg/api/dtos"
	"github.com/xformation/sdp/pkg/bus"
	"github.com/xformation/sdp/pkg/middleware"
	m "github.com/xformation/sdp/pkg/models"
	"github.com/xformation/sdp/pkg/services/dashboards"
	"github.com/xformation/sdp/pkg/services/guardian"
)

func GetFolderPermissionList(c *middleware.Context) Response {
	s := dashboards.NewFolderService(c.OrgId, c.SignedInUser)
	folder, err := s.GetFolderByUid(c.Params(":uid"))

	if err != nil {
		return toFolderError(err)
	}

	g := guardian.New(folder.Id, c.OrgId, c.SignedInUser)

	if canAdmin, err := g.CanAdmin(); err != nil || !canAdmin {
		return toFolderError(m.ErrFolderAccessDenied)
	}

	acl, err := g.GetAcl()
	if err != nil {
		return ApiError(500, "Failed to get folder permissions", err)
	}

	for _, perm := range acl {
		perm.FolderId = folder.Id
		perm.DashboardId = 0

		if perm.Slug != "" {
			perm.Url = m.GetDashboardFolderUrl(perm.IsFolder, perm.Uid, perm.Slug)
		}
	}

	return Json(200, acl)
}

func UpdateFolderPermissions(c *middleware.Context, apiCmd dtos.UpdateDashboardAclCommand) Response {
	s := dashboards.NewFolderService(c.OrgId, c.SignedInUser)
	folder, err := s.GetFolderByUid(c.Params(":uid"))

	if err != nil {
		return toFolderError(err)
	}

	g := guardian.New(folder.Id, c.OrgId, c.SignedInUser)
	canAdmin, err := g.CanAdmin()
	if err != nil {
		return toFolderError(err)
	}

	if !canAdmin {
		return toFolderError(m.ErrFolderAccessDenied)
	}

	cmd := m.UpdateDashboardAclCommand{}
	cmd.DashboardId = folder.Id

	for _, item := range apiCmd.Items {
		cmd.Items = append(cmd.Items, &m.DashboardAcl{
			OrgId:       c.OrgId,
			DashboardId: folder.Id,
			UserId:      item.UserId,
			TeamId:      item.TeamId,
			Role:        item.Role,
			Permission:  item.Permission,
			Created:     time.Now(),
			Updated:     time.Now(),
		})
	}

	if okToUpdate, err := g.CheckPermissionBeforeUpdate(m.PERMISSION_ADMIN, cmd.Items); err != nil || !okToUpdate {
		if err != nil {
			if err == guardian.ErrGuardianPermissionExists ||
				err == guardian.ErrGuardianOverride {
				return ApiError(400, err.Error(), err)
			}

			return ApiError(500, "Error while checking folder permissions", err)
		}

		return ApiError(403, "Cannot remove own admin permission for a folder", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrDashboardAclInfoMissing {
			err = m.ErrFolderAclInfoMissing
		}
		if err == m.ErrDashboardPermissionDashboardEmpty {
			err = m.ErrFolderPermissionFolderEmpty
		}

		if err == m.ErrFolderAclInfoMissing || err == m.ErrFolderPermissionFolderEmpty {
			return ApiError(409, err.Error(), err)
		}

		return ApiError(500, "Failed to create permission", err)
	}

	return ApiSuccess("Folder permissions updated")
}
