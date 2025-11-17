package services

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/chandan/trading-dashboard/mockdata"
	"github.com/chandan/trading-dashboard/models"
	"github.com/chandan/trading-dashboard/websocket"
)

func init() {
	// seed randomness
	rand.Seed(time.Now().UnixNano())
	// start websocket hub
	go websocket.RunHub()
}

// StartPriceSimulator runs indefinitely; interval controls how often prices change
func StartPriceSimulator(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		updates := make([]models.PriceUpdate, 0, len(mockdata.Stocks))

		for _, s := range mockdata.Stocks {
			// save previous price
			// prev := s.Price
			// percent change between -2.0% and +2.0% but avoid 0; typical 0.5-2%
			sign := 1.0
			if rand.Intn(2) == 0 {
				sign = -1.0
			}
			pct := (0.5 + rand.Float64()*1.5) * sign // between -2.0 and -0.5 or +0.5 and +2.0
			factor := 1 + pct/100.0
			newPrice := s.Price * factor
			// round to 2 decimals
			newPrice = math.Round(newPrice*100) / 100
			if newPrice <= 0 {
				newPrice = s.Price // safeguard
			}
			s.PreviousPrice = s.Price
			s.Price = newPrice
			s.ChangePct = math.Round(((s.Price-s.PreviousPrice)/s.PreviousPrice)*100*100) / 100 // 2 decimals

			updates = append(updates, models.PriceUpdate{
				Symbol: s.Symbol,
				Price:  s.Price,
			})
		}

		// broadcast updates as JSON array
		b, err := json.Marshal(updates)
		if err != nil {
			log.Println("failed to marshal updates:", err)
			continue
		}
		websocket.BroadcastPriceUpdate(b)
	}
}
