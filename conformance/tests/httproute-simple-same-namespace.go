/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tests

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/gateway-api/conformance/utils/http"
	"sigs.k8s.io/gateway-api/conformance/utils/kubernetes"
	"sigs.k8s.io/gateway-api/conformance/utils/roundtripper"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, HTTPRouteSimpleSameNamespace)
}

var HTTPRouteSimpleSameNamespace = suite.ConformanceTest{
	ShortName:   "HTTPRouteSimpleSameNamespace",
	Description: "A single HTTPRoute in the gateway-conformance-infra namespace attaches to a Gateway in the same namespace",
	Manifests:   []string{"tests/httproute-simple-same-namespace.yaml"},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		ns := v1alpha2.Namespace("gateway-conformance-infra")
		routeNN := types.NamespacedName{Name: "gateway-conformance-infra-test", Namespace: string(ns)}
		gwNN := types.NamespacedName{Name: "same-namespace", Namespace: string(ns)}
		gwAddr := kubernetes.GatewayAndHTTPRoutesMustBeReady(t, suite.Client, suite.ControllerName, gwNN, routeNN)

		t.Run("Simple HTTP request should reach infra-backend", func(t *testing.T) {
			t.Logf("Making request to http://%s", gwAddr)
			cReq, cRes, err := suite.RoundTripper.CaptureRoundTrip(roundtripper.Request{
				URL:      url.URL{Scheme: "http", Host: gwAddr},
				Protocol: "HTTP",
			})

			require.NoErrorf(t, err, "error making request")

			http.ExpectResponse(t, cReq, cRes, http.ExpectedResponse{
				Request: http.ExpectedRequest{
					Method: "GET",
					Path:   "/",
				},
				StatusCode: 200,
				Backend:    "infra-backend-v1",
				Namespace:  "gateway-conformance-infra",
			})
		})
	},
}
