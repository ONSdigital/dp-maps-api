package api_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ONSdigital/dp-maps-api/api"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	from          = "from"
	to            = "to"
	targetAPIURL  = "https://somewhere/"
	testApiKey    = "somekeyvalue"
	samplePayload = "fromsomething"
)

func TestCreateOSMapsProxy(t *testing.T) {
	Convey("Given an OS Max Proxy", t, func() {
		targetURL, err := url.Parse(targetAPIURL)
		So(err, ShouldBeNil)

		var testConfig = api.Config{
			OrdnanceSurveyAPIURL: targetURL,
			OrdnanceSurveyAPIKey: testApiKey,
			CacheMaxAge:          2 * time.Minute,
		}

		Convey("without a response modifier", func() {
			rp := api.CreateOSMapsProxy(targetURL, testConfig, nil)
			So(rp, ShouldNotBeNil)
			So(rp.Director, ShouldNotBeNil)
			So(rp.ModifyResponse, ShouldNotBeNil)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			Convey("then a valid request will be transformed correctly", func() {
				rp.Director(req)
				So(req, ShouldNotBeNil)
				So(req.Host, ShouldEqual, targetURL.Host)
			})

			Convey("then the modified response should be transformed", func() {
				res := createResponse(req, io.NopCloser(bytes.NewReader([]byte(samplePayload))))
				rp.ModifyResponse(&res)
				So(res.Header, ShouldNotBeNil)
				So(res.Header.Get("Cache-Control"), ShouldEqual, "max-age=120")
				So(res.Body, ShouldNotBeNil)
				body, err := io.ReadAll(res.Body)
				So(err, ShouldBeNil)
				So(string(body), ShouldResemble, samplePayload)
			})
		})

		Convey("with a successful response modifier", func() {
			modifiedResponse := false
			responseModifier := func(r *http.Response) error {
				modifiedResponse = true
				return nil
			}

			rp := api.CreateOSMapsProxy(targetURL, testConfig, responseModifier)
			So(rp, ShouldNotBeNil)
			So(rp.Director, ShouldNotBeNil)
			So(rp.ModifyResponse, ShouldNotBeNil)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			Convey("then a the request will ????", func() {
				rp.Director(req)
				So(req, ShouldNotBeNil)
				So(req.Header, ShouldNotBeNil)
				So(req.Header.Get("Accept-Encoding"), ShouldResemble, "")
			})

			Convey("then the modified response should be transformed", func() {
				res := createResponse(req, io.NopCloser(bytes.NewReader([]byte(samplePayload))))
				err = rp.ModifyResponse(&res)
				So(err, ShouldBeNil)
				So(res.Header, ShouldNotBeNil)
				So(res.Header.Get("Cache-Control"), ShouldEqual, "max-age=120")
				So(res.Body, ShouldNotBeNil)
				body, err := io.ReadAll(res.Body)
				So(err, ShouldBeNil)
				So(string(body), ShouldResemble, samplePayload)
				So(modifiedResponse, ShouldBeTrue)
			})
		})

		Convey("with an errorring response modifier", func() {
			fakeErr := errors.New("someerror")
			responseModifier := func(r *http.Response) error {
				return fakeErr
			}

			rp := api.CreateOSMapsProxy(targetURL, testConfig, responseModifier)
			So(rp, ShouldNotBeNil)
			So(rp.Director, ShouldNotBeNil)
			So(rp.ModifyResponse, ShouldNotBeNil)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			Convey("then the modified response should be transformed", func() {
				res := createResponse(req, io.NopCloser(bytes.NewReader([]byte(samplePayload))))
				err := rp.ModifyResponse(&res)
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, fakeErr)
			})
		})
	})
}

func TestStringReplaceResponseModifier(t *testing.T) {
	Convey("Given a StringReplaceResponseModifier", t, func() {
		srrm := api.StringReplaceResponseModifier(from, to)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		Convey("when we have a response containing the from text", func() {
			body := io.NopCloser(bytes.NewReader([]byte(samplePayload)))
			r := createResponse(req, body)
			err := srrm(&r)
			So(err, ShouldBeNil)

			newBody, err := ioutil.ReadAll(r.Body)
			So(err, ShouldBeNil)
			So(newBody, ShouldResemble, []byte("tosomething"))
		})

		Convey("when we have a response that errors when read", func() {
			expectedErr := errors.New("Blah")
			body := io.NopCloser(&erroringResponse{expectedErr})
			r := createResponse(req, body)
			err := srrm(&r)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedErr)
		})
	})
}

func createResponse(req *http.Request, body io.ReadCloser) http.Response {
	req = req.WithContext(context.Background())
	return http.Response{
		Body:    body,
		Request: req,
	}
}

type erroringResponse struct {
	Err error
}

func (er *erroringResponse) Read(p []byte) (n int, err error) {
	return 0, er.Err
}
