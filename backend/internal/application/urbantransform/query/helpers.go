package query

import "github.com/masterfabric-go/masterfabric/internal/shared/pagination"

// normalizePage clamps page/perPage into valid pagination parameters.
func normalizePage(page, perPage int) pagination.Params {
	if page < 1 {
		page = pagination.DefaultPage
	}
	if perPage < 1 {
		perPage = pagination.DefaultPerPage
	}
	if perPage > pagination.MaxPerPage {
		perPage = pagination.MaxPerPage
	}
	return pagination.Params{Page: page, PerPage: perPage}
}
