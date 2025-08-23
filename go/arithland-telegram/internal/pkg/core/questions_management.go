package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/avistopia/monorepo/go/arithland-telegram/internal/models"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/compact"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/components"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/flows"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/texts"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/values"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func (s *Service) questionsManagementFlow() flows.Flow {
	return flows.Flow{
		KeyboardButtonActions: map[string]components.Action{
			texts.QuestionsManagement_Show: func(user *models.User, message *tgbotapi.Message) error {
				response, err := s.questionsManagementMessage(user)
				if err != nil {
					return fmt.Errorf("failed to prepare message: %w", err)
				}

				if err = s.telegram.SendOrEdit(*response, message.Chat.ID, nil); err != nil {
					return fmt.Errorf("failed to send question management message: %w", err)
				}

				return nil
			},
		},
		InlineButtonActions: map[components.InlineButtonActionName]components.InlineButtonAction{
			actionName_questionManagement_showQuestion: func(
				user *models.User, query *tgbotapi.CallbackQuery, rawData string,
			) (string, error) {
				data := &actionData_questionManagement_showQuestion{}

				if err := compact.Unmarshal(rawData, data); err != nil {
					return "", fmt.Errorf("invalid data %q received: %w", rawData, err)
				}

				userQuestion, err := s.userQuestionRepo.Last(models.UserQuestionFilter{
					UserID:   values.Ptr(user.ID),
					Answered: values.Ptr(data.Answered),
				})
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return texts.Err_NoQuestionsFound, nil
					}

					return "", fmt.Errorf("failed to find a question: %w", err)
				}

				if err = s.sendOrEditUserQuestionDetailsMessage(
					user, data.Answered, userQuestion, query.Message.Chat.ID, nil,
				); err != nil {
					return "", fmt.Errorf("failed to send user question details message: %w", err)
				}

				return "", nil
			},
			actionName_questionManagement_buyQuestion: func(
				user *models.User, query *tgbotapi.CallbackQuery, data string,
			) (string, error) {
				level := models.QuestionLevel(data)
				price := float64(models.QuestionPrices[level])
				errInsufficientFunds := errors.New("insufficient funds")
				errNoQuestionsFound := errors.New("no questions found")

				var userQuestion *models.UserQuestion

				err := s.db.Transaction(func(tx *gorm.DB) error {
					user, err := s.userRepo.GetByID(user.ID)
					if err != nil {
						return fmt.Errorf("failed to get user: %w", err)
					}

					if user.Balance < price {
						return errInsufficientFunds
					}

					question, err := s.questionRepo.GetRandom(models.QuestionFilter{
						Level:             values.Ptr(level),
						Inactive:          values.Ptr(false),
						UserIDDoesNotHave: values.Ptr(user.ID),
					})
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return errNoQuestionsFound
						}

						return fmt.Errorf("failed to get a random question: %w", err)
					}

					user.Balance -= price

					if err = s.userRepo.Save(user); err != nil {
						return fmt.Errorf("failed to save user: %w", err)
					}

					userQuestion = &models.UserQuestion{
						UserID:     user.ID,
						QuestionID: question.ID,
						User:       *user,
						Question:   *question,
					}

					if err = s.userQuestionRepo.Save(userQuestion); err != nil {
						return fmt.Errorf("failed to save user question: %w", err)
					}

					return nil
				})
				if err != nil {
					if errors.Is(err, errInsufficientFunds) {
						return texts.Err_InsufficientFunds, nil
					}

					if errors.Is(err, errNoQuestionsFound) {
						return texts.Err_NoQuestionsFound, nil
					}

					return "", fmt.Errorf("failed to commit purchase atomic transaction: %w", err)
				}

				if err = s.sendOrEditUserQuestionDetailsMessage(
					user, false, userQuestion, query.Message.Chat.ID, nil,
				); err != nil {
					return "", fmt.Errorf("failed to send user question details message: %w", err)
				}

				return "", nil
			},
			actionName_questionManagement_navigate: func(
				user *models.User, query *tgbotapi.CallbackQuery, rawData string,
			) (string, error) {
				data := &actionData_questionManagement_showQuestion{}

				if err := compact.Unmarshal(rawData, data); err != nil {
					return "", fmt.Errorf("invalid data %q received: %w", rawData, err)
				}

				userQuestion, err := s.userQuestionRepo.Navigate(
					data.UserQuestionID,
					models.UserQuestionFilter{
						UserID:   values.Ptr(user.ID),
						Answered: values.Ptr(data.Answered),
					},
					models.Direction(data.Direction),
				)
				if err != nil {
					return "", fmt.Errorf("failed to get previous user question: %w", err)
				}

				if err = s.sendOrEditUserQuestionDetailsMessage(
					user, data.Answered, userQuestion, query.Message.Chat.ID, values.Ptr(query.Message.MessageID),
				); err != nil {
					return "", fmt.Errorf("failed to send user question details message: %w", err)
				}
				return "", nil
			},
		},
		MessageActions: map[models.StateName]components.Action{
			models.StateName_WaitingForQuestionAnswer: func(user *models.User, message *tgbotapi.Message) error {
				userQuestionID := user.State.Data.WaitingForQuestionAnswer.UserQuestionID

				userQuestion, err := s.userQuestionRepo.GetByID(userQuestionID)
				if err != nil {
					return fmt.Errorf("failed to get user question: %w", err)
				}

				if userQuestion.Answered {
					if err = s.telegram.SendOrEdit(
						components.Message{Text: texts.QuestionsManagement_AlreadyAnswered}, message.Chat.ID, nil,
					); err != nil {
						return fmt.Errorf("failed to send incorrect answer message: %w", err)
					}

					return nil
				}

				userAnswer := texts.NormalizeValue(message.Text)
				correctAnswer := texts.NormalizeValue(userQuestion.Question.Answer)

				if userAnswer != correctAnswer {
					if err = s.telegram.SendOrEdit(
						components.Message{Text: texts.QuestionsManagement_IncorrectAnswer}, message.Chat.ID, nil,
					); err != nil {
						return fmt.Errorf("failed to send incorrect answer message: %w", err)
					}

					return nil
				}

				reward := float64(models.QuestionPrices[userQuestion.Question.Level] * 2)

				err = s.db.Transaction(func(tx *gorm.DB) error {
					user, err := s.userRepo.GetByID(user.ID)
					if err != nil {
						return fmt.Errorf("failed to get user: %w", err)
					}

					user.Balance += reward
					user.State = models.NewDefaultState()

					if err = s.userRepo.Save(user); err != nil {
						return fmt.Errorf("failed to save user: %w", err)
					}

					userQuestion.Answered = true
					userQuestion.AnsweredAt = time.Now()

					if err = s.userQuestionRepo.Save(userQuestion); err != nil {
						return fmt.Errorf("failed to save user question: %w", err)
					}

					return nil
				})
				if err != nil {
					return fmt.Errorf("failed to commit answer atomic transaction: %w", err)
				}

				if err = s.telegram.SendOrEdit(
					components.Message{
						Text: texts.Format(
							texts.QuestionsManagement_CorrectAnswer,
							map[string]string{
								"price": texts.FormatFloat(reward),
							},
						),
					},
					message.Chat.ID,
					nil,
				); err != nil {
					return fmt.Errorf("failed to send incorrect answer message: %w", err)
				}

				return nil
			},
		},
	}
}

