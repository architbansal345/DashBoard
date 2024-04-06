package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type TeamData struct {
	Teams []Team `json:"teams"`
}

type PlayerData struct {
	Players []Player `json:"players"`
}

type Player struct {
	Name               string `json:"Name"`
	Wickets            int    `json:"Wickets"`
	Runs               int    `json:"Runs"`
	Catches            int    `json:"Catches"`
	ManOfMatchesAwards int    `json:"ManOfMatchesAwards"`
}

type Team struct {
	Name          string `json:"Name"`
	MatchesPlayed int    `json:"MatchesPlayed"`
	Wins          int    `json:"Wins"`
	Losses        int    `json:"Losses"`
	Points        int    `json:"Points"`
}

func sortByWickets(players []Player) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Wickets > players[j].Wickets
	})
}

func sortByRuns(players []Player) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Runs > players[j].Runs
	})
}

func sortByValue(players []Player) {
	sort.Slice(players, func(i, j int) bool {
		return calculateValue(players[i]) > calculateValue(players[j])
	})
}

func calculateValue(player Player) int {
	return player.Runs + 20*player.Wickets + 5*player.Catches + 100*player.ManOfMatchesAwards
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("static/*.html")
	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// Read team data from file
	teamDataFile, err := ioutil.ReadFile("team.json")
	if err != nil {
		fmt.Println("Error reading team data:", err)
		return
	}
	var teamData TeamData
	err = json.Unmarshal(teamDataFile, &teamData)
	if err != nil {
		fmt.Println("Error unmarshalling team data:", err)
		return
	}

	// Read player data from file
	playerDataFile, err := ioutil.ReadFile("player.json")
	if err != nil {
		fmt.Println("Error reading player data:", err)
		return
	}
	var playerData PlayerData
	err = json.Unmarshal(playerDataFile, &playerData)
	if err != nil {
		fmt.Println("Error unmarshalling player data:", err)
		return
	}

	// Routing
	r.GET("/top-players-by-wickets", func(c *gin.Context) {
		sortByWickets(playerData.Players)
		c.HTML(http.StatusOK, "wickets.html", gin.H{
			"Player": playerData.Players[:10], // Pass top 10 players
		})
	})

	r.GET("/top-players-by-runs", func(c *gin.Context) {
		sortByRuns(playerData.Players)
		c.HTML(http.StatusOK, "runs.html", gin.H{
			"Player": playerData.Players[:10], // Pass top 10 players
		})

	})

	r.GET("/top-players-by-value", func(c *gin.Context) {
		sortByValue(playerData.Players)
		c.HTML(http.StatusOK, "value.html", gin.H{
			"Player": playerData.Players[:10], // Pass top 10 players
		})
	})
	r.GET("/teams", func(c *gin.Context) {
		sort.Slice(teamData.Teams, func(i, j int) bool {
			return teamData.Teams[i].Points > teamData.Teams[j].Points
		})
		c.HTML(http.StatusOK, "teams.html", gin.H{
			"Teams": teamData.Teams[:10], // Pass top 10 players
		})
	})

	// Run the server
	r.Run(":8080")
}
