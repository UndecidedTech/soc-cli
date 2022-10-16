package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Status struct {
	Type string `json:"type"`
}

type Event struct {
	Id     int    `json:"id"`
	Slug   string `json:"slug"`
	Status Status
}

type League struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Events []Event
}

type Choice struct {
	Id   int
	Name string
}

type GameObj struct {
	Id    int
	Teams string
}

// responsible for returning an array of choices for our multiselect

func GetJson() []League {
	url := "https://sportscentral.io/new-api/matches?timeZone=240&date=2022-10-15"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
	}

	defer resp.Body.Close()

	body, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Printf("%s", err2.Error())
	}

	leaguesJson := body
	var leagues []League
	json.Unmarshal([]byte(leaguesJson), &leagues)

	// fmt.Printf("League : %+v", leagues[0])
	return leagues
}

func GetActiveLeagues(leagues []League) []Choice {
	var activeLeagues []Choice
	// loop over leagues
	for i, v := range leagues {
		// fmt.Println(leagues[i].Name, leagues[i].Id, v)
		_ = v
		for o, p := range leagues[i].Events {
			// fmt.Println(leagues[i].Events[o].Status.Type, p)
			_ = p
			if leagues[i].Events[o].Status.Type == "inprogress" {
				var newChoice Choice
				newChoice.Id = leagues[i].Id
				newChoice.Name = leagues[i].Name
				activeLeagues = append(activeLeagues, newChoice)
			}
		}
	}
	// fmt.Println(activeLeagues)
	return activeLeagues
}

func GetGames(leagueId int, leagues []League) []Event {
	var games []Event
	for i, v := range leagues {
		_ = v
		if leagues[i].Id == leagueId {
			for o, p := range leagues[i].Events {
				if leagues[i].Events[o].Status.Type == "" {
					_ = p
					var newEvent Event
					newEvent.Id = leagues[i].Events[o].Id
					newEvent.Slug = leagues[i].Events[o].Slug
					newEvent.Status.Type = leagues[i].Events[o].Status.Type
					games = append(games, newEvent)
					fmt.Println(games, newEvent)
				}

			}
		}
	}

	// fmt.Println(games)
	return games
}
