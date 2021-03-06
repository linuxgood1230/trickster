/*
 * Copyright 2018 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package irondb

import (
	"testing"

	"github.com/tricksterproxy/trickster/pkg/backends"
	"github.com/tricksterproxy/trickster/pkg/backends/irondb/model"
	oo "github.com/tricksterproxy/trickster/pkg/backends/options"
	cr "github.com/tricksterproxy/trickster/pkg/cache/registration"
	"github.com/tricksterproxy/trickster/pkg/config"
	tl "github.com/tricksterproxy/trickster/pkg/logging"
)

func TestIRONdbClientInterfacing(t *testing.T) {

	// this test ensures the client will properly conform to the
	// Client and TimeseriesClient interfaces

	c := &Client{name: "test"}
	var oc backends.Client = c
	var tc backends.TimeseriesClient = c

	if oc.Name() != "test" {
		t.Errorf("expected %s got %s", "test", oc.Name())
	}

	if tc.Name() != "test" {
		t.Errorf("expected %s got %s", "test", tc.Name())
	}
}

var testModeler = model.NewModeler()

func TestNewClient(t *testing.T) {
	conf, _, err := config.Load("trickster", "test",
		[]string{"-origin-url", "http://example.com", "-provider", "TEST_CLIENT"})
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}

	caches := cr.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer cr.CloseCaches(caches)
	cache, ok := caches["default"]
	if !ok {
		t.Errorf("Could not find default configuration")
	}

	oc := &oo.Options{Provider: "TEST_CLIENT"}
	c, err := NewClient("default", oc, nil, cache, testModeler)
	if err != nil {
		t.Error(err)
	}
	if c.Name() != "default" {
		t.Errorf("expected %s got %s", "default", c.Name())
	}

	if c.Cache().Configuration().Provider != "memory" {
		t.Errorf("expected %s got %s", "memory", c.Cache().Configuration().Provider)
	}

	if c.Configuration().Provider != "TEST_CLIENT" {
		t.Errorf("expected %s got %s", "TEST_CLIENT", c.Configuration().Provider)
	}
}

func TestConfiguration(t *testing.T) {
	oc := &oo.Options{Provider: "TEST"}
	client := Client{config: oc}
	c := client.Configuration()
	if c.Provider != "TEST" {
		t.Errorf("expected %s got %s", "TEST", c.Provider)
	}
}

func TestCache(t *testing.T) {
	conf, _, err := config.Load("trickster", "test",
		[]string{"-origin-url", "http://example.com", "-provider", "TEST_CLIENT"})
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}

	caches := cr.LoadCachesFromConfig(conf, tl.ConsoleLogger("error"))
	defer cr.CloseCaches(caches)
	cache, ok := caches["default"]
	if !ok {
		t.Errorf("Could not find default configuration")
	}

	client := Client{cache: cache}
	c := client.Cache()
	if c.Configuration().Provider != "memory" {
		t.Errorf("expected %s got %s", "memory", c.Configuration().Provider)
	}
}

func TestName(t *testing.T) {
	client := Client{name: "TEST"}
	c := client.Name()
	if c != "TEST" {
		t.Errorf("expected %s got %s", "TEST", c)
	}
}

func TestRouter(t *testing.T) {
	client := Client{name: "TEST"}
	r := client.Router()
	if r != nil {
		t.Error("expected nil router")
	}
}

func TestHTTPClient(t *testing.T) {
	oc := &oo.Options{Provider: "TEST"}
	client, err := NewClient("test", oc, nil, nil, testModeler)
	if err != nil {
		t.Error(err)
	}
	if client.HTTPClient() == nil {
		t.Errorf("missing http client")
	}
}

func TestSetCache(t *testing.T) {
	c, err := NewClient("test", oo.New(), nil, nil, testModeler)
	if err != nil {
		t.Error(err)
	}
	c.SetCache(nil)
	if c.Cache() != nil {
		t.Errorf("expected nil cache for client named %s", "test")
	}
}
