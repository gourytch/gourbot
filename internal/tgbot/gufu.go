package tgbot

import "github.com/go-telegram/bot/models"

// GetUserFromUpdate extracts the user information from an update.
// It checks various fields in the update to find a non-nil user and returns it.
// If no user information is found, it returns nil.
func GetUserFromUpdate(update *models.Update) *models.User {
	switch {
	// update.ID, int64, no From

	// update.Message, *models.Message, From as pointer
	case update.Message != nil && update.Message.From != nil:
		return update.Message.From

	// update.EditedMessage, *models.Message, From as pointer
	case update.EditedMessage != nil && update.EditedMessage.From != nil:
		return update.EditedMessage.From

	// update.ChannelPost, *models.Message, From as pointer
	case update.ChannelPost != nil && update.ChannelPost.From != nil:
		return update.ChannelPost.From

	// update.EditedChannelPost, *models.Message, From as pointer
	case update.EditedChannelPost != nil && update.EditedChannelPost.From != nil:
		return update.EditedChannelPost.From

	// update.BusinessConnection, *models.BusinessConnection, From as value
	case update.BusinessConnection != nil:
		return &update.BusinessConnection.User

	// update.BusinessMessage, *models.Message, From as pointer
	case update.BusinessMessage != nil && update.BusinessMessage.From != nil:
		return update.BusinessMessage.From

	// update.EditedBusinessMessage, *models.Message, From as pointer
	case update.EditedBusinessMessage != nil && update.EditedBusinessMessage.From != nil:
		return update.EditedBusinessMessage.From

	// update.DeletedBusinessMessages, *models.BusinessMessagesDeleted, no From
	// update.MessageReaction, *models.MessageReactionUpdated, From as pointer
	case update.MessageReaction != nil && update.MessageReaction.User != nil:
		return update.MessageReaction.User

	// update.MessageReactionCount, *models.MessageReactionCountUpdated, no From
	// update.InlineQuery, *models.InlineQuery, From as pointer
	case update.InlineQuery != nil && update.InlineQuery.From != nil:
		return update.InlineQuery.From

	// update.ChosenInlineResult, *models.ChosenInlineResult, From as value
	case update.ChosenInlineResult != nil:
		return &update.ChosenInlineResult.From

	// update.CallbackQuery, *models.CallbackQuery, From as value
	case update.CallbackQuery != nil:
		return &update.CallbackQuery.From

	// update.ShippingQuery, *models.ShippingQuery, From as pointer
	case update.ShippingQuery != nil && update.ShippingQuery.From != nil:
		return update.ShippingQuery.From

	// update.PreCheckoutQuery, *models.PreCheckoutQuery, From as pointer
	case update.PreCheckoutQuery != nil && update.PreCheckoutQuery.From != nil:
		return update.PreCheckoutQuery.From

	// update.PurchasedPaidMedia, *models.PaidMediaPurchased, no From
	// update.Poll, *models.Poll, no From
	// update.PollAnswer, *models.PollAnswer, From as pointer
	case update.PollAnswer != nil && update.PollAnswer.User != nil:
		return update.PollAnswer.User

	// update.MyChatMember, *models.ChatMemberUpdated, From as value
	case update.MyChatMember != nil:
		return &update.MyChatMember.From

	// update.ChatMember, *models.ChatMemberUpdated, From as value
	case update.ChatMember != nil:
		return &update.ChatMember.From

	// update.ChatJoinRequest, *models.ChatJoinRequest, From as value
	case update.ChatJoinRequest != nil:
		return &update.ChatJoinRequest.From

	// update.ChatBoost, *models.ChatBoostUpdated, no From
	// update.RemovedChatBoost, *models.ChatBoostRemoved, no From

	default:
		return nil
	}
}
