package api

import (
	"time"

	"github.com/xformation/sdp/pkg/api/dtos"
	"github.com/xformation/sdp/pkg/bus"
	"github.com/xformation/sdp/pkg/metrics"
	"github.com/xformation/sdp/pkg/middleware"
	m "github.com/xformation/sdp/pkg/models"
	"github.com/xformation/sdp/pkg/services/guardian"
	"github.com/xformation/sdp/pkg/setting"
	"github.com/xformation/sdp/pkg/util"
)

func GetSharingOptions(c *middleware.Context) {
	c.JSON(200, util.DynMap{
		"externalSnapshotURL":  setting.ExternalSnapshotUrl,
		"externalSnapshotName": setting.ExternalSnapshotName,
		"externalEnabled":      setting.ExternalEnabled,
	})
}

func CreateDashboardSnapshot(c *middleware.Context, cmd m.CreateDashboardSnapshotCommand) {
	if cmd.Name == "" {
		cmd.Name = "Unnamed snapshot"
	}

	if cmd.External {
		// external snapshot ref requires key and delete key
		if cmd.Key == "" || cmd.DeleteKey == "" {
			c.JsonApiErr(400, "Missing key and delete key for external snapshot", nil)
			return
		}

		cmd.OrgId = -1
		cmd.UserId = -1
		metrics.M_Api_Dashboard_Snapshot_External.Inc()
	} else {
		cmd.Key = util.GetRandomString(32)
		cmd.DeleteKey = util.GetRandomString(32)
		cmd.OrgId = c.OrgId
		cmd.UserId = c.UserId
		metrics.M_Api_Dashboard_Snapshot_Create.Inc()
	}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Failed to create snaphost", err)
		return
	}

	c.JSON(200, util.DynMap{
		"key":       cmd.Key,
		"deleteKey": cmd.DeleteKey,
		"url":       setting.ToAbsUrl("dashboard/snapshot/" + cmd.Key),
		"deleteUrl": setting.ToAbsUrl("api/snapshots-delete/" + cmd.DeleteKey),
	})
}

// GET /api/snapshots/:key
func GetDashboardSnapshot(c *middleware.Context) {
	key := c.Params(":key")
	query := &m.GetDashboardSnapshotQuery{Key: key}

	err := bus.Dispatch(query)
	if err != nil {
		c.JsonApiErr(500, "Failed to get dashboard snapshot", err)
		return
	}

	snapshot := query.Result

	// expired snapshots should also be removed from db
	if snapshot.Expires.Before(time.Now()) {
		c.JsonApiErr(404, "Dashboard snapshot not found", err)
		return
	}

	dto := dtos.DashboardFullWithMeta{
		Dashboard: snapshot.Dashboard,
		Meta: dtos.DashboardMeta{
			Type:       m.DashTypeSnapshot,
			IsSnapshot: true,
			Created:    snapshot.Created,
			Expires:    snapshot.Expires,
		},
	}

	metrics.M_Api_Dashboard_Snapshot_Get.Inc()

	c.Resp.Header().Set("Cache-Control", "public, max-age=3600")
	c.JSON(200, dto)
}

// GET /api/snapshots-delete/:key
func DeleteDashboardSnapshot(c *middleware.Context) Response {
	key := c.Params(":key")

	query := &m.GetDashboardSnapshotQuery{DeleteKey: key}

	err := bus.Dispatch(query)
	if err != nil {
		return ApiError(500, "Failed to get dashboard snapshot", err)
	}

	if query.Result == nil {
		return ApiError(404, "Failed to get dashboard snapshot", nil)
	}
	dashboard := query.Result.Dashboard
	dashboardId := dashboard.Get("id").MustInt64()

	guardian := guardian.New(dashboardId, c.OrgId, c.SignedInUser)
	canEdit, err := guardian.CanEdit()
	if err != nil {
		return ApiError(500, "Error while checking permissions for snapshot", err)
	}

	if !canEdit && query.Result.UserId != c.SignedInUser.UserId {
		return ApiError(403, "Access denied to this snapshot", nil)
	}

	cmd := &m.DeleteDashboardSnapshotCommand{DeleteKey: key}

	if err := bus.Dispatch(cmd); err != nil {
		return ApiError(500, "Failed to delete dashboard snapshot", err)
	}

	return Json(200, util.DynMap{"message": "Snapshot deleted. It might take an hour before it's cleared from a CDN cache."})
}

// GET /api/dashboard/snapshots
func SearchDashboardSnapshots(c *middleware.Context) Response {
	query := c.Query("query")
	limit := c.QueryInt("limit")

	if limit == 0 {
		limit = 1000
	}

	searchQuery := m.GetDashboardSnapshotsQuery{
		Name:         query,
		Limit:        limit,
		OrgId:        c.OrgId,
		SignedInUser: c.SignedInUser,
	}

	err := bus.Dispatch(&searchQuery)
	if err != nil {
		return ApiError(500, "Search failed", err)
	}

	dtos := make([]*m.DashboardSnapshotDTO, len(searchQuery.Result))
	for i, snapshot := range searchQuery.Result {
		dtos[i] = &m.DashboardSnapshotDTO{
			Id:          snapshot.Id,
			Name:        snapshot.Name,
			Key:         snapshot.Key,
			DeleteKey:   snapshot.DeleteKey,
			OrgId:       snapshot.OrgId,
			UserId:      snapshot.UserId,
			External:    snapshot.External,
			ExternalUrl: snapshot.ExternalUrl,
			Expires:     snapshot.Expires,
			Created:     snapshot.Created,
			Updated:     snapshot.Updated,
		}
	}

	return Json(200, dtos)
}
