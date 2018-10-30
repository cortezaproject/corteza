package config

// @todo need to decide on settings format & structure...
// type (
// 	// MessagingSettings holds configuration settings for mesaging service
// 	MessagingSettings struct {
// 		Messages struct {
// 			Body struct {
// 				// How long can a message be
// 				MaxLength uint `json:"maxLength,omitempty"`
//
// 				// On render, convert textual [:)] emoji to graph. empji
// 				EmojiConvertToPicture bool `json:"emojiConvertToPicture,omitempty"`
// 			} `json:"body,omitempty"`
//
// 			Avatars struct {
// 				// Display image that users use for the avatar
// 				// if false, it will use initials only
// 				DisplaySelectedImage bool `json:"displaySelectedImage,omitempty"`
//
// 				// Enable avatar
// 				Enabled bool `json:"enabled,omitempty"`
// 			} `json:"avatars,omitempty"`
// 		} `json:"messages,omitempty"`
//
// 		Channels struct {
// 			Name struct {
// 				// How long can a name be
// 				MaxLength uint `json:"maxLength,omitempty"`
//
// 				// Can we have spaces in channel's name
// 				AllowSpaces bool `json:"allowSpaces,omitempty"`
// 			} `json:"name,omitempty"`
//
// 			Topic struct {
// 				// How long can a name be
// 				MaxLength uint `json:"maxLength,omitempty"`
//
// 				// Can we have spaces in channel's name
// 				Enabled bool `json:"enabled,omitempty"`
// 			} `json:"topic,omitempty"`
// 		} `json:"channels,omitempty"`
// 	}
// )
//
// func DefaultMessagingSettings() (s MessagingSettings) {
// 	s.Messages.Body.MaxLength = 10000
// 	s.Messages.Body.EmojiConvertToPicture = true
// 	s.Messages.Avatars.DisplaySelectedImage = true
// 	s.Messages.Avatars.Enabled = true
// 	s.Channels.Name.MaxLength = 40
// 	s.Channels.Name.AllowSpaces = true
// 	s.Channels.Topic.MaxLength = 200
// 	s.Channels.Topic.Enabled = true
//
// 	return s
// }
