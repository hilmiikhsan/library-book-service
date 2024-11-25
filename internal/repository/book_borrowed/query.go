package book_borrowed

const (
	queryInsertNewBookBorrowed = `
		INSERT INTO borrowed_books
		(
			user_id,
			book_id,
			due_date
		) VALUES (?, ?, ?)
	`

	queryValidateBookBorrowed = `
		SELECT COUNT(id) 
		FROM borrowed_books 
		WHERE book_id = ? AND user_id = ? AND returned_date IS NULL
	`

	queryUpdateBookReturned = `
		UPDATE borrowed_books 
		SET 
			returned_date = ?,
			updated_at = NOW()
		WHERE book_id = ? AND user_id = ?
	`

	queryValidateBookReturned = `
		SELECT COUNT(id)
		FROM borrowed_books
		WHERE book_id = ? AND user_id = ? AND returned_date IS NOT NULL
	`
)
