/*
Copyright 2022 The KubeService-Stack Authors.

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

package framework

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

type TestCtx struct {
	id         string
	cleanUpFns []FinalizerFn
}

type FinalizerFn func() error

func (f *Framework) NewTestCtx(t *testing.T) *TestCtx {
	// TestCtx is used among others for namespace names where '/' is forbidden
	prefix := strings.TrimPrefix(
		strings.ReplaceAll(
			strings.ToLower(t.Name()),
			"/",
			"-",
		),
		"test",
	)

	id := prefix + "-" + strconv.FormatInt(time.Now().Unix(), 36)
	return &TestCtx{
		id: id,
	}
}

func (ctx *TestCtx) ID() string {
	return ctx.id + "-" + strconv.Itoa(len(ctx.cleanUpFns))
}

func (ctx *TestCtx) Cleanup(t *testing.T) {
	var eg errgroup.Group

	for i := len(ctx.cleanUpFns) - 1; i >= 0; i-- {
		eg.Go(ctx.cleanUpFns[i])
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}

func (ctx *TestCtx) AddFinalizerFn(fn FinalizerFn) {
	ctx.cleanUpFns = append(ctx.cleanUpFns, fn)
}
