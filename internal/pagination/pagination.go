package pagination

import (
	"encoding/json"
	"fmt"
	"strconv"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
)

func hasNextPage(pageInfo map[string]any) bool {
	v, ok := pageInfo["has_next_page"]
	if !ok {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true"
	default:
		return false
	}
}

func PaginateCRM(client *zohttp.Client, url string, params map[string]string, maxPages int) ([]json.RawMessage, error) {
	if maxPages == 0 {
		maxPages = 20
	}
	var all []json.RawMessage
	p := make(map[string]string)
	for k, v := range params {
		p[k] = v
	}
	if _, ok := p["per_page"]; !ok {
		p["per_page"] = "200"
	}
	p["page"] = "1"

	for range maxPages {
		raw, err := client.Request("GET", url, &zohttp.RequestOpts{Params: p})
		if err != nil {
			return all, err
		}

		var envelope struct {
			Data []json.RawMessage `json:"data"`
			Info struct {
				MoreRecords   bool   `json:"more_records"`
				NextPageToken string `json:"next_page_token"`
			} `json:"info"`
		}
		if err := json.Unmarshal(raw, &envelope); err != nil {
			return all, nil
		}
		all = append(all, envelope.Data...)

		if !envelope.Info.MoreRecords || len(envelope.Data) == 0 {
			break
		}
		if envelope.Info.NextPageToken != "" {
			p["page_token"] = envelope.Info.NextPageToken
			delete(p, "page")
		} else {
			page, _ := strconv.Atoi(p["page"])
			p["page"] = strconv.Itoa(page + 1)
		}
	}
	return all, nil
}

func PaginateProjects(client *zohttp.Client, url string, itemsKey string, params map[string]string, maxPages int) ([]json.RawMessage, error) {
	if maxPages == 0 {
		maxPages = 20
	}
	var all []json.RawMessage
	p := make(map[string]string)
	for k, v := range params {
		p[k] = v
	}
	page := 1

	for range maxPages {
		p["page"] = strconv.Itoa(page)
		p["per_page"] = "100"

		raw, err := client.Request("GET", url, &zohttp.RequestOpts{Params: p})
		if err != nil {
			return all, err
		}

		var rawList []json.RawMessage
		if err := json.Unmarshal(raw, &rawList); err == nil {
			all = append(all, rawList...)
			break
		}

		var envelope map[string]json.RawMessage
		if err := json.Unmarshal(raw, &envelope); err != nil {
			break
		}

		if itemsKey != "" {
			if items, ok := envelope[itemsKey]; ok {
				var list []json.RawMessage
				if err := json.Unmarshal(items, &list); err != nil {
					break
				}
				all = append(all, list...)
				if len(list) == 0 {
					break
				}
			}
		}

		var pageInfo map[string]any
		if pi, ok := envelope["page_info"]; ok {
			json.Unmarshal(pi, &pageInfo)
		}
		if !hasNextPage(pageInfo) {
			break
		}
		page++
	}
	return all, nil
}

func PaginateWorkDrive(client *zohttp.Client, url string, params map[string]string, maxPages int) ([]json.RawMessage, error) {
	if maxPages == 0 {
		maxPages = 10
	}
	perPage := 50
	var all []json.RawMessage
	p := make(map[string]string)
	for k, v := range params {
		p[k] = v
	}
	if _, ok := p["page[limit]"]; !ok {
		p["page[limit]"] = fmt.Sprintf("%d", perPage)
	}

	for range maxPages {
		p["page[offset]"] = fmt.Sprintf("%d", len(all))

		raw, err := client.Request("GET", url, &zohttp.RequestOpts{Params: p})
		if err != nil {
			return all, err
		}

		var envelope struct {
			Data []json.RawMessage `json:"data"`
			Meta struct {
				HasNext bool `json:"has_next"`
			} `json:"meta"`
		}
		if err := json.Unmarshal(raw, &envelope); err != nil {
			break
		}
		all = append(all, envelope.Data...)

		if !envelope.Meta.HasNext || len(envelope.Data) == 0 {
			break
		}
	}
	return all, nil
}
