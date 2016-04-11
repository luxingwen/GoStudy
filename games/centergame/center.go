package centergame

import (
	"GoStudy/games/myipc"
	"encoding/json"
	"errors"
	"sync"
)

//var  ipc,Server=&CenterServer{}
type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type CenterServer struct {
	servers map[string]myipc.Server
	players []*Player
	rooms   []*Room
	mutex   sync.RWMutex
}

func NewCenterServer() *CenterServer {
	servers := make(map[string]*myipc.Server)
	players := make([]*Player, 0)
	return &CenterServer{servers: servers, players: players}
}
func (server *CenterServer) addPlaer(params string) error {
	player := NewPlayer()
	err := json.Unmarshal([]byte(params), &player)
	if err != nil {
		return err
	}
	server.mutex.Lock()
	defer server.mutex.Unlock()

	server.players = append(server.players, player)
	return nil
}

func (server *CenterServer) removePlayer(params string) error {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	for i, v := range server.players {
		if v.Name == params {
			if len(server.players) == 1 {
				server.players = make([]*Player, 0)
			} else if i == len(server.players)-1 {
				server.players = server.players[:i]
			} else if i == 0 {
				server.players = server.players[i:]
			} else {
				server.players = append(server.players[:i-1], server.players[i+1:]...)
			}
			return nil
		}
	}
	return errors.New("player not found.")
}

func (server *CenterServer) listplayer(params string) (players string, err error) {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	if len(server.players) > 0 {
		b, _ := json.Marshal(server.players)
		players = string(b)
	} else {
		err = errors.New("No player online.")
	}
	return
}

func (server *CenterServer) brodcast(params string) error {
	var message Message
	err := json.Unmarshal([]byte(params), &message)
	if err != nil {
		return err
	}
	server.mutex.Lock()
	defer server.mutex.Unlock()

	if len(server.players) > 0 {
		for _, player := range server.players {
			player.mq <- &message
		}
	} else {
		err = errors.New("No player online.")
	}
	return err
}

func (server *CenterServer) Handle(method, params string) *myipc.Response {
	switch method {
	case "addplayer":
		err := server.addPlaer(params)
		if err != nil {
			return &myipc.Response{Code: err.Error()}
		}
	case "removeplayer":
		err := server.removePlayer(params)
		if err != nil {
			return &myipc.Response{Code: err.Error()}
		}
	case "listplayer":
		err := server.listplayer(params)
		if err != nil {
			return &myipc.Response{Code: err.Error()}
		}
	case "broadcast":
		if err := server.brodcast(params); err != nil {
			return &myipc.Respoense{Code: err.Error()}
		}
	default:
		return &myipc.Response{Code: "404", Body: method + ":" + params}
	}
	return &myipc.Response{Code: "200"}
}
func (server *CenterServer) Name() string {
	return "CenterServer"
}
