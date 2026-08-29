package screener

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"sampleproject/internal/model"
)

// ParseCompanyPage turns a raw screener.in company page into structured
// data. It deliberately extracts only raw, as-reported figures: the
// financial statement line items and the bare market facts (price, face
// value) that cannot be derived from those statements. It does not use
// any of screener.in's own pre-computed ratios or growth figures - those
// are computed independently in the analysis package.
func ParseCompanyPage(html []byte, sourceURL string) (*model.CompanyData, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	data := &model.CompanyData{
		SourceURL: sourceURL,
		Name:      cleanText(doc.Find("h1").First().Text()),
		About:     cleanText(doc.Find(".company-profile .show-more-box.about").First().Text()),
		Market:    parseMarketFacts(doc),

		Quarterly:    parseFinancialTable(doc, "#quarters", "Quarterly Results"),
		ProfitLoss:   parseFinancialTable(doc, "#profit-loss", "Profit & Loss"),
		BalanceSheet: parseFinancialTable(doc, "#balance-sheet", "Balance Sheet"),
		CashFlow:     parseFinancialTable(doc, "#cash-flow", "Cash Flow"),
		Shareholding: parseFinancialTable(doc, "#shareholding #quarterly-shp", "Shareholding Pattern (Quarterly)"),
	}

	return data, nil
}

// parseMarketFacts reads only the handful of top-ratios entries that are
// raw market data rather than a computed ratio: current price, 52-week
// high/low, and face value. Screener's own P/E, ROE, ROCE, Book Value,
// Market Cap and Dividend Yield entries in that same box are intentionally
// skipped - this program computes its own versions of those.
func parseMarketFacts(doc *goquery.Document) model.MarketFacts {
	var facts model.MarketFacts
	doc.Find("#top-ratios > li").Each(func(_ int, li *goquery.Selection) {
		name := cleanText(li.Find(".name").First().Text())
		valueText := cleanText(li.Find(".value").First().Text())
		switch name {
		case "Current Price":
			facts.CurrentPrice = firstNumber(valueText)
		case "High / Low":
			nums := allNumbers(valueText)
			if len(nums) >= 2 {
				facts.High52W, facts.Low52W = nums[0], nums[1]
			}
		case "Face Value":
			facts.FaceValue = firstNumber(valueText)
		}
	})
	return facts
}

// parseFinancialTable reads a screener "data-table": a header row of
// period labels followed by one row per line item.
func parseFinancialTable(doc *goquery.Document, sectionSelector, title string) model.FinancialTable {
	table := model.FinancialTable{Title: title}

	sel := doc.Find(sectionSelector).Find("table.data-table").First()
	if sel.Length() == 0 {
		return table
	}

	sel.Find("thead th").Each(func(i int, th *goquery.Selection) {
		if i == 0 {
			return // first header cell is the blank label column
		}
		table.Periods = append(table.Periods, cleanText(th.Text()))
		key, _ := th.Attr("data-date-key")
		table.PeriodKeys = append(table.PeriodKeys, key)
	})

	sel.Find("tbody > tr").Each(func(_ int, tr *goquery.Selection) {
		tds := tr.Find("td")
		if tds.Length() == 0 {
			return
		}
		row := model.DataRow{Label: cleanRowLabel(tds.Eq(0).Text())}
		tds.Each(func(i int, td *goquery.Selection) {
			if i == 0 {
				return
			}
			row.Values = append(row.Values, cleanText(td.Text()))
		})
		if row.Label != "" {
			table.Rows = append(table.Rows, row)
		}
	})

	return table
}

// cleanText collapses whitespace/newlines from goquery's raw text
// extraction into a single readable string.
func cleanText(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// cleanRowLabel additionally strips the "+" expand affordance screener
// renders next to row labels that have a drill-down schedule (e.g. "Sales +").
func cleanRowLabel(s string) string {
	s = cleanText(s)
	s = strings.TrimSuffix(s, "+")
	return strings.TrimSpace(s)
}

// firstNumber extracts the first numeric token from a formatted value
// string like "₹ 2,302" -> 2302.
func allNumbers(s string) []float64 {
	var out []float64
	var cur strings.Builder
	flush := func() {
		if cur.Len() == 0 {
			return
		}
		if v, err := strconv.ParseFloat(cur.String(), 64); err == nil {
			out = append(out, v)
		}
		cur.Reset()
	}
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9', r == '.':
			cur.WriteRune(r)
		case r == ',':
			// thousands separator, skip
		default:
			flush()
		}
	}
	flush()
	return out
}

func firstNumber(s string) float64 {
	nums := allNumbers(s)
	if len(nums) == 0 {
		return 0
	}
	return nums[0]
}
