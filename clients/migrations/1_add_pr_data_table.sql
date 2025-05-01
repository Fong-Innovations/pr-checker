-- +goose Up
CREATE TABLE pr_metrics (
    id INT AUTO_INCREMENT PRIMARY KEY, -- Unique identifier for the pull request
    repo VARCHAR(255) NOT NULL, -- The name of the repository where the PR was made
    target_branch VARCHAR(255) NOT NULL, -- The branch being merged into (e.g., master)
    source_branch VARCHAR(255) NOT NULL, -- The branch with feature changes
    merged BOOLEAN NOT NULL DEFAULT FALSE, -- Whether the PR has been merged
    comments INT NOT NULL DEFAULT 0, -- Count of comments made on the PR
    changed_files INT NOT NULL DEFAULT 0, -- Count of files changed in the PR
    opened_at DATETIME NOT NULL, -- Timestamp when the PR was opened
    merged_at DATETIME DEFAULT NULL, -- Timestamp when the PR was merged
    closed_at DATETIME DEFAULT NULL, -- Timestamp when the PR was closed
    issue_url VARCHAR(2083) NOT NULL, -- URL to fetch the issue
    user VARCHAR(255) NOT NULL -- Username of the PR author
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- +goose Down
DROP TABLE pr_metrics;