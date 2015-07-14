//  Copyright (c) 2015 Rackspace
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
//  implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package proxyserver

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openstack/swift/go/hummingbird"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	expectedBody := "OK"
	conf := "/etc/swift/proxy-server.conf"
	ip, port, server, _, _ := GetServer(conf, &flag.FlagSet{})
	handler := server.GetHandler()
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("http://%s:%d/healthcheck", ip, port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Error("Unable to create new Request")
	}
	handler.ServeHTTP(recorder, req)

	if recorder.Body.String() != expectedBody {
		t.Error("Excepting ", expectedBody, " got ", recorder.Body.String())
	}
}

func TestGetServer(t *testing.T) {
	tests := []struct {
		conf    string
		err_msg string
	}{
		{"/etc/swift/proxy-server.conf", ""},
		{"/tmp/asdf", "Unable to load /tmp/asdf"},
	}
	for _, test := range tests {
		if test.err_msg != "" {
			_, _, _, _, err := GetServer(test.conf, &flag.FlagSet{})
			assert.Equal(t, test.err_msg, err.Error())
			continue
		}
		ip, port, server, _, _ := GetServer(test.conf, &flag.FlagSet{})
		assert.NotNil(t, net.ParseIP(ip))
		assert.Equal(t, port, 8080)
		assert.NotPanics(t, func() {
			_ = server.(hummingbird.Server)
		})

	}
}
