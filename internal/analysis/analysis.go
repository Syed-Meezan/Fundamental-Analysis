// Package analysis performs fundamental analysis entirely from raw,
// as-reported financial statement figures (Profit & Loss, Balance Sheet)
// plus the bare market facts (share price, face value). It does not read
// or trust any pre-computed ratio from screener.in - every number here
// (shares outstanding, market cap, EPS, P/E, book value, ROE, ROCE,
// debt-to-equity, operating margin, dividend yield, and all growth/CAGR
// trends) is derived in Go using standard formulas.
package analysis

import (
	"math"
	"strconv"
	"strings"

	"sampleproject/internal/model"
)

// AnnotateRowChanges fills in ChangePercents for every row of a raw
// statement table (Quarterly Results, P&L, Balance Sheet, Cash Flow,
// Shareholding), computing the period-over-period percentage change
// directly from that row's own values. This is applied uniformly to every
// row across every section, not just the headline metrics Compute focuses
// on, so e.g. "how much did Sales go up vs. the previous year" is visible
// straight in the table rather than only for a curated few line items.
func AnnotateRowChanges(table *model.FinancialTable) {
	for i := range table.Rows {
		row := &table.Rows[i]
		row.ChangePercents = make([]*float64, len(row.Values))

		var prev float64
		havePrev := false
		for j, raw := range row.Values {
			v, ok := parseNumber(raw)
			if ok && havePrev {
				if pct, ok2 := pctChange(prev, v); ok2 {
					pctCopy := pct
					row.ChangePercents[j] = &pctCopy
				}
			}
			// A gap (missing/non-numeric period) breaks the chain rather
			// than comparing across it, so a change is only ever shown
			// between two genuinely adjacent reported periods.
			prev, havePrev = v, ok
		}
	}
}

