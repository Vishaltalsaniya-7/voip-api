package response

import "time"

type CDRResponse struct {
	XMLCDRUUID             string     `json:"xml_cdr_uuid"`
	DomainUUID             string     `json:"domain_uuid"`
	ProviderUUID           string     `json:"provider_uuid"`
	ExtensionUUID          string     `json:"extension_uuid"`
	SIPCallID              string     `json:"sip_call_id"`
	DomainName             string     `json:"domain_name"`
	AccountCode            string     `json:"account_code"`
	Direction              string     `json:"direction"`
	DefaultLanguage        string     `json:"default_language"`
	Context                string     `json:"context"`
	CallerIDName           string     `json:"caller_id_name"`
	CallerIDNumber         string     `json:"caller_id_number"`
	CallerDestination      string     `json:"caller_destination"`
	SourceNumber           string     `json:"source_number"`
	DestinationNumber      string     `json:"destination_number"`
	StartEpoch             int64      `json:"start_epoch"`
	StartStamp             time.Time  `json:"start_stamp"`
	AnswerStamp            *time.Time `json:"answer_stamp"`
	AnswerEpoch            int64      `json:"answer_epoch"`
	EndEpoch               int64      `json:"end_epoch"`
	EndStamp               *time.Time `json:"end_stamp"`
	Duration               int64      `json:"duration"`
	MDuration              int64      `json:"mduration"`
	BillSec                int64      `json:"billsec"`
	BillMsec               int64      `json:"billmsec"`
	HoldAccumSeconds       int64      `json:"hold_accum_seconds"`
	BridgeUUID             string     `json:"bridge_uuid"`
	ReadCodec              string     `json:"read_codec"`
	ReadRate               int64      `json:"read_rate"`
	WriteCodec             string     `json:"write_codec"`
	WriteRate              int64      `json:"write_rate"`
	RemoteMediaIP          string     `json:"remote_media_ip"`
	NetworkAddr            string     `json:"network_addr"`
	RecordPath             string     `json:"record_path"`
	RecordName             string     `json:"record_name"`
	RecordLength           int64      `json:"record_length"`
	RecordTranscription    string     `json:"record_transcription"`
	Leg                    string     `json:"leg"`
	OriginatingLegUUID     string     `json:"originating_leg_uuid"`
	PDDMs                  int64      `json:"pdd_ms"`
	RTPAudioInMOS          float64    `json:"rtp_audio_in_mos"`
	LastApp                string     `json:"last_app"`
	LastArg                string     `json:"last_arg"`
	VoicemailMessage       string     `json:"voicemail_message"`
	MissedCall             bool       `json:"missed_call"`
	CallCenterQueueUUID    string     `json:"call_center_queue_uuid"`
	CCSide                 string     `json:"cc_side"`
	CCMemberUUID           string     `json:"cc_member_uuid"`
	CCQueueJoinedEpoch     int64      `json:"cc_queue_joined_epoch"`
	CCQueue                string     `json:"cc_queue"`
	CCMemberSessionUUID    string     `json:"cc_member_session_uuid"`
	CCAgentUUID            string     `json:"cc_agent_uuid"`
	CCAgent                string     `json:"cc_agent"`
	CCAgentType            string     `json:"cc_agent_type"`
	CCAgentBridged         bool       `json:"cc_agent_bridged"`
	CCQueueAnsweredEpoch   int64      `json:"cc_queue_answered_epoch"`
	CCQueueTerminatedEpoch int64      `json:"cc_queue_terminated_epoch"`
	CCQueueCanceledEpoch   int64      `json:"cc_queue_canceled_epoch"`
	CCCancelReason         string     `json:"cc_cancel_reason"`
	CCCause                string     `json:"cc_cause"`
	WaitSec                int64      `json:"waitsec"`
	ConferenceName         string     `json:"conference_name"`
	ConferenceUUID         string     `json:"conference_uuid"`
	ConferenceMemberID     string     `json:"conference_member_id"`
	DigitsDialed           string     `json:"digits_dialed"`
	PINNumber              string     `json:"pin_number"`
	Status                 string     `json:"status"`
	HangupCause            string     `json:"hangup_cause"`
	HangupCauseQ850        string     `json:"hangup_cause_q850"`
	SIPHangupDisposition   string     `json:"sip_hangup_disposition"`
	RingGroupUUID          string     `json:"ring_group_uuid"`
	IVRMenuUUID            string     `json:"ivr_menu_uuid"`
	CallFlow               string     `json:"call_flow"`
	XML                    string     `json:"xml"`
	JSON                   string     `json:"json"`
	InsertDate             time.Time  `json:"insert_date"`
	InsertUser             string     `json:"insert_user"`
	UpdateDate             time.Time  `json:"update_date"`
	UpdateUser             string     `json:"update_user"`
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