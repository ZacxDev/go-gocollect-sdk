# go-gocollect-sdk

A Go SDK for the GoCollect API. This SDK provides a convenient way to interact with the GoCollect API for managing collectibles, insights, sold examples, and staged sales.

## Installation

```bash
go get github.com/ZacxDev/go-gocollect-sdk
```

## Usage

### Creating a Client

```go
import "github.com/ZacxDev/go-gocollect-sdk"

// Create a new client with your API token
client, err := gocollect.NewClient("your-api-token")
if err != nil {
    log.Fatal(err)
}

// Optionally, customize the client
client, err = gocollect.NewClient(
    "your-api-token",
    gocollect.WithBaseURL("https://gocollect.com"),
    gocollect.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
)
```

### Searching for Collectibles

```go
// Search for items
items, err := client.Collectibles.SearchItems(gocollect.SearchItemsOptions{
    Query: "Incredible Hulk #181",
    CAM:   "Comics",
    Limit: 10,
})
if err != nil {
    log.Fatal(err)
}

for _, item := range items {
    fmt.Printf("Item: %s (ID: %d)\n", item.Name, item.ItemID)
}
```

### Getting Item Insights

```go
// Get insights for an item
insights, err := client.Insights.GetItemInsights(
    223124,          // item ID
    "9.8",          // grade
    "CGC",          // company
    "Universal",     // label
)
if err != nil {
    log.Fatal(err)
}

// Access metrics
fmt.Printf("FMV: $%.2f\n", *insights.FMV)
fmt.Printf("30-day sales count: %d\n", insights.Metrics["30"].SoldCount)
```

### Managing Sold Examples

```go
// Create a sold example
soldExample := &gocollect.SoldExample{
    PartnerSaleID:       "12345",
    CAM:                 "Comics",
    Title:               "Amazing Spider-Man #300",
    CertificationCompany: "CGC",
    SoldPrice:           999.99,
    SoldAt:              time.Now(),
    URL:                 "https://example.com/sale/12345",
    Format:              gocollect.SaleFormatFixedPrice,
}

err := client.SoldExamples.CreateSoldExample(soldExample)
if err != nil {
    log.Fatal(err)
}

// Get a sold example
example, err := client.SoldExamples.GetSoldExample("12345")
if err != nil {
    log.Fatal(err)
}
```

### Working with Staged Sales

```go
// Create a staged sale
stagedSale := &gocollect.StagedSale{
    PartnerSaleID:       "67890",
    CAM:                 "Comics",
    Title:               "X-Men #1",
    IsActive:            true,
    IsGraded:            true,
    CertificationCompany: "CGC",
    SoldAt:              time.Now(),
    URL:                 "https://example.com/sale/67890",
    Format:              gocollect.SaleFormatAuction,
}

err := client.StagedSales.CreateStagedSale(stagedSale)
if err != nil {
    log.Fatal(err)
}

// Get a staged sale
sale, err := client.StagedSales.GetStagedSale("67890")
if err != nil {
    log.Fatal(err)
}
```

## API Documentation

### Services

The SDK provides four main services:

1. **CollectiblesService**
   - `SearchItems(opts SearchItemsOptions) ([]SearchItem, error)`

2. **InsightsService**
   - `GetItemInsights(itemID int, grade string, company string, label string) (*ItemInsights, error)`
   - `GetItemInsightsByCGCID(cgcID string, grade string, company string, label string) (*ItemInsights, error)`

3. **SoldExamplesService**
   - `CreateSoldExample(example *SoldExample) error`
   - `GetSoldExample(partnerSaleID string) (*SoldExample, error)`

4. **StagedSalesService**
   - `CreateStagedSale(sale *StagedSale) error`
   - `GetStagedSale(id string) (*StagedSale, error)`

### Rate Limits

The GoCollect API has the following rate limits:

- Collectibles API: 100 requests per day for subscribers, 50 for non-subscribers
- Insights API: 100 requests per day for subscribers, 50 for non-subscribers
- Sold Examples API: 500 requests per hour
- Staged Sales API: 500 requests per hour

When rate limits are exceeded, the API will return a 429 status code.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This SDK is distributed under the MIT license. See the [LICENSE](LICENSE) file for more information.
