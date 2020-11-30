// Copyright 2020 Ye Zi Jie. All Rights Reserved.
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
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/11/30 22:17:07

package logit

import (
	"os"
	"sync"
)

// Roller is an interface of rolling a file writer.
// File writer will call TimeToRoll() to know if time to roll before writing,
// and if true, the Roll() will be called.
type Roller interface {

	// TimeToRoll returns true if need rolling or false.
	// Although file is a pointer, you shouldn't change it in this method.
	// Remember, file in this method should be read only.
	TimeToRoll(fw *FileWriter) bool

	// Roll will roll this file and returns an error if failed.
	// Although file is a pointer, you shouldn't change it in this method.
	// Return an os.File instance will be fine.
	Roll(fw *FileWriter) (*os.File, error)
}

// =================================== file writer ===================================

// FileWriter writes logs to one or more files.
// We provide a powerful writer of file, which can roll to several files
// automatically in some conditions.
type FileWriter struct {

	// name is the name of log file.
	name string

	// file is a pointer to the real file in os.
	file *os.File

	// seq is the serial number of rolling file.
	// If name is "test.log" and seq is 1, then the file rolled will be "test.log.0000000001"
	seq uint32

	// rollers stores all rollers will be used.
	// If one of them say: "Time to roll!", then this file writer will start to roll.
	// After rolling, the rollers after it will not be used this loop.
	rollers []Roller

	// lock is for safe-concurrency.
	lock *sync.RWMutex
}

// NewFileWriter returns a new file writer with given name and rollers.
func NewFileWriter(name string, rollers ...Roller) (*FileWriter, error) {

	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		name:    name,
		file:    file,
		seq:     0,
		rollers: rollers,
		lock:    &sync.RWMutex{},
	}, nil
}

// Write writes len(p) bytes from p to file and returns an error if failed.
// The precise count of written bytes is n.
func (fw *FileWriter) Write(p []byte) (n int, err error) {

	fw.lock.RLock()
	defer fw.lock.RUnlock()

	// Check rolling condition first and replace to newFile only if roll returns nil error
	for _, roller := range fw.rollers {
		if roller.TimeToRoll(fw) {
			if newFile, err := roller.Roll(fw); err == nil {
				fw.file = newFile
				break
			}
		}
	}
	return fw.file.Write(p)
}

// Close closes this file writer and returns an error if failed.
func (fw *FileWriter) Close() error {
	fw.lock.Lock()
	defer fw.lock.Unlock()
	fw.rollers = nil
	return fw.file.Close()
}
