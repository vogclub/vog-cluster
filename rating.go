package vogcluster

import (
	"errors"
	"time"
)

// RatingUpdated is published by vog-games on
// SubjectRatingUpdated(game_type) whenever a player's rating changes
// after a finished match. Game instances subscribed to the game type
// update their in-memory rating display.
type RatingUpdated struct {
	UserID     string    `json:"user_id"`
	GameType   string    `json:"game_type"`
	RatingType string    `json:"rating_type"`
	OldRating  int       `json:"old_rating"`
	NewRating  int       `json:"new_rating"`
	GameID     string    `json:"game_id"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// Validate reports whether the message is well-formed.
func (m RatingUpdated) Validate() error {
	if m.UserID == "" {
		return errors.New("vogcluster: RatingUpdated.user_id is required")
	}
	if m.GameType == "" {
		return errors.New("vogcluster: RatingUpdated.game_type is required")
	}
	if m.RatingType == "" {
		return errors.New("vogcluster: RatingUpdated.rating_type is required")
	}
	if m.GameID == "" {
		return errors.New("vogcluster: RatingUpdated.game_id is required")
	}
	return nil
}
