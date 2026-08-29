// Package model holds the shared data types passed between the screener
// scraper, the analysis engine, and the HTTP API.
package model // this file belongs to the "model" package, so other files import it as "model.SomeType"

// SearchResult is one match returned by screener.in's company search.
type SearchResult struct { // defines a new struct (a labeled box of fields) called SearchResult
	ID   int    `json:"id"`   // a whole number id, shown as "id" when converted to JSON
	Name string `json:"name"` // the company's full name, shown as "name" in JSON
	URL  string `json:"url"`  // the path to that company's page on screener.in, shown as "url" in JSON
} // end of the SearchResult struct

// DataRow is one labeled row of a financial statement table, e.g. "Sales"
// with one value per reporting period, plus the period-over-period
// percentage change for that same row (ChangePercents[i] is the change
// from period i-1 to period i; entry 0 and any entry that couldn't be
// computed - e.g. a missing/non-numeric value - is nil).
type DataRow struct { // one row inside a bigger table, e.g. the "Sales" row
	Label          string     `json:"label"`                     // the row's name, e.g. "Sales"
	Values         []string   `json:"values"`                    // a list of text values, one per column/period
	ChangePercents []*float64 `json:"change_percents,omitempty"` // matching list of % changes; omitted from JSON if empty
} // end of the DataRow struct

// FinancialTable is a raw statement table (Quarterly Results, Profit &
// Loss, Balance Sheet, Cash Flow, or Shareholding Pattern) exactly as
// reported - no derived ratios, just the line items and their values per
// period, kept as the display strings screener.in renders.
type FinancialTable struct { // a whole table, e.g. "Profit & Loss"
	Title      string    `json:"title"`       // the table's display name
	Periods    []string  `json:"periods"`     // column headers shown to the user, e.g. "Mar 2024"
	PeriodKeys []string  `json:"period_keys"` // hidden machine-readable version of each column, e.g. "2024-03-31" or "TTM"
	Rows       []DataRow `json:"rows"`        // every row in this table, in display order
} // end of the FinancialTable struct

// MarketFacts are the few bare market data points that cannot be derived
// from the financial statements themselves (they come from the stock
// market, not from company filings): the current share price, its
// 52-week range, and the face (nominal) value of one share.
type MarketFacts struct { // raw stock-market facts, not calculated by us
	CurrentPrice float64 `json:"current_price"` // the share's last traded price, in rupees
	High52W      float64 `json:"high_52w"`      // highest price traded in the last 52 weeks
	Low52W       float64 `json:"low_52w"`       // lowest price traded in the last 52 weeks
	FaceValue    float64 `json:"face_value"`    // the nominal value of one share, e.g. ₹1, ₹2, ₹10
} // end of the MarketFacts struct

// SeriesPoint is one (period, numeric value) sample in a computed trend.
type SeriesPoint struct { // one dot on a graph: which period, and what number
	Period string  `json:"period"` // the period label, e.g. "Mar 2024"
	Value  float64 `json:"value"`  // the numeric value for that period
} // end of the SeriesPoint struct

// GrowthMetric is a trend for one financial line item across the annual
// (non-TTM) periods available, computed entirely from the raw statement
// numbers - not sourced from any pre-computed figure on screener.in.
type GrowthMetric struct { // a full growth trend for one thing, e.g. "Sales"
	Label   string        `json:"label"`              // what this trend is of, e.g. "Sales"
	Unit    string        `json:"unit"`               // the unit of the values, e.g. "Rs. Cr."
	Series  []SeriesPoint `json:"series"`             // the year-by-year numbers making up this trend
	Found   bool          `json:"found"`              // true if we actually found this row in the source data
	CAGR3Y  *float64      `json:"cagr_3y,omitempty"`  // compound annual growth rate over the last 3 years, if computable
	CAGR5Y  *float64      `json:"cagr_5y,omitempty"`  // compound annual growth rate over the last 5 years, if computable
	CAGR10Y *float64      `json:"cagr_10y,omitempty"` // compound annual growth rate over the last 10 years, if computable
	// LatestYoYPct is the most recent year-over-year percentage change,
	// used for flow metrics (Sales, Net Profit, EPS).
	LatestYoYPct *float64 `json:"latest_yoy_pct,omitempty"` // % change from last year to this year, for money-amount metrics
	// LatestYoYChange is the most recent year-over-year change in
	// percentage points, used for margin metrics (OPM %).
	LatestYoYChange *float64 `json:"latest_yoy_change,omitempty"` // point change from last year to this year, for margin/% metrics
} // end of the GrowthMetric struct

