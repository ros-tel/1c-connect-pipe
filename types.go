package pipe

import (
	"encoding/xml"
	"net"
)

type (
	Client struct {
		conn    net.Conn
		login   string
		debug   bool
		Event   chan *Event
		command chan *Command
		Result  chan *CommandResult
	}

	Command struct {
		XMLName xml.Name `xml:"Command"`

		Action    string `xml:"Action,attr"`
		Mode      string `xml:"Mode,attr"`
		ID        string `xml:"ID,attr,omitempty"`
		Object    string `xml:"Object,attr,omitempty"`
		Initiator string `xml:"Initiator,attr,omitempty"`

		ColleagueID  string `xml:"ColleagueID,omitempty"`
		ServiceID    string `xml:"ServiceID,omitempty"`
		UserID       string `xml:"UserID,omitempty"`
		MessageBody  string `xml:"MessageBody,omitempty"`
		Status       string `xml:"Status,omitempty"`
		CallTo       string `xml:"CallTo,omitempty"`
		TreatmentID  string `xml:"TreatmentID,omitempty"`
		SpecialistID string `xml:"SpecialistID,omitempty"`
		Citation     string `xml:"Citation,omitempty"`
	}

	CommandResult struct {
		Action string `xml:"Action,attr"`
		ID     string `xml:"ID,attr"`
		Result string `xml:"Result"`
	}

	Event struct {
		Mode      string `xml:"Mode,attr"`
		Object    string `xml:"Object,attr"`
		Initiator string `xml:"Initiator,attr"`
		Time      string `xml:"Time,attr"`

		CallID      string `xml:"CallID"`
		StartTime   string `xml:"StartTime"`
		CallFrom    string `xml:"CallFrom"`
		State       string `xml:"State"`
		ColleagueID string `xml:"ColleagueID"`
		AcceptTime  string `xml:"AcceptTime"`
		BillSec     string `xml:"BillSec"`
		CallTo      string `xml:"CallTo"`
		Duration    string `xml:"Duration"`
		EndTime     string `xml:"EndTime"`
		CallResult  string `xml:"CallResult"`

		Sended      string `xml:"Sended"`
		MessageID   string `xml:"MessageID"`
		MessageBody string `xml:"MessageBody"`
		AuthorID    string `xml:"AuthorID"`

		Status        string `xml:"Status"`
		UserID        string `xml:"UserID"`
		ClientID      string `xml:"ClientID"`
		CounterpartID string `xml:"CounterpartID"`

		SessionID string `xml:"SessionID"`

		FileDownloadedFiles []EventFile `xml:"DownloadedFiles>File"`
		FileUploadedFiles   []EventFile `xml:"UploadedFiles>File"`
	}

	EventFile struct {
		StartTime string `xml:"StartTime,attr"`
		Duration  string `xml:"Duration,attr"`
		Name      string `xml:"Name,attr"`
		SourceDir string `xml:"SourceDir,attr"`
		DestinDir string `xml:"DestinDir,attr"`
		Size      string `xml:"Size,attr"`
	}
)
