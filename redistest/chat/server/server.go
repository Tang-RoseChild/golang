package server

import (
	"log"

	"chat/msgproto"
	proto "github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"

	"fmt"
	"net/http"
	"strconv"
)

type MsgType int32

const (
	_             = iota
	LOGIN MsgType = iota
	CONNECT
	DISCONNECT
	MSG_CONTENT
	MSG_CHATROOM
)

var (
	clients      map[int32]*websocket.Conn           // for store all clients
	topicclients map[string]map[*websocket.Conn]bool // for store the same topic clients
	msgCodec     websocket.Codec
)

func init() {
	clients = make(map[int32]*websocket.Conn, 1<<5)
	topicclients = make(map[string]map[*websocket.Conn]bool, 1<<5) // if not init,it will be nil
	msgCodec = websocket.Codec{
		Marshal:   PbMarshal,
		Unmarshal: PbUnmarshal,
	}

}

func Run(addr string) error {
	// welcome web
	http.HandleFunc("/", index)
	// socket entry
	http.HandleFunc("/sock", sock)
	// server run
	return http.ListenAndServe(addr, nil)

}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "welcome")
}

func sock(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(wsserver).ServeHTTP(w, r)
}

func wsserver(ws *websocket.Conn) {

	// read msg
	var (
		msg msgproto.Msg
		err error
		id  int32
		// topic string
	)
	for {
		err = msgCodec.Receive(ws, &msg)
		checkErr("Receive", err)
		if err != nil {
			break
		}
		// add client and topic client
		id = msg.GetId()
		if _, ok := clients[id]; !ok {
			clients[id] = ws
		}

		handleMsg(msg, ws)
	}

	defer func() {
		delete(clients, id)
		ws.Close()
	}()
}

func handleMsg(msg msgproto.Msg, ws *websocket.Conn) {
	log.Printf("msg received : id : %d, type : %d \n", msg.GetId(), msg.GetType())
	switch MsgType(msg.GetType()) {
	case LOGIN:
		// send all clients
		var msgcontent = "clients: \n"
		for key, _ := range clients {
			msgcontent += strconv.Itoa(int(key)) + "\n"
		}
		// pbMsg := &msgproto.Msg{
		// 	Id:      proto.Int32(msg.GetId()),
		// 	Topic:   proto.String(msg.GetTopic()),
		// 	Content: proto.String(msgcontent),
		// 	Type:    proto.Int32(msg.GetType()),
		// }
		pbMsg := PbMsgFactory(msg.GetId(), msg.GetTopic(), msgcontent, msg.GetType())

		err := msgCodec.Send(ws, pbMsg)
		checkErr("Send at LOGIN in handleMsg", err)
		log.Printf("Login send : %s %v\n", msgcontent, clients)
	case CONNECT:
		pbMsg := PbMsgFactory(msg.GetId(), msg.GetTopic(), "connected, start to talk", msg.GetType())

		err := msgCodec.Send(ws, pbMsg)
		checkErr("Send at CONNECT in handleMsg", err)
	case MSG_CONTENT:
		// here topic is another id
		id, err := strconv.Atoi(msg.GetTopic())
		checkErr("Atoi at MSG_CONTENT case of handleMsg", err)
		// find the client in the clients and send the data
		// proto no setter,just get
		toClient := clients[int32(id)]
		// pbMsg := &msgproto.Msg{
		// 	Id:      proto.Int32(int32(id)),
		// 	Topic:   proto.String(strconv.Itoa(int(msg.GetId()))),
		// 	Content: proto.String(msg.GetContent()),
		// 	Type:    proto.Int32(msg.GetType()),
		// }
		pbMsg := PbMsgFactory(int32(id), strconv.Itoa(int(msg.GetId())), msg.GetContent(), msg.GetType())

		err = msgCodec.Send(toClient, pbMsg)
		checkErr("Send at MSG_CONTENT case of handleMsg", err)
	case MSG_CHATROOM:
		topic := msg.GetTopic()
		if v, ok := topicclients[topic]; ok {
			// if not exists
			if _, ok := v[ws]; !ok {
				// add
				v[ws] = true
			}

		} else {
			topicclients[topic] = map[*websocket.Conn]bool{ws: true}
		}

		// check topic exist or not
		if v, ok := topicclients[topic]; ok {
			for c, _ := range v {

				err := msgCodec.Send(c, &msg)
				checkErr("Send at MSG_CHATROOM case of handleMsg", err)
			}

		} else {
			pbMsg := PbMsgFactory(int32(msg.GetId()), strconv.Itoa(int(msg.GetId())), "chatroom doesn't exist", msg.GetType())
			err := msgCodec.Send(clients[msg.GetId()], pbMsg)
			checkErr("Send at MSG_CHATROOM case of handleMsg", err)
		}
		fmt.Println("topicclients : ", topicclients)
		// send back msg
	default:
		log.Printf("msg type not supported : %d \n", msg.GetType())
	}

}

// checkErr: key string and err
func checkErr(key string, err error) {
	if err != nil {
		fmt.Printf("%s err occur : %s ", key, err.Error())
	}
}

/*
// Codec defined in websocket
type Codec struct {
    Marshal   func(v interface{}) (data []byte, payloadType byte, err error)
    Unmarshal func(data []byte, payloadType byte, v interface{}) (err error)
}

// Marshal in protobuf
func Marshal(pb Message) ([]byte, error)

// Unmarshal in protobuf
func Unmarshal(buf []byte, pb Message) error

*/

// pbMarshal : wrap protobuf marshal as websocket marshal
func PbMarshal(v interface{}) (data []byte, payloadType byte, err error) {
	// type assertion
	pb, ok := v.(proto.Message)

	if !ok {
		data, payloadType, err = make([]byte, 0), websocket.BinaryFrame, fmt.Errorf("type wrong, need proto.Message but getting %T ", v)
		return
	}

	data, err = proto.Marshal(pb)
	payloadType = websocket.BinaryFrame
	return
}

// pbUnmarshal : wrap protobuf Unmarshal as websocket marshal
func PbUnmarshal(data []byte, payloadType byte, v interface{}) (err error) {
	// type assertion
	pb, ok := v.(proto.Message)

	if !ok {
		data, payloadType, err = make([]byte, 0), websocket.BinaryFrame, fmt.Errorf("type wrong, need proto.Message but getting %T ", v)
		return
	}

	err = proto.Unmarshal(data, pb)
	return
}

// PbMsgFactory: create msg by given params
func PbMsgFactory(Id int32, Topic string, Content string, Type int32) *msgproto.Msg {
	return &msgproto.Msg{
		Id:      proto.Int32(Id),
		Topic:   proto.String(Topic),
		Content: proto.String(Content),
		Type:    proto.Int32(Type),
	}
}
