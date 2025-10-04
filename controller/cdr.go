package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vishaltalsaniya-7/voip-api/models"
	"github.com/vishaltalsaniya-7/voip-api/response"
)

type CDRController struct {
	db *sql.DB
}

func NewCDRController(db *sql.DB) *CDRController {
	return &CDRController{
		db: db,
	}
}

func (cdc *CDRController) GetCDRs(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int
	err = cdc.db.QueryRow(`SELECT COUNT(*) FROM v_xml_cdr`).Scan(&total)
	if err != nil {
		log.Printf("Failed to count CDRs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count CDRs"})
		return
	}

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

	rows, err := cdc.db.Query(query)
	if err != nil {
		log.Printf("Failed to fetch CDRs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch CDRs"})
		return
	}
	defer rows.Close()

	var cdrs []response.CDRResponse
	for rows.Next() {
		var cdr models.CDR
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
			log.Printf("Failed to scan CDR row: %v", err)
			continue
		}

		cdrResponse := cdc.mapCDRToResponse(cdr)
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

func (cdc *CDRController) mapCDRToResponse(cdr models.CDR) response.CDRResponse {
	resp := response.CDRResponse{
		XMLCDRUUID:             cdr.XMLCDRUUID.String,
		DomainUUID:             cdr.DomainUUID.String,
		ProviderUUID:           cdr.ProviderUUID.String,
		ExtensionUUID:          cdr.ExtensionUUID.String,
		SIPCallID:              cdr.SIPCallID.String,
		DomainName:             cdr.DomainName.String,
		AccountCode:            cdr.AccountCode.String,
		Direction:              cdr.Direction.String,
		DefaultLanguage:        cdr.DefaultLanguage.String,
		Context:                cdr.Context.String,
		CallerIDName:           cdr.CallerIDName.String,
		CallerIDNumber:         cdr.CallerIDNumber.String,
		CallerDestination:      cdr.CallerDestination.String,
		SourceNumber:           cdr.SourceNumber.String,
		DestinationNumber:      cdr.DestinationNumber.String,
		StartEpoch:             cdr.StartEpoch.Int64,
		StartStamp:             cdr.StartStamp,
		AnswerEpoch:            cdr.AnswerEpoch.Int64,
		EndEpoch:               cdr.EndEpoch.Int64,
		Duration:               cdr.Duration.Int64,
		MDuration:              cdr.MDuration.Int64,
		BillSec:                cdr.BillSec.Int64,
		BillMsec:               cdr.BillMsec.Int64,
		HoldAccumSeconds:       cdr.HoldAccumSeconds.Int64,
		BridgeUUID:             cdr.BridgeUUID.String,
		ReadCodec:              cdr.ReadCodec.String,
		ReadRate:               cdr.ReadRate.Int64,
		WriteCodec:             cdr.WriteCodec.String,
		WriteRate:              cdr.WriteRate.Int64,
		RemoteMediaIP:          cdr.RemoteMediaIP.String,
		NetworkAddr:            cdr.NetworkAddr.String,
		RecordPath:             cdr.RecordPath.String,
		RecordName:             cdr.RecordName.String,
		RecordLength:           cdr.RecordLength.Int64,
		RecordTranscription:    cdr.RecordTranscription.String,
		Leg:                    cdr.Leg.String,
		OriginatingLegUUID:     cdr.OriginatingLegUUID.String,
		PDDMs:                  cdr.PDDMs.Int64,
		RTPAudioInMOS:          cdr.RTPAudioInMOS.Float64,
		LastApp:                cdr.LastApp.String,
		LastArg:                cdr.LastArg.String,
		VoicemailMessage:       cdr.VoicemailMessage.String,
		MissedCall:             cdr.MissedCall.Bool,
		CallCenterQueueUUID:    cdr.CallCenterQueueUUID.String,
		CCSide:                 cdr.CCSide.String,
		CCMemberUUID:           cdr.CCMemberUUID.String,
		CCQueueJoinedEpoch:     cdr.CCQueueJoinedEpoch.Int64,
		CCQueue:                cdr.CCQueue.String,
		CCMemberSessionUUID:    cdr.CCMemberSessionUUID.String,
		CCAgentUUID:            cdr.CCAgentUUID.String,
		CCAgent:                cdr.CCAgent.String,
		CCAgentType:            cdr.CCAgentType.String,
		CCAgentBridged:         cdr.CCAgentBridged.Bool,
		CCQueueAnsweredEpoch:   cdr.CCQueueAnsweredEpoch.Int64,
		CCQueueTerminatedEpoch: cdr.CCQueueTerminatedEpoch.Int64,
		CCQueueCanceledEpoch:   cdr.CCQueueCanceledEpoch.Int64,
		CCCancelReason:         cdr.CCCancelReason.String,
		CCCause:                cdr.CCCause.String,
		WaitSec:                cdr.WaitSec.Int64,
		ConferenceName:         cdr.ConferenceName.String,
		ConferenceUUID:         cdr.ConferenceUUID.String,
		ConferenceMemberID:     cdr.ConferenceMemberID.String,
		DigitsDialed:           cdr.DigitsDialed.String,
		PINNumber:              cdr.PINNumber.String,
		Status:                 cdr.Status.String,
		HangupCause:            cdr.HangupCause.String,
		HangupCauseQ850:        cdr.HangupCauseQ850.String,
		SIPHangupDisposition:   cdr.SIPHangupDisposition.String,
		RingGroupUUID:          cdr.RingGroupUUID.String,
		IVRMenuUUID:            cdr.IVRMenuUUID.String,
		CallFlow:               cdr.CallFlow.String,
		XML:                    cdr.XML.String,
		JSON:                   cdr.JSON.String,
		InsertDate:             cdr.InsertDate.Time,
		InsertUser:             cdr.InsertUser.String,
		UpdateDate:             cdr.UpdateDate.Time,
		UpdateUser:             cdr.UpdateUser.String,
	}

	if cdr.AnswerStamp.Valid {
		resp.AnswerStamp = &cdr.AnswerStamp.Time
	}

	if cdr.EndStamp.Valid {
		resp.EndStamp = &cdr.EndStamp.Time
	}

	return resp
}
