package tests

import (
	"os"
	"reflect"
	"testing"

	"github.com/L11D/avito-review-assign-service/pkg/api/dto"
)

func TestCreateTeam(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "u1",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "u2",
		Username: "Alice",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "TeamA",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, body := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)
	var createdTeam dto.TeamDTO
	ParseJSONResponse(t, body, &createdTeam)

	if createdTeam.Name != team.Name {
		t.Fatalf("expected team name %s, got %s", team.Name, createdTeam.Name)
	}

	if !reflect.DeepEqual(createdTeam.Members, team.Members) {
		t.Fatal("team members do not match")
	}
}

func TestCreateTeam_SameName(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "u11",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "u21",
		Username: "Alice",
		IsActive: GetBoolPtr(true),
	}

	team1 := dto.TeamDTO{
		Name:    "Team1",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team1)
	AssertStatusCode(t, resp, 201)

	resp, _ = MakeJSONRequest(t, "POST", url, team1)
	AssertStatusCode(t, resp, 400)
}

func TestCreateTeam_SameUser(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "s1",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "s2",
		Username: "Alice",
		IsActive: GetBoolPtr(true),
	}

	team1 := dto.TeamDTO{
		Name:    "Team88",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team1)
	AssertStatusCode(t, resp, 201)

	team2 := dto.TeamDTO{
		Name:    "Team2",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	resp, _ = MakeJSONRequest(t, "POST", url, team2)
	AssertStatusCode(t, resp, 400)
}

func TestGetTeam(t *testing.T) {
	teamMember1 := dto.TeamMemberDTO{
		ID:       "d1",
		Username: "Bob",
		IsActive: GetBoolPtr(true),
	}

	teamMember2 := dto.TeamMemberDTO{
		ID:       "d2",
		Username: "Alice",
		IsActive: GetBoolPtr(true),
	}

	team := dto.TeamDTO{
		Name:    "TeamD",
		Members: []dto.TeamMemberDTO{teamMember1, teamMember2},
	}

	url := os.Getenv("API_URL") + "/team/add"
	resp, _ := MakeJSONRequest(t, "POST", url, team)
	AssertStatusCode(t, resp, 201)

	url = os.Getenv("API_URL") + "/team/get"
	resp, body := MakeQueryRequest(t, "GET", url,
		map[string]string{"name": team.Name},
	)
	AssertStatusCode(t, resp, 200)

	var fetchedTeam dto.TeamDTO
	ParseJSONResponse(t, body, &fetchedTeam)

	if fetchedTeam.Name != team.Name {
		t.Fatalf("expected team name %s, got %s", team.Name, fetchedTeam.Name)
	}

	if !reflect.DeepEqual(fetchedTeam.Members, team.Members) {
		t.Fatal("team members do not match")
	}
}

func TestGetTeam_NotFound(t *testing.T) {
	url := os.Getenv("API_URL") + "/team/get"
	resp, _ := MakeQueryRequest(t, "GET", url,
		map[string]string{"name": "nonexistent_team"},
	)
	AssertStatusCode(t, resp, 404)
}
