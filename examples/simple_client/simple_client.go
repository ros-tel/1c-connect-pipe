package main

import (
	"flag"
	"log"

	"github.com/ros-tel/1c-connect-pipe"
)

var (
	login = flag.String("login", "", "Set Login")
	debug = flag.Bool("debug", false, "Print debug information on stderr")
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	flag.Parse()

	if *login == "" {
		flag.PrintDefaults()
		log.Fatal("Required parameter \"login\" not set")
	}

	// Инициируем клиент
	client := pipe.InitClient(*login, *debug)

	/*
	 *  \/ Подписываемся на все возможные события: \/
	 */
	// ... измения статусов
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "ServicesClients",
		Object:    "AgentOnlineStatus",
		Initiator: "Incoming",
	})

	// ... текстовых сообщений
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "ServicesClients",
		Object:    "Message",
		Initiator: "Incoming",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "ServicesClients",
		Object:    "Message",
		Initiator: "Self",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "Colleagues",
		Object:    "Message",
		Initiator: "Incoming",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "Colleagues",
		Object:    "Message",
		Initiator: "Self",
	})

	// ... звонков
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "ServicesClients",
		Object:    "Call",
		Initiator: "Incoming",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "ServicesClients",
		Object:    "Call",
		Initiator: "Self",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "Colleagues",
		Object:    "Call",
		Initiator: "Incoming",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "Colleagues",
		Object:    "Call",
		Initiator: "Self",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "Softphone",
		Object:    "Call",
		Initiator: "Incoming",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "Softphone",
		Object:    "Call",
		Initiator: "Self",
	})

	// ... сеансов удаленного доступа
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "ServicesClients",
		Object:    "RemoteAccessSession",
		Initiator: "Incoming",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "ServicesClients",
		Object:    "RemoteAccessSession",
		Initiator: "Self",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "Colleagues",
		Object:    "RemoteAccessSession",
		Initiator: "Incoming",
	})
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "Colleagues",
		Object:    "RemoteAccessSession",
		Initiator: "Self",
	})

	// ... обращений
	client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "ServicesClients",
		Object:    "Treatment",
		Initiator: "Incoming",
	})
	/*
	 *  /\ Подписались на все возможные события /\
	 */

	// В отдельной рутине читаем события из каналов
	go func() {
		for {
			select {
			case e := <-client.Event:
				log.Printf("CommandResult received %+v", e)
			case r := <-client.Result:
				log.Printf("Event received %+v", r)
			}
		}
	}()

	// Запускаем клиент
	client.Start()
}
