package mockdata

import "github.com/chandan/trading-dashboard/models"

// initial mock stocks
var Stocks = map[string]*models.Stock{
	"AAPL":  {Symbol: "AAPL", Price: 175.00, PreviousPrice: 175.00},
	"TSLA":  {Symbol: "TSLA", Price: 240.00, PreviousPrice: 240.00},
	"AMZN":  {Symbol: "AMZN", Price: 145.00, PreviousPrice: 145.00},
	"INFY":  {Symbol: "INFY", Price: 20.50, PreviousPrice: 20.50},
	"TCS":   {Symbol: "TCS", Price: 34.75, PreviousPrice: 34.75},
}
