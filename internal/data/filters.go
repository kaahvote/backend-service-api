package data

import (
	"strings"
	"unicode"
)

type Metadata struct {
	CurrentPage  int `json:"currentPage"`
	PageSize     int `json:"pageSize"`
	FirstPage    int `json:"firstPage"`
	LastPage     int `json:"lastPage"`
	TotalRecords int `json:"totalRecords"`
}

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

type FlowFilters struct {
	SessionID int64
	Filters
}

type SessionFilters struct {
	Name              string
	VotingPolicyID    int64
	VotersPolicyID    int64
	CandidatePolicyID int64
	CreatedBy         int64
	ExpiresAtFrom     *string
	ExpiresAtTo       *string
	CreatedAtFrom     *string
	CreatedAtTo       *string
	Filters
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	if pageSize == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) sortColumn() string {
	for _, saveValue := range f.SortSafeList {
		if f.Sort == saveValue {
			col := strings.TrimPrefix(f.Sort, "-")
			hasUpper, char := HasUpperCase(col)
			if hasUpper {
				newChar := "_" + strings.ToLower(char)
				return strings.ReplaceAll(col, char, newChar)
			}
			return col
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f FlowFilters) sortColumn() string {
	col := f.Filters.sortColumn()
	if col == "state" {
		return "state_id"
	}
	return col
}

func HasUpperCase(s string) (bool, string) {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true, string(r)
		}
	}
	return false, ""
}
