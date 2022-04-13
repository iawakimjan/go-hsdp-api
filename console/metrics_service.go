package console

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hasura/go-graphql-client"
)

type MetricsService struct {
	client *Client
}

type Instance struct {
	CreatedAt    time.Time `json:"createdAt"`
	GUID         string    `json:"guid"`
	Name         string    `json:"name"`
	Organization string    `json:"organization"`
	Space        string    `json:"space"`
	Details      Details   `json:"details"`
}

type Details struct {
	Hostname string `json:"hostname"`
}

type MetricsResponse struct {
	Data struct {
		Instances []Instance `json:"instances"`
	} `json:"data"`
	Status string `json:"status"`
	Error  Error  `json:"error,omitempty"`
}

type Group struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type RuleResponse struct {
	Data struct {
		Groups []Group `json:"groups"`
	} `json:"data"`
	Status string `json:"status"`
	Error  Error  `json:"error,omitempty"`
}

type Result struct {
	Metric json.RawMessage `json:"metric"`
	Values [][]any         `json:"values"`
}

type Data struct {
	Result     []Result `json:"result"`
	ResultType string   `json:"resultType"`
}

type DataResponse struct {
	Data   Data   `json:"data"`
	Status string `json:"status"`
}

type Threshold struct {
	Default int      `json:"default,omitempty"`
	Enabled bool     `json:"enabled"`
	Max     float64  `json:"max"`
	Min     float64  `json:"min"`
	Name    string   `json:"name"`
	Type    string   `json:"type,omitempty"`
	Unit    []string `json:"unit,omitempty"`
}

type Application struct {
	Enabled      bool        `json:"enabled"`
	MaxInstances int         `json:"maxInstances"`
	MinInstances int         `json:"minInstances"`
	Name         string      `json:"name"`
	Thresholds   []Threshold `json:"thresholds,omitempty"`
}

type AutoscalersResponse struct {
	Data struct {
		Applications []Application `json:"applications"`
	} `json:"data"`
	Status string `json:"status"`
	Error  Error  `json:"error,omitempty"`
}

type Rule struct {
	Annotations struct {
		Description string `json:"description"`
		Resolved    string `json:"resolved"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Metric      string `json:"metric"`
	Rule        struct {
		ExtraFor []struct {
			Name         string   `json:"name"`
			Options      []string `json:"options"`
			Type         string   `json:"type"`
			VariableName string   `json:"variableName"`
		} `json:"extraFor,omitempty"`
		Extras []struct {
			Name         string   `json:"name"`
			Options      []string `json:"options"`
			Type         string   `json:"type"`
			VariableName string   `json:"variableName"`
		} `json:"extras"`
		Operators []string  `json:"operators"`
		Subject   string    `json:"subject"`
		Threshold Threshold `json:"threshold"`
	} `json:"rule"`
	Template string `json:"template"`
}

// GetGroupedRules looks up available rules
func (c *MetricsService) GetGroupedRules(options ...OptionFunc) (*[]Group, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/rules", nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response RuleResponse

	resp, err := c.client.do(req, &response)
	if err != nil {
		if resp != nil {
			resp.Error = response.Error
		}
		return nil, resp, err
	}
	return &response.Data.Groups, resp, err
}

// GetRuleByID retrievs a rule by ID
func (c *MetricsService) GetRuleByID(id string, options ...OptionFunc) (*Rule, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/rules/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response struct {
		Data   Rule   `json:"data"`
		Status string `json:"status"`
		Error  Error  `json:"error,omitempty"`
	}

	resp, err := c.client.do(req, &response)
	if err != nil {
		if resp != nil {
			resp.Error = response.Error
		}
		return nil, resp, err
	}
	return &response.Data, resp, err
}

func (c *MetricsService) GQLGetInstances(ctx context.Context) (*[]Instance, error) {
	var query struct {
		Instances []Instance `graphql:"instances"`
	}
	err := c.client.gql.Query(ctx, &query, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return &query.Instances, nil
}

func (c *MetricsService) GQLGetInstanceByID(ctx context.Context, guid string) (*Instance, error) {
	var query struct {
		Instance Instance `graphql:"instance(guid: $guid)"`
	}
	err := c.client.gql.Query(ctx, &query, map[string]interface{}{
		"guid": graphql.String(guid),
	})
	if err != nil {
		return nil, err
	}
	return &query.Instance, nil
}

func (c *MetricsService) PrometheusGetData(ctx context.Context, host, query string, options ...OptionFunc) (*DataResponse, error) {
	var dataResponse DataResponse

	return &dataResponse, nil
}

// GetInstances looks up available instances
func (c *MetricsService) GetInstances(options ...OptionFunc) (*[]Instance, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/instances", nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response MetricsResponse

	resp, err := c.client.do(req, &response)
	if err != nil {
		if resp != nil {
			resp.Error = response.Error
		}
		return nil, resp, err
	}
	return &response.Data.Instances, resp, err
}

// GetInstanceByID looks up an instance by ID
func (c *MetricsService) GetInstanceByID(id string, options ...OptionFunc) (*Instance, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/instances/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response struct {
		Data   Instance `json:"data"`
		Status string   `json:"status"`
		Error  Error    `json:"error,omitempty"`
	}

	resp, err := c.client.do(req, &response)
	if err != nil {
		if resp != nil {
			resp.Error = response.Error
		}
		return nil, resp, err
	}
	return &response.Data, resp, err
}

// GetApplicationAutoscalers looks up all available autoscalers
func (c *MetricsService) GetApplicationAutoscalers(id string, options ...OptionFunc) (*[]Application, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/"+id+"/autoscalers", nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response AutoscalersResponse

	resp, err := c.client.do(req, &response)
	if err != nil {
		if resp != nil {
			resp.Error = response.Error
		}
		return nil, resp, err
	}
	return &response.Data.Applications, resp, err
}

// GetApplicationAutoscaler looks up a specific application autoscaler settings
func (c *MetricsService) GetApplicationAutoscaler(id, app string, options ...OptionFunc) (*Application, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "GET", "v3/metrics/"+id+"/autoscalers/"+app, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response struct {
		Data struct {
			Application Application `json:"application"`
		} `json:"data"`
		Status string `json:"status"`
		Error  Error  `json:"error,omitempty"`
	}

	resp, err := c.client.do(req, &response)
	if err != nil {
		if resp != nil {
			resp.Error = response.Error
		}
		return nil, resp, err
	}
	return &response.Data.Application, resp, err
}

// GetApplicationAutoscaler looks up a specific application autoscaler settings
func (c *MetricsService) UpdateApplicationAutoscaler(id string, settings Application, options ...OptionFunc) (*Application, *Response, error) {
	req, err := c.client.newRequest(CONSOLE, "PUT", "v3/metrics/"+id+"/autoscalers", &settings, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response struct {
		Data struct {
			Application Application `json:"application"`
		} `json:"data,omitempty"`
		Status string `json:"status,omitempty"`
		Error  Error  `json:"error,omitempty"`
	}

	resp, err := c.client.do(req, &response)
	if err != nil {
		if resp != nil {
			resp.Error = response.Error
		}
		return nil, resp, err
	}
	return &response.Data.Application, resp, err
}