// Analysis is this program's own fundamental-analysis layer: every ratio
// and trend here is computed in Go from the raw statement line items and
// the bare market facts - none of it is copied from screener.in's own
// ratio boxes.
type Analysis struct { // everything we personally calculated for one company
	SharesOutstandingCr   *float64 `json:"shares_outstanding_cr,omitempty"`     // total shares in circulation, in crores
	MarketCapCr           *float64 `json:"market_cap_cr,omitempty"`             // price x shares outstanding, in crores of rupees
	EPS                   *float64 `json:"eps,omitempty"`                       // latest annual earnings per share, in rupees
	PE                    *float64 `json:"pe,omitempty"`                        // price divided by EPS
	BookValuePerShare     *float64 `json:"book_value_per_share,omitempty"`      // net worth divided by shares outstanding
	ROEPercent            *float64 `json:"roe_percent,omitempty"`               // net profit as a % of net worth
	ROCEPercent           *float64 `json:"roce_percent,omitempty"`              // EBIT as a % of capital employed
	DebtToEquity          *float64 `json:"debt_to_equity,omitempty"`            // borrowings divided by net worth
	OPMPercent            *float64 `json:"opm_percent,omitempty"`               // latest operating profit as a % of sales
	NetMarginPercent      *float64 `json:"net_margin_percent,omitempty"`        // net profit as a % of sales
	DepreciationPercent   *float64 `json:"depreciation_percent,omitempty"`      // depreciation as a % of sales
	DividendYieldPercent  *float64 `json:"dividend_yield_percent,omitempty"`    // dividend per share as a % of current price
	OCFToSalesPercent     *float64 `json:"ocf_to_sales_percent,omitempty"`      // operating cash flow as a % of sales
	OCFToNetProfitPercent *float64 `json:"ocf_to_net_profit_percent,omitempty"` // operating cash flow as a % of net profit

	SalesGrowth            GrowthMetric `json:"sales_growth"`             // the full Sales (or Revenue) trend
	ProfitGrowth           GrowthMetric `json:"profit_growth"`            // the full Net Profit trend
	EPSGrowth              GrowthMetric `json:"eps_growth"`               // the full EPS trend
	OPMTrend               GrowthMetric `json:"opm_trend"`                // the full Operating Margin trend
	SharesOutstandingTrend GrowthMetric `json:"shares_outstanding_trend"` // the full Shares Outstanding trend (dilution/buybacks)

	Notes []string `json:"notes"` // short plain-English observations generated from the numbers above
} // end of the Analysis struct

// CompanyData is the full payload returned by the /api/analyze endpoint.
type CompanyData struct { // everything sent back to the browser for one company
	Name         string         `json:"name"`          // the company's display name
	About        string         `json:"about"`         // a short description of what the company does
	SourceURL    string         `json:"source_url"`    // the screener.in page this data came from
	Market       MarketFacts    `json:"market"`        // raw market facts (price, face value, etc.)
	Quarterly    FinancialTable `json:"quarterly"`     // the Quarterly Results table
	ProfitLoss   FinancialTable `json:"profit_loss"`   // the Profit & Loss table
	BalanceSheet FinancialTable `json:"balance_sheet"` // the Balance Sheet table
	CashFlow     FinancialTable `json:"cash_flow"`     // the Cash Flow table
	Shareholding FinancialTable `json:"shareholding"`  // the Shareholding Pattern table
	Analysis     Analysis       `json:"analysis"`      // everything we calculated ourselves
} // end of the CompanyData struct
