package tests

import (
	"os"
	"testing"

	"github.com/L11D/avito-review-assign-service/pkg/api/dto"
)

func TestCreatePullRequest(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "1323u1",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "124sft",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "Team444",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/create"

	createPRDTO := dto.PullRequestCreateDTO{
		ID:       "SomePRID123",
		Name:     "pull req",
		AuthorID: teamMember1.ID,
	}

	resp, body := MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 201)

	var fetchedFullPR dto.FullPullRequestDTO
	ParseJSONResponse(t, body, &fetchedFullPR)
	fetchedPR := fetchedFullPR.PullRequest

	if fetchedPR.ID != createPRDTO.ID ||
		fetchedPR.Name != createPRDTO.Name ||
		fetchedPR.AuthorID != createPRDTO.AuthorID ||
		fetchedPR.Status != dto.StatusOpen ||
		fetchedPR.Reviewers[0] != teamMember2.ID {
		t.Fatalf("Pull Request data does not match expected values")
	}
}

func TestCreatePullRequest_MemberNotActive(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "1323u1a13",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "124sfst321",
		Username: "Bob",
		IsActive: GetBoolPtr(false),
	}

	team := dto.TeamDTO{
		Name:    "Team423411",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, body := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/create"

	createPRDTO := dto.PullRequestCreateDTO{
		ID:       "SomePRID123sda",
		Name:     "pull req",
		AuthorID: teamMember1.ID,
	}

	resp, body = MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 201)

	var fetchedFullPR dto.FullPullRequestDTO
	ParseJSONResponse(t, body, &fetchedFullPR)
	fetchedPR := fetchedFullPR.PullRequest

	if fetchedPR.ID != createPRDTO.ID ||
		fetchedPR.Name != createPRDTO.Name ||
		fetchedPR.AuthorID != createPRDTO.AuthorID ||
		fetchedPR.Status != dto.StatusOpen ||
		len(fetchedPR.Reviewers) > 0 {
		t.Fatalf("Pull Request data does not match expected values")
	}
}

func TestCreatePullRequest_UserNotFound(t *testing.T) {
	url := os.Getenv("API_URL") + "/pullRequest/create"

	createPRDTO := dto.PullRequestCreateDTO{
		ID:       "SomePweeRID123",
		Name:     "pull req",
		AuthorID: "nonexistent_user_id",
	}

	resp, _ := MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 404)
}

func TestCreatePullRequest_PRAlreadyExists(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "66hj1a13",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "Teamfgg4411",
		Members: []dto.TeamMemberDTO{teamMember1},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/create"

	createPRDTO := dto.PullRequestCreateDTO{
		ID:       "SomePRID123sdadfsfd",
		Name:     "pull req",
		AuthorID: teamMember1.ID,
	}

	resp, _ = MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 201)

	resp, _ = MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 409)
}

func TestMergePullRequest(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "fgfgsgrtt",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "fsgsfdgft22",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "Teamffggg555",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/create"

	createPRDTO := dto.PullRequestCreateDTO{
		ID:       "SomedffD123d",
		Name:     "pull req",
		AuthorID: teamMember1.ID,
	}

	resp, _ = MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/merge"
	mergePRDTO := dto.PullRequestMergeDTO{
		ID: createPRDTO.ID,
	}

	resp, body := MakeJSONRequest(t, "POST", url, mergePRDTO)
	AssertStatusCode(t, resp, 200)

	var fetchedFullPR dto.FullPullRequestDTO
	ParseJSONResponse(t, body, &fetchedFullPR)
	fetchedPR := fetchedFullPR.PullRequest

	if fetchedPR.ID != createPRDTO.ID ||
		fetchedPR.Name != createPRDTO.Name ||
		fetchedPR.AuthorID != createPRDTO.AuthorID ||
		fetchedPR.Status != dto.StatusMerged {
		t.Fatalf("Pull Request data does not match expected values")
	}
}

