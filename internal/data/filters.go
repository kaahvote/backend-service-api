package data

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
