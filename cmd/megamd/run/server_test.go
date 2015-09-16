package run

import (
//	"fmt"
//	"net/http"
//	"net/url"
//	"strconv"
//	"strings"
//	"testing"
//	"time"
)


/*
// Ensure that HTTP responses include the InfluxDB version.
func TestServer_HTTPResponseVersion(t *testing.T) {
	version := "v1234"
	s := OpenServerWithVersion(NewConfig(), version)
	defer s.Close()

	resp, _ := http.Get(s.URL() + "/query")
	got := resp.Header.Get("X-Influxdb-Version")
	if got != version {
		t.Errorf("Server responded with incorrect version, exp %s, got %s", version, got)
	}
}

// Ensure the database commands work.
func TestServer_DatabaseCommands(t *testing.T) {
	t.Parallel()
	s := OpenServer(NewConfig(), "")
	defer s.Close()

	test := Test{
		queries: []*Query{
			&Query{
				name:    "create database should succeed",
				command: `CREATE DATABASE db0`,
				exp:     `{"results":[{}]}`,
			},
			&Query{
				name:    "create database should error with bad name",
				command: `CREATE DATABASE 0xdb0`,
				exp:     `{"error":"error parsing query: found 0, expected identifier at line 1, char 17"}`,
			},
			&Query{
				name:    "show database should succeed",
				command: `SHOW DATABASES`,
				exp:     `{"results":[{"series":[{"name":"databases","columns":["name"],"values":[["db0"]]}]}]}`,
			},
			&Query{
				name:    "create database should error if it already exists",
				command: `CREATE DATABASE db0`,
				exp:     `{"results":[{"error":"database already exists"}]}`,
			},
			&Query{
				name:    "drop database should succeed",
				command: `DROP DATABASE db0`,
				exp:     `{"results":[{}]}`,
			},
			&Query{
				name:    "show database should have no results",
				command: `SHOW DATABASES`,
				exp:     `{"results":[{"series":[{"name":"databases","columns":["name"]}]}]}`,
			},
			&Query{
				name:    "drop database should error if it doesn't exist",
				command: `DROP DATABASE db0`,
				exp:     `{"results":[{"error":"database not found: db0"}]}`,
			},
		},
	}

	for _, query := range test.queries {
		if query.skip {
			t.Logf("SKIP:: %s", query.name)
			continue
		}
		if err := query.Execute(s); err != nil {
			t.Error(query.Error(err))
		} else if !query.success() {
			t.Error(query.failureMessage())
		}
	}
}

func TestServer_Query_DropAndRecreateDatabase(t *testing.T) {
	t.Parallel()
	s := OpenServer(NewConfig(), "")
	defer s.Close()

	if err := s.CreateDatabaseAndRetentionPolicy("db0", newRetentionPolicyInfo("rp0", 1, 0)); err != nil {
		t.Fatal(err)
	}
	if err := s.MetaStore.SetDefaultRetentionPolicy("db0", "rp0"); err != nil {
		t.Fatal(err)
	}

	writes := []string{
		fmt.Sprintf(`cpu,host=serverA,region=uswest val=23.2 %d`, mustParseTime(time.RFC3339Nano, "2000-01-01T00:00:00Z").UnixNano()),
	}

	test := NewTest("db0", "rp0")
	test.write = strings.Join(writes, "\n")

	test.addQueries([]*Query{
		&Query{
			name:    "Drop database after data write",
			command: `DROP DATABASE db0`,
			exp:     `{"results":[{}]}`,
		},
		&Query{
			name:    "Recreate database",
			command: `CREATE DATABASE db0`,
			exp:     `{"results":[{}]}`,
		},
		&Query{
			name:    "Recreate retention policy",
			command: `CREATE RETENTION POLICY rp0 ON db0 DURATION 365d REPLICATION 1 DEFAULT`,
			exp:     `{"results":[{}]}`,
		},
		&Query{
			name:    "Show measurements after recreate",
			command: `SHOW MEASUREMENTS`,
			exp:     `{"results":[{}]}`,
			params:  url.Values{"db": []string{"db0"}},
		},
		&Query{
			name:    "Query data after recreate",
			command: `SELECT * FROM cpu`,
			exp:     `{"results":[{}]}`,
			params:  url.Values{"db": []string{"db0"}},
		},
	}...)

	for i, query := range test.queries {
		if i == 0 {
			if err := test.init(s); err != nil {
				t.Fatalf("test init failed: %s", err)
			}
		}
		if query.skip {
			t.Logf("SKIP:: %s", query.name)
			continue
		}
		if err := query.Execute(s); err != nil {
			t.Error(query.Error(err))
		} else if !query.success() {
			t.Error(query.failureMessage())
		}
	}
}

*/
