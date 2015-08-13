/**
 * Copyright (c) 2015 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {

	assert.NotNil(t, args, "nil config args")
	assert.NotNil(t, args.Pub, "nil pub")
	assert.NotNil(t, args.Pub.URI, "nil uri")
	assert.True(t, len(args.Pub.URI) > 0, "no uris")
	assert.NotNil(t, args.Server, "nil server")
	assert.NotNil(t, args.Server.Host, "nil host")
	assert.NotEqual(t, args.Server.Port, 0, "nil Port")
	assert.NotNil(t, args.Server.Root, "nil root")
	assert.NotNil(t, args.Server.Token, "nil root")

}