func TestReassignPullRequest(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "nmmn5cds",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "dfhht34",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember3 := dto.TeamMemberDTO{
		ID:       "jkkllli7",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember4 := dto.TeamMemberDTO{
		ID:       "vbvbb4455",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "Teamddsfs234556gg",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2, teamMember3, teamMember4},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/create"

	createPRDTO := dto.PullRequestCreateDTO{
		ID:       "Someddmmmmqq11f",
		Name:     "pull req",
		AuthorID: teamMember1.ID,
	}

	resp, body := MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 201)

	var fetchedFullPR dto.FullPullRequestDTO
	ParseJSONResponse(t, body, &fetchedFullPR)
	fetchedPR := fetchedFullPR.PullRequest

	assignRewiever1 := fetchedPR.Reviewers[0]

	url = os.Getenv("API_URL") + "/pullRequest/reassign"
	reassignPRDTO := dto.PullRequestReassignDTO{
		PullRequestID: createPRDTO.ID,
		OldReviewerID: assignRewiever1,
	}

	resp, body = MakeJSONRequest(t, "POST", url, reassignPRDTO)
	AssertStatusCode(t, resp, 200)

	ParseJSONResponse(t, body, &fetchedFullPR)
	fetchedPR = fetchedFullPR.PullRequest

	if fetchedPR.Reviewers[0] == assignRewiever1 ||
		fetchedPR.Reviewers[1] == assignRewiever1 {
		t.Fatalf("Pull Request reassignment did not work as expected")
	}
}

func TestReassignPullRequest_OnMerged(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "gfdgjty5",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "nmnbmb544",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "Teamddd346yy",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/create"

	createPRDTO := dto.PullRequestCreateDTO{
		ID:       "Somedddfs341",
		Name:     "pull req",
		AuthorID: teamMember1.ID,
	}

	resp, _ = MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/merge"
	mergePRDTO := dto.PullRequestMergeDTO{
		ID: createPRDTO.ID,
	}

	resp, _ = MakeJSONRequest(t, "POST", url, mergePRDTO)
	AssertStatusCode(t, resp, 200)

	url = os.Getenv("API_URL") + "/pullRequest/reassign"
	reassignPRDTO := dto.PullRequestReassignDTO{
		PullRequestID: createPRDTO.ID,
		OldReviewerID: teamMember2.ID,
	}

	resp, body := MakeJSONRequest(t, "POST", url, reassignPRDTO)
	AssertStatusCode(t, resp, 409)

	var errorMessage dto.FullErrorDTO
	ParseJSONResponse(t, body, &errorMessage)

	if errorMessage.Error.Code != "PR_MERGED" {
		t.Fatalf("expected error code PR_MERGED, got %s", errorMessage.Error.Code)
	}
}

func TestReassignPullRequest_NoCandidate(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "dfsf33444",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "dfabnyjtub",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "Teamddfsdkjuy656",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/create"

	createPRDTO := dto.PullRequestCreateDTO{
		ID:       "Somedddf22113mm",
		Name:     "pull req",
		AuthorID: teamMember1.ID,
	}

	resp, _ = MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/reassign"
	reassignPRDTO := dto.PullRequestReassignDTO{
		PullRequestID: createPRDTO.ID,
		OldReviewerID: teamMember2.ID,
	}

	resp, body := MakeJSONRequest(t, "POST", url, reassignPRDTO)
	AssertStatusCode(t, resp, 409)

	var errorMessage dto.FullErrorDTO
	ParseJSONResponse(t, body, &errorMessage)

	if errorMessage.Error.Code != "NO_CANDIDATE" {
		t.Fatalf("expected error code NO_CANDIDATE, got %s", errorMessage.Error.Code)
	}
}

func TestReassignPullRequest_NotAssigned(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "hhh556",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "nmbnmbmn554",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "Teamddsfs234556",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/create"

	createPRDTO := dto.PullRequestCreateDTO{
		ID:       "Someddmmmmqq11",
		Name:     "pull req",
		AuthorID: teamMember1.ID,
	}

	resp, _ = MakeJSONRequest(t, "POST", url, createPRDTO)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/pullRequest/reassign"
	reassignPRDTO := dto.PullRequestReassignDTO{
		PullRequestID: createPRDTO.ID,
		OldReviewerID: teamMember1.ID,
	}

	resp, body := MakeJSONRequest(t, "POST", url, reassignPRDTO)
	AssertStatusCode(t, resp, 409)

	var errorMessage dto.FullErrorDTO
	ParseJSONResponse(t, body, &errorMessage)

	if errorMessage.Error.Code != "NOT_ASSIGNED" {
		t.Fatalf("expected error code NOT_ASSIGNED, got %s", errorMessage.Error.Code)
	}
}
