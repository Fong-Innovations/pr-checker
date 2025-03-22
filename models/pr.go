package models

import "time"

type PullRequestRequest struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	RepoID  string `json:"repo_id"`
}

type PullRequestResponse struct {
	PR PullRequest
}

type Comment struct {
	ID       string
	Category string
	Content  string
}

type PullRequest struct {
	URL                 string      `json:"url"`
	ID                  int         `json:"id"`
	NodeID              string      `json:"node_id"`
	HTMLURL             string      `json:"html_url"`
	DiffURL             string      `json:"diff_url"`
	PatchURL            string      `json:"patch_url"`
	IssueURL            string      `json:"issue_url"`
	CommitsURL          string      `json:"commits_url"`
	ReviewCommentsURL   string      `json:"review_comments_url"`
	ReviewCommentURL    string      `json:"review_comment_url"`
	CommentsURL         string      `json:"comments_url"`
	StatusesURL         string      `json:"statuses_url"`
	Number              int         `json:"number"`
	State               string      `json:"state"`
	Locked              bool        `json:"locked"`
	Title               string      `json:"title"`
	User                User        `json:"user"`
	Body                string      `json:"body"`
	Labels              []Label     `json:"labels"`
	Milestone           Milestone   `json:"milestone"`
	ActiveLockReason    string      `json:"active_lock_reason"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
	ClosedAt            time.Time   `json:"closed_at"`
	MergedAt            time.Time   `json:"merged_at"`
	MergeCommitSHA      string      `json:"merge_commit_sha"`
	Assignee            User        `json:"assignee"`
	Assignees           []User      `json:"assignees"`
	RequestedReviewers  []User      `json:"requested_reviewers"`
	RequestedTeams      []Team      `json:"requested_teams"`
	Head                Branch      `json:"head"`
	Base                Branch      `json:"base"`
	Links               Links       `json:"_links"`
	AuthorAssociation   string      `json:"author_association"`
	AutoMerge           interface{} `json:"auto_merge"`
	Draft               bool        `json:"draft"`
	Merged              bool        `json:"merged"`
	Mergeable           bool        `json:"mergeable"`
	Rebaseable          bool        `json:"rebaseable"`
	MergeableState      string      `json:"mergeable_state"`
	MergedBy            User        `json:"merged_by"`
	Comments            int         `json:"comments"`
	ReviewComments      int         `json:"review_comments"`
	MaintainerCanModify bool        `json:"maintainer_can_modify"`
	Commits             int         `json:"commits"`
	Additions           int         `json:"additions"`
	Deletions           int         `json:"deletions"`
	ChangedFiles        int         `json:"changed_files"`
}

type User struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Label struct {
	ID          int    `json:"id"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

type Milestone struct {
	URL          string    `json:"url"`
	HTMLURL      string    `json:"html_url"`
	LabelsURL    string    `json:"labels_url"`
	ID           int       `json:"id"`
	NodeID       string    `json:"node_id"`
	Number       int       `json:"number"`
	State        string    `json:"state"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Creator      User      `json:"creator"`
	OpenIssues   int       `json:"open_issues"`
	ClosedIssues int       `json:"closed_issues"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ClosedAt     time.Time `json:"closed_at"`
	DueOn        time.Time `json:"due_on"`
}

type Team struct {
	ID                  int    `json:"id"`
	NodeID              string `json:"node_id"`
	URL                 string `json:"url"`
	HTMLURL             string `json:"html_url"`
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	Description         string `json:"description"`
	Privacy             string `json:"privacy"`
	NotificationSetting string `json:"notification_setting"`
	Permission          string `json:"permission"`
	MembersURL          string `json:"members_url"`
	RepositoriesURL     string `json:"repositories_url"`
}

type Branch struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	SHA   string `json:"sha"`
	User  User   `json:"user"`
	Repo  Repo   `json:"repo"`
}

type Repo struct {
	ID               int         `json:"id"`
	NodeID           string      `json:"node_id"`
	Name             string      `json:"name"`
	FullName         string      `json:"full_name"`
	Owner            User        `json:"owner"`
	Private          bool        `json:"private"`
	HTMLURL          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ArchiveURL       string      `json:"archive_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BlobsURL         string      `json:"blobs_url"`
	BranchesURL      string      `json:"branches_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	CommentsURL      string      `json:"comments_url"`
	CommitsURL       string      `json:"commits_url"`
	CompareURL       string      `json:"compare_url"`
	ContentsURL      string      `json:"contents_url"`
	ContributorsURL  string      `json:"contributors_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	DownloadsURL     string      `json:"downloads_url"`
	EventsURL        string      `json:"events_url"`
	ForksURL         string      `json:"forks_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitURL           string      `json:"git_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	IssuesURL        string      `json:"issues_url"`
	KeysURL          string      `json:"keys_url"`
	LabelsURL        string      `json:"labels_url"`
	LanguagesURL     string      `json:"languages_url"`
	MergesURL        string      `json:"merges_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	PullsURL         string      `json:"pulls_url"`
	ReleasesURL      string      `json:"releases_url"`
	SSHURL           string      `json:"ssh_url"`
	StargazersURL    string      `json:"stargazers_url"`
	StatusesURL      string      `json:"statuses_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	TagsURL          string      `json:"tags_url"`
	TeamsURL         string      `json:"teams_url"`
	TreesURL         string      `json:"trees_url"`
	CloneURL         string      `json:"clone_url"`
	MirrorURL        string      `json:"mirror_url"`
	HooksURL         string      `json:"hooks_url"`
	SVNURL           string      `json:"svn_url"`
	Homepage         string      `json:"homepage"`
	Language         interface{} `json:"language"`
	ForksCount       int         `json:"forks_count"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Size             int         `json:"size"`
	DefaultBranch    string      `json:"default_branch"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	Topics           []string    `json:"topics"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	HasDownloads     bool        `json:"has_downloads"`
	HasDiscussions   bool        `json:"has_discussions"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	PushedAt         time.Time   `json:"pushed_at"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Permissions      Permissions `json:"permissions"`
	AllowRebaseMerge bool        `json:"allow_rebase_merge"`
	TempCloneToken   string      `json:"temp_clone_token"`
	AllowSquashMerge bool        `json:"allow_squash_merge"`
	AllowMergeCommit bool        `json:"allow_merge_commit"`
	AllowForking     bool        `json:"allow_forking"`
	Forks            int         `json:"forks"`
	OpenIssues       int         `json:"open_issues"`
	License          License     `json:"license"`
	Watchers         int         `json:"watchers"`
}

type Permissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	SPDXID string `json:"spdx_id"`
	NodeID string `json:"node_id"`
}

type Links struct {
	Self           Link `json:"self"`
	HTML           Link `json:"html"`
	Issue          Link `json:"issue"`
	Comments       Link `json:"comments"`
	ReviewComments Link `json:"review_comments"`
	ReviewComment  Link `json:"review_comment"`
	Commits        Link `json:"commits"`
	Statuses       Link `json:"statuses"`
}

type Link struct {
	Href string `json:"href"`
}
