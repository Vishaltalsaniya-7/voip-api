package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fiorix/go-eventsocket/eventsocket"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type CallRequest struct {
	Caller string `json:"caller" binding:"required"`
	Callee string `json:"callee" binding:"required"`
}

type CDR struct {
	XMLCDRUUID             sql.NullString
	DomainUUID             sql.NullString
	ProviderUUID           sql.NullString
	ExtensionUUID          sql.NullString
	SIPCallID              sql.NullString
	DomainName             sql.NullString
	AccountCode            sql.NullString
	Direction              sql.NullString
	DefaultLanguage        sql.NullString
	Context                sql.NullString
	CallerIDName           sql.NullString
	CallerIDNumber         sql.NullString
	CallerDestination      sql.NullString
	SourceNumber           sql.NullString
	DestinationNumber      sql.NullString
	StartEpoch             sql.NullInt64
	StartStamp             time.Time
	AnswerStamp            sql.NullTime
	AnswerEpoch            sql.NullInt64
	EndEpoch               sql.NullInt64
	EndStamp               sql.NullTime
	Duration               sql.NullInt64
	MDuration              sql.NullInt64
	BillSec                sql.NullInt64
	BillMsec               sql.NullInt64
	HoldAccumSeconds       sql.NullInt64
	BridgeUUID             sql.NullString
	ReadCodec              sql.NullString
	ReadRate               sql.NullInt64
	WriteCodec             sql.NullString
	WriteRate              sql.NullInt64
	RemoteMediaIP          sql.NullString
	NetworkAddr            sql.NullString
	RecordPath             sql.NullString
	RecordName             sql.NullString
	RecordLength           sql.NullInt64
	RecordTranscription    sql.NullString
	Leg                    sql.NullString
	OriginatingLegUUID     sql.NullString
	PDDMs                  sql.NullInt64
	RTPAudioInMOS          sql.NullFloat64
	LastApp                sql.NullString
	LastArg                sql.NullString
	VoicemailMessage       sql.NullString
	MissedCall             sql.NullBool
	CallCenterQueueUUID    sql.NullString
	CCSide                 sql.NullString
	CCMemberUUID           sql.NullString
	CCQueueJoinedEpoch     sql.NullInt64
	CCQueue                sql.NullString
	CCMemberSessionUUID    sql.NullString
	CCAgentUUID            sql.NullString
	CCAgent                sql.NullString
	CCAgentType            sql.NullString
	CCAgentBridged         sql.NullBool
	CCQueueAnsweredEpoch   sql.NullInt64
	CCQueueTerminatedEpoch sql.NullInt64
	CCQueueCanceledEpoch   sql.NullInt64
	CCCancelReason         sql.NullString
	CCCause                sql.NullString
	WaitSec                sql.NullInt64
	ConferenceName         sql.NullString
	ConferenceUUID         sql.NullString
	ConferenceMemberID     sql.NullString
	DigitsDialed           sql.NullString
	PINNumber              sql.NullString
	Status                 sql.NullString
	HangupCause            sql.NullString
	HangupCauseQ850        sql.NullString
	SIPHangupDisposition   sql.NullString
	RingGroupUUID          sql.NullString
	IVRMenuUUID            sql.NullString
	CallFlow               sql.NullString
	XML                    sql.NullString
	JSON                   sql.NullString
	InsertDate             sql.NullTime
	InsertUser             sql.NullString
	UpdateDate             sql.NullTime
	UpdateUser             sql.NullString
}
type CDRResponse struct {
	XMLCDRUUID             string
	DomainUUID             string
	ProviderUUID           string
	ExtensionUUID          string
	SIPCallID              string
	DomainName             string
	AccountCode            string
	Direction              string
	DefaultLanguage        string
	Context                string
	CallerIDName           string
	CallerIDNumber         string
	CallerDestination      string
	SourceNumber           string
	DestinationNumber      string
	StartEpoch             int64 // Changed from int to int64
	StartStamp             time.Time
	AnswerStamp            *time.Time
	AnswerEpoch            int64 // Changed from int to int64
	EndEpoch               int64 // Changed from int to int64
	EndStamp               *time.Time
	Duration               int64 // Changed from int to int64
	MDuration              int64 // Changed from int to int64
	BillSec                int64 // Changed from int to int64
	BillMsec               int64 // Changed from int to int64
	HoldAccumSeconds       int64 // Changed from int to int64
	BridgeUUID             string
	ReadCodec              string
	ReadRate               int64 // Changed from int to int64
	WriteCodec             string
	WriteRate              int64 // Changed from int to int64
	RemoteMediaIP          string
	NetworkAddr            string
	RecordPath             string
	RecordName             string
	RecordLength           int64 // Changed from int to int64
	RecordTranscription    string
	Leg                    string
	OriginatingLegUUID     string
	PDDMs                  int64   // Changed from int to int64
	RTPAudioInMOS          float64 // Changed from float to float64
	LastApp                string
	LastArg                string
	VoicemailMessage       string
	MissedCall             bool
	CallCenterQueueUUID    string
	CCSide                 string
	CCMemberUUID           string
	CCQueueJoinedEpoch     int64 // Changed from int to int64
	CCQueue                string
	CCMemberSessionUUID    string
	CCAgentUUID            string
	CCAgent                string
	CCAgentType            string
	CCAgentBridged         bool
	CCQueueAnsweredEpoch   int64 // Changed from int to int64
	CCQueueTerminatedEpoch int64 // Changed from int to int64
	CCQueueCanceledEpoch   int64 // Changed from int to int64
	CCCancelReason         string
	CCCause                string
	WaitSec                int64 // Changed from int to int64
	ConferenceName         string
	ConferenceUUID         string
	ConferenceMemberID     string
	DigitsDialed           string
	PINNumber              string
	Status                 string
	HangupCause            string
	HangupCauseQ850        string
	SIPHangupDisposition   string
	RingGroupUUID          string
	IVRMenuUUID            string
	CallFlow               string
	XML                    string
	JSON                   string
	InsertDate             time.Time
	InsertUser             string
	UpdateDate             time.Time
	UpdateUser             string
}

