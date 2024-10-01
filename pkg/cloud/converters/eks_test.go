/*
Copyright 2020 The Kubernetes Authors.

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

package converters

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/stretchr/testify/assert"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestTaintsFromSDK(t *testing.T) {
	tests := []struct {
		name      string
		input     []*eks.Taint
		expected  expinfrav1.Taints
		expectErr bool
	}{
		{
			name: "Single taint conversion",
			input: []*eks.Taint{
				{
					Key:    awsString("key1"),
					Value:  awsString("value1"),
					Effect: awsString(eks.TaintEffectNoSchedule),
				},
			},
			expected: expinfrav1.Taints{
				{
					Key:    "key1",
					Value:  "value1",
					Effect: expinfrav1.TaintEffectNoSchedule,
				},
			},
			expectErr: false,
		},
		{
			name: "Multiple taints conversion",
			input: []*eks.Taint{
				{
					Key:    awsString("key1"),
					Value:  awsString("value1"),
					Effect: awsString(eks.TaintEffectNoExecute),
				},
				{
					Key:    awsString("key2"),
					Value:  awsString("value2"),
					Effect: awsString(eks.TaintEffectPreferNoSchedule),
				},
			},
			expected: expinfrav1.Taints{
				{
					Key:    "key1",
					Value:  "value1",
					Effect: expinfrav1.TaintEffectNoExecute,
				},
				{
					Key:    "key2",
					Value:  "value2",
					Effect: expinfrav1.TaintEffectPreferNoSchedule,
				},
			},
			expectErr: false,
		},
		{
			name: "Unknown taint effect",
			input: []*eks.Taint{
				{
					Key:    awsString("key1"),
					Value:  awsString("value1"),
					Effect: awsString("UnknownEffect"),
				},
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Nil taint value",
			input: []*eks.Taint{
				{
					Key:    awsString("key1"),
					Value:  nil,
					Effect: awsString(eks.TaintEffectNoSchedule),
				},
			},
			expected: expinfrav1.Taints{
				{
					Key:    "key1",
					Value:  "",
					Effect: expinfrav1.TaintEffectNoSchedule,
				},
			},
			expectErr: false,
		},
		{
			name:      "Empty input slice",
			input:     []*eks.Taint{},
			expected:  expinfrav1.Taints{},
			expectErr: false,
		},
		{
			name:      "Nil input",
			input:     nil,
			expected:  expinfrav1.Taints{},
			expectErr: false,
		},
		{
			name: "Mixed valid and invalid taints",
			input: []*eks.Taint{
				{
					Key:    awsString("key1"),
					Value:  awsString("value1"),
					Effect: awsString(eks.TaintEffectNoExecute),
				},
				{
					Key:    awsString("key2"),
					Value:  awsString("value2"),
					Effect: awsString("InvalidEffect"),
				},
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Empty key",
			input: []*eks.Taint{
				{
					Key:    awsString(""),
					Value:  awsString("value1"),
					Effect: awsString(eks.TaintEffectNoSchedule),
				},
			},
			expected: expinfrav1.Taints{
				{
					Key:    "",
					Value:  "value1",
					Effect: expinfrav1.TaintEffectNoSchedule,
				},
			},
			expectErr: false,
		},
		{
			name: "Empty value",
			input: []*eks.Taint{
				{
					Key:    awsString("key1"),
					Value:  awsString(""),
					Effect: awsString(eks.TaintEffectNoSchedule),
				},
			},
			expected: expinfrav1.Taints{
				{
					Key:    "key1",
					Value:  "",
					Effect: expinfrav1.TaintEffectNoSchedule,
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TaintsFromSDK(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result, "Expected result to be nil on error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result, "Converted taints do not match expected")
			}
		})
	}
}

// Helper function to create AWS String pointers.
func awsString(value string) *string {
	return &value
}
