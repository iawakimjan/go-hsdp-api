package iam

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServicesCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	serviceName := "testservice"
	serviceDescription := "Service description"
	applicationID := "b0889958-3762-4427-af07-2d6268c26988"
	serviceID := serviceName + ".testapp.testprop@testdev.devorg.1e100.io"
	id := "2c266886-f918-4223-941d-437cb3cd09e8"
	orgID := "bda40124-54fa-4967-b2fb-23dcc4e0ad1a"

	muxIDM.HandleFunc("/authorize/identity/Service", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
				"id": "`+id+`",
				"serviceId": "`+serviceID+`",
				"organizationId": "`+orgID+`",
				"privateKey": "-----BEGIN RSA PRIVATE KEY-----MIIEowIBAAKCAQEAh2kjyJRO/rJuwMOECMcR5ZoSCwIsq3205gABA8BNGAZk1CGxxbduIX3peZMdMC5nyQ2gVVHj97Bsx9sYHXs3ihUVzeCth9SEipmrIuktrIzMMbN1Hd0DBMYEVfgxER/F5yYvt3D2RwnydBOg9QSc7GoFa6zzJMfplyVQyCiDgCYN9qf+BHTPk7+x/7ev835ylBsYOFHebN0WsXUlEBWc8nCSJm0Z0hFcmUeLcVSsB2X7iMqGq4VpgSsRfP9sloCX8lNzDo4ujwvTIwsCpHwokQhxCV2avhmjH2hMQXHGlgUfyPm3RjQZ0waSGRRCEuWDhm7xLq9hRzad4cD2auVwhQIDAQABAoIBAD3YwrRBMNdZzgYTBsIvkjgJJ8aJZrepAa+vPsdk1JFtki3ledmxTwbTCIkzrTgtac/Ffn6ZmYKuvPCHXDtS5OoXeU8AGKIaabMYPrcCQ480+6qTqaFLKa7Ldn2Bj3+fwHcz1MV3PbTykR99O53NTpMYVYN5idA50rHrJDtXbcBgdc8KDA79keN6Fv6pggm//Vms7/E9/bgYYY1W8FmjngNHHDgH2EEnERMkyDp2Ng2/2jJIAa8uKPPFxjJAtApQjlLEv82MOcq+Yeq2/VXGRbvgxWlhdEUQcfCf2/ncFEhslau2llcgKTyvcqbaB+jI9Fj8AeXkVyAGVA4X3fuodp0CgYEA0p3oTonVovr+hXELd8hshElAQ+xP4KtPRFH00VPBkB5m3rAALVK+brDAzzZ7W/5X/QTpELcIfIbx5tbFOcQXjuakxVrbof0wZEDGT9KGPSqAAScp7GqYKMSNWUdZ3G5bLisnHT0q7jdWw9se5HGYPEBxYGW8NgbSX+ejsd84ydsCgYEApJax/f+h5RkTUx9TBSRoTz7jMmOls8bbC3IiAUSlVd3wyJX98etsiIsGFh0mrG16Z8Kg0C597woFoUEWTSv7vzN+SLH2CJAt7JzDdlbgmj7EtB0WJwOCPTRUCbQBeMvy3G34NwAlQxGcQqwbo9JMx/8haxqSgCpfrKaD0/CUrR8CgYBgbwqpwzR9Lj0RbkQY8Ty2iS+SqgWc0fM2TewxWA8dZL4nIiDCn8svtWBiwAhVg6xX3kK0c4nAMq1Zy2Z8X4uF05cIAeTkU6AvlvT2IWdzZB096eepJtlKeUxa32+GnUTEa9+55ILelZn1jUOkx1oz5DHFOG+nsRHr9Yye6Zz/1wKBgCAQF0aS6Rf3RZN407R5vjRJ3Pqw/NPD1mIpbsRuegL7RG/fAGSDZ1ZGNv5R2XnXrfPOr4M+u1u4yRX71vtbqSQ7RMuml3ZdmASzGUTRcdm6hplL3UfmYBXKPuDRB0Rf/sTAS41zYs7o/FbkrlHAoyKG6hyyRX3gQ1kf6yh7gosjAoGBAJv0O9x0oY7HZ5QF1PNlDLZUSF/8UtdmyKnm+6VLBmZBaOBJ7MWOJUzzJOHgMRZyzSbH1Z7aBfDXEdXuPihzel/m6TTt/sBv6P8UbVz8cA/uK5BC6nYVwhhu+/wQT/pWNxfqERDKcS7LJ+F3XWkngYA1fBVFmwwDeDWhzgcLQmU6-----END RSA PRIVATE KEY-----",
				"expiresOn": "2019-08-15T17:38:06.322Z",
				"name": "`+serviceName+`",
				"applicationId": "`+applicationID+`",
				"defaultScopes": [
					"openid"
				],
				"scopes": [
					"openid"
				]
			}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
				"total": 1,
				"entry": [
					{
						"id": "`+id+`",
						"serviceId": "`+serviceID+`",
						"organizationId": "`+orgID+`",
						"expiresOn": "2019-08-15T17:38:06.322Z",
						"name": "testservice",
						"applicationId": "`+applicationID+`",
						"defaultScopes": [
							"openid"
						],
						"scopes": [
							"openid"
						]
					}
				]
			}`)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/Service/"+id, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	var r Service
	r.Name = serviceName
	r.Description = serviceDescription
	r.ApplicationID = applicationID

	service, resp, err := client.Services.CreateService(r)
	if ok := assert.Nil(t, err); !ok {
		return
	}
	if ok := assert.NotNil(t, resp); !ok {
		return
	}
	if ok := assert.NotNil(t, service); !ok {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, serviceName, service.Name)
	assert.Equal(t, applicationID, service.ApplicationID)
	assert.True(t, len(service.PrivateKey) > 0)

	foundService, resp, err := client.Services.GetServiceByID(service.ID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.NotNil(t, foundService) {
		return
	}
	assert.Equal(t, service.ID, foundService.ID)

	services, resp, err := client.Services.GetServicesByApplicationID(applicationID)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, services) {
		return
	}
	assert.Equal(t, 1, len(*services))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ok, resp, err := client.Services.DeleteService(*foundService)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.True(t, ok)
}

func TestScopes(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	serviceName := "testservice"
	serviceDescription := "Service description"
	applicationID := "b0889958-3762-4427-af07-2d6268c26988"
	id := "2c266886-f918-4223-941d-437cb3cd09e8"

	muxIDM.HandleFunc("/authorize/identity/Service/"+id+"/$scopes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "PUT":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	var r Service
	r.ID = id
	r.Name = serviceName
	r.Description = serviceDescription
	r.ApplicationID = applicationID

	ok, resp, err := client.Services.AddScopes(r, []string{"foo"}, []string{"foo"})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, ok) {
		return
	}
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	ok, resp, err = client.Services.RemoveScopes(r, []string{"foo"}, []string{"foo"})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, ok) {
		return
	}
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
