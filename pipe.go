package pipe

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net"
	"time"
	"unicode/utf16"

	"github.com/Microsoft/go-winio"
)

// Инициирует структуру клиента
func InitClient(login string, debug bool) *Client {
	c := new(Client)
	c.login = login
	c.debug = debug

	c.command = make(chan *Command, 100)
	c.Result = make(chan *CommandResult, 100)
	c.Event = make(chan *Event, 1000)

	return c
}

// Запускает клиент
func (c *Client) Start() {
	// Ждем пока подключимся
	c.waitConnect()

	// Цикл отправки команд из очереди
	go func() {
		for {
			select {
			case command := <-c.command:
				err := c.sendCommand(command)
				if err != nil {
					log.Fatalf("SendCommand error %s", err)
				}
			}
		}
	}()

	// Цикл получения событий
	for {
		val, err := next(c.conn, c.debug)
		if err != nil {
			log.Printf("Bad message received %+v %+v", val, err)
			continue
		}
		switch val.(type) {
		case *CommandResult:
			if c.debug {
				log.Printf("CommandResult received %+v", val)
			}
			if len(c.Result) < cap(c.Result) {
				c.Result <- val.(*CommandResult)
				continue
			}
			log.Println(`Channel "result" overloaded`)
		case *Event:
			if c.debug {
				log.Printf("Event received %+v", val)
			}
			if len(c.Event) < cap(c.Event) {
				c.Event <- val.(*Event)
				continue
			}
			log.Println(`Channel "event" overloaded`)
		}
	}
}

// Ожидает некоторое время подключения к pipe (агент может запускаться в это время)
func (c *Client) waitConnect() {
	for i := 0; i < 10; i++ {
		conn, err := winio.DialPipe(`\\.\pipe\BuhphoneAgentAPI2_`+c.login, nil)
		if err != nil {
			log.Printf("Connect error %+v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		c.conn = conn

		log.Println("Connect established")
		return
	}
	log.Fatalln("Not connect to PIPE")
}

// Добавляет команду в очередь
func (c *Client) SendCommand(command Command) {
	if len(c.command) < cap(c.command) {
		c.command <- &command
		return
	}
	log.Println(`Channel "command" overloaded`)
}

// Маршалит и отправлет команду в сокет
func (c *Client) sendCommand(command *Command) error {
	bc, err := xml.Marshal(command)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(bc)
	return err
}

// Два первых байта в int (порядок байтов Little-Endian)
func byteToInt(d []byte) (int, error) {
	var res int
	res = int(uint(d[0]) + uint(d[1])<<8)
	if res < 50 || res > 65535 {
		return 0, errors.New("Bad Packet")
	}
	return res, nil
}

// Сканирует поток токенов XML, чтобы найти следующий стартовый элемент
func nextStart(p *xml.Decoder) (xml.StartElement, error) {
	for {
		t, err := p.Token()
		if err != nil && err != io.EOF || t == nil {
			return xml.StartElement{}, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			return t, nil
		}
	}
}

// Получает следующее сообщение и декодирует в нужный тип
func next(conn net.Conn, debug bool) (interface{}, error) {
	// Выделяем буфер и читаем первые 4 байта
	bytePacketSize := make([]byte, 4)
	n, err := conn.Read(bytePacketSize)
	if err != nil || n != 4 {
		log.Fatalf("%+v", err)
	}
	/*
		if debug {
			log.Printf("Read %x -> %s\n", bytePacketSize, bytePacketSize)
		}
	*/

	// Первые два байта кодируют длину сообщения, два последующих всегда 0x00
	packetSize, err := byteToInt(bytePacketSize)
	if err != nil {
		log.Printf("Read PacketLen", err)
		return nil, errors.New("Read PacketLen")

	}

	if packetSize%2 != 0 {
		return nil, errors.New("Error PacketLen")
	}

	// Выделяем буфер и читаем всё сообщение
	body := make([]uint16, packetSize/2)
	err = binary.Read(conn, binary.LittleEndian, &body)
	if err != nil {
		log.Println("failed to Read:", err, len(body))
		return nil, errors.New("Read bad data")
	}

	// Сообщения имеют кодировку UTF-16
	byteMessage := utf16.Decode(body)

	if debug {
		log.Printf("Read %x -> %s\n", byteMessage, string(byteMessage))
	}

	// Вешаем XML-декодер (bytes.NewReader принимает []byte, а byteMessage у нас []rune, соотв. конвертим)
	p := xml.NewDecoder(bytes.NewReader([]byte(string(byteMessage))))

	// Получаем первый элемент XML для определения типа сообщения
	se, err := nextStart(p)
	if err != nil {
		return nil, err
	}

	if debug {
		log.Printf("Read XML %+v", se)
	}

	// Выясняем тип события
	var nv interface{}
	switch se.Name.Local {
	case "CommandResult":
		nv = &CommandResult{}
	case "Event":
		nv = &Event{}
	default:
		return nil, errors.New("Read unexpected message " + " <" + se.Name.Local + "/>")
	}

	// Анмаршалим в нужную структуру
	if err = p.DecodeElement(nv, &se); err != nil {
		return nil, err
	}

	return nv, err
}
