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
	"path/filepath"
	"sort"
	"time"
)

const (
	backupSeparator = "."
)

type backup struct {
	path string
	t    time.Time
}

func (b backup) BeforeTime(t time.Time) bool {
	return b.t.Before(t)
}

func (b backup) Before(other backup) bool {
	return b.t.Before(other.t)
}

func sortBackups(backups []backup) {
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Before(backups[j])
	})
}

func backupPrefixAndExt(path string) (string, string) {
	ext := filepath.Ext(path)
	prefix := path[:len(path)-len(ext)] + backupSeparator
	return prefix, ext
}

func backupPath(path string, timeFormat string) string {
	name, ext := backupPrefixAndExt(path)
	now := now().Format(timeFormat)
	return name + now + ext
}

func parseBackupTime(filename string, prefix string, ext string, timeFormat string) (time.Time, error) {
	ts := filename[len(prefix) : len(filename)-len(ext)]
	return time.Parse(timeFormat, ts)
}
