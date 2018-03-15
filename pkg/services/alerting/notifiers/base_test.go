package notifiers

import (
	"context"
	"testing"

	m "github.com/xformation/sdp/pkg/models"
	"github.com/xformation/sdp/pkg/services/alerting"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBaseNotifier(t *testing.T) {
	Convey("Base notifier tests", t, func() {
		Convey("should notify", func() {
			Convey("pending -> ok", func() {
				context := alerting.NewEvalContext(context.TODO(), &alerting.Rule{
					State: m.AlertStatePending,
				})
				context.Rule.State = m.AlertStateOK
				So(defaultShouldNotify(context), ShouldBeFalse)
			})

			Convey("ok -> alerting", func() {
				context := alerting.NewEvalContext(context.TODO(), &alerting.Rule{
					State: m.AlertStateOK,
				})
				context.Rule.State = m.AlertStateAlerting
				So(defaultShouldNotify(context), ShouldBeTrue)
			})
		})
	})
}
