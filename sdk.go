package gocollect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://gocollect.com"
)

// Client manages communication with the GoCollect API
type Client struct {
	client  *http.Client
	baseURL *url.URL
	token   string

	// Services
	Collectibles *CollectiblesService
	Insights     *InsightsService
	SoldExamples *SoldExamplesService
	StagedSales  *StagedSalesService
}

// ClientOption is a function that modifies the client
type ClientOption func(*Client) error

// NewClient creates a new GoCollect API client
func NewClient(token string, opts ...ClientOption) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:  http.DefaultClient,
		baseURL: baseURL,
		token:   token,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Initialize services
	c.Collectibles = &CollectiblesService{client: c}
	c.Insights = &InsightsService{client: c}
	c.SoldExamples = &SoldExamplesService{client: c}
	c.StagedSales = &StagedSalesService{client: c}

	return c, nil
}

// WithBaseURL sets a custom base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.baseURL = parsedURL
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		c.client = httpClient
		return nil
	}
}

// newRequest creates a new API request
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do sends an API request and returns the response
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// CollectiblesService handles communication with the collectible related endpoints
type CollectiblesService struct {
	client *Client
}

// SearchItemsOptions represents the parameters for searching items
type SearchItemsOptions struct {
	Query string
	CAM   string
	Limit int
}

// SearchItem represents a collectible item in search results
type SearchItem struct {
	ItemID             int     `json:"item_id"`
	UUID               string  `json:"uuid"`
	Slug               string  `json:"slug"`
	Name               string  `json:"name"`
	VariantOfItemID    *int    `json:"variant_of_item_id"`
	VariantDescription *string `json:"variant_description"`
}

// SearchItems searches for collectible items
func (s *CollectiblesService) SearchItems(opts SearchItemsOptions) ([]SearchItem, error) {
	params := url.Values{}
	params.Add("query", opts.Query)
	if opts.CAM != "" {
		params.Add("cam", opts.CAM)
	}
	if opts.Limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", opts.Limit))
	}

	path := fmt.Sprintf("/api/collectibles/v1/item/search?%s", params.Encode())
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var items []SearchItem
	_, err = s.client.do(req, &items)
	return items, err
}

// InsightsService handles communication with the insights related endpoints
type InsightsService struct {
	client *Client
}

// Metrics represents sales metrics for a specific time period
type Metrics struct {
	SoldCount    int     `json:"sold_count"`
	LowPrice     float64 `json:"low_price"`
	HighPrice    float64 `json:"high_price"`
	AveragePrice float64 `json:"average_price"`
}

// ItemInsights represents insights for a collectible item
type ItemInsights struct {
	ItemID      int                `json:"item_id"`
	Title       string             `json:"title"`
	IssueNumber string             `json:"issue_number"`
	CAM         string             `json:"cam"`
	Company     string             `json:"company"`
	Label       string             `json:"label"`
	Grade       string             `json:"grade"`
	Metrics     map[string]Metrics `json:"metrics"`
	FMV         *float64           `json:"fmv"`
}

// GetItemInsights retrieves insights for a specific item
func (s *InsightsService) GetItemInsights(itemID int, grade string, company string, label string) (*ItemInsights, error) {
	params := url.Values{}
	params.Add("grade", grade)
	if company != "" {
		params.Add("company", company)
	}
	if label != "" {
		params.Add("label", label)
	}

	path := fmt.Sprintf("/api/insights/v1/item/%d?%s", itemID, params.Encode())
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	insights := new(ItemInsights)
	_, err = s.client.do(req, insights)
	return insights, err
}

// GetItemInsightsByCGCID retrieves insights for a specific CGC item
func (s *InsightsService) GetItemInsightsByCGCID(cgcID string, grade string, company string, label string) (*ItemInsights, error) {
	params := url.Values{}
	params.Add("grade", grade)
	if company != "" {
		params.Add("company", company)
	}
	if label != "" {
		params.Add("label", label)
	}

	path := fmt.Sprintf("/api/insights/v1/item/cgc-id/%s?%s", cgcID, params.Encode())
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	insights := new(ItemInsights)
	_, err = s.client.do(req, insights)
	return insights, err
}

// Common types for both SoldExamples and StagedSales
type SaleFormat string

const (
	SaleFormatAuction    SaleFormat = "auction"
	SaleFormatFixedPrice SaleFormat = "fixed_price"
)

// SoldExamplesService handles communication with the sold examples related endpoints
type SoldExamplesService struct {
	client *Client
}

// SoldExample represents a sold collectible
type SoldExample struct {
	PartnerSaleID        string     `json:"partner_sale_id"`
	CAM                  string     `json:"cam"`
	Title                string     `json:"title"`
	ImageURLs            []string   `json:"image_urls"`
	GocollectItemID      *int       `json:"gocollect_item_id"`
	CertificationCompany string     `json:"certification_company"`
	CertificationKey     *string    `json:"certification_key"`
	ListedPrice          *float64   `json:"listed_price"`
	ListedAt             time.Time  `json:"listed_at"`
	SoldPrice            float64    `json:"sold_price"`
	SoldAt               time.Time  `json:"sold_at"`
	URL                  string     `json:"url"`
	Format               SaleFormat `json:"format"`
	AuctionName          *string    `json:"auction_name"`
	BidCount             *int       `json:"bid_count"`
}

// CreateSoldExample creates a new sold example
func (s *SoldExamplesService) CreateSoldExample(example *SoldExample) error {
	req, err := s.client.newRequest("POST", "/api/resources/v1/sold-examples", example)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}

// GetSoldExample retrieves a specific sold example
func (s *SoldExamplesService) GetSoldExample(partnerSaleID string) (*SoldExample, error) {
	path := fmt.Sprintf("/api/resources/v1/sold-examples/%s", partnerSaleID)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data SoldExample `json:"data"`
	}
	_, err = s.client.do(req, &response)
	return &response.Data, err
}

// StagedSalesService handles communication with the staged sales related endpoints
type StagedSalesService struct {
	client *Client
}

// StagedSale represents a staged sale
type StagedSale struct {
	PartnerSaleID        string     `json:"partner_sale_id"`
	CAM                  string     `json:"cam"`
	Title                string     `json:"title"`
	IsActive             bool       `json:"is_active"`
	ImageURLs            []string   `json:"image_urls"`
	GocollectItemID      *int       `json:"gocollect_item_id"`
	IsGraded             bool       `json:"is_graded"`
	CertificationCompany string     `json:"certification_company"`
	CertificationKey     *string    `json:"certification_key"`
	ListedPrice          *float64   `json:"listed_price"`
	Price                *float64   `json:"price"`
	SoldAt               time.Time  `json:"sold_at"`
	URL                  string     `json:"url"`
	Format               SaleFormat `json:"format"`
	AuctionName          *string    `json:"auction_name"`
	EndsAt               *time.Time `json:"ends_at"`
}

// CreateStagedSale creates a new staged sale
func (s *StagedSalesService) CreateStagedSale(sale *StagedSale) error {
	req, err := s.client.newRequest("POST", "/api/resources/v1/staged-sales", sale)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}

// GetStagedSale retrieves a specific staged sale
func (s *StagedSalesService) GetStagedSale(id string) (*StagedSale, error) {
	path := fmt.Sprintf("/api/resources/v1/staged-sales/%s", id)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data StagedSale `json:"data"`
	}
	_, err = s.client.do(req, &response)
	return &response.Data, err
}
