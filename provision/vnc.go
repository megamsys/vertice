package provision

import (
	//"encoding/json"
	"fmt"
	kvnc "github.com/kward/go-vnc"
	"golang.org/x/net/context"
	"log"
	"net"
	"time"
	//log "github.com/Sirupsen/logrus"
	//nsqc "github.com/crackcomm/nsqueue/consumer"
	//nsqp "github.com/crackcomm/nsqueue/producer"
	//"github.com/megamsys/vertice/meta"
)

type VncListener struct {
	B <-chan Boxlog
}

func NewVNCListener(a *Box) (*VncListener, error) {
	b := make(chan Boxlog, maxInFlight)

	go func() {
	//	defer close(b)
		nc, err := net.Dial("tcp", "136.243.49.217:6643")
		fmt.Println("*************************")
    fmt.Printf("%#v",nc)
		if err != nil {
			log.Fatalf("Error connecting to VNC host. %v", err)
		}

		// Negotiate connection with the server.
		vcc := kvnc.NewClientConfig("")
		fmt.Println("**********clinent***************")
   fmt.Printf("%#v",vcc)
		vc, err := kvnc.Connect(context.Background(), nc, vcc)
		fmt.Println("*********conn****************")
     fmt.Printf("%#v",vc)
		 fmt.Println("***********conn err*****************")
		 fmt.Println(err)
		if err != nil {
			log.Fatalf("Error negotiating connection to VNC host. %v", err)
		}

		// Periodically request framebuffer updates.
	//go func() {
			w, h := vc.FramebufferWidth(), vc.FramebufferHeight()
			fmt.Println("***************buffer************************************")
     fmt.Println(vc.FramebufferWidth())
      fmt.Println(vc.FramebufferHeight())
			fmt.Println(vc)

			for {
				if err := vc.FramebufferUpdateRequest(kvnc.RFBTrue, 0, 0, w, h); err != nil {
					log.Printf("error requesting framebuffer update: %v", err)
				}
					fmt.Println("*******************sleep")
}
				fmt.Println("*******************sleep")
				time.Sleep(1 * time.Second)

		//}()
   fmt.Println("************server**********")
		// Listen and handle server messages.
		go vc.ListenAndHandle()

		// Process messages coming in on the ServerMessage channel.
		for {
			msg := <-vcc.ServerMessageCh

			switch msg.Type() {
			case kvnc.FramebufferUpdateMsg:
				log.Println("Received FramebufferUpdate message.")
			default:
				log.Printf("Received message type:%v msg:%v\n", msg.Type(), msg)
			}
		}
	}()
	l := VncListener{B: b}
	return &l, nil

}
/*
func (l *VncListener) Close() (err error) {
	l.Stop()
	return nil
}
*/