// Compute fills data.Analysis from data.ProfitLoss, data.BalanceSheet and
// data.Market. It never fails outright: a missing or unparseable input
// (not every company reports every line item) just leaves the
// corresponding output as nil/not-found rather than erroring the request.
func Compute(data *model.CompanyData) {
	pl := data.ProfitLoss
	bs := data.BalanceSheet
	mkt := data.Market

	// Banks, NBFCs and insurers report on a different statement template:
	// "Revenue"/"Interest"/"Financing Profit" instead of
	// "Sales"/"Expenses"/"Operating Profit", Interest is their core cost of
	// funds (not a financing charge to add back), and their real funding
	// source is Deposits rather than Borrowings. The Operating
	// Profit/ROCE/Debt-to-Equity formulas below assume the non-financial
	// template and would be misleading here, so they are skipped for such
	// companies (detected by the presence of a "Financing Profit" row).
	isFinancial := hasRow(pl, "Financing Profit")

	topLineLabel := "Sales"
	topLine := annualSeries(pl, topLineLabel)
	if len(topLine) == 0 {
		topLineLabel = "Revenue"
		topLine = annualSeries(pl, topLineLabel)
	}

	expenses := annualSeries(pl, "Expenses")
	netProfit := annualSeries(pl, "Net Profit")
	pbt := annualSeries(pl, "Profit before tax")
	interest := annualSeries(pl, "Interest")
	depreciation := annualSeries(pl, "Depreciation")
	dividendPayoutPct := annualSeries(pl, "Dividend Payout %")
	operatingCashFlow := annualSeries(data.CashFlow, "Cash from Operating Activity")

	equityCapital := annualSeries(bs, "Equity Capital")
	reserves := annualSeries(bs, "Reserves")
	borrowings := annualSeries(bs, "Borrowings")

	netMargin := combine(netProfit, topLine, func(np, s float64) (float64, bool) {
		if s == 0 {
			return 0, false
		}
		return np / s * 100, true
	})
	depreciationPct := combine(depreciation, topLine, func(d, s float64) (float64, bool) {
		if s == 0 {
			return 0, false
		}
		return d / s * 100, true
	})
	ocfToSales := combine(operatingCashFlow, topLine, func(ocf, s float64) (float64, bool) {
		if s == 0 {
			return 0, false
		}
		return ocf / s * 100, true
	})
	ocfToNetProfit := combine(operatingCashFlow, netProfit, func(ocf, np float64) (float64, bool) {
		if np == 0 {
			return 0, false
		}
		return ocf / np * 100, true
	})

	// Operating Profit = Sales - Expenses, computed here rather than taken
	// from screener's own "Operating Profit" row. Not meaningful under the
	// financial-company template (see isFinancial above), so left empty
	// there.
	var opm []model.SeriesPoint
	if !isFinancial {
		operatingProfit := combine(topLine, expenses, func(s, e float64) (float64, bool) { return s - e, true })
		opm = combine(operatingProfit, topLine, func(op, s float64) (float64, bool) {
			if s == 0 {
				return 0, false
			}
			return op / s * 100, true
		})
	}

	netWorth := combine(equityCapital, reserves, func(e, r float64) (float64, bool) { return e + r, true })

	// Shares outstanding per period = Equity Capital (Rs Cr) / Face Value,
	// scaled to a share count. This assumes face value is constant across
	// the reported history, which holds unless the company has done a
	// stock split or consolidation.
	var shares []model.SeriesPoint
	if mkt.FaceValue > 0 {
		for _, p := range equityCapital {
			shares = append(shares, model.SeriesPoint{Period: p.Period, Value: p.Value * 1e7 / mkt.FaceValue})
		}
	}

	eps := combine(netProfit, shares, func(np, sh float64) (float64, bool) {
		if sh <= 0 {
			return 0, false
		}
		return np * 1e7 / sh, true
	})

	var sharesCr []model.SeriesPoint
	for _, p := range shares {
		sharesCr = append(sharesCr, model.SeriesPoint{Period: p.Period, Value: p.Value / 1e7})
	}

	a := model.Analysis{
		SalesGrowth:            growthMetric(topLineLabel, "Rs. Cr.", topLine, true),
		ProfitGrowth:           growthMetric("Net Profit", "Rs. Cr.", netProfit, true),
		EPSGrowth:              growthMetric("EPS", "Rs.", eps, true),
		OPMTrend:               growthMetric("Operating Margin", "%", opm, false),
		SharesOutstandingTrend: growthMetric("Shares Outstanding", "Cr.", sharesCr, true),
	}

	if v, ok := last(netMargin); ok {
		a.NetMarginPercent = &v
	}
	if v, ok := last(depreciationPct); ok {
		a.DepreciationPercent = &v
	}
	if v, ok := last(ocfToSales); ok {
		a.OCFToSalesPercent = &v
	}
	if v, ok := last(ocfToNetProfit); ok {
		a.OCFToNetProfitPercent = &v
	}

	if v, ok := last(shares); ok && v > 0 {
		sharesCr := v / 1e7
		a.SharesOutstandingCr = &sharesCr
		if mkt.CurrentPrice > 0 {
			mc := mkt.CurrentPrice * v / 1e7
			a.MarketCapCr = &mc
		}
	}

	if v, ok := last(eps); ok {
		a.EPS = &v
		if mkt.CurrentPrice > 0 && v > 0 {
			pe := mkt.CurrentPrice / v
			a.PE = &pe
		}
		if payout, ok2 := last(dividendPayoutPct); ok2 && mkt.CurrentPrice > 0 {
			dps := payout / 100 * v
			dy := dps / mkt.CurrentPrice * 100
			a.DividendYieldPercent = &dy
		}
	}

	latestNetWorth, hasNetWorth := last(netWorth)
	if hasNetWorth {
		if v, ok := last(shares); ok && v > 0 {
			bv := latestNetWorth * 1e7 / v
			a.BookValuePerShare = &bv
		}
		if v, ok := last(netProfit); ok && latestNetWorth > 0 {
			roe := v / latestNetWorth * 100
			a.ROEPercent = &roe
		}
	}

	if !isFinancial {
		if latestPBT, ok := last(pbt); ok && hasNetWorth {
			latestInterest, _ := last(interest) // absent/zero for debt-free companies
			ebit := latestPBT + latestInterest
			latestBorrowings, _ := last(borrowings)
			// Capital employed approximated as Net Worth + Borrowings, since
			// screener's summarized balance sheet does not break out current
			// vs. non-current liabilities separately.
			capitalEmployed := latestNetWorth + latestBorrowings
			if capitalEmployed > 0 {
				roce := ebit / capitalEmployed * 100
				a.ROCEPercent = &roce
			}
		}

		if latestBorrowings, ok := last(borrowings); ok && hasNetWorth && latestNetWorth > 0 {
			de := latestBorrowings / latestNetWorth
			a.DebtToEquity = &de
		}
	}

	if v, ok := last(opm); ok {
		a.OPMPercent = &v
	}

	a.Notes = buildNotes(a)
	if isFinancial {
		a.Notes = append(a.Notes, "This looks like a bank/NBFC/insurer: Operating Margin, ROCE and Debt-to-Equity are omitted because their standard non-financial formulas (which treat Interest as a financing charge and Borrowings as the funding base) don't fit a business funded mainly by deposits, where interest is the core cost.")
	}
	data.Analysis = a
}

// hasRow reports whether a row with this label exists in the table, used
// to detect which statement template (standard vs. bank/NBFC) a company
// was reported under.
func hasRow(table model.FinancialTable, label string) bool {
	_, ok := findRow(table, label)
	return ok
}

