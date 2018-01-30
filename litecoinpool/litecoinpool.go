package litecoinpool

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/json-iterator/go"
)

// Why isn't this part of the stdlib?!
// Oh right: https://github.com/json-iterator/go/issues/231
var json = jsoniter.ConfigDefault

func fetchData(url string, timeout time.Duration) (*[]byte, error) {
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "litecoinpool_exporter")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &body, err
}

func parseJson(body *[]byte) (*PoolData, error) {
	poolData := PoolData{}

	err := json.Unmarshal(*body, &poolData)
	if err != nil {
		return nil, err
	}

	return &poolData, nil
}

type LCPClient struct {
	urlPrefix string
	apiKey    string
	timeout   time.Duration
}

func (c *LCPClient) Fetch() (*PoolData, error) {
	body, err := fetchData(c.urlPrefix+c.apiKey, c.timeout)
	if err != nil {
		return nil, err
	}

	return parseJson(body)
}

func NewClient(apiKey string, timeout time.Duration) *LCPClient {
	return &LCPClient{
		urlPrefix: "https://www.litecoinpool.org/api?api_key=",
		apiKey:    apiKey,
		timeout:   timeout,
	}
}
