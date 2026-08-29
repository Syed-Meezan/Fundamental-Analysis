// Package httpapi wires the screener client and analysis engine into HTTP
// endpoints consumed by the web frontend.
package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"sampleproject/internal/analysis"
	"sampleproject/internal/model"
	"sampleproject/internal/screener"
)

// Server holds the shared dependencies for the API handlers.
type Server struct {
	screener *screener.Client
}

// NewServer builds a Server ready to register routes.
func NewServer() *Server {
	return &Server{screener: screener.NewClient()}
}

// Routes registers the API endpoints on mux.
func (s *Server) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/search", s.handleSearch)
	mux.HandleFunc("GET /api/analyze", s.handleAnalyze)
}

// handleSearch proxies screener.in's company search, for a frontend
// autocomplete dropdown.
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) < 2 {
		writeJSON(w, http.StatusOK, []model.SearchResult{})
		return
	}

	results, err := s.screener.Search(q)
	if err != nil {
		log.Printf("search %q: %v", q, err)
		writeError(w, http.StatusBadGateway, "could not search screener.in right now")
		return
	}
	writeJSON(w, http.StatusOK, results)
}

// handleAnalyze resolves a free-text company name to a screener.in company
// page, scrapes its financial statements, runs the growth-trend analysis,
// and returns everything as one JSON payload.
func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "missing required query param: name")
		return
	}

	data, err := s.analyzeByName(name)
	if err != nil {
		log.Printf("analyze %q: %v", name, err)
		status := http.StatusBadGateway
		if errors.Is(err, errNoMatch) {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, data)
}

var errNoMatch = errors.New("no matching company found on screener.in")

func (s *Server) analyzeByName(name string) (*model.CompanyData, error) {
	results, err := s.screener.Search(name)
	if err != nil {
		return nil, fmt.Errorf("searching for %q: %w", name, err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("%q: %w", name, errNoMatch)
	}

	best := results[0]
	html, err := s.screener.FetchCompanyPage(best.URL)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", best.Name, err)
	}

	data, err := screener.ParseCompanyPage(html, "https://www.screener.in"+best.URL)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", best.Name, err)
	}

	for _, t := range []*model.FinancialTable{&data.Quarterly, &data.ProfitLoss, &data.BalanceSheet, &data.CashFlow, &data.Shareholding} {
		analysis.AnnotateRowChanges(t)
	}
	analysis.Compute(data)
	return data, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