func (s *Service) questionsManagementMessage(user *models.User) (*components.Message, error) {
	questionsCount, err := s.userQuestionRepo.Count(models.UserQuestionFilter{
		UserID: values.Ptr(user.ID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get count user questions: %w", err)
	}

	answeredQuestionsCount, err := s.userQuestionRepo.Count(models.UserQuestionFilter{
		UserID:   values.Ptr(user.ID),
		Answered: values.Ptr(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get count user answered questions: %w", err)
	}

	showUnansweredQuestionsData, err := compact.Marshal(actionData_questionManagement_showQuestion{Answered: false})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal showUnansweredQuestionsData: %w", err)
	}

	showAnsweredQuestionsData, err := compact.Marshal(actionData_questionManagement_showQuestion{Answered: true})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal showAnsweredQuestionsData: %w", err)
	}

	inlineKeyboard := components.InlineKeyboard{
		{components.NewInlineButton(
			texts.QuestionsManagement_ShowUnansweredQuestions,
			actionName_questionManagement_showQuestion,
			showUnansweredQuestionsData,
		)},
		{components.NewInlineButton(
			texts.QuestionsManagement_ShowAnsweredQuestions,
			actionName_questionManagement_showQuestion,
			showAnsweredQuestionsData,
		)},
	}

	for _, level := range models.QuestionLevels {
		inlineKeyboard = append(inlineKeyboard, []components.InlineButton{components.NewInlineButton(
			texts.Format(
				texts.QuestionsManagement_BuyQuestion,
				map[string]string{
					"levelEmoji": level.Emoji(),
					"level":      level.Render(),
					"price":      texts.ReplaceDigitsToFarsi(texts.FormatInt(models.QuestionPrices[level])),
				},
			),
			actionName_questionManagement_buyQuestion,
			string(level),
		)})
	}

	return &components.Message{
		Text: texts.Format(texts.QuestionsManagement, map[string]string{
			"questionsCount":         texts.FormatInt(questionsCount),
			"answeredQuestionsCount": texts.FormatInt(answeredQuestionsCount),
			"balance":                texts.FormatFloat(user.Balance),
		}),
		InlineKeyboard: inlineKeyboard,
	}, nil
}

func (s *Service) sendOrEditUserQuestionDetailsMessage(
	user *models.User, answered bool, userQuestion *models.UserQuestion, chatID int64, originalMessageID *int,
) error {
	inlineKeyboard := make(components.InlineKeyboard, 0, 2)

	nextData, err := compact.Marshal(actionData_questionManagement_showQuestion{
		UserQuestionID: userQuestion.ID,
		Answered:       answered,
		Direction:      string(models.Direction_Older),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal nextData: %w", err)
	}

	backData, err := compact.Marshal(actionData_questionManagement_showQuestion{
		UserQuestionID: userQuestion.ID,
		Answered:       answered,
		Direction:      string(models.Direction_Newer),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal nextData: %w", err)
	}

	inlineKeyboard = append(inlineKeyboard, []components.InlineButton{
		components.NewInlineButton(texts.Common_Next, actionName_questionManagement_navigate, nextData),
		components.NewInlineButton(texts.Common_Back, actionName_questionManagement_navigate, backData),
	})

	answer := userQuestion.Question.Answer
	answeredAt := texts.FormatTime(userQuestion.AnsweredAt)
	footer := ""

	if !userQuestion.Answered {
		answer = texts.QuestionsManagement_Unanswered
		answeredAt = texts.QuestionsManagement_Unanswered
		footer = texts.QuestionsManagement_DetailsFooterUnanswered

		user.State = models.NewWaitingForQuestionAnswerState(userQuestion.ID)
		if err = s.userRepo.Save(user); err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}
	}

	header := texts.QuestionsManagement_DetailsHeaderUnanswered

	if answered {
		header = texts.QuestionsManagement_DetailsHeaderAnswered
	}

	photoID := userQuestion.Question.PhotoID
	if photoID == "" {
		photoID = s.defaultPhotoID
	}

	message := components.Message{
		Text: texts.Format(texts.QuestionsManagement_Details, map[string]string{
			"header":     header,
			"id":         texts.FormatInt(int64(userQuestion.Question.ID)),
			"levelEmoji": userQuestion.Question.Level.Emoji(),
			"level":      userQuestion.Question.Level.Render(),
			"text":       userQuestion.Question.Text,
			"answer":     answer,
			"createdAt":  texts.FormatTime(userQuestion.CreatedAt),
			"answeredAt": answeredAt,
			"footer":     footer,
		}),
		PhotoID:        photoID,
		InlineKeyboard: inlineKeyboard,
	}

	if err := s.telegram.SendOrEdit(message, chatID, originalMessageID); err != nil {
		return fmt.Errorf("failed to send or edit user question details message: %w", err)
	}

	return nil
}
