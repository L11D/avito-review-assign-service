package tests

import (
	"os"
	"testing"

	"github.com/L11D/avito-review-assign-service/pkg/api/dto"
)

func TestSetIsActive(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "ddd112",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "TeamDfF",
		Members: []dto.TeamMemberDTO{teamMember1},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/users/setIsActive"

	active := GetBoolPtr(false)
	setActiveDTO := dto.UserSetIsActiveDTO{
		UserID:   teamMember1.ID,
		IsActive: active,
	}

	resp, body := MakeJSONRequest(t, "POST", url, setActiveDTO)
	AssertStatusCode(t, resp, 200)
	var updatedUser dto.UserDTO
	ParseJSONResponse(t, body, &updatedUser)

	if updatedUser.ID != teamMember1.ID ||
		updatedUser.IsActive != *active ||
		updatedUser.Username != teamMember1.Username ||
		updatedUser.TeamName != team.Name {
		t.Fatalf("User data does not match expected values")
	}
}

func TestSetIsActive_NotFound(t *testing.T) {
	url := os.Getenv("API_URL") + "/users/setIsActive"

	setActiveDTO := dto.UserSetIsActiveDTO{
		UserID:   "nonexistent_user_id",
		IsActive: GetBoolPtr(false),
	}

	resp, _ := MakeJSONRequest(t, "POST", url, setActiveDTO)
	AssertStatusCode(t, resp, 404)
}

func TestGetPr(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "13dff23u1",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "124sftsss",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "Team12b444",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/create"

	createPRDTO := dto.PullRequestCreateDTO{
		ID:       "sfd112",
		Name:     "pull req",
		AuthorID: teamMember1.ID,
	}

	resp, _ = MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/users/getReview"
	resp, body := MakeQueryRequest(t, "GET", url,
		map[string]string{"user_id": teamMember2.ID},
	)
	AssertStatusCode(t, resp, 200)

	var userPRs dto.UserPRsDTO
	ParseJSONResponse(t, body, &userPRs)

	if userPRs.UserID != teamMember2.ID ||
		len(userPRs.PullRequests) != 1 ||
		userPRs.PullRequests[0].ID != createPRDTO.ID ||
		userPRs.PullRequests[0].Name != createPRDTO.Name ||
		userPRs.PullRequests[0].AuthorID != createPRDTO.AuthorID {
		t.Fatalf("User PRs data does not match expected values")
	}
}
