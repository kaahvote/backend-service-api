package data

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
