package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type VoterPolicy struct {
	BaseFields
}

type VoterPolicyModel struct {
	DB *sql.DB
}

func (m *VoterPolicyModel) ListFiltering(filters *SettingsFilters) ([]*VoterPolicy, Metadata, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) OVER(), id, name, created_at 
			  FROM voter_policies
			  WHERE 1=1
			  AND (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
			  AND (created_at >= $2 OR $2 IS NULL)
			  AND (created_at <= $3 OR $3 IS NULL)
			  ORDER BY %s %s, id ASC
			  LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	args := []any{strings.ToLower(filters.Name), filters.CreatedAtFrom, filters.CreatedAtTo, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	policies := make([]*VoterPolicy, 0)
	totalRecords := 0
	defer rows.Close()

	for rows.Next() {
		var p VoterPolicy
		err = rows.Scan(&totalRecords, &p.ID, &p.Name, &p.CreatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}

		policies = append(policies, &p)
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return policies, metadata, nil
}
