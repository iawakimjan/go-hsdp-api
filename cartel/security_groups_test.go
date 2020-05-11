package cartel

import (
	"testing"

	"net/http"

	"github.com/stretchr/testify/assert"
)

func TestSecurityGroups(t *testing.T) {
	teardown, err := setup(t, Config{
		Token:  sharedToken,
		Secret: sharedSecret,
		Host:   "foo",
		NoTLS:  true,
	})

	muxCartel.HandleFunc("/v3/api/get_security_groups", endpointMocker(sharedSecret,
		`[
    "foo",
    "bar",
    "baz"
]`))

	defer teardown()

	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.GetSecurityGroups()
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
