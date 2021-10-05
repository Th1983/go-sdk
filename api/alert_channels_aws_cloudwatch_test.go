//
// Author:: Vatasha White (<vatasha.white@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestAlertChannelsService_GetCloudwatchEb(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		apiPath     = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer  = lacework.MockServer()
		eventBusArn = "arn:aws:events:YourRegion:YourAccountID:event-bus/default"
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetCloudwatchEb() should be a GET method")
		fmt.Fprintf(w, generateAlertChannelResponse(singleAWSCloudwatchAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.GetCloudwatchEb(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, eventBusArn, response.Data.Data.EventBusArn)
	assert.Equal(t, "Events", response.Data.Data.IssueGrouping)
}

func TestAlertChannelsService_UpdateCloudwatchEb(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		apiPath     = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer  = lacework.MockServer()
		eventBusArn = "arn:aws:events:YourRegion:YourAccountID:event-bus/default"
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateCloudwatchEb() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "CloudwatchEb", "wrong cloud account type")
			assert.Contains(t, body, eventBusArn, "missing event bus arn")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateAlertChannelResponse(singleAWSCloudwatchAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	emailAlertChan := api.NewAlertChannel("integration_name",
		api.CloudwatchEbAlertChannelType,
		api.CloudwatchEbDataV2{EventBusArn: eventBusArn},
	)
	assert.Equal(t, "integration_name", emailAlertChan.Name)
	assert.Equal(t, "CloudwatchEb", emailAlertChan.Type)
	assert.Equal(t, 1, emailAlertChan.Enabled)
	emailAlertChan.IntgGuid = intgGUID

	response, err := c.V2.AlertChannels.UpdateCloudwatchEb(emailAlertChan)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, eventBusArn, response.Data.Data.EventBusArn)
	assert.Equal(t, "Events", response.Data.Data.IssueGrouping)
}

func singleAWSCloudwatchAlertChannel(id string) string {
	return fmt.Sprintf(`
	{
		"createdOrUpdatedBy": "vatasha.white@lacework.net",
		"createdOrUpdatedTime": "2021-09-29T117:55:47.277316",
		"data": {
			"eventBusArn": "arn:aws:events:YourRegion:YourAccountID:event-bus/default",
			"issueGrouping": "Events"
		},
		"enabled": 1,
		"intgGuid": %q,
		"isOrg": 0,
		"name": "integration_name",
		"state": {
		"details": {},
		"lastSuccessfulTime": 1632932665892,
			"lastUpdatedTime": 1632932665892,
			"ok": true
	},
		"type": "CloudwatchEb"
	}
	`, id)
}