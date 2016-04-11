package centergame

import (
	"GoStudy/games/myipc"
	"encoding/json"
	"errors"
)

type CenterClient struct {
	*myipc.IpcClient
}

func (client *CenterClient) AddPlayer(player *Player) (err error) {
	b, err := json.Marshal(*Player)
	if err != nil {
		return
	}
	resp, err := client.Call("addplayer", string(b))
	if err == nil && resp.Code == "200" {
		return
	}
	return
}

func (client *CenterClient) RemovePlayer(name string) (err error) {
	ret, err := client.Call("removeplayer", name)
	if err == nil && ret.Code == "200" {
		return
	}
	return
}

func (client *CenterClient) ListPlayer(params string) (ps []*Player, err error) {
	resp, _ := client.Call("listplayer", params)
	if resp.Code != "200" {
		err = errors.New(resp.Code)
		return
	}
	err = json.Unmarshal([]byte(resp.Body), &ps)
	return
}

func (client *CenterClient) Broadcast(message string) (err error) {
	m := &Message{Content: message}
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	resp, err := client.Call("broadcast", string(b))
	if err == nil && resp.Code == "200" {
		return
	}
	return errors.New(resp.Code)
}
