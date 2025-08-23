package core

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/avistopia/monorepo/go/arithland-telegram/internal/models"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/compact"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/components"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/flows"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/texts"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/values"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *Service) AdminQuestionsFlow() flows.Flow {
	return flows.Flow{
		KeyboardButtonActions: map[string]components.Action{
			texts.AdminQuestions_Show: func(user *models.User, message *tgbotapi.Message) error {
				if !s.isAdmin(message.Chat.UserName) {
					if err := s.telegram.SendOrEdit(s.forbiddenMessage(), message.Chat.ID, nil); err != nil {
						return fmt.Errorf("failed to send forbidden message: %w", err)
					}

					return nil
				}

				user.State = models.NewWaitingForAdminQuestionsIDState()

				if err := s.userRepo.Save(user); err != nil {
					return fmt.Errorf("failed to save user: %w", err)
				}

				response, err := s.AdminQuestionsMessage()
				if err != nil {
					return fmt.Errorf("failed to prepare message: %w", err)
				}

				if err = s.telegram.SendOrEdit(*response, message.Chat.ID, nil); err != nil {
					return fmt.Errorf("failed to send questions admin message: %w", err)
				}

				return nil
			},
		},
		InlineButtonActions: map[components.InlineButtonActionName]components.InlineButtonAction{
			actionName_adminQuestions_create: func(
				user *models.User, query *tgbotapi.CallbackQuery, data string,
			) (string, error) {
				if !s.isAdmin(query.From.UserName) {
					return texts.Err_Forbidden, nil
				}

				level := models.QuestionLevel(data)

				question := &models.Question{Level: level}

				if err := s.questionRepo.Save(question); err != nil {
					return "", fmt.Errorf("failed to create question: %w", err)
				}

				user.State = models.NewWaitingForAdminQuestionFieldState(
					question.ID, models.QuestionFieldName_TextAndImage, false,
				)

				if err := s.userRepo.Save(user); err != nil {
					return "", fmt.Errorf("failed to save user: %w", err)
				}

				if err := s.telegram.SendOrEdit(
					components.Message{
						Text: texts.Format(
							texts.AdminQuestions_Created,
							map[string]string{
								"id":    texts.FormatInt(int64(question.ID)),
								"level": level.Render(),
								"field": models.QuestionFieldName_TextAndImage.Render(),
							},
						),
					},
					query.Message.Chat.ID,
					nil,
				); err != nil {
					return "", fmt.Errorf("failed to send question created message: %w", err)
				}

				return "", nil
			},
			actionName_adminQuestions_updateField: func(
				user *models.User, query *tgbotapi.CallbackQuery, rawData string,
			) (string, error) {
				data := &actionData_adminQuestions_updateField{}

				if err := compact.Unmarshal(rawData, data); err != nil {
					return "", fmt.Errorf("invalid data %q received: %w", rawData, err)
				}

				fieldName := models.QuestionFieldName(data.FieldName)

				user.State = models.NewWaitingForAdminQuestionFieldState(
					data.QuestionID, fieldName, true,
				)

				if err := s.userRepo.Save(user); err != nil {
					return "", fmt.Errorf("failed to save user: %w", err)
				}

				if err := s.telegram.SendOrEdit(
					components.Message{Text: texts.Format(
						texts.AdminQuestions_SendField, map[string]string{"field": fieldName.Render()},
					)},
					query.Message.Chat.ID,
					nil,
				); err != nil {
					return "", fmt.Errorf("failed to send text and image message: %w", err)
				}

				return "", nil
			},
			actionName_adminQuestions_toggleInactive: func(
				user *models.User, query *tgbotapi.CallbackQuery, data string,
			) (string, error) {
				questionID, err := strconv.Atoi(data)
				if err != nil {
					return "", fmt.Errorf("invalid data %q received: %w", data, err)
				}

				question, err := s.questionRepo.GetByID(uint(questionID))
				if err != nil {
					return "", fmt.Errorf("failed to get by ID question: %w", err)
				}

				question.Inactive = !question.Inactive

				if err = s.questionRepo.Save(question); err != nil {
					return "", fmt.Errorf("failed to save question: %w", err)
				}

				message, err := s.adminQuestionDetailsMessage(question)
				if err != nil {
					return "", fmt.Errorf("failed to prepare admin question details message: %w", err)
				}

				if err = s.telegram.SendOrEdit(
					*message, query.Message.Chat.ID, values.Ptr(query.Message.MessageID),
				); err != nil {
					return "", fmt.Errorf("failed to send admin question details message: %w", err)
				}

				response := texts.AdminQuestions_Deleted
				if !question.Inactive {
					response = texts.AdminQuestions_UndoDeleted
				}

				return response, nil
			},
			actionName_adminQuestions_refresh: func(
				user *models.User, query *tgbotapi.CallbackQuery, data string,
			) (string, error) {
				questionID, err := strconv.Atoi(data)
				if err != nil {
					return "", fmt.Errorf("invalid data %q received: %w", data, err)
				}

				question, err := s.questionRepo.GetByID(uint(questionID))
				if err != nil {
					return "", fmt.Errorf("failed to get by ID question: %w", err)
				}

				message, err := s.adminQuestionDetailsMessage(question)
				if err != nil {
					return "", fmt.Errorf("failed to prepare question details message: %w", err)
				}

				if err = s.telegram.SendOrEdit(
					*message, query.Message.Chat.ID, values.Ptr(query.Message.MessageID),
				); err != nil {
					logrus.WithError(err).Error("failed to send question details message")
				}

				return "", nil
			},
			actionName_adminQuestions_updateLevel: func(
				user *models.User, query *tgbotapi.CallbackQuery, rawData string,
			) (string, error) {
				data := &actionData_adminQuestions_updateLevel{}

				if err := compact.Unmarshal(rawData, data); err != nil {
					return "", fmt.Errorf("invalid data %q received: %w", rawData, err)
				}

				question, err := s.questionRepo.GetByID(data.QuestionID)
				if err != nil {
					return "", fmt.Errorf("failed to get by ID question: %w", err)
				}

				question.Level = models.QuestionLevel(data.Level)

				if err = s.questionRepo.Save(question); err != nil {
					return "", fmt.Errorf("failed to save question: %w", err)
				}

				message, err := s.adminQuestionDetailsMessage(question)
				if err != nil {
					return "", fmt.Errorf("failed to prepare question details message: %w", err)
				}

				if err = s.telegram.SendOrEdit(
					*message, query.Message.Chat.ID, values.Ptr(query.Message.MessageID),
				); err != nil {
					logrus.WithError(err).Error("failed to send question details message")
				}

				return texts.Format(
					texts.AdminQuestions_FieldUpdated, map[string]string{"field": texts.AdminQuestions_Level},
				), nil
			},
			actionName_adminQuestions_navigate: func(
				user *models.User, query *tgbotapi.CallbackQuery, rawData string,
			) (string, error) {
				data := &actionData_adminQuestions_navigate{}

				if err := compact.Unmarshal(rawData, data); err != nil {
					return "", fmt.Errorf("invalid data %q received: %w", rawData, err)
				}

				question, err := s.questionRepo.Navigate(data.QuestionID, models.Direction(data.Direction))
				if err != nil {
					return "", fmt.Errorf("failed to get next question: %w", err)
				}

				message, err := s.adminQuestionDetailsMessage(question)
				if err != nil {
					return "", fmt.Errorf("failed to prepare question details message: %w", err)
				}

				if err = s.telegram.SendOrEdit(
					*message, query.Message.Chat.ID, values.Ptr(query.Message.MessageID),
				); err != nil {
					logrus.WithError(err).Error("failed to send question details message")
				}

				return "", nil
			},
		},
		MessageActions: map[models.StateName]components.Action{
			models.StateName_WaitingForAdminQuestionsID: func(user *models.User, message *tgbotapi.Message) error {
				fmt.Println(message.Text)
				fmt.Println(texts.NormalizeValue(message.Text))
				questionID, err := strconv.Atoi(texts.NormalizeValue(message.Text))
				if err != nil {
					if err = s.telegram.SendOrEdit(
						components.Message{Text: texts.Err_InvalidInteger}, message.Chat.ID, nil,
					); err != nil {
						return fmt.Errorf("failed to send invalid integer message: %w", err)
					}

					return nil
				}

				question, err := s.questionRepo.GetByIDOrNext(uint(questionID))
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						if err = s.telegram.SendOrEdit(
							components.Message{Text: texts.AdminQuestions_NotFound}, message.Chat.ID, nil,
						); err != nil {
							return fmt.Errorf("failed to send questions not found message: %w", err)
						}

						return nil
					}

					return fmt.Errorf("failed to get by ID or next question: %w", err)
				}

				response, err := s.adminQuestionDetailsMessage(question)
				if err != nil {
					return fmt.Errorf("failed to prepare question details message: %w", err)
				}

				if err = s.telegram.SendOrEdit(*response, message.Chat.ID, nil); err != nil {
					return fmt.Errorf("failed to send question details message: %w", err)
				}

				return nil
			},
			models.StateName_WaitingForAdminQuestionField: func(user *models.User, message *tgbotapi.Message) error {
				stateData := user.State.Data.WaitingForAdminQuestionField

				question, err := s.questionRepo.GetByID(stateData.QuestionID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						if err = s.telegram.SendOrEdit(
							components.Message{Text: texts.AdminQuestions_NotFound}, message.Chat.ID, nil,
						); err != nil {
							return fmt.Errorf("failed to send questions not found message: %w", err)
						}

						return nil
					}

					return fmt.Errorf("failed to get question by ID: %w", err)
				}

				responseText := ""
				var next models.QuestionFieldName

				switch stateData.FieldName {
				case models.QuestionFieldName_TextAndImage:
					question.Text = texts.NormalizeDescription(message.Text)
					if question.Text == "" {
						question.Text = texts.NormalizeDescription(message.Caption)
					}

					question.PhotoID = largestPhotoID(message.Photo)

					if !stateData.Terminate {
						next = models.QuestionFieldName_Answer
					}
				case models.QuestionFieldName_Answer:
					question.Answer = texts.NormalizeValue(message.Text)
				default:
					return fmt.Errorf("invalid question field name %s", stateData.FieldName)
				}

				if err = s.questionRepo.Save(question); err != nil {
					return fmt.Errorf("failed to save question: %w", err)
				}

				if next != "" {
					user.State = models.NewWaitingForAdminQuestionFieldState(question.ID, next, false)
					responseText = texts.Format(texts.AdminQuestions_FieldUpdatedSendField, map[string]string{
						"updatedField": stateData.FieldName.Render(),
						"nextField":    next.Render(),
					})
				} else {
					user.State = models.NewDefaultState()
					responseText = texts.Format(texts.AdminQuestions_FieldUpdated, map[string]string{
						"field": stateData.FieldName.Render(),
					})
				}

				if err = s.userRepo.Save(user); err != nil {
					return fmt.Errorf("failed to save user: %w", err)
				}

				if err = s.telegram.SendOrEdit(
					components.Message{Text: responseText}, message.Chat.ID, nil,
				); err != nil {
					return fmt.Errorf("failed to send response message: %w", err)
				}

				return nil
			},
		},
	}
}

