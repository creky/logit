// Copyright 2022 FishGoddess. All Rights Reserved.
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

package file

import (
	"os"
	"time"

	"github.com/go-logit/logit/support/size"
)

const (
	day = 24 * time.Hour
)

type config struct {
	mode       os.FileMode
	timeFormat string
	maxSize    size.ByteSize
	maxAge     time.Duration
	maxBackups uint
}

func newDefaultConfig() config {
	return config{
		mode:       0644,
		timeFormat: "20060102-150405.000",
		maxSize:    256 * size.MB,
		maxAge:     14 * day,
		maxBackups: 14,
	}
}
