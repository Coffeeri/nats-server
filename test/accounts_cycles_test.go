// Copyright 2020 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"os"
	"strings"
	"testing"

	"github.com/nats-io/nats-server/v2/server"
)

func TestAccountCycleService(t *testing.T) {
	conf := createConfFile(t, []byte(`
		accounts {
		  A {
		    exports [ { service: help } ]
			imports [ { service { subject: help, account: B } } ]
		  }
		  B {
		    exports [ { service: help } ]
			imports [ { service { subject: help, account: A } } ]
		  }
		}
	`))
	defer os.Remove(conf)

	if _, err := server.ProcessConfigFile(conf); err == nil || !strings.Contains(err.Error(), server.ErrImportFormsCycle.Error()) {
		t.Fatalf("Expected an error on cycle service import, got none")
	}

	conf = createConfFile(t, []byte(`
		accounts {
		  A {
		    exports [ { service: * } ]
			imports [ { service { subject: help, account: B } } ]
		  }
		  B {
		    exports [ { service: help } ]
			imports [ { service { subject: *, account: A } } ]
		  }
		}
	`))
	defer os.Remove(conf)

	if _, err := server.ProcessConfigFile(conf); err == nil || !strings.Contains(err.Error(), server.ErrImportFormsCycle.Error()) {
		t.Fatalf("Expected an error on cycle service import, got none")
	}

	conf = createConfFile(t, []byte(`
		accounts {
		  A {
		    exports [ { service: * } ]
			imports [ { service { subject: help, account: B } } ]
		  }
		  B {
		    exports [ { service: help } ]
			imports [ { service { subject: help, account: C } } ]
		  }
		  C {
		    exports [ { service: * } ]
			imports [ { service { subject: *, account: A } } ]
		  }
		}
	`))
	defer os.Remove(conf)

	if _, err := server.ProcessConfigFile(conf); err == nil || !strings.Contains(err.Error(), server.ErrImportFormsCycle.Error()) {
		t.Fatalf("Expected an error on cycle service import, got none")
	}
}

func TestAccountCycleStream(t *testing.T) {
	conf := createConfFile(t, []byte(`
		accounts {
		  A {
		    exports [ { stream: strm } ]
			imports [ { stream { subject: strm, account: B } } ]
		  }
		  B {
		    exports [ { stream: strm } ]
			imports [ { stream { subject: strm, account: A } } ]
		  }
		}
	`))
	defer os.Remove(conf)
	if _, err := server.ProcessConfigFile(conf); err == nil || !strings.Contains(err.Error(), server.ErrImportFormsCycle.Error()) {
		t.Fatalf("Expected an error on cyclic import, got none")
	}
}

func TestAccountCycleStreamWithMapping(t *testing.T) {
	conf := createConfFile(t, []byte(`
		accounts {
		  A {
		    exports [ { stream: * } ]
			imports [ { stream { subject: bar, account: B } } ]
		  }
		  B {
		    exports [ { stream: bar } ]
			imports [ { stream { subject: foo, account: A }, to: bar } ]
		  }
		}
	`))
	defer os.Remove(conf)
	if _, err := server.ProcessConfigFile(conf); err == nil || !strings.Contains(err.Error(), server.ErrImportFormsCycle.Error()) {
		t.Fatalf("Expected an error on cyclic import, got none")
	}
}

func TestAccountCycleNonCycleStreamWithMapping(t *testing.T) {
	conf := createConfFile(t, []byte(`
		accounts {
		  A {
		    exports [ { stream: foo } ]
			imports [ { stream { subject: bar, account: B } } ]
		  }
		  B {
		    exports [ { stream: bar } ]
			imports [ { stream { subject: baz, account: C }, to: bar } ]
		  }
		  C {
		    exports [ { stream: baz } ]
			imports [ { stream { subject: foo, account: A }, to: bar } ]
		  }
		}
	`))
	defer os.Remove(conf)
	if _, err := server.ProcessConfigFile(conf); err != nil {
		t.Fatalf("Expected no error but got %s", err)
	}
}

func TestAccountCycleServiceCycleWithMapping(t *testing.T) {
	conf := createConfFile(t, []byte(`
		accounts {
		  A {
		    exports [ { service: a } ]
			imports [ { service { subject: b, account: B }, to: a } ]
		  }
		  B {
		    exports [ { service: b } ]
			imports [ { service { subject: a, account: A }, to: b } ]
		  }
		}
	`))
	defer os.Remove(conf)
	if _, err := server.ProcessConfigFile(conf); err == nil || !strings.Contains(err.Error(), server.ErrImportFormsCycle.Error()) {
		t.Fatalf("Expected an error on cycle service import, got none")
	}
}

func TestAccountCycleServiceNonCycle(t *testing.T) {
	conf := createConfFile(t, []byte(`
		accounts {
		  A {
		    exports [ { service: * } ]
			imports [ { service { subject: help, account: B } } ]
		  }
		  B {
		    exports [ { service: help } ]
			imports [ { service { subject: nohelp, account: C } } ]
		  }
		  C {
		    exports [ { service: * } ]
			imports [ { service { subject: *, account: A } } ]
		  }
		}
	`))
	defer os.Remove(conf)

	if _, err := server.ProcessConfigFile(conf); err != nil {
		t.Fatalf("Expected no error but got %s", err)
	}
}

func TestAccountCycleServiceNonCycleChain(t *testing.T) {
	conf := createConfFile(t, []byte(`
		accounts {
		  A {
		    exports [ { service: help } ]
			imports [ { service { subject: help, account: B } } ]
		  }
		  B {
		    exports [ { service: help } ]
			imports [ { service { subject: help, account: C } } ]
		  }
		  C {
		    exports [ { service: help } ]
			imports [ { service { subject: help, account: D } } ]
		  }
		  D {
		    exports [ { service: help } ]
		  }
		}
	`))
	defer os.Remove(conf)

	if _, err := server.ProcessConfigFile(conf); err != nil {
		t.Fatalf("Expected no error but got %s", err)
	}
}
