package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"syscall"

	"github.com/pions/webrtc/pkg/media"
	"golang.org/x/net/websocket"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func listenUnix() {
	path := "/tmp/opus"
	addr, err := net.ResolveUnixAddr("unixgram", path)

	// unlink if exists
	err = syscall.Unlink(path)
	if err != nil {
		// not really important if it fails
		log.Println("Unlink()", err)
	}

	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//simple read
	buf := make([]byte, 4000)

	for {
		n, _, err := conn.ReadFrom(buf)
		sample := media.RTCSample{Data: buf[0:n], Samples: 960}
		log.Println("read: ", n, "sample: ", sample)
		if err != nil {
			continue
		}
	}
}

// In a real project, these would be defined in a common file
type Args struct {
	A int
	B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	log.Println("listening...")

	rpc.Register(new(Arith))

	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)
	http.Handle("/ws", websocket.Handler(serve))
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", nil))
}

func serve(ws *websocket.Conn) {
	log.Printf("Handler starting")
	jsonrpc.ServeConn(ws)
	log.Printf("Handler exiting")
}

/*
	// Register codecs
	// Opus has to be 48000, 2 other rates and channels are signaled within Opus.
	webrtc.RegisterCodec(webrtc.NewRTCRtpOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000, 2))
	// Prepare the configuration
	config := webrtc.RTCConfiguration{
		IceServers: []webrtc.RTCIceServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.New(config)
	check(err)

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange = func(connectionState ice.ConnectionState) {
		log.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	}

	// Create a audio track
	opusTrack, err := peerConnection.NewRTCTrack(webrtc.DefaultPayloadTypeOpus, "audio", "SDRserver")
	check(err)

	_, err = peerConnection.AddTrack(opusTrack)
	check(err)

	// Channels to handle incoming offers and outgoing answers
	offerChan := make(chan webrtc.RTCSessionDescription)
	answerChan := make(chan webrtc.RTCSessionDescription)
	// Offer HTTP handler
	http.HandleFunc("/offer", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		//var offer webrtc.RTCSessionDescription
		// Decode r.Body into offer
		offer := webrtc.RTCSessionDescription{}
		err := json.NewDecoder(r.Body).Decode(&offer)
		check(err)

		log.Println("got offer: ", offer)

		offerChan <- offer // send the offer to the
		answer := <-answerChan

		err = json.NewEncoder(w).Encode(answer)
		check(err)

	})

	// Static files
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	go http.ListenAndServe(":7373", nil)

	// Loop
	for {
		log.Println("I'm waiting for an offer")

		// Wait for the remote SessionDescription
		offer := <-offerChan

		err = peerConnection.SetRemoteDescription(offer)
		check(err)

		// Sets the LocalDescription, and starts our UDP listeners
		answer, err := peerConnection.CreateAnswer(nil)
		check(err)

		// Send the answer
		answerChan <- answer
	}
	/*
			go func() {
			// Get packets from gnuradio
			pc, err := net.ListenPacket("udp", ":1234")
			if err != nil {
				log.Fatal(err)
			}
			defer pc.Close()

			buf := make([]byte, 4000)
			for {

				n, addr, err := pc.ReadFrom(buf)
				if err != nil {
					continue
				}
				log.Println("read: ", n, "from: ", addr)

				sample := media.RTCSample{Data: buf[0:n], Samples: 960}
				opusTrack.Samples <- sample
			}
		}()
		}
*/
