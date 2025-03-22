package handlers

// import (
// 	"ai-api/models"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // MockPRService is a mock implementation of the PRService interface
// type MockPRService struct {
// 	mock.Mock
// }

// func (m *MockPRService) GetPRsFromGitHub(req models.PullRequestRequest) ([]models.PullRequestResponse, error) {
// 	args := m.Called(req)
// 	// return args.Get(0), args.Error(1)
// }

// func TestGetPR(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	tests := []struct {
// 		name           string
// 		owner          string
// 		repo           string
// 		id             string
// 		mockReturn     interface{}
// 		mockError      error
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name:  "Success",
// 			owner: "owner",
// 			repo:  "repo",
// 			id:    "1",
// 			mockReturn: map[string]interface{}{
// 				"id":    1,
// 				"title": "Test PR",
// 			},
// 			mockError:      nil,
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   `{"message": "PR details", "pr": {"id": 1, "title": "Test PR"}}`,
// 		},
// 		{
// 			name:           "BadRequest_MissingOwner",
// 			owner:          "",
// 			repo:           "repo",
// 			id:             "1",
// 			mockReturn:     nil,
// 			mockError:      nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   `{"error": "missing field in FetchPullRequest body"}`,
// 		},
// 		{
// 			name:           "BadRequest_ServiceError",
// 			owner:          "owner",
// 			repo:           "repo",
// 			id:             "1",
// 			mockReturn:     nil,
// 			mockError:      errors.New("service error"),
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   `{"error": "service error"}`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockService := new(MockPRService)
// 			handler := NewPRHandler(mockService)

// 			reqBody := models.PullRequestRequest{
// 				OwnerID: tt.owner,
// 				RepoID:  tt.repo,
// 				ID:      tt.id,
// 			}

// 			mockService.On("GetPRsFromGitHub", reqBody).Return(tt.mockReturn, tt.mockError)

// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)
// 			c.Params = gin.Params{
// 				{Key: "owner", Value: tt.owner},
// 				{Key: "repo", Value: tt.repo},
// 				{Key: "id", Value: tt.id},
// 			}

// 			handler.GetPR(c)

// 			assert.Equal(t, tt.expectedStatus, w.Code)
// 			assert.JSONEq(t, tt.expectedBody, w.Body.String())
// 			mockService.AssertExpectations(t)
// 		})
// 	}
// }
