package core

import "github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/components"

const (
	actionName_adminQuestions_create         components.InlineButtonActionName = "adminQuestions_create"
	actionName_adminQuestions_updateField    components.InlineButtonActionName = "adminQuestions_updateField"
	actionName_adminQuestions_toggleInactive components.InlineButtonActionName = "adminQuestions_toggleInactive"
	actionName_adminQuestions_refresh        components.InlineButtonActionName = "adminQuestions_refresh"
	actionName_adminQuestions_updateLevel    components.InlineButtonActionName = "adminQuestions_updateLevel"
	actionName_adminQuestions_navigate       components.InlineButtonActionName = "adminQuestions_navigate"

	actionName_profileManagement_changeDisplayName           components.InlineButtonActionName = "profileManagement_changeDisplayName"           //nolint: lll
	actionName_profileManagement_backToShowProfileManagement components.InlineButtonActionName = "profileManagement_backToShowProfileManagement" //nolint: lll

	actionName_questionManagement_showQuestion components.InlineButtonActionName = "questionManagement_showQuestion"
	actionName_questionManagement_buyQuestion  components.InlineButtonActionName = "questionManagement_buyQuestion"
	actionName_questionManagement_navigate     components.InlineButtonActionName = "questionManagement_navigate"
)

type actionData_adminQuestions_navigate struct {
	QuestionID uint
	Direction  string
}

type actionData_adminQuestions_updateField struct {
	QuestionID uint
	FieldName  string
}

type actionData_adminQuestions_updateLevel struct {
	QuestionID uint
	Level      string
}

type actionData_questionManagement_showQuestion struct {
	UserQuestionID uint
	Answered       bool
	Direction      string
}
