package data

import "time"

type Session struct {
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	PublicID           string    `json:"publicId"`
	ExpiresAt          time.Time `json:"expiresAt"`
	VotingPolicyID     int64     `json:"votingPolicyId"`
	VotersPolicyID     int64     `json:"votersPolicyId"`
	CandidatesPolicyID int64     `json:"candidatesPolicyId"`
	CreatedBy          int64     `json:"createdBy"`
	CreatedAt          time.Time `json:"createdAt"`
}
