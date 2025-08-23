package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

//nolint:recvcheck
type State struct {
	Name StateName `json:"name"`
	Data StateData `json:"data"`
}

func (s *State) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to assert []byte type")
	}

	return json.Unmarshal(bytes, s)
}

func (s State) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type StateName string

const (
	StateName_Default                      StateName = "Default"
	StateName_WaitingForUserField          StateName = "WaitingForUserField"
	StateName_WaitingForQuestionAnswer     StateName = "WaitingForQuestionAnswer"
	StateName_WaitingForAdminQuestionsID   StateName = "WaitingForAdminQuestionsID"
	StateName_WaitingForAdminQuestionField StateName = "WaitingForAdminQuestionField"
)

type StateData struct {
	Default                      DefaultStateData                      `json:"default"`
	WaitingForUserField          WaitingForUserFieldStateData          `json:"waiting_for_user_field"`
	WaitingForQuestionAnswer     WaitingForQuestionAnswerStateData     `json:"waiting_for_question_answer"`
	WaitingForAdminQuestionsID   WaitingForAdminQuestionsIDStateData   `json:"waiting_for_admin_questions_id"`
	WaitingForAdminQuestionField WaitingForAdminQuestionFieldStateData `json:"waiting_for_question_field"`
}

type DefaultStateData struct{}

func NewDefaultState() State {
	return State{
		Name: StateName_Default,
		Data: StateData{
			Default: DefaultStateData{},
		},
	}
}

type WaitingForUserFieldStateData struct {
	FieldName UserFieldName `json:"field_name"`
}

type UserFieldName string

const (
	UserFieldName_DisplayName = "DisplayName"
)

func NewWaitingForUserFieldState(fieldName UserFieldName) State {
	return State{
		Name: StateName_WaitingForUserField,
		Data: StateData{
			WaitingForUserField: WaitingForUserFieldStateData{
				FieldName: fieldName,
			},
		},
	}
}

type WaitingForQuestionAnswerStateData struct {
	UserQuestionID uint `json:"user_question_id"`
}

func NewWaitingForQuestionAnswerState(userQuestionID uint) State {
	return State{
		Name: StateName_WaitingForQuestionAnswer,
		Data: StateData{
			WaitingForQuestionAnswer: WaitingForQuestionAnswerStateData{
				UserQuestionID: userQuestionID,
			},
		},
	}
}

type WaitingForAdminQuestionsIDStateData struct {
	UserQuestionID string `json:"user_question_id"`
}

func NewWaitingForAdminQuestionsIDState() State {
	return State{
		Name: StateName_WaitingForAdminQuestionsID,
		Data: StateData{
			WaitingForAdminQuestionsID: WaitingForAdminQuestionsIDStateData{},
		},
	}
}

type WaitingForAdminQuestionFieldStateData struct {
	QuestionID uint              `json:"question_id"`
	FieldName  QuestionFieldName `json:"field_name"`
	Terminate  bool              `json:"terminate"`
}

type QuestionFieldName string

func (v QuestionFieldName) Render() string {
	switch v {
	case QuestionFieldName_TextAndImage:
		return "صورت"
	case QuestionFieldName_Answer:
		return "پاسخ"
	default:
		return "نامشخص"
	}
}

const (
	QuestionFieldName_TextAndImage QuestionFieldName = "TextAndImage"
	QuestionFieldName_Answer       QuestionFieldName = "Answer"
)

var (
	QuestionFieldNames = []QuestionFieldName{QuestionFieldName_TextAndImage, QuestionFieldName_Answer}
)

func NewWaitingForAdminQuestionFieldState(questionID uint, fieldName QuestionFieldName, terminate bool) State {
	return State{
		Name: StateName_WaitingForAdminQuestionField,
		Data: StateData{
			WaitingForAdminQuestionField: WaitingForAdminQuestionFieldStateData{
				QuestionID: questionID,
				FieldName:  fieldName,
				Terminate:  terminate,
			},
		},
	}
}
