package middleware

import (
	"strings"
	"testing"

	"github.com/xformation/sdp/pkg/bus"
	m "github.com/xformation/sdp/pkg/models"
	"github.com/xformation/sdp/pkg/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMiddlewareDashboardRedirect(t *testing.T) {
	Convey("Given the dashboard redirect middleware", t, func() {
		bus.ClearBusHandlers()
		redirectFromLegacyDashboardUrl := RedirectFromLegacyDashboardUrl()
		redirectFromLegacyDashboardSoloUrl := RedirectFromLegacyDashboardSoloUrl()

		fakeDash := m.NewDashboard("Child dash")
		fakeDash.Id = 1
		fakeDash.FolderId = 1
		fakeDash.HasAcl = false
		fakeDash.Uid = util.GenerateShortUid()

		bus.AddHandler("test", func(query *m.GetDashboardQuery) error {
			query.Result = fakeDash
			return nil
		})

		middlewareScenario("GET dashboard by legacy url", func(sc *scenarioContext) {
			sc.m.Get("/dashboard/db/:slug", redirectFromLegacyDashboardUrl, sc.defaultHandler)

			sc.fakeReqWithParams("GET", "/dashboard/db/dash?orgId=1&panelId=2", map[string]string{}).exec()

			Convey("Should redirect to new dashboard url with a 301 Moved Permanently", func() {
				So(sc.resp.Code, ShouldEqual, 301)
				redirectUrl, _ := sc.resp.Result().Location()
				So(redirectUrl.Path, ShouldEqual, m.GetDashboardUrl(fakeDash.Uid, fakeDash.Slug))
				So(len(redirectUrl.Query()), ShouldEqual, 2)
			})
		})

		middlewareScenario("GET dashboard solo by legacy url", func(sc *scenarioContext) {
			sc.m.Get("/dashboard-solo/db/:slug", redirectFromLegacyDashboardSoloUrl, sc.defaultHandler)

			sc.fakeReqWithParams("GET", "/dashboard-solo/db/dash?orgId=1&panelId=2", map[string]string{}).exec()

			Convey("Should redirect to new dashboard url with a 301 Moved Permanently", func() {
				So(sc.resp.Code, ShouldEqual, 301)
				redirectUrl, _ := sc.resp.Result().Location()
				expectedUrl := m.GetDashboardUrl(fakeDash.Uid, fakeDash.Slug)
				expectedUrl = strings.Replace(expectedUrl, "/d/", "/d-solo/", 1)
				So(redirectUrl.Path, ShouldEqual, expectedUrl)
				So(len(redirectUrl.Query()), ShouldEqual, 2)
			})
		})
	})
}
