package notifiers

import (
	"github.com/xformation/sdp/pkg/components/simplejson"
	m "github.com/xformation/sdp/pkg/models"
	"github.com/xformation/sdp/pkg/services/alerting"
)

type NotifierBase struct {
	Name        string
	Type        string
	Id          int64
	IsDeault    bool
	UploadImage bool
}

func NewNotifierBase(id int64, isDefault bool, name, notifierType string, model *simplejson.Json) NotifierBase {
	uploadImage := model.Get("uploadImage").MustBool(false)

	return NotifierBase{
		Id:          id,
		Name:        name,
		IsDeault:    isDefault,
		Type:        notifierType,
		UploadImage: uploadImage,
	}
}

func defaultShouldNotify(context *alerting.EvalContext) bool {
	if context.PrevAlertState == context.Rule.State {
		return false
	}
	if (context.PrevAlertState == m.AlertStatePending) && (context.Rule.State == m.AlertStateOK) {
		return false
	}
	return true
}

func (n *NotifierBase) GetType() string {
	return n.Type
}

func (n *NotifierBase) NeedsImage() bool {
	return n.UploadImage
}

func (n *NotifierBase) GetNotifierId() int64 {
	return n.Id
}

func (n *NotifierBase) GetIsDefault() bool {
	return n.IsDeault
}