func (s *Service) AdminQuestionsMessage() (*components.Message, error) {
	count, err := s.questionRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("failed to get questions count: %w", err)
	}

	inlineKeyboard := make(components.InlineKeyboard, 0, len(models.QuestionLevels))

	// TODO add show last question inline keyboard button

	for _, level := range models.QuestionLevels {
		inlineKeyboard = append(inlineKeyboard, []components.InlineButton{
			components.NewInlineButton(
				texts.Format(texts.AdminQuestions_Create, map[string]string{
					"levelEmoji": level.Emoji(),
					"level":      level.Render(),
				}),
				actionName_adminQuestions_create,
				string(level),
			),
		})
	}

	return &components.Message{
		Text: texts.Format(texts.AdminQuestions, map[string]string{
			"questionsCount": texts.FormatInt(count),
		}),
		InlineKeyboard: inlineKeyboard,
	}, nil
}

func (s *Service) adminQuestionDetailsMessage(question *models.Question) (*components.Message, error) {
	updateFieldKeyboardRow := make([]components.InlineButton, 0, len(models.QuestionFieldNames))

	for _, fieldName := range models.QuestionFieldNames {
		data, err := compact.Marshal(actionData_adminQuestions_updateField{
			QuestionID: question.ID,
			FieldName:  string(fieldName),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal update question field data: %w", err)
		}

		updateFieldKeyboardRow = append(updateFieldKeyboardRow, components.NewInlineButton(
			texts.Format(texts.AdminQuestions_UpdateField, map[string]string{"field": fieldName.Render()}),
			actionName_adminQuestions_updateField,
			data,
		))
	}

	slices.Reverse(updateFieldKeyboardRow)

	deleteText := texts.AdminQuestions_Delete
	if question.Inactive {
		deleteText = texts.AdminQuestions_UndoDelete
	}

	inlineKeyboard := components.InlineKeyboard{
		{components.NewInlineButton(
			texts.AdminQuestions_Refresh, actionName_adminQuestions_refresh, fmt.Sprintf("%d", question.ID),
		)},
		updateFieldKeyboardRow,
		{components.NewInlineButton(
			deleteText, actionName_adminQuestions_toggleInactive, fmt.Sprintf("%d", question.ID),
		)},
	}

	for _, level := range models.QuestionLevels {
		if question.Level == level {
			continue
		}

		data, err := compact.Marshal(actionData_adminQuestions_updateLevel{
			QuestionID: question.ID,
			Level:      string(level),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal update question level data: %w", err)
		}

		inlineKeyboard = append(inlineKeyboard, []components.InlineButton{
			components.NewInlineButton(
				texts.Format(texts.AdminQuestions_UpdateLevel, map[string]string{
					"levelEmoji": level.Emoji(),
					"level":      level.Render(),
				}),
				actionName_adminQuestions_updateLevel,
				data,
			),
		})
	}

	nextData, err := compact.Marshal(actionData_adminQuestions_navigate{
		QuestionID: question.ID,
		Direction:  string(models.Direction_Older),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal nextData: %w", err)
	}

	backData, err := compact.Marshal(actionData_adminQuestions_navigate{
		QuestionID: question.ID,
		Direction:  string(models.Direction_Newer),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal backData: %w", err)
	}

	inlineKeyboard = append(inlineKeyboard, []components.InlineButton{
		components.NewInlineButton(texts.Common_Next, actionName_adminQuestions_navigate, nextData),
		components.NewInlineButton(texts.Common_Back, actionName_adminQuestions_navigate, backData),
	})

	photoID := question.PhotoID
	if photoID == "" {
		photoID = s.defaultPhotoID
	}

	return &components.Message{
		Text: texts.Format(texts.AdminQuestions_Details, map[string]string{
			"active":     texts.FormatBoolAsEmoji(!question.Inactive),
			"id":         texts.FormatInt(int64(question.ID)),
			"levelEmoji": question.Level.Emoji(),
			"level":      question.Level.Render(),
			"text":       question.Text,
			"answer":     question.Answer,
		}),
		PhotoID:        photoID,
		InlineKeyboard: inlineKeyboard,
	}, nil
}

func largestPhotoID(photos []tgbotapi.PhotoSize) string {
	maxID := ""
	maxSize := 0

	for i := range photos {
		size := photos[i].Width * photos[i].Height

		if size > maxSize {
			maxSize = size
			maxID = photos[i].FileID
		}
	}

	return maxID
}
