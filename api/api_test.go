package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ONSdigital/dp-maps-api/api"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const invalidURL = "123:32423:invalid"

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()

		apiURL, err := url.Parse("https://api.example:123/")
		So(err, ShouldBeNil)
		cfg := api.Config{OrdnanceSurveyAPIURL: apiURL}
		api, err := api.Setup(ctx, cfg, r)
		So(err, ShouldBeNil)

		Convey("When created the following routes should have been added", func() {
			// Add to check below with any newly added api endpoints
			So(hasRoute(api.Router, "/maps/vector/v1/vts/resources/styles", http.MethodGet), ShouldBeTrue)
			So(hasRoute(api.Router, "/maps/vector/v1/vts", http.MethodGet), ShouldBeTrue)
			So(hasRoute(api.Router, "/maps/vector/v1/vts/something", http.MethodGet), ShouldBeTrue)
		})

		Convey("When created invalid routes should bot have been added", func() {
			So(hasRoute(api.Router, "/something/invalid", http.MethodGet), ShouldBeFalse)
		})
	})
}

func TestSetup_BadConfig(t *testing.T) {
	Convey("Given an empty API config", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		cfg := api.Config{}
		Convey("When setup is called it should error due to an invalid url", func() {
			api, err := api.Setup(ctx, cfg, r)
			So(api, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "nil OrdnanceSurveyAPIURL supplied in api setup config")

		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
