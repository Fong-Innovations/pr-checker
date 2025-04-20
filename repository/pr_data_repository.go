package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pr-checker/models"
)

type PRDataRepository struct {
	db *sql.DB
}

const (
	insertPRDataQuery = `INSERT INTO pr_metrics (
            repo,
            target_branch,
            source_branch,
            merged,
            comments,
            changed_files,
            opened_at,
            merged_at,
            closed_at,
            issue_url,
            user
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
)

// NewPRDataRepository creates a new instance of PRDataRepository
func NewPRDataRepository(db *sql.DB) *PRDataRepository {
	return &PRDataRepository{db: db}
}

// GetPRData retrieves pull request data from the database
// func (r *PRDataRepository) GetPRData(id int) (string, error) {
// 	var data string
// 	query := "SELECT data FROM pr_data WHERE id = ?"
// 	err := r.db.QueryRow(query, id).Scan(&data)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get PR data: %w", err)
// 	}
// 	return data, nil
// }

// SavePRData saves pull request data to the database
func (r *PRDataRepository) InsertPRDataEntry(ctx context.Context, data models.PullRequestDBEntry) error {
	query := insertPRDataQuery
	_, err := r.db.Exec(query, data.Repo, data.TargetBranch, data.SourceBranch, data.Merged, data.Comments, data.ChangedFiles, data.OpenedAt, data.MergedAt, data.ClosedAt, data.IssueUrl, data.User)
	if err != nil {
		return fmt.Errorf("failed to save PR data: %w", err)
	}
	return nil
}
