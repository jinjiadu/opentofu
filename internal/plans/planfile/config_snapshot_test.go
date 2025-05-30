// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package planfile

import (
	"archive/zip"
	"bytes"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/opentofu/opentofu/internal/configs"
	"github.com/opentofu/opentofu/internal/configs/configload"
)

func TestConfigSnapshotRoundtrip(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "test-config")
	loader, err := configload.NewLoader(&configload.Config{
		ModulesDir: filepath.Join(fixtureDir, ".terraform", "modules"),
	})
	if err != nil {
		t.Fatal(err)
	}

	_, snapIn, diags := loader.LoadConfigWithSnapshot(t.Context(), fixtureDir, configs.RootModuleCallForTesting())
	if diags.HasErrors() {
		t.Fatal(diags.Error())
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	err = writeConfigSnapshot(snapIn, zw)
	if err != nil {
		t.Fatalf("failed to write snapshot: %s", err)
	}
	zw.Close()

	raw := buf.Bytes()
	r := bytes.NewReader(raw)
	zr, err := zip.NewReader(r, int64(len(raw)))
	if err != nil {
		t.Fatal(err)
	}

	snapOut, err := readConfigSnapshot(zr)
	if err != nil {
		t.Fatalf("failed to read snapshot: %s", err)
	}

	if !reflect.DeepEqual(snapIn, snapOut) {
		t.Errorf("result does not match input\nresult: %sinput: %s", spew.Sdump(snapOut), spew.Sdump(snapIn))
	}
}
