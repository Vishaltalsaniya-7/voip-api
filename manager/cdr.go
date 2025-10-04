package manager

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fiorix/go-eventsocket/eventsocket"
	"github.com/vishaltalsaniya-7/voip-api/config"
)

type ESLManager struct {
	config config.FreeSWITCHConfig
}

func NewESLManager(cfg config.FreeSWITCHConfig) *ESLManager {
	return &ESLManager{
		config: cfg,
	}
}

func (e *ESLManager) OriginateCall(caller, callee string) (string, error) {
	addr := fmt.Sprintf("%s:%s", e.config.Host, e.config.Port)
	client, err := eventsocket.Dial(addr, e.config.Password)
	if err != nil {
		return "", fmt.Errorf("failed to connect to FreeSWITCH: %w", err)
	}
	defer client.Close()

	cmd := fmt.Sprintf("originate {origination_caller_id_number=%s}user/%s@172.27.191.2 &bridge(user/%s@172.27.191.2)", caller, caller, callee)
	resp, err := client.Send(fmt.Sprintf("api %s", cmd))
	if err != nil {
		return "", fmt.Errorf("failed to originate call: %w", err)
	}

	callID := e.parseCallID(resp)
	if callID == "" {
		return "", fmt.Errorf("no Call-ID returned from FreeSWITCH")
	}

	return callID, nil
}

func (e *ESLManager) ListenEvents() {
	addr := fmt.Sprintf("%s:%s", e.config.Host, e.config.Port)
	
	for {
		conn, err := eventsocket.Dial(addr, e.config.Password)
		if err != nil {
			log.Printf("Failed to connect to FreeSWITCH for events: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		conn.Send("event plain CHANNEL_HANGUP")
		log.Println("ESL event listener connected")

		for {
			ev, err := conn.ReadEvent()
			if err != nil {
				log.Printf("ESL read error: %v", err)
				conn.Close()
				break
			}

			if ev.Get("Event-Name") == "CHANNEL_HANGUP" {
				go e.handleHangupEvent(ev)
			}
		}

		log.Println("ESL connection lost, reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

func (e *ESLManager) handleHangupEvent(ev *eventsocket.Event) {
	callID := ev.Get("Unique-ID")
	if callID == "" {
		log.Println("No Unique-ID in hangup event")
		return
	}

	durationStr := ev.Get("variable_duration")
	billSecStr := ev.Get("variable_billsec")
	hangupCause := ev.Get("Hangup-Cause")

	var duration, billSec int
	if durationStr != "" {
		fmt.Sscanf(durationStr, "%d", &duration)
	}
	if billSecStr != "" {
		fmt.Sscanf(billSecStr, "%d", &billSec)
	}

	status := "COMPLETED"
	if hangupCause != "NORMAL_CLEARING" {
		status = fmt.Sprintf("FAILED_%s", hangupCause)
	}

	log.Printf("Call %s ended: duration=%d, billsec=%d, status=%s, hangup_cause=%s", 
		callID, duration, billSec, status, hangupCause)
}

func (e *ESLManager) parseCallID(resp *eventsocket.Event) string {
	body := strings.TrimSpace(resp.Body)
	if strings.HasPrefix(body, "+OK") {
		parts := strings.Fields(body)
		if len(parts) >= 2 {
			return strings.TrimSpace(parts[1])
		}
	}
	if body != "" && body != "-ERR" && !strings.HasPrefix(body, "-ERR") {
		return body
	}
	return ""
}




type CallStatus struct {
	UUID          string `json:"uuid"`
	Caller        string `json:"caller"`
	Callee        string `json:"callee"`
	State         string `json:"state"`
	CallDirection string `json:"callDirection"`
	CreatedTime   string `json:"createdTime"`
	AnsweredTime  string `json:"answeredTime"`
	Duration      int    `json:"duration"`
}

func (e *ESLManager) GetCallStatus(uuid string) (*CallStatus, error) {
		addr := fmt.Sprintf("%s:%s", e.config.Host, e.config.Port)
		conn, err := eventsocket.Dial(addr, e.config.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to FreeSWITCH: %v", err)
	}
	defer conn.Close()

cmd := fmt.Sprintf("api uuid_dump %s", uuid)
resp, err := conn.Send(cmd)
	// resp, err := conn.Send(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to send ESL command: %v", err)
	}
	if resp.Body == "" {
		return nil, fmt.Errorf("no call found with uuid: %s", uuid)
	}

	lines := strings.Split(resp.Body, "\n")
	status := &CallStatus{UUID: uuid}

	for _, line := range lines {
		if strings.Contains(line, "Caller-Caller-ID-Number") {
			status.Caller = getValue(line)
		}
		if strings.Contains(line, "Caller-Destination-Number") {
			status.Callee = getValue(line)
		}
		if strings.Contains(line, "Channel-State") {
			status.State = getValue(line)
		}
		if strings.Contains(line, "Call-Direction") {
			status.CallDirection = getValue(line)
		}
		if strings.Contains(line, "Caller-Channel-Created-Time") {
			status.CreatedTime = getValue(line)
		}
		if strings.Contains(line, "Caller-Channel-Answered-Time") {
			status.AnsweredTime = getValue(line)
		}
		if strings.Contains(line, "variable_billsec") {
			fmt.Sscanf(getValue(line), "%d", &status.Duration)
		}
	}

	return status, nil
}

func getValue(line string) string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}