(() => {
  // Plain-English descriptions for financial shorthand and jargon used
  // throughout the page. Looked up by the exact label text as rendered, so
  // an "i" info icon can be attached next to any term that has an entry.
  const GLOSSARY = {
    // Market facts
    'Current Price': 'The price at which the stock last traded on the exchange.',
    '52W High / Low': 'The highest and lowest price the stock has traded at over the past 52 weeks.',
    'Face Value': 'The nominal value assigned to one share by the company (often ₹1, ₹2, ₹5 or ₹10). Used here to work out how many shares are outstanding from the Equity Capital figure.',

    // Computed fundamentals
    'Market Cap': 'Total market value of the company: Current Price × Shares Outstanding.',
    'Shares Outstanding': 'The total number of shares currently held by all shareholders.',
    'EPS': 'Earnings Per Share — Net Profit divided by Shares Outstanding. How much profit the company made per share.',
    'P/E': 'Price-to-Earnings ratio — Current Price divided by EPS. Roughly, how many years of current profit it would take to earn back the share price; a common (rough) gauge of whether a stock is cheap or expensive.',
    'Book Value / Share': 'Net Worth (Equity Capital + Reserves) divided by Shares Outstanding — the accounting value backing each share.',
    'ROE': 'Return on Equity — Net Profit as a percentage of Net Worth. How efficiently the company turns shareholders’ money into profit.',
    'ROCE': 'Return on Capital Employed — EBIT as a percentage of total capital employed (equity + debt). How efficiently the company uses all the capital it has, not just shareholders’ equity.',
    'Debt / Equity': 'Total Borrowings divided by Net Worth. Higher means the company relies more on debt to fund itself; lower is generally safer.',
    'Operating Margin': 'Operating Profit (Sales minus operating Expenses) as a percentage of Sales. Profitability from core operations, before interest and tax.',
    'Net Margin': 'Net Profit as a percentage of Sales. How much of every rupee of sales ends up as actual bottom-line profit.',
    'Depreciation %': 'Depreciation expense as a percentage of Sales. A rough gauge of how capital-intensive the business is (how much it relies on wearing-down physical assets).',
    'Dividend Yield': 'Annual dividend per share as a percentage of the current price — the cash return you get just from holding the stock, separate from any price change.',
    'OCF / Sales': 'Operating Cash Flow as a percentage of Sales. How much of revenue is actually converted into real cash rather than staying on paper.',
    'OCF / Net Profit': 'Operating Cash Flow as a percentage of Net Profit — an "earnings quality" check. A ratio well below 100% can mean reported profit isn’t translating into real cash.',

    // Growth / trend terms
    'CAGR': 'Compound Annual Growth Rate — the steady annual growth rate that would take a value from its starting point to its ending point over that many years.',
    'YoY': 'Year-over-Year — the percentage (or point) change compared to the same period one year earlier.',
    'TTM': 'Trailing Twelve Months — the most recent four quarters added together, a rolling annual figure.',

    // Raw statement line items
    'Sales': 'Total income earned from the company’s core business, before any expenses are deducted.',
    'Revenue': 'Total income earned (used here for banks/NBFCs in place of "Sales") — mainly interest income for a bank, before any expenses are deducted.',
    'Expenses': 'Total costs incurred to run the business, excluding interest, depreciation and tax.',
    'Operating Profit': 'Sales minus Expenses — profit generated from core operations, before interest and tax.',
    'OPM %': 'Operating Profit Margin — Operating Profit as a percentage of Sales.',
    'Other Income': 'Income from non-core activities, e.g. interest earned on surplus cash or gains on investments.',
    'Interest': 'The cost of borrowed money (finance costs) paid during the period. For a bank, this is interest paid to depositors and lenders — a core cost, not a side expense.',
    'Depreciation': 'The accounting charge that spreads the cost of a fixed asset over its useful life, reflecting wear-and-tear/aging.',
    'Profit before tax': 'Profit after all expenses and interest, but before income tax is deducted.',
    'Tax %': 'The effective income tax rate paid on Profit before tax.',
    'Net Profit': 'The company’s final profit after all expenses, interest, depreciation and tax — the "bottom line".',
    'EPS in Rs': 'Earnings Per Share as reported by the company, in rupees.',
    'Dividend Payout %': 'The percentage of Net Profit paid out to shareholders as dividends; the rest is retained in the business.',
    'Equity Capital': 'The total face value of all shares the company has issued.',
    'Reserves': 'Accumulated profits and other equity retained in the business over the years.',
    'Borrowings': 'Total debt owed by the company to lenders.',
    'Borrowing': 'Total debt owed by the company to lenders (excludes customer deposits for a bank — see Deposits).',
    'Deposits': 'Money held by a bank on behalf of its customers — a bank’s main source of funding, distinct from Borrowings.',
    'Fixed Assets': 'Long-term physical assets — property, plant and equipment — used to run the business.',
    'CWIP': 'Capital Work in Progress — money spent on fixed assets that are still under construction and not yet in use.',
    'Investments': 'Money placed in financial instruments (shares, bonds, mutual funds) rather than the core business.',
    'Other Assets': 'Assets not separately broken out elsewhere, e.g. cash, receivables, inventory.',
    'Other Liabilities': 'Liabilities not separately broken out elsewhere, e.g. trade payables, provisions.',
    'Total Assets': 'Everything the company owns. Always equal to Total Liabilities (Assets = Liabilities + Equity).',
    'Total Liabilities': 'Everything the company owes plus shareholders’ equity. Always equal to Total Assets.',
    'Cash from Operating Activity': 'Cash generated or used by the company’s core, day-to-day business operations.',
    'Cash from Investing Activity': 'Cash spent on or received from buying/selling long-term assets and investments.',
    'Cash from Financing Activity': 'Cash raised from or paid to lenders and shareholders — loans, dividends, buybacks, share issues.',
    'Net Cash Flow': 'The total change in the company’s cash balance over the period (Operating + Investing + Financing cash flows).',
    'Free Cash Flow': 'Cash left over after the business covers its operating and capital needs — cash genuinely free to pay down debt, pay dividends, or reinvest.',

    // Financial-company specific
    'Financing Profit': 'Bank/NBFC equivalent of Operating Profit: Revenue minus Interest expended minus operating Expenses.',
    'Financing Margin %': 'Bank/NBFC equivalent of Operating Margin: Financing Profit as a percentage of Revenue.',
    'Gross NPA %': 'Gross Non-Performing Assets — the percentage of a bank’s loans where the borrower has stopped repaying, before accounting for provisions.',
    'Net NPA %': 'Net Non-Performing Assets — the percentage of a bank’s loans that are non-performing after provisions set aside for expected losses are subtracted.',

    // Shareholding pattern
    'Promoters': 'Shares held by the company’s founders and promoter group.',
    'FIIs': 'Foreign Institutional Investors — shares held by overseas investment institutions.',
    'DIIs': 'Domestic Institutional Investors — shares held by Indian institutions such as mutual funds and insurance companies.',
    'Government': 'Shares held by government bodies.',
    'Public': 'Shares held by individual retail investors.',
    'No. of Shareholders': 'The total count of people/entities holding at least one share of the company.',

    'Growth Analysis': 'CAGR (Compound Annual Growth Rate) is the steady yearly growth rate over that many years. YoY (Year-over-Year) is the most recent single-year change. TTM (Trailing Twelve Months) is the latest four quarters combined.',
  };

  // withInfo returns label text with a small "i" info button appended when
  // a glossary entry exists for it. The description is shown in a
  // click-triggered popover (see the info-popover wiring below), not a
  // native title-attribute tooltip.
  function withInfo(label) {
    const desc = GLOSSARY[label];
    const safeLabel = escapeHtml(label);
    if (!desc) return safeLabel;
    return safeLabel + ' <span class="info-btn" tabindex="0" role="button" aria-haspopup="true" aria-expanded="false" data-tip="' + escapeHtml(desc) + '">i</span>';
  }

  // --- Click-triggered info popover -----------------------------------
  // A single shared popover element, positioned next to whichever
  // .info-btn was last activated. Kept as one reused element (rather than
  // one per button) to keep z-index/positioning/outside-click handling simple.
  const infoPopover = document.createElement('div');
  infoPopover.className = 'info-popover hidden';
  infoPopover.setAttribute('role', 'tooltip');
  document.body.appendChild(infoPopover);
  let openInfoBtn = null;

  function closeInfoPopover() {
    infoPopover.classList.add('hidden');
    if (openInfoBtn) openInfoBtn.setAttribute('aria-expanded', 'false');
    openInfoBtn = null;
  }

  function openInfoPopoverFor(btn) {
    infoPopover.textContent = btn.getAttribute('data-tip') || '';
    infoPopover.classList.remove('hidden');
    btn.setAttribute('aria-expanded', 'true');
    openInfoBtn = btn;

    // Position below the button by default, flipping above if there is
    // not enough room, and clamping horizontally to stay on-screen.
    const btnRect = btn.getBoundingClientRect();
    const popRect = infoPopover.getBoundingClientRect();
    let top = btnRect.bottom + 8;
    if (top + popRect.height > window.innerHeight - 8) {
      top = btnRect.top - popRect.height - 8;
    }
    let left = btnRect.left + btnRect.width / 2 - popRect.width / 2;
    left = Math.max(8, Math.min(left, window.innerWidth - popRect.width - 8));
    infoPopover.style.top = top + 'px';
    infoPopover.style.left = left + 'px';
  }

  document.addEventListener('click', (e) => {
    const btn = e.target.closest('.info-btn');
    if (btn) {
      e.stopPropagation();
      if (openInfoBtn === btn) {
        closeInfoPopover();
      } else {
        openInfoPopoverFor(btn);
      }
      return;
    }
    if (!e.target.closest('.info-popover')) closeInfoPopover();
  });

  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
      closeInfoPopover();
      return;
    }
    if ((e.key === 'Enter' || e.key === ' ') && e.target.classList && e.target.classList.contains('info-btn')) {
      e.preventDefault();
      if (openInfoBtn === e.target) {
        closeInfoPopover();
      } else {
        openInfoPopoverFor(e.target);
      }
    }
  });

  window.addEventListener('scroll', closeInfoPopover, true);
  window.addEventListener('resize', closeInfoPopover);

  const input = document.getElementById('search-input');
  const suggestionsEl = document.getElementById('suggestions');
  const statusEl = document.getElementById('status');
  const resultEl = document.getElementById('result');

  let debounceTimer = null;
  let activeIndex = -1;
  let currentSuggestions = [];

  input.addEventListener('input', () => {
    clearTimeout(debounceTimer);
    const q = input.value.trim();
    if (q.length < 2) {
      hideSuggestions();
      return;
    }
    debounceTimer = setTimeout(() => fetchSuggestions(q), 250);
  });

  input.addEventListener('keydown', (e) => {
    if (suggestionsEl.classList.contains('hidden')) {
      if (e.key === 'Enter') runAnalysis(input.value.trim());
      return;
    }
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setActive(activeIndex + 1);
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setActive(activeIndex - 1);
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (activeIndex >= 0 && currentSuggestions[activeIndex]) {
        selectSuggestion(currentSuggestions[activeIndex]);
      } else {
        runAnalysis(input.value.trim());
      }
    } else if (e.key === 'Escape') {
      hideSuggestions();
    }
  });

  document.addEventListener('click', (e) => {
    if (!e.target.closest('.search-wrap')) hideSuggestions();
  });

  async function fetchSuggestions(q) {
    try {
      const res = await fetch('/api/search?q=' + encodeURIComponent(q));
      if (!res.ok) return;
      currentSuggestions = await res.json();
      renderSuggestions();
    } catch (_) {
      /* ignore transient autocomplete errors */
    }
  }

  function renderSuggestions() {
    activeIndex = -1;
    if (!currentSuggestions.length) {
      hideSuggestions();
      return;
    }
    suggestionsEl.innerHTML = '';
    currentSuggestions.forEach((c, i) => {
      const li = document.createElement('li');
      li.textContent = c.name;
      li.addEventListener('click', () => selectSuggestion(c));
      suggestionsEl.appendChild(li);
    });
    suggestionsEl.classList.remove('hidden');
  }

  function setActive(idx) {
    const items = [...suggestionsEl.children];
    if (!items.length) return;
    activeIndex = (idx + items.length) % items.length;
    items.forEach((el, i) => el.classList.toggle('active', i === activeIndex));
    items[activeIndex].scrollIntoView({ block: 'nearest' });
  }

  function hideSuggestions() {
    suggestionsEl.classList.add('hidden');
    suggestionsEl.innerHTML = '';
    currentSuggestions = [];
    activeIndex = -1;
  }

  function selectSuggestion(c) {
    input.value = c.name;
    hideSuggestions();
    runAnalysis(c.name);
  }

  async function runAnalysis(name) {
    if (!name) return;
    hideSuggestions();
    resultEl.classList.add('hidden');
    statusEl.classList.remove('hidden', 'error');
    statusEl.textContent = 'Fetching and analyzing ' + name + '...';

    try {
      const res = await fetch('/api/analyze?name=' + encodeURIComponent(name));
      const body = await res.json();
      if (!res.ok) throw new Error(body.error || 'request failed');
      render(body);
      statusEl.classList.add('hidden');
    } catch (err) {
      statusEl.textContent = 'Could not load data: ' + err.message;
      statusEl.classList.add('error');
    }
  }

  function render(data) {
    resultEl.innerHTML = '';
    resultEl.appendChild(renderHeader(data));
    resultEl.appendChild(renderComputedFundamentals(data));
    resultEl.appendChild(renderGrowthAnalysis(data));
    for (const table of [data.quarterly, data.profit_loss, data.balance_sheet, data.cash_flow, data.shareholding]) {
      if (table && table.rows && table.rows.length) {
        resultEl.appendChild(renderTable(table));
      }
    }
    resultEl.classList.remove('hidden');
  }

  function el(tag, className, html) {
    const e = document.createElement(tag);
    if (className) e.className = className;
    if (html !== undefined) e.innerHTML = html;
    return e;
  }

  function renderHeader(data) {
    const card = el('section', 'card company-header');
    card.appendChild(el('h2', null, escapeHtml(data.name || 'Unknown company')));
    if (data.source_url) {
      card.appendChild(el('div', 'source-link', 'Raw financial statements sourced from <a href="' + data.source_url + '" target="_blank" rel="noopener">' + data.source_url + '</a> — every ratio below is calculated by this app, not copied from the source.'));
    }

    const m = data.market || {};
    const grid = el('div', 'ratio-grid');
    const item = (label, value) => {
      const d = el('div', 'ratio-item');
      d.appendChild(el('div', 'label', withInfo(label)));
      d.appendChild(el('div', 'value', escapeHtml(value)));
      grid.appendChild(d);
    };
    if (m.current_price) item('Current Price', '₹ ' + fmtNum(m.current_price));
    if (m.high_52w) item('52W High / Low', '₹ ' + fmtNum(m.high_52w) + ' / ' + fmtNum(m.low_52w));
    if (m.face_value) item('Face Value', '₹ ' + fmtNum(m.face_value));
    card.appendChild(grid);

    if (data.about) {
      card.appendChild(el('div', 'about-text', escapeHtml(data.about)));
    }
    return card;
  }

  function renderComputedFundamentals(data) {
    const a = data.analysis || {};
    const card = el('section', 'card');
    card.appendChild(el('h2', 'table-title', 'Computed Fundamentals'));
    card.appendChild(el('div', 'about-text', 'Calculated in Go from the raw Profit & Loss and Balance Sheet figures above — not sourced from any pre-computed ratio.'));

    const grid = el('div', 'ratio-grid');
    const item = (label, value) => {
      if (value === null || value === undefined) return;
      const d = el('div', 'ratio-item');
      d.appendChild(el('div', 'label', withInfo(label)));
      d.appendChild(el('div', 'value', escapeHtml(value)));
      grid.appendChild(d);
    };
    item('Market Cap', a.market_cap_cr != null ? '₹ ' + fmtNum(a.market_cap_cr) + ' Cr.' : undefined);
    item('Shares Outstanding', a.shares_outstanding_cr != null ? fmtNum(a.shares_outstanding_cr) + ' Cr.' : undefined);
    item('EPS', a.eps != null ? '₹ ' + fmtNum(a.eps) : undefined);
    item('P/E', a.pe != null ? fmtNum(a.pe) : undefined);
    item('Book Value / Share', a.book_value_per_share != null ? '₹ ' + fmtNum(a.book_value_per_share) : undefined);
    item('ROE', a.roe_percent != null ? fmtNum(a.roe_percent) + ' %' : undefined);
    item('ROCE', a.roce_percent != null ? fmtNum(a.roce_percent) + ' %' : undefined);
    item('Debt / Equity', a.debt_to_equity != null ? fmtNum(a.debt_to_equity) : undefined);
    item('Operating Margin', a.opm_percent != null ? fmtNum(a.opm_percent) + ' %' : undefined);
    item('Net Margin', a.net_margin_percent != null ? fmtNum(a.net_margin_percent) + ' %' : undefined);
    item('Depreciation %', a.depreciation_percent != null ? fmtNum(a.depreciation_percent) + ' %' : undefined);
    item('Dividend Yield', a.dividend_yield_percent != null ? fmtNum(a.dividend_yield_percent) + ' %' : undefined);
    item('OCF / Sales', a.ocf_to_sales_percent != null ? fmtNum(a.ocf_to_sales_percent) + ' %' : undefined);
    item('OCF / Net Profit', a.ocf_to_net_profit_percent != null ? fmtNum(a.ocf_to_net_profit_percent) + ' %' : undefined);
    card.appendChild(grid);

    return card;
  }

  function fmtNum(v) {
    if (v === null || v === undefined) return '-';
    return Number(v).toLocaleString('en-IN', { maximumFractionDigits: 2 });
  }

  function renderGrowthAnalysis(data) {
    const card = el('section', 'card');
    card.appendChild(el('h2', 'table-title', withInfo('Growth Analysis')));

    const grid = el('div', 'growth-grid');
    const metrics = [
      data.analysis && data.analysis.sales_growth,
      data.analysis && data.analysis.profit_growth,
      data.analysis && data.analysis.eps_growth,
      data.analysis && data.analysis.opm_trend,
      data.analysis && data.analysis.shares_outstanding_trend,
    ].filter(m => m && m.found);
    metrics.forEach(m => grid.appendChild(renderGrowthCard(m)));
    if (metrics.length) card.appendChild(grid);

    if (data.analysis && data.analysis.notes && data.analysis.notes.length) {
      const ul = el('ul', 'notes-list');
      data.analysis.notes.forEach(n => ul.appendChild(el('li', null, escapeHtml(n))));
      card.appendChild(ul);
    }

    return card;
  }

  function renderGrowthCard(m) {
    const card = el('div', 'growth-card');
    card.appendChild(el('h3', null, withInfo(m.label) + (m.unit ? ' (' + escapeHtml(m.unit) + ')' : '')));

    const spark = el('div', 'sparkline');
    spark.innerHTML = sparklineSvg(m.series);
    card.appendChild(spark);

    if (m.latest_yoy_pct !== undefined) {
      card.appendChild(cagrRow('Latest YoY', m.latest_yoy_pct, '%'));
    }
    if (m.latest_yoy_change !== undefined) {
      card.appendChild(cagrRow('Latest YoY', m.latest_yoy_change, ' pts'));
    }
    if (m.cagr_3y !== undefined) card.appendChild(cagrRow('3Y CAGR', m.cagr_3y, '%'));
    if (m.cagr_5y !== undefined) card.appendChild(cagrRow('5Y CAGR', m.cagr_5y, '%'));
    if (m.cagr_10y !== undefined) card.appendChild(cagrRow('10Y CAGR', m.cagr_10y, '%'));

    return card;
  }

  function cagrRow(label, value, suffix) {
    const row = el('div', 'cagr-row');
    const cls = value >= 0 ? 'pos' : 'neg';
    row.innerHTML = '<span>' + label + '</span><b class="' + cls + '">' + value.toFixed(1) + suffix + '</b>';
    return row;
  }

  function sparklineSvg(series) {
    if (!series || series.length < 2) return '';
    const w = 200, h = 46, pad = 4;
    const values = series.map(p => p.value);
    const min = Math.min(...values), max = Math.max(...values);
    const range = max - min || 1;
    const step = (w - pad * 2) / (series.length - 1);
    const points = values.map((v, i) => {
      const x = pad + i * step;
      const y = h - pad - ((v - min) / range) * (h - pad * 2);
      return x.toFixed(1) + ',' + y.toFixed(1);
    }).join(' ');
    const last = values[values.length - 1] >= values[0] ? '#3ecf8e' : '#ff6b6b';
    return '<svg viewBox="0 0 ' + w + ' ' + h + '" width="100%" height="' + h + '" preserveAspectRatio="none">' +
      '<polyline fill="none" stroke="' + last + '" stroke-width="2" points="' + points + '" /></svg>';
  }

  // dataCell builds one statement-table cell: the reported value, plus a
  // small colored line underneath showing that period's percentage change
  // from the previous period, when one was computable.
  function dataCell(value, changePct) {
    const td = document.createElement('td');
    td.appendChild(el('div', 'cell-value', escapeHtml(value)));
    if (changePct !== null && changePct !== undefined) {
      const cls = changePct >= 0 ? 'pos' : 'neg';
      const arrow = changePct >= 0 ? '▲' : '▼';
      td.appendChild(el('div', 'cell-delta ' + cls, arrow + ' ' + Math.abs(changePct).toFixed(1) + '%'));
    }
    return td;
  }

  function renderTable(table) {
    const card = el('section', 'card');
    card.appendChild(el('h2', 'table-title', escapeHtml(table.title)));
    card.appendChild(el('div', 'about-text', 'Each ▲/▼ is that period’s % change from the previous one.'));
    const scroll = el('div', 'table-scroll');
    const el2 = document.createElement('table');
    el2.className = 'data-table';

    const thead = document.createElement('thead');
    const headRow = document.createElement('tr');
    headRow.appendChild(document.createElement('th'));
    (table.periods || []).forEach(p => headRow.appendChild(el('th', null, escapeHtml(p))));
    thead.appendChild(headRow);
    el2.appendChild(thead);

    const tbody = document.createElement('tbody');
    (table.rows || []).forEach(row => {
      const tr = document.createElement('tr');
      tr.appendChild(el('td', null, withInfo(row.label)));
      (row.values || []).forEach((v, i) => tr.appendChild(dataCell(v, row.change_percents && row.change_percents[i])));
      tbody.appendChild(tr);
    });
    el2.appendChild(tbody);

    scroll.appendChild(el2);
    card.appendChild(scroll);
    return card;
  }

  function escapeHtml(s) {
    if (s === undefined || s === null) return '';
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }
})();