var db *sql.DB

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Connect to PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	log.Printf("conected !!!!!!!!!!!!!!!!!!!!!!!!")
	// Test the DB connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping PostgreSQL:", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	// Initialize Gin router
	r := gin.Default()

	// POST /call to trigger SIP call
	r.POST("/call", initiateCall)

	// GET /cdrs to retrieve CDRs
	r.GET("/cdrs", getCDRs)

	// Start ESL listener in a goroutine
	go listenESLEvents()

	// Start server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initiateCall(c *gin.Context) {
	var req CallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Connect to FreeSWITCH ESL
	client, err := eventsocket.Dial("127.0.0.1:8021", "ClueCon")
	if err != nil {
		log.Printf("Failed to connect to FreeSWITCH: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to FreeSWITCH"})
		return
	}
	defer client.Close()

	// Originate call
	cmd := fmt.Sprintf("originate {origination_caller_id_number=%s}user/%s@192.168.1.246 &bridge(user/%s@192.168.1.246)", req.Caller, req.Caller, req.Callee)
	resp, err := client.Send(fmt.Sprintf("api %s", cmd))
	if err != nil {
		log.Printf("Failed to originate call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to originate call"})
		return
	}

	// Parse response for Call-ID
	callID := parseCallID(resp)
	if callID == "" {
		log.Printf("No Call-ID returned from FreeSWITCH. Response: %v", resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No Call-ID returned"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"call_id": callID, "status": "Call initiated"})
}

func listenESLEvents() {
	for {
		// Connect to FreeSWITCH for event listening
		conn, err := eventsocket.Dial("127.0.0.1:8021", "ClueCon")
		if err != nil {
			log.Printf("Failed to connect to FreeSWITCH for events: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Subscribe to CHANNEL_HANGUP events
		conn.Send("event plain CHANNEL_HANGUP")
		log.Println("ESL event listener connected")

		// Listen for events
		for {
			ev, err := conn.ReadEvent()
			if err != nil {
				log.Printf("ESL read error: %v", err)
				conn.Close()
				break
			}

			if ev.Get("Event-Name") == "CHANNEL_HANGUP" {
				go handleHangupEvent(ev)
			}
		}

		// Reconnect after a delay
		log.Println("ESL connection lost, reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

func handleHangupEvent(ev *eventsocket.Event) {
	callID := ev.Get("Unique-ID")
	if callID == "" {
		log.Println("No Unique-ID in hangup event")
		return
	}

	// endTime := time.Now()
	durationStr := ev.Get("variable_duration")
	billSecStr := ev.Get("variable_billsec")
	// answerStampStr := ev.Get("variable_answer_stamp")
	hangupCause := ev.Get("Hangup-Cause")

	var duration, billSec int
	// var answerStamp *time.Time

	// Parse duration and billsec
	if durationStr != "" {
		fmt.Sscanf(durationStr, "%d", &duration)
	}
	if billSecStr != "" {
		fmt.Sscanf(billSecStr, "%d", &billSec)
	}

	// Set status based on Hangup Cause
	status := "COMPLETED"
	if hangupCause != "NORMAL_CLEARING" {
		status = fmt.Sprintf("FAILED_%s", hangupCause)
	}

	// Log the hangup event
	log.Printf("Call %s ended: duration=%d, billsec=%d, status=%s, hangup_cause=%s", callID, duration, billSec, status, hangupCause)
}

// func getCDRs(c *gin.Context) {
// 	pageStr := c.DefaultQuery("page", "1")
// 	limitStr := c.DefaultQuery("limit", "10")

// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil || limit < 1 || limit > 100 {
// 		limit = 10 // max page size = 100
// 	}

// 	offset := (page - 1) * limit
// 	var total int
//     err = db.QueryRow(`SELECT COUNT(*) FROM v_xml_cdr`).Scan(&total)
//     if err != nil {
//         log.Printf("Failed to count CDRs: %v", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count CDRs"})
//         return
//     }


// 	query:=`
//         SELECT 
//             xml_cdr_uuid, domain_uuid, provider_uuid, extension_uuid, sip_call_id, domain_name, accountcode, direction, 
//             default_language, context, caller_id_name, caller_id_number, caller_destination, source_number, destination_number, 
//             start_epoch, start_stamp, answer_stamp, answer_epoch, end_epoch, end_stamp, duration, mduration, billsec, billmsec, 
//             hold_accum_seconds, bridge_uuid, read_codec, read_rate, write_codec, write_rate, remote_media_ip, network_addr, 
//             record_path, record_name, record_length, record_transcription, leg, originating_leg_uuid, pdd_ms, rtp_audio_in_mos, 
//             last_app, last_arg, voicemail_message, missed_call, call_center_queue_uuid, cc_side, cc_member_uuid, cc_queue_joined_epoch, 
//             cc_queue, cc_member_session_uuid, cc_agent_uuid, cc_agent, cc_agent_type, cc_agent_bridged, cc_queue_answered_epoch, 
//             cc_queue_terminated_epoch, cc_queue_canceled_epoch, cc_cancel_reason, cc_cause, waitsec, conference_name, conference_uuid, 
//             conference_member_id, digits_dialed, pin_number, status, hangup_cause, hangup_cause_q850, sip_hangup_disposition, 
//             ring_group_uuid, ivr_menu_uuid, call_flow, xml, json, insert_date, insert_user, update_date, update_user
//         FROM v_xml_cdr 
//         ORDER BY start_stamp DESC  
// 		LIMIT ? OFFSET ?`

// 	rows, err := db.Query(query, limit, offset)

// 	if err != nil {
// 		log.Printf("Failed to fetch CDRs: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch CDRs"})
// 		return
// 	}
// 	defer rows.Close()

// 	var cdrs []CDRResponse
// 	scanErrors := 0
// 	rowCount := 0
// 	for rows.Next() {
// 		rowCount++
// 		var cdr CDR
// 		if err := rows.Scan(
// 			&cdr.XMLCDRUUID, &cdr.DomainUUID, &cdr.ProviderUUID, &cdr.ExtensionUUID, &cdr.SIPCallID, &cdr.DomainName,
// 			&cdr.AccountCode, &cdr.Direction, &cdr.DefaultLanguage, &cdr.Context, &cdr.CallerIDName, &cdr.CallerIDNumber,
// 			&cdr.CallerDestination, &cdr.SourceNumber, &cdr.DestinationNumber, &cdr.StartEpoch, &cdr.StartStamp, &cdr.AnswerStamp,
// 			&cdr.AnswerEpoch, &cdr.EndEpoch, &cdr.EndStamp, &cdr.Duration, &cdr.MDuration, &cdr.BillSec, &cdr.BillMsec,
// 			&cdr.HoldAccumSeconds, &cdr.BridgeUUID, &cdr.ReadCodec, &cdr.ReadRate, &cdr.WriteCodec, &cdr.WriteRate,
// 			&cdr.RemoteMediaIP, &cdr.NetworkAddr, &cdr.RecordPath, &cdr.RecordName, &cdr.RecordLength, &cdr.RecordTranscription,
// 			&cdr.Leg, &cdr.OriginatingLegUUID, &cdr.PDDMs, &cdr.RTPAudioInMOS, &cdr.LastApp, &cdr.LastArg, &cdr.VoicemailMessage,
// 			&cdr.MissedCall, &cdr.CallCenterQueueUUID, &cdr.CCSide, &cdr.CCMemberUUID, &cdr.CCQueueJoinedEpoch, &cdr.CCQueue,
// 			&cdr.CCMemberSessionUUID, &cdr.CCAgentUUID, &cdr.CCAgent, &cdr.CCAgentType, &cdr.CCAgentBridged,
// 			&cdr.CCQueueAnsweredEpoch, &cdr.CCQueueTerminatedEpoch, &cdr.CCQueueCanceledEpoch, &cdr.CCCancelReason,
// 			&cdr.CCCause, &cdr.WaitSec, &cdr.ConferenceName, &cdr.ConferenceUUID, &cdr.ConferenceMemberID, &cdr.DigitsDialed,
// 			&cdr.PINNumber, &cdr.Status, &cdr.HangupCause, &cdr.HangupCauseQ850, &cdr.SIPHangupDisposition, &cdr.RingGroupUUID,
// 			&cdr.IVRMenuUUID, &cdr.CallFlow, &cdr.XML, &cdr.JSON, &cdr.InsertDate, &cdr.InsertUser, &cdr.UpdateDate, &cdr.UpdateUser,
// 		); err != nil {
// 			log.Printf("Failed to scan CDR row %d: %v", rowCount, err)
// 			scanErrors++
// 			continue
// 		}

// 		// Initialize CDRResponse
// 		cdrResponse := CDRResponse{
// 			XMLCDRUUID:             cdr.XMLCDRUUID.String,
// 			DomainUUID:             cdr.DomainUUID.String,
// 			ProviderUUID:           cdr.ProviderUUID.String,
// 			ExtensionUUID:          cdr.ExtensionUUID.String,
// 			SIPCallID:              cdr.SIPCallID.String,
// 			DomainName:             cdr.DomainName.String,
// 			AccountCode:            cdr.AccountCode.String,
// 			Direction:              cdr.Direction.String,
// 			DefaultLanguage:        cdr.DefaultLanguage.String,
// 			Context:                cdr.Context.String,
// 			CallerIDName:           cdr.CallerIDName.String,
// 			CallerIDNumber:         cdr.CallerIDNumber.String,
// 			CallerDestination:      cdr.CallerDestination.String,
// 			SourceNumber:           cdr.SourceNumber.String,
// 			DestinationNumber:      cdr.DestinationNumber.String,
// 			StartEpoch:             cdr.StartEpoch.Int64,
// 			StartStamp:             cdr.StartStamp,
// 			AnswerEpoch:            cdr.AnswerEpoch.Int64,
// 			EndEpoch:               cdr.EndEpoch.Int64,
// 			Duration:               cdr.Duration.Int64,
// 			MDuration:              cdr.MDuration.Int64,
// 			BillSec:                cdr.BillSec.Int64,
// 			BillMsec:               cdr.BillMsec.Int64,
// 			HoldAccumSeconds:       cdr.HoldAccumSeconds.Int64,
// 			BridgeUUID:             cdr.BridgeUUID.String,
// 			ReadCodec:              cdr.ReadCodec.String,
// 			ReadRate:               cdr.ReadRate.Int64,
// 			WriteCodec:             cdr.WriteCodec.String,
// 			WriteRate:              cdr.WriteRate.Int64,
// 			RemoteMediaIP:          cdr.RemoteMediaIP.String,
// 			NetworkAddr:            cdr.NetworkAddr.String,
// 			RecordPath:             cdr.RecordPath.String,
// 			RecordName:             cdr.RecordName.String,
// 			RecordLength:           cdr.RecordLength.Int64,
// 			RecordTranscription:    cdr.RecordTranscription.String,
// 			Leg:                    cdr.Leg.String,
// 			OriginatingLegUUID:     cdr.OriginatingLegUUID.String,
// 			PDDMs:                  cdr.PDDMs.Int64,
// 			RTPAudioInMOS:          cdr.RTPAudioInMOS.Float64,
// 			LastApp:                cdr.LastApp.String,
// 			LastArg:                cdr.LastArg.String,
// 			VoicemailMessage:       cdr.VoicemailMessage.String,
// 			MissedCall:             cdr.MissedCall.Bool,
// 			CallCenterQueueUUID:    cdr.CallCenterQueueUUID.String,
// 			CCSide:                 cdr.CCSide.String,
// 			CCMemberUUID:           cdr.CCMemberUUID.String,
// 			CCQueueJoinedEpoch:     cdr.CCQueueJoinedEpoch.Int64,
// 			CCQueue:                cdr.CCQueue.String,
// 			CCMemberSessionUUID:    cdr.CCMemberSessionUUID.String,
// 			CCAgentUUID:            cdr.CCAgentUUID.String,
// 			CCAgent:                cdr.CCAgent.String,
// 			CCAgentType:            cdr.CCAgentType.String,
// 			CCAgentBridged:         cdr.CCAgentBridged.Bool,
// 			CCQueueAnsweredEpoch:   cdr.CCQueueAnsweredEpoch.Int64,
// 			CCQueueTerminatedEpoch: cdr.CCQueueTerminatedEpoch.Int64,
// 			CCQueueCanceledEpoch:   cdr.CCQueueCanceledEpoch.Int64,
// 			CCCancelReason:         cdr.CCCancelReason.String,
// 			CCCause:                cdr.CCCause.String,
// 			WaitSec:                cdr.WaitSec.Int64,
// 			ConferenceName:         cdr.ConferenceName.String,
// 			ConferenceUUID:         cdr.ConferenceUUID.String,
// 			ConferenceMemberID:     cdr.ConferenceMemberID.String,
// 			DigitsDialed:           cdr.DigitsDialed.String,
// 			PINNumber:              cdr.PINNumber.String,
// 			Status:                 cdr.Status.String,
// 			HangupCause:            cdr.HangupCause.String,
// 			HangupCauseQ850:        cdr.HangupCauseQ850.String,
// 			SIPHangupDisposition:   cdr.SIPHangupDisposition.String,
// 			RingGroupUUID:          cdr.RingGroupUUID.String,
// 			IVRMenuUUID:            cdr.IVRMenuUUID.String,
// 			CallFlow:               cdr.CallFlow.String,
// 			XML:                    cdr.XML.String,
// 			JSON:                   cdr.JSON.String,
// 			InsertDate:             cdr.InsertDate.Time,
// 			InsertUser:             cdr.InsertUser.String,
// 			UpdateDate:             cdr.UpdateDate.Time,
// 			UpdateUser:             cdr.UpdateUser.String,
// 		}

// 		// Handle nullable time fields
// 		if cdr.AnswerStamp.Valid {
// 			cdrResponse.AnswerStamp = &cdr.AnswerStamp.Time
// 		} else {
// 			cdrResponse.AnswerStamp = nil
// 		}
// 		if cdr.EndStamp.Valid {
// 			cdrResponse.EndStamp = &cdr.EndStamp.Time
// 		} else {
// 			cdrResponse.EndStamp = nil
// 		}

// 		cdrs = append(cdrs, cdrResponse)
// 	}

// 	log.Printf("Queried %d rows, %d scan errors, %d CDRs successfully scanned", rowCount, scanErrors, len(cdrs))
// 	if err := rows.Err(); err != nil {
// 		log.Printf("Row iteration error: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing CDR rows"})
// 		return
// 	}

// 	if len(cdrs) == 0 {
// 		log.Println("No CDRs found in the database")
// 		c.JSON(http.StatusOK, gin.H{"message": "No CDRs found", "cdrs": []CDRResponse{}})
// 		return
// 	}

// 	log.Printf("Fetched %d CDRs", len(cdrs))
// 	c.JSON(http.StatusOK, gin.H{
//         "cdrs": cdrs,
//         "meta": gin.H{
//             "page":  page,
//             "limit": limit,
//             "total": total,
//         },
//     })
// }

func getCDRs(c *gin.Context) {
    // Parse query params
    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "10")

    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 || limit > 100 {
        limit = 10 // reasonable default & max limit
    }

    offset := (page - 1) * limit

    // Get total CDR count
    var total int
    err = db.QueryRow(`SELECT COUNT(*) FROM v_xml_cdr`).Scan(&total)
    if err != nil {
        log.Printf("Failed to count CDRs: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count CDRs"})
        return
    }

    // Fetch paginated CDRs
    // query := `
    //     SELECT 
    //         xml_cdr_uuid, domain_uuid, provider_uuid, extension_uuid, sip_call_id, domain_name, accountcode, direction, 
    //         default_language, context, caller_id_name, caller_id_number, caller_destination, source_number, destination_number, 
    //         start_epoch, start_stamp, answer_stamp, answer_epoch, end_epoch, end_stamp, duration, mduration, billsec, billmsec, 
    //         hold_accum_seconds, bridge_uuid, read_codec, read_rate, write_codec, write_rate, remote_media_ip, network_addr, 
    //         record_path, record_name, record_length, record_transcription, leg, originating_leg_uuid, pdd_ms, rtp_audio_in_mos, 
    //         last_app, last_arg, voicemail_message, missed_call, call_center_queue_uuid, cc_side, cc_member_uuid, cc_queue_joined_epoch, 
    //         cc_queue, cc_member_session_uuid, cc_agent_uuid, cc_agent, cc_agent_type, cc_agent_bridged, cc_queue_answered_epoch, 
    //         cc_queue_terminated_epoch, cc_queue_canceled_epoch, cc_cancel_reason, cc_cause, waitsec, conference_name, conference_uuid, 
    //         conference_member_id, digits_dialed, pin_number, status, hangup_cause, hangup_cause_q850, sip_hangup_disposition, 
    //         ring_group_uuid, ivr_menu_uuid, call_flow, xml, json, insert_date, insert_user, update_date, update_user
    //     FROM v_xml_cdr 
    //     ORDER BY start_stamp DESC 
    //     LIMIT ? OFFSET ?`

    // rows, err := db.Query(query, limit, offset)
    query := fmt.Sprintf(`
    SELECT 
        xml_cdr_uuid, domain_uuid, provider_uuid, extension_uuid, sip_call_id, domain_name, accountcode, direction, 
        default_language, context, caller_id_name, caller_id_number, caller_destination, source_number, destination_number, 
        start_epoch, start_stamp, answer_stamp, answer_epoch, end_epoch, end_stamp, duration, mduration, billsec, billmsec, 
        hold_accum_seconds, bridge_uuid, read_codec, read_rate, write_codec, write_rate, remote_media_ip, network_addr, 
        record_path, record_name, record_length, record_transcription, leg, originating_leg_uuid, pdd_ms, rtp_audio_in_mos, 
        last_app, last_arg, voicemail_message, missed_call, call_center_queue_uuid, cc_side, cc_member_uuid, cc_queue_joined_epoch, 
        cc_queue, cc_member_session_uuid, cc_agent_uuid, cc_agent, cc_agent_type, cc_agent_bridged, cc_queue_answered_epoch, 
        cc_queue_terminated_epoch, cc_queue_canceled_epoch, cc_cancel_reason, cc_cause, waitsec, conference_name, conference_uuid, 
        conference_member_id, digits_dialed, pin_number, status, hangup_cause, hangup_cause_q850, sip_hangup_disposition, 
        ring_group_uuid, ivr_menu_uuid, call_flow, xml, json, insert_date, insert_user, update_date, update_user
    FROM v_xml_cdr 
    ORDER BY start_stamp DESC 
    LIMIT %d OFFSET %d`, limit, offset)
	rows, err := db.Query(query)

	if err != nil {
        log.Printf("Failed to fetch CDRs: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch CDRs"})
        return
    }
    defer rows.Close()

    var cdrs []CDRResponse
    scanErrors := 0
    rowCount := 0
    for rows.Next() {
        rowCount++
        var cdr CDR
        if err := rows.Scan(
            &cdr.XMLCDRUUID, &cdr.DomainUUID, &cdr.ProviderUUID, &cdr.ExtensionUUID, &cdr.SIPCallID, &cdr.DomainName,
            &cdr.AccountCode, &cdr.Direction, &cdr.DefaultLanguage, &cdr.Context, &cdr.CallerIDName, &cdr.CallerIDNumber,
            &cdr.CallerDestination, &cdr.SourceNumber, &cdr.DestinationNumber, &cdr.StartEpoch, &cdr.StartStamp, &cdr.AnswerStamp,
            &cdr.AnswerEpoch, &cdr.EndEpoch, &cdr.EndStamp, &cdr.Duration, &cdr.MDuration, &cdr.BillSec, &cdr.BillMsec,
            &cdr.HoldAccumSeconds, &cdr.BridgeUUID, &cdr.ReadCodec, &cdr.ReadRate, &cdr.WriteCodec, &cdr.WriteRate,
            &cdr.RemoteMediaIP, &cdr.NetworkAddr, &cdr.RecordPath, &cdr.RecordName, &cdr.RecordLength, &cdr.RecordTranscription,
            &cdr.Leg, &cdr.OriginatingLegUUID, &cdr.PDDMs, &cdr.RTPAudioInMOS, &cdr.LastApp, &cdr.LastArg, &cdr.VoicemailMessage,
            &cdr.MissedCall, &cdr.CallCenterQueueUUID, &cdr.CCSide, &cdr.CCMemberUUID, &cdr.CCQueueJoinedEpoch, &cdr.CCQueue,
            &cdr.CCMemberSessionUUID, &cdr.CCAgentUUID, &cdr.CCAgent, &cdr.CCAgentType, &cdr.CCAgentBridged,
            &cdr.CCQueueAnsweredEpoch, &cdr.CCQueueTerminatedEpoch, &cdr.CCQueueCanceledEpoch, &cdr.CCCancelReason,
            &cdr.CCCause, &cdr.WaitSec, &cdr.ConferenceName, &cdr.ConferenceUUID, &cdr.ConferenceMemberID, &cdr.DigitsDialed,
            &cdr.PINNumber, &cdr.Status, &cdr.HangupCause, &cdr.HangupCauseQ850, &cdr.SIPHangupDisposition, &cdr.RingGroupUUID,
            &cdr.IVRMenuUUID, &cdr.CallFlow, &cdr.XML, &cdr.JSON, &cdr.InsertDate, &cdr.InsertUser, &cdr.UpdateDate, &cdr.UpdateUser,
        ); err != nil {
            log.Printf("Failed to scan CDR row %d: %v", rowCount, err)
            scanErrors++
            continue
        }

        cdrResponse := CDRResponse{
            XMLCDRUUID:          cdr.XMLCDRUUID.String,
            DomainUUID:          cdr.DomainUUID.String,
            ProviderUUID:        cdr.ProviderUUID.String,
            ExtensionUUID:       cdr.ExtensionUUID.String,
            SIPCallID:           cdr.SIPCallID.String,
            DomainName:          cdr.DomainName.String,
            AccountCode:         cdr.AccountCode.String,
            Direction:           cdr.Direction.String,
            DefaultLanguage:     cdr.DefaultLanguage.String,
            Context:             cdr.Context.String,
            CallerIDName:        cdr.CallerIDName.String,
            CallerIDNumber:      cdr.CallerIDNumber.String,
            CallerDestination:   cdr.CallerDestination.String,
            SourceNumber:        cdr.SourceNumber.String,
            DestinationNumber:   cdr.DestinationNumber.String,
            StartEpoch:          cdr.StartEpoch.Int64,
            StartStamp:          cdr.StartStamp,
            AnswerEpoch:         cdr.AnswerEpoch.Int64,
            EndEpoch:            cdr.EndEpoch.Int64,
            Duration:            cdr.Duration.Int64,
            MDuration:           cdr.MDuration.Int64,
            BillSec:             cdr.BillSec.Int64,
            BillMsec:            cdr.BillMsec.Int64,
            HoldAccumSeconds:    cdr.HoldAccumSeconds.Int64,
            BridgeUUID:          cdr.BridgeUUID.String,
            ReadCodec:           cdr.ReadCodec.String,
            ReadRate:            cdr.ReadRate.Int64,
            WriteCodec:          cdr.WriteCodec.String,
            WriteRate:           cdr.WriteRate.Int64,
            RemoteMediaIP:       cdr.RemoteMediaIP.String,
            NetworkAddr:         cdr.NetworkAddr.String,
            RecordPath:          cdr.RecordPath.String,
            RecordName:          cdr.RecordName.String,
            RecordLength:        cdr.RecordLength.Int64,
            RecordTranscription: cdr.RecordTranscription.String,
            Leg:                 cdr.Leg.String,
            OriginatingLegUUID:  cdr.OriginatingLegUUID.String,
            PDDMs:               cdr.PDDMs.Int64,
            RTPAudioInMOS:       cdr.RTPAudioInMOS.Float64,
            LastApp:             cdr.LastApp.String,
            LastArg:             cdr.LastArg.String,
            VoicemailMessage:    cdr.VoicemailMessage.String,
            MissedCall:          cdr.MissedCall.Bool,
            CallCenterQueueUUID: cdr.CallCenterQueueUUID.String,
            CCSide:              cdr.CCSide.String,
            CCMemberUUID:        cdr.CCMemberUUID.String,
            CCQueueJoinedEpoch:  cdr.CCQueueJoinedEpoch.Int64,
            CCQueue:             cdr.CCQueue.String,
            CCMemberSessionUUID: cdr.CCMemberSessionUUID.String,
            CCAgentUUID:         cdr.CCAgentUUID.String,
            CCAgent:             cdr.CCAgent.String,
            CCAgentType:         cdr.CCAgentType.String,
            CCAgentBridged:      cdr.CCAgentBridged.Bool,
            CCQueueAnsweredEpoch:   cdr.CCQueueAnsweredEpoch.Int64,
            CCQueueTerminatedEpoch: cdr.CCQueueTerminatedEpoch.Int64,
            CCQueueCanceledEpoch:   cdr.CCQueueCanceledEpoch.Int64,
            CCCancelReason:      cdr.CCCancelReason.String,
            CCCause:             cdr.CCCause.String,
            WaitSec:             cdr.WaitSec.Int64,
            ConferenceName:      cdr.ConferenceName.String,
            ConferenceUUID:      cdr.ConferenceUUID.String,
            ConferenceMemberID:  cdr.ConferenceMemberID.String,
            DigitsDialed:        cdr.DigitsDialed.String,
            PINNumber:           cdr.PINNumber.String,
            Status:              cdr.Status.String,
            HangupCause:         cdr.HangupCause.String,
            HangupCauseQ850:     cdr.HangupCauseQ850.String,
            SIPHangupDisposition: cdr.SIPHangupDisposition.String,
            RingGroupUUID:       cdr.RingGroupUUID.String,
            IVRMenuUUID:         cdr.IVRMenuUUID.String,
            CallFlow:            cdr.CallFlow.String,
            XML:                 cdr.XML.String,
            JSON:                cdr.JSON.String,
            InsertDate:          cdr.InsertDate.Time,
            InsertUser:          cdr.InsertUser.String,
            UpdateDate:          cdr.UpdateDate.Time,
            UpdateUser:          cdr.UpdateUser.String,
        }

        if cdr.AnswerStamp.Valid {
            cdrResponse.AnswerStamp = &cdr.AnswerStamp.Time
        } else {
            cdrResponse.AnswerStamp = nil
        }

        if cdr.EndStamp.Valid {
            cdrResponse.EndStamp = &cdr.EndStamp.Time
        } else {
            cdrResponse.EndStamp = nil
        }

        cdrs = append(cdrs, cdrResponse)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Row iteration error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing CDR rows"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "cdrs": cdrs,
        "meta": gin.H{
            "page":  page,
            "limit": limit,
            "total": total,
        },
    })
}

func parseCallID(resp *eventsocket.Event) string {
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

// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/fiorix/go-eventsocket/eventsocket"
// 	"github.com/gin-gonic/gin"
// 	_ "github.com/lib/pq"
// 	// _ "github.com/go-sql-driver/mysql"
// )

// type CallRequest struct {
// 	Caller string `json:"caller" binding:"required"`
// 	Callee string `json:"callee" binding:"required"`
// }
// type CDR struct {
//     XMLCDRUUID               sql.NullString `json:"xml_cdr_uuid"`
//     DomainUUID               sql.NullString `json:"domain_uuid"`
//     ProviderUUID             sql.NullString `json:"provider_uuid"`
//     ExtensionUUID            sql.NullString `json:"extension_uuid"`
//     SIPCallID                sql.NullString `json:"sip_call_id"`
//     DomainName               sql.NullString `json:"domain_name"`
//     AccountCode              sql.NullString `json:"accountcode"`
//     Direction                sql.NullString `json:"direction"`
//     DefaultLanguage          sql.NullString `json:"default_language"`
//     Context                  sql.NullString `json:"context"`
//     CallerIDName             sql.NullString `json:"caller_id_name"`
//     CallerIDNumber           sql.NullString `json:"caller_id_number"`
//     CallerDestination        sql.NullString `json:"caller_destination"`
//     SourceNumber             sql.NullString `json:"source_number"`
//     DestinationNumber        sql.NullString `json:"destination_number"`
//     StartEpoch               int            `json:"start_epoch"`
//     StartStamp               time.Time      `json:"start_stamp"`
//     AnswerStamp              *time.Time     `json:"answer_stamp"`
//     AnswerEpoch              int            `json:"answer_epoch"`
//     EndEpoch                 int            `json:"end_epoch"`
//     EndStamp                 *time.Time     `json:"end_stamp"`
//     Duration                 int            `json:"duration"`
//     MDuration                int            `json:"mduration"`
//     BillSec                  int            `json:"billsec"`
//     BillMsec                 int            `json:"billmsec"`
//     HoldAccumSeconds         int            `json:"hold_accum_seconds"`
//     BridgeUUID               sql.NullString `json:"bridge_uuid"`
//     ReadCodec                sql.NullString `json:"read_codec"`
//     ReadRate                 int            `json:"read_rate"`
//     WriteCodec               sql.NullString `json:"write_codec"`
//     WriteRate                int            `json:"write_rate"`
//     RemoteMediaIP            sql.NullString `json:"remote_media_ip"`
//     NetworkAddr              sql.NullString `json:"network_addr"`
//     RecordPath               sql.NullString `json:"record_path"`
//     RecordName               sql.NullString `json:"record_name"`
//     RecordLength             int            `json:"record_length"`
//     RecordTranscription      sql.NullString `json:"record_transcription"`
//     Leg                      sql.NullString `json:"leg"`
//     OriginatingLegUUID       sql.NullString `json:"originating_leg_uuid"`
//     PDDMs                    int            `json:"pdd_ms"`
//     RTPAudioInMOS            float64        `json:"rtp_audio_in_mos"`
//     LastApp                  sql.NullString `json:"last_app"`
//     LastArg                  sql.NullString `json:"last_arg"`
//     VoicemailMessage         sql.NullString `json:"voicemail_message"`
//     MissedCall               bool           `json:"missed_call"`
//     CallCenterQueueUUID      sql.NullString `json:"call_center_queue_uuid"`
//     CCSide                   sql.NullString `json:"cc_side"`
//     CCMemberUUID             sql.NullString `json:"cc_member_uuid"`
//     CCQueueJoinedEpoch       int            `json:"cc_queue_joined_epoch"`
//     CCQueue                  sql.NullString `json:"cc_queue"`
//     CCMemberSessionUUID      sql.NullString `json:"cc_member_session_uuid"`
//     CCAgentUUID              sql.NullString `json:"cc_agent_uuid"`
//     CCAgent                  sql.NullString `json:"cc_agent"`
//     CCAgentType              sql.NullString `json:"cc_agent_type"`
//     CCAgentBridged           bool           `json:"cc_agent_bridged"`
//     CCQueueAnsweredEpoch     int            `json:"cc_queue_answered_epoch"`
//     CCQueueTerminatedEpoch   int            `json:"cc_queue_terminated_epoch"`
//     CCQueueCanceledEpoch     int            `json:"cc_queue_canceled_epoch"`
//     CCCancelReason           sql.NullString `json:"cc_cancel_reason"`
//     CCCause                  sql.NullString `json:"cc_cause"`
//     WaitSec                  int            `json:"waitsec"`
//     ConferenceName           sql.NullString `json:"conference_name"`
//     ConferenceUUID           sql.NullString `json:"conference_uuid"`
//     ConferenceMemberID       sql.NullString `json:"conference_member_id"`
//     DigitsDialed             sql.NullString `json:"digits_dialed"`
//     PINNumber                sql.NullString `json:"pin_number"`
//     Status                   sql.NullString `json:"status"`
//     HangupCause              sql.NullString `json:"hangup_cause"`
//     HangupCauseQ850          sql.NullString `json:"hangup_cause_q850"`
//     SIPHangupDisposition     sql.NullString `json:"sip_hangup_disposition"`
//     RingGroupUUID            sql.NullString `json:"ring_group_uuid"`
//     IVRMenuUUID              sql.NullString `json:"ivr_menu_uuid"`
//     CallFlow                 sql.NullString `json:"call_flow"`
//     XML                      sql.NullString `json:"xml"`
//     JSON                     sql.NullString `json:"json"`
//     InsertDate               time.Time      `json:"insert_date"`
//     InsertUser               sql.NullString `json:"insert_user"`
//     UpdateDate               time.Time      `json:"update_date"`
//     UpdateUser               sql.NullString `json:"update_user"`
// }

// var db *sql.DB

// func main() {
// 	// Connect to MySQL
// 	var err error
// 	db, err = sql.Open("postgres", "host=127.0.0.1 port=5432 user=fusionpbx password=75AxREM3XRhmyVuWcpTt9O8Hg dbname=fusionpbx sslmode=disable")
// 	if err != nil {
// 		log.Fatal("Failed to connect to PostgreSQL:", err)
// 	}

// 	// Test the DB connection
// 	if err = db.Ping(); err != nil {
// 		log.Fatal("Failed to ping PostgreSQL:", err)
// 	}
// 	log.Println("Successfully connected to PostgreSQL")
// 	// Initialize Gin router
// 	r := gin.Default()

// 	// POST /call to trigger SIP call
// 	r.POST("/call", initiateCall)

// 	// GET /cdrs to retrieve CDRs
// 	r.GET("/cdrs", getCDRs)

// 	// Start ESL listener in a goroutine
// 	go listenESLEvents()

// 	// Start server
// 	log.Println("Starting server on :8080")
// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatal("Failed to start server:", err)
// 	}
// }

// func initiateCall(c *gin.Context) {
// 	var req CallRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Connect to FreeSWITCH ESL
// 	client, err := eventsocket.Dial("127.0.0.1:8021", "ClueCon")
// 	if err != nil {
// 		log.Printf("Failed to connect to FreeSWITCH: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to FreeSWITCH"})
// 		return
// 	}
// 	defer client.Close()

// 	// // Originate call
// 	// cmd := fmt.Sprintf("originate {origination_caller_id_number=%s}user/%s &bridge(user/%s)",
// 	// 	req.Caller, req.Caller, req.Callee)
// 	// cmd := fmt.Sprintf("originate {origination_caller_id_number=%s}sofia/internal/%s@192.168.1.246 &bridge(sofia/internal/%s@192.168.1.246)",
// 	// req.Caller, req.Caller, req.Callee)

// 	cmd := fmt.Sprintf("originate {origination_caller_id_number=%s}user/%s@192.168.1.246 &bridge(user/%s@192.168.1.246)", req.Caller, req.Caller, req.Callee)

// 	resp, err := client.Send(fmt.Sprintf("api %s", cmd))
// 	if err != nil {
// 		log.Printf("Failed to originate call: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to originate call"})
// 		return
// 	}

// 	// Parse response for Call-ID
// 	callID := parseCallID(resp)
// 	if callID == "" {
// 		log.Printf("No Call-ID returned from FreeSWITCH. Response: %v", resp.Body)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "No Call-ID returned"})
// 		return
// 	}

// 	// Store initial CDR
// 	_, err = db.Exec("INSERT INTO cdrs (call_id, caller, callee, start_time, status) VALUES (?, ?, ?, ?, ?)",
// 		callID, req.Caller, req.Callee, time.Now(), "INITIATED")
// 	if err != nil {
// 		log.Println("Failed to store CDR:", err)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"call_id": callID, "status": "Call initiated"})
// }

// func listenESLEvents() {
// 	for {
// 		// Connect to FreeSWITCH for event listening
// 		conn, err := eventsocket.Dial("127.0.0.1:8021", "ClueCon")
// 		if err != nil {
// 			log.Printf("Failed to connect to FreeSWITCH for events: %v", err)
// 			time.Sleep(5 * time.Second)
// 			continue
// 		}

// 		// Subscribe to CHANNEL_HANGUP events
// 		conn.Send("event plain CHANNEL_HANGUP")

// 		log.Println("ESL event listener connected")

// 		// Listen for events
// 		for {
// 			ev, err := conn.ReadEvent()
// 			if err != nil {
// 				log.Printf("ESL read error: %v", err)
// 				conn.Close()
// 				break
// 			}

// 			if ev.Get("Event-Name") == "CHANNEL_HANGUP" {
// 				go handleHangupEvent(ev)
// 			}
// 		}

// 		// Reconnect after a delay
// 		log.Println("ESL connection lost, reconnecting in 5 seconds...")
// 		time.Sleep(5 * time.Second)
// 	}
// }

// func handleHangupEvent(ev *eventsocket.Event) {
// 	callID := ev.Get("Unique-ID")
// 	if callID == "" {
// 		log.Println("No Unique-ID in hangup event")
// 		return
// 	}

// 	endTime := time.Now()

// 	durationStr := ev.Get("variable_duration")
// 	billSecStr := ev.Get("variable_billsec")
// 	answerStampStr := ev.Get("variable_answer_stamp")
// 	hangupCause := ev.Get("Hangup-Cause")

// 	var duration, billSec int
// 	var answerStamp *time.Time

// 	// Parse duration and billsec
// 	if durationStr != "" {
// 		fmt.Sscanf(durationStr, "%d", &duration)
// 	}
// 	if billSecStr != "" {
// 		fmt.Sscanf(billSecStr, "%d", &billSec)
// 	}

// 	// Parse answer_stamp
// 	if answerStampStr != "" {
// 		t, err := time.Parse("2006-01-02 15:04:05", answerStampStr)
// 		if err == nil {
// 			answerStamp = &t
// 		}
// 	}

// 	// Set status based on Hangup Cause
// 	status := "COMPLETED"
// 	if hangupCause != "NORMAL_CLEARING" {
// 		status = fmt.Sprintf("FAILED_%s", hangupCause)
// 	}

// 	// Update the CDR row in the database
// 	_, err := db.Exec(`
// 		UPDATE cdrs
// 		SET end_stamp = ?, end_epoch = ?, duration = ?, billsec = ?, answer_stamp = ?, hangup_cause = ?, status = ?
// 		WHERE uuid = ?`,
// 		endTime, int(endTime.Unix()), duration, billSec, answerStamp, hangupCause, status, callID)

// 	if err != nil {
// 		log.Printf("Failed to update CDR for call %s: %v", callID, err)
// 	} else {
// 		log.Printf("Updated CDR for call %s: duration=%d, billsec=%d, status=%s", callID, duration, billSec, status)
// 	}
// }

// func getCDRs(c *gin.Context) {
//     rows, err := db.Query(`
//         SELECT
//             xml_cdr_uuid, domain_uuid, provider_uuid, extension_uuid, sip_call_id, domain_name, accountcode, direction,
//             default_language, context, caller_id_name, caller_id_number, caller_destination, source_number, destination_number,
//             start_epoch, start_stamp, answer_stamp, answer_epoch, end_epoch, end_stamp, duration, mduration, billsec, billmsec,
//             hold_accum_seconds, bridge_uuid, read_codec, read_rate, write_codec, write_rate, remote_media_ip, network_addr,
//             record_path, record_name, record_length, record_transcription, leg, originating_leg_uuid, pdd_ms, rtp_audio_in_mos,
//             last_app, last_arg, voicemail_message, missed_call, call_center_queue_uuid, cc_side, cc_member_uuid, cc_queue_joined_epoch,
//             cc_queue, cc_member_session_uuid, cc_agent_uuid, cc_agent, cc_agent_type, cc_agent_bridged, cc_queue_answered_epoch,
//             cc_queue_terminated_epoch, cc_queue_canceled_epoch, cc_cancel_reason, cc_cause, waitsec, conference_name, conference_uuid,
//             conference_member_id, digits_dialed, pin_number, status, hangup_cause, hangup_cause_q850, sip_hangup_disposition,
//             ring_group_uuid, ivr_menu_uuid, call_flow, xml, json, insert_date, insert_user, update_date, update_user
//         FROM v_xml_cdr
//         ORDER BY start_stamp DESC`)
//     if err != nil {
//         log.Printf("Failed to fetch CDRs: %v", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch CDRs"})
//         return
//     }
//     defer rows.Close()

//     var cdrs []CDR
//     for rows.Next() {
//         var cdr CDR
//         var answerStamp, endStamp sql.NullTime

//         if err := rows.Scan(
//             &cdr.XMLCDRUUID, &cdr.DomainUUID, &cdr.ProviderUUID, &cdr.ExtensionUUID, &cdr.SIPCallID, &cdr.DomainName,
//             &cdr.AccountCode, &cdr.Direction, &cdr.DefaultLanguage, &cdr.Context, &cdr.CallerIDName, &cdr.CallerIDNumber,
//             &cdr.CallerDestination, &cdr.SourceNumber, &cdr.DestinationNumber, &cdr.StartEpoch, &cdr.StartStamp, &answerStamp,
//             &cdr.AnswerEpoch, &cdr.EndEpoch, &endStamp, &cdr.Duration, &cdr.MDuration, &cdr.BillSec, &cdr.BillMsec,
//             &cdr.HoldAccumSeconds, &cdr.BridgeUUID, &cdr.ReadCodec, &cdr.ReadRate, &cdr.WriteCodec, &cdr.WriteRate,
//             &cdr.RemoteMediaIP, &cdr.NetworkAddr, &cdr.RecordPath, &cdr.RecordName, &cdr.RecordLength, &cdr.RecordTranscription,
//             &cdr.Leg, &cdr.OriginatingLegUUID, &cdr.PDDMs, &cdr.RTPAudioInMOS, &cdr.LastApp, &cdr.LastArg, &cdr.VoicemailMessage,
//             &cdr.MissedCall, &cdr.CallCenterQueueUUID, &cdr.CCSide, &cdr.CCMemberUUID, &cdr.CCQueueJoinedEpoch, &cdr.CCQueue,
//             &cdr.CCMemberSessionUUID, &cdr.CCAgentUUID, &cdr.CCAgent, &cdr.CCAgentType, &cdr.CCAgentBridged,
//             &cdr.CCQueueAnsweredEpoch, &cdr.CCQueueTerminatedEpoch, &cdr.CCQueueCanceledEpoch, &cdr.CCCancelReason,
//             &cdr.CCCause, &cdr.WaitSec, &cdr.ConferenceName, &cdr.ConferenceUUID, &cdr.ConferenceMemberID, &cdr.DigitsDialed,
//             &cdr.PINNumber, &cdr.Status, &cdr.HangupCause, &cdr.HangupCauseQ850, &cdr.SIPHangupDisposition, &cdr.RingGroupUUID,
//             &cdr.IVRMenuUUID, &cdr.CallFlow, &cdr.XML, &cdr.JSON, &cdr.InsertDate, &cdr.InsertUser, &cdr.UpdateDate, &cdr.UpdateUser); err != nil {
//             log.Printf("Failed to scan CDR: %v", err)
//             continue
//         }

//         if answerStamp.Valid {
//             cdr.AnswerStamp = &answerStamp.Time
//         }
//         if endStamp.Valid {
//             cdr.EndStamp = &endStamp.Time
//         }

//         cdrs = append(cdrs, cdr)
//     }

//     if err := rows.Err(); err != nil {
//         log.Printf("Row iteration error: %v", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing CDR rows"})
//         return
//     }

//     if len(cdrs) == 0 {
//         log.Println("No CDRs found in the database")
//         c.JSON(http.StatusOK, gin.H{"message": "No CDRs found", "cdrs": []CDR{}})
//         return
//     }

//     log.Printf("Fetched %d CDRs", len(cdrs))
//     c.JSON(http.StatusOK, cdrs)
// }

// func parseCallID(resp *eventsocket.Event) string {
// 	// The response body contains the UUID
// 	body := strings.TrimSpace(resp.Body)

// 	// FreeSWITCH returns "+OK <uuid>" for successful originate
// 	if strings.HasPrefix(body, "+OK") {
// 		parts := strings.Fields(body)
// 		if len(parts) >= 2 {
// 			return strings.TrimSpace(parts[1])
// 		}
// 	}

// 	// Sometimes it just returns the UUID directly
// 	if body != "" && body != "-ERR" && !strings.HasPrefix(body, "-ERR") {
// 		return body
// 	}

// 	return ""
// }
