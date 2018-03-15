package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncoding(t *testing.T) {
	Convey("When generating base64 header", t, func() {
		result := GetBasicAuthHeader("sdp", "1234")

		So(result, ShouldEqual, "Basic Z3JhZmFuYToxMjM0")
	})

	Convey("When decoding basic auth header", t, func() {
		header := GetBasicAuthHeader("sdp", "1234")
		username, password, err := DecodeBasicAuthHeader(header)
		So(err, ShouldBeNil)

		So(username, ShouldEqual, "sdp")
		So(password, ShouldEqual, "1234")
	})

	Convey("When encoding password", t, func() {
		encodedPassword := EncodePassword("iamgod", "pepper")
		So(encodedPassword, ShouldEqual, "e59c568621e57756495a468f47c74e07c911b037084dd464bb2ed72410970dc849cabd71b48c394faf08a5405dae53741ce9")
	})
}
