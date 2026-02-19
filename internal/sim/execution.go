package sim

import (
	"log"
	"tbankbot/internal/models"
)

func ExecuteOrder(o *models.Order, accountID, token, baseURL string) {
	switch o.Side {

	case models.Buy:
		postErr := models.PostOrder(accountID, token, baseURL, "BBG004730N88", models.Buy)
		if postErr != nil {
			log.Println(postErr)
		}

	case models.Sell:
		postErr := models.PostOrder(accountID, token, baseURL, "BBG004730N88", models.Buy)
		if postErr != nil {
			log.Println(postErr)
		}
	}
}