// growthMetric computes YoY change and 3/5/10-year CAGR over an annual
// series. isFlow distinguishes flow items (Sales, Net Profit, EPS - where
// a percentage CAGR is meaningful) from margin items like OPM %  (where
// only a point-change YoY is meaningful).
func growthMetric(label, unit string, series []model.SeriesPoint, isFlow bool) model.GrowthMetric {
	metric := model.GrowthMetric{Label: label, Unit: unit, Series: series}
	n := len(series)
	if n == 0 {
		return metric
	}
	metric.Found = true
	last := series[n-1].Value

	if n >= 2 {
		prev := series[n-2].Value
		if isFlow {
			if pct, ok := pctChange(prev, last); ok {
				metric.LatestYoYPct = &pct
			}
		} else {
			change := last - prev
			metric.LatestYoYChange = &change
		}
	}

	if isFlow {
		metric.CAGR3Y = cagrOver(series, 3)
		metric.CAGR5Y = cagrOver(series, 5)
		metric.CAGR10Y = cagrOver(series, 10)
	}

	return metric
}

// cagrOver computes the compound annual growth rate over the trailing
// `years` of the series, or nil if there isn't enough history or the
// starting value isn't positive (CAGR is undefined for a non-positive base).
func cagrOver(series []model.SeriesPoint, years int) *float64 {
	n := len(series)
	if n < years+1 {
		return nil
	}
	base := series[n-1-years].Value
	final := series[n-1].Value
	if base <= 0 {
		return nil
	}
	cagr := (math.Pow(final/base, 1/float64(years)) - 1) * 100
	return &cagr
}

// pctChange returns (final-base)/|base| as a percentage; undefined for a
// zero base.
func pctChange(base, final float64) (float64, bool) {
	if base == 0 {
		return 0, false
	}
	return (final - base) / math.Abs(base) * 100, true
}

// annualSeries extracts one labeled row from a raw statement table as a
// numeric series over its annual (non-TTM) periods, in chronological order.
func annualSeries(table model.FinancialTable, label string) []model.SeriesPoint {
	row, ok := findRow(table, label)
	if !ok {
		return nil
	}
	var out []model.SeriesPoint
	for i, raw := range row.Values {
		if i >= len(table.PeriodKeys) || table.PeriodKeys[i] == "TTM" {
			continue
		}
		v, ok := parseNumber(raw)
		if !ok {
			continue
		}
		period := label
		if i < len(table.Periods) {
			period = table.Periods[i]
		}
		out = append(out, model.SeriesPoint{Period: period, Value: v})
	}
	return out
}

func findRow(table model.FinancialTable, label string) (model.DataRow, bool) {
	for _, row := range table.Rows {
		if strings.EqualFold(row.Label, label) {
			return row, true
		}
	}
	return model.DataRow{}, false
}

// combine element-wise merges two series matched by period label (not by
// index, since two tables can start their history in different years),
// preserving a's chronological order. fn's second return value indicates
// whether that point is defined (e.g. division by zero is skipped).
func combine(a, b []model.SeriesPoint, fn func(x, y float64) (float64, bool)) []model.SeriesPoint {
	bByPeriod := make(map[string]float64, len(b))
	for _, p := range b {
		bByPeriod[p.Period] = p.Value
	}
	var out []model.SeriesPoint
	for _, p := range a {
		bv, ok := bByPeriod[p.Period]
		if !ok {
			continue
		}
		v, ok := fn(p.Value, bv)
		if !ok {
			continue
		}
		out = append(out, model.SeriesPoint{Period: p.Period, Value: v})
	}
	return out
}

func last(series []model.SeriesPoint) (float64, bool) {
	if len(series) == 0 {
		return 0, false
	}
	return series[len(series)-1].Value, true
}

// parseNumber cleans screener's display formatting (thousands separators,
// currency symbols, percent signs) and parses the underlying number.
func parseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	s = strings.NewReplacer(",", "", "%", "", "₹", "", "Cr.", "", "Cr", "").Replace(s)
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return 0, false
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func buildNotes(a model.Analysis) []string {
	var notes []string
	note := func(m model.GrowthMetric, subject string) {
		if !m.Found || m.LatestYoYPct == nil {
			return
		}
		dir := "grew"
		if *m.LatestYoYPct < 0 {
			dir = "declined"
		}
		notes = append(notes, subject+" "+dir+" "+trimTrailingZeros(*m.LatestYoYPct)+"% year-over-year in the latest reported year.")
	}
	note(a.SalesGrowth, a.SalesGrowth.Label)
	note(a.ProfitGrowth, "Net profit")
	note(a.EPSGrowth, "EPS")

	if a.OPMTrend.Found && a.OPMTrend.LatestYoYChange != nil {
		dir := "expanded"
		if *a.OPMTrend.LatestYoYChange < 0 {
			dir = "contracted"
		}
		notes = append(notes, "Operating margin "+dir+" by "+trimTrailingZeros(math.Abs(*a.OPMTrend.LatestYoYChange))+" percentage points year-over-year.")
	}

	if a.ROCEPercent != nil {
		notes = append(notes, "ROCE approximates capital employed as Net Worth + Borrowings, since the summarized balance sheet does not separate current and non-current liabilities.")
	}

	return notes
}

func trimTrailingZeros(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64)
}
