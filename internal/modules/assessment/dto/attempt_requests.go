package dto

type StartAttemptReq struct {
	// เผื่ออนาคต เช่น proctoring token ฯลฯ
}

type UpsertAnswerReq struct {
	QuestionID string `json:"question_id" validate:"required"`
	// ใช้ก้านใดก้านหนึ่งตามประเภทคำถาม:
	SelectedChoiceIDs []string `json:"selected_choice_ids"` // single/multiple/true_false
	TextAnswer        *string  `json:"text_answer"`         // short_text
}

type SubmitAttemptReq struct {
	// เผื่ออนาคต เช่น confirm=true
}
