package book_user_preferences

const (
	queryInsertNewBookUserPreferences = `
		INSERT INTO book_user_preferences 
		(
			user_id, 
			preferred_category
		) VALUES (?, ?)
	`

	queryCountUserPreferences = `
		SELECT COUNT(1)
		FROM book_user_preferences
		WHERE user_id = ? AND preferred_category = ?
	`

	queryDeleteExistingUserPreferences = `
		DELETE FROM book_user_preferences
		WHERE user_id = ? AND preferred_category = ?
	`
)
