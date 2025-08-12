package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Sessions        SessionModel
	Flows           FlowModel
	Users           UserModel
	VotingPolicy    VotingPolicyModel
	VoterPolicy     VoterPolicyModel
	CandidatePolicy CandidatePolicyModel
	State           StateModel
	Candidate       CandidateModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Sessions:        SessionModel{DB: db},
		Flows:           FlowModel{DB: db},
		Users:           UserModel{DB: db},
		VotingPolicy:    VotingPolicyModel{DB: db},
		VoterPolicy:     VoterPolicyModel{DB: db},
		CandidatePolicy: CandidatePolicyModel{DB: db},
		State:           StateModel{DB: db},
		Candidate:       CandidateModel{DB: db},
	}
}
