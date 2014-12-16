//
// DEPRECIATED. REPLACED WITH worker.go AND workmanager.go
// LEFT IN PROJECT FOR REFERENCE ONLY.
//
/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes, Daniel Sparks
 Created: 8/29/2014

 Module: distribute
 Description: Host fills a channel with games to be played. Clients connect and
 	request a game to play. The client machine submits control over itself to
 	the host.

 TODO:
 	-Break it up so the user types 'host' then 'start'
 	-Take a good look at passing around net.Conn vs pointers. Same with Game.

*******************************************************************************/

package main

//import (
//"bufio"
//"bytes"
//"encoding/gob" //TODO: change to gob
//"fmt"
//"io"
//"net"
//"strings"

//)

// Wrapper to send commands and objects back and forth:
/*
type NetMessage struct {
	Command string
	Object  interface{}
}

func Send(what *NetMessage, where net.Conn) {
	// Note: this works very easily because net.Conn impliments Reader and Writer.

	// Define what types of data we can encode:
	gob.Register(Game{})
	// Encode and write:
	//encoder := gob.NewEncoder(where)
	//err := encoder.Encode(what)

	writer := bufio.NewWriter(where)
	writer.WriteString("testing")
	writer.Flush()

	//if err != nil {
	//	fmt.Println("Error encoding:", err.Error())
	//}
}
*/
/*******************************************************************************

	Hosting:

*******************************************************************************/
/*
// Host object:
type ClientManager struct {
	List       []net.Conn
	Connecting chan net.Conn
	Incoming   chan *NetMessage
	Stop       chan struct{}
}

func NewClientManager() *ClientManager {
	CM := &ClientManager{
		Stop:       make(chan struct{}),
		Incoming:   make(chan *NetMessage),
		Connecting: make(chan net.Conn),
	}
	return CM
}
*/
// Primary function:

/*
func HostTourney(T *Tourney) error {

	// Generate the list of games to pull from:
	T.GameQue = make(chan Game, len(T.GameList))
	T.CompletedGameQue = make(chan Game)
	for i, _ := range T.GameList {
		T.GameQue <- T.GameList[i]
	}

	// Setup the client manager:
	CM := NewClientManager()
	go CM.ReadAndExec() // wait for and process commands from clients

	go CM.ListenAndServe(T)
	// TODO: go PlayLocally()

	// Wait for games to be completed:
	for {
		// TODO: think about deadlocks
		select {

		// Check for a force stop:
		case <-T.Done:
			return nil

		// Check for a completed game to process:
		case g := <-T.CompletedGameQue:
			// place it in the proper location:
			i := g.Round - 1
			T.GameList[i] = g

			// TODO: Broadcast result.
		}

		// Stop condition:
		if T.Complete() {
			// TODO: if the GameQue is empty but there are incomplete games,
			//		 add the noncompleted games to the que again.
			if CM.ForcedStop() {
				close(T.Done) // TODO: this may cause problems with user typing stop?
			}
			break
		}
	}

	return nil
}
*/

/*
// Wait for clients to connect then assign them a game to play.
func (CM *ClientManager) ListenAndServe(T *Tourney) {

	// Listen:
	fmt.Println("\nListening on port 9000...")
	server, err := net.Listen("tcp", ":9000") //TODO: user chosen port.
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer server.Close()

	// Serve:
	for {

		// Wait for a connection:
		conn, err := server.Accept() // TODO: add a timeout
		if err != nil {
			fmt.Println("Client connection error:", err.Error())
			continue
		}

		// Establish the incomming connection:
		go func() {
			CM.Connecting <- conn
		}()
	}
}

// Wait for clients to give commands, then execute them:
func (CM *ClientManager) ReadAndExec() {
	for {
		select {

		// Forced stop from somewhere else in the pipeline:
		case <-CM.Stop:
			return

		// Finish establishing connections to waiting clients:
		case client := <-CM.Connecting:
			fmt.Println(client.RemoteAddr(), "connected.")
			CM.List = append(CM.List, client)
			go CM.Read(&client)

		// Process incoming messages:
		case inc := <-CM.Incoming:
			CM.Exec(inc)
		}
	}
}

// Listen to a particular socket:
func (CM *ClientManager) Read(Client *net.Conn) error {
	// TODO: is this loop taxing? YES!!!!!!!!
	gob.Register(Game{})
	//decoder := gob.NewDecoder(*Client)
	reader := bufio.NewReader(*Client)
	for {
		fmt.Print(".")
		if CM.ForcedStop() {
			return nil
		}
		//message := NetMessage{}
		//err := decoder.Decode(&message)
		s, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				//continue
				//return nil
			}
			fmt.Println("Decoding error:", err.Error())
		}
		//CM.Incoming <- &message
		fmt.Println(s)
	}
	return nil
}

// Host executes a command coming in from a client:
func (CM *ClientManager) Exec(m *NetMessage) error {
	fmt.Println(*m) // temporary
	return nil
}

// Helper:
func (CM *ClientManager) ForcedStop() bool {
	select {
	case <-CM.Stop:
		return true
	default:
	}
	return false
}

*/

/*******************************************************************************

	Client:

*******************************************************************************/

// Establish a connection to the host, then do what the host says:
/*
func ConnectAndPlay(address string) {

		// First connect to the host:
		fmt.Print("\nConnecting to " + address + "... ")
		host, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer host.Close()
		fmt.Println("Success.")

		// Second start playing:

		//testing:
		var g Game
		g.initialize()
		m := NetMessage{Command: "NewGame", Object: g}

*/
//fmt.Println(m.Object)

//fmt.Println("ENCODED:")
//encoder := gob.NewEncoder(os.Stdout)
//encoder.Encode(m)

/*
	var network bytes.Buffer //dummy
	enc := gob.NewEncoder(&network)
	dec := gob.NewDecoder(&network)
	gob.Register(Game{})
	err = enc.Encode(m)
	if err != nil {
		fmt.Println(err.Error())
	}
	var m2 NetMessage
	err = dec.Decode(&m2)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(m2)
*/

//Send(&m, host)
//netEncoder := gob.NewEncoder(host)
//netEncoder.Encode(m)

//}
