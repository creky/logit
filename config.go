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
// Author: fish
// Email: fishinlove@163.com
// Created at 2020/03/25 23:06:59

package logit

import (
    "bytes"
    "encoding/json"
    "errors"
    "io/ioutil"
)

var (
    HandlerDoesNotHaveNameError = errors.New("handler should provide a name which has been used to register to logit")
)

// Config is for configuring a Logger.
// You can use a config to create a Logger to use.
// See NewLoggerFrom(config Config).
type Config struct {

    // Level is the level of Logger.
    // If the level of log is smaller than this Level, this log will be ignored.
    Level Level

    Handlers []Handler
}

// fileConfig is the config mapping a file.
type fileConfig struct {
    Level string `json:"level"`

    Handlers []map[string]string `json:"handlers"`
}

func removeComments(fileInBytes []byte) []byte {
    // 注释只能是单独起一行，并且以 # 开头
    var buffer []byte
    lines := bytes.Split(fileInBytes, []byte("\n"))
    for _, line := range lines {
        if !bytes.HasPrefix(bytes.TrimSpace(line), []byte("#")) {
            buffer = append(buffer, line...)
        }
    }
    return buffer
}

func parseConfig(fileConfig fileConfig) (Config, error) {

    // 解析日志级别
    level, err := ParseLevel(fileConfig.Level)
    if err != nil {
        return Config{}, err
    }

    // 解析日志处理器
    var handlers []Handler
    for _, handler := range fileConfig.Handlers {
        name, ok := handler["name"]
        if !ok {
            return Config{}, HandlerDoesNotHaveNameError
        }
        delete(handler, "name")
        handlers = append(handlers, HandlerOf(name, handler))
    }

    return Config{
        Level:    level,
        Handlers: handlers,
    }, nil
}

func ParseConfigFile(configFile string) (Config, error) {

    // 配置文件一般不会太大，直接全部读取进内存
    fileInBytes, err := ioutil.ReadFile(configFile)
    if err != nil {
        return Config{}, err
    }

    // 由于配置文件使用 Json 格式，而 Json 规范要求使用 {} 包裹内容，但这个 {} 不方便配置文件的阅读，
    // 所以我们设定在配置文件中不使用 {} 包裹，而是交给我们读取出来之后进行包裹
    // 另外，我们的配置文件是支持注释的，而 Json 规范中并没有对注释的支持，所以我们需要对注释进行擦除
    configInBytes := bytes.Join([][]byte{[]byte("{"), removeComments(fileInBytes), []byte("}")}, []byte(""))
    fileConfig := fileConfig{
        Level: "debug",
    }
    err = json.Unmarshal(configInBytes, &fileConfig)
    if err != nil {
        return Config{}, err
    }

    return parseConfig(fileConfig)
}
