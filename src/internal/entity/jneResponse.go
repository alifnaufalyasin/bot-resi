package entity

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aprilahijriyan/gojne/lib"
	"github.com/gocolly/colly"
)

func GetResiJNEHistory(trackingNumber string) lib.DetailTracking {
	// Get tracking
	var token string
	var receivers = []string{}
	var Output = lib.DetailTracking{}
	var is_delivered = false

	c := colly.NewCollector()

	// Get csrf token
	c.OnHTML("#traceForm", func(e *colly.HTMLElement) {
		token = lib.ParseCsrfToken(e)
	})

	c.Visit("https://www.jne.co.id/id/beranda")
	c.Wait()

	var cc = colly.NewCollector()
	cc.OnHTML("html", func(h *colly.HTMLElement) {
		// Detail delivery
		delivery := lib.ParseDeliveryDetail(h)
		// fmt.Println("Delivery:", delivery)
		Output.Data.Delivery = delivery
		// Detail history
		history := lib.ParseDetailHistory(h)
		// fmt.Println("history:", history)
		receivers = append(receivers, history.Receiver.Name)
		Output.Data.History = history
		// Detail shipment
		shipment := lib.ParseDetailShipment(h)
		// fmt.Println("shipment:", shipment)
		receiverName := strings.TrimSpace(h.DOM.Find("div.col-md-12:nth-child(11) > h4:nth-child(2)").Text())
		receivers = append(receivers, receiverName)
		Output.Data.Shipment = shipment
		if len(history.Data) > 0 {
			re := regexp.MustCompile(`(?i)delivered to \[(?P<name>[a-zA-Z0-9 ]+)`)
			subexp := re.SubexpNames()
			target := history.Data[len(history.Data)-1].Title
			m := re.FindAllStringSubmatch(target, -1)[0]
			for i, n := range m {
				pattern := subexp[i]
				n = strings.TrimSpace(n)
				if pattern == "name" && lib.ArrayContains(receivers, n) {
					is_delivered = true
				}
			}
		}
		// Validasi status pengiriman
		var code string
		var message string

		if is_delivered {
			code = "060101"
			message = "Delivery has been received by the recipient"
		} else if len(delivery) > 0 {
			code = "060102"
			message = "Delivery is still in transit"
		} else {
			code = "060103"
			message = "Unable to retrieve delivery information"
		}
		Output.Status = lib.DetailStatus{
			Code:    code,
			Message: message,
		}
	})
	fmt.Println("Get tracking information...")
	var url = fmt.Sprintf("https://cekresi.jne.co.id/%s/", trackingNumber)
	var requestData = map[string]string{
		"_token":   token,
		"code":     trackingNumber,
		"tracking": "",
	}
	cc.Request("POST", url, lib.CreateFormReader(requestData), nil, http.Header{"Referer": []string{"https://www.jne.co.id/id/beranda"}})
	return Output
}
