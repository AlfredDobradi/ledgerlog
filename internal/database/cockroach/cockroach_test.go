package cockroach

import (
	"testing"
)

func TestBuildConnectionString(t *testing.T) {
	SetUser("test")
	SetPassword("pass")
	SetHost("1.1.1.1")
	SetPort("33333")
	SetDatabase("testingdb")
	SetSSLMode("verify-full")
	SetSSLRootCert("asd.crt")
	SetCluster("test-cluster")
	expectedConnectionString := "postgresql://test:pass@1.1.1.1:33333/test-cluster.testingdb?sslmode=verify-full&sslrootcert=asd.crt"

	if actual, expected := buildConnectionString(), expectedConnectionString; actual != expected {
		t.Fatalf("Fail: Connection strings don't match.\nExpected: %s\nGot: %s\n", expected, actual)
	}

	t.Log("Pass")
}
