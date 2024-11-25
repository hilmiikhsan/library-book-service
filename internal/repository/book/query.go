package book

const (
	queryInsertNewBook = `
		INSERT INTO books
		(
			title,
			author_id,
			category_id,
			isbn,
			description,
			published_date
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	queryFindBookByID = `
		SELECT
			id,
			author_id,
			category_id,
			isbn,
			description,
			published_date,
			created_at,
			updated_at
		FROM books
		WHERE id = ?
	`

	queryFindAllBook = `
		SELECT
			id,
			title,
			description,
			isbn,
			published_date
		FROM books
		ORDER BY updated_at DESC
		LIMIT ?
		OFFSET ?
	`

	queryDeleteBookByID = `
		DELETE FROM books WHERE id = ?
	`

	queryGetRecommendations = `
		SELECT 
			b.id, 
			b.title, 
			b.author_id, 
			b.category_id, 
			b.description, 
			b.published_date
		FROM books b
		INNER JOIN book_user_preferences p ON b.category_id = p.preferred_category
		WHERE p.user_id = ?
		ORDER BY b.published_date DESC
		LIMIT ? OFFSET ?
	`
)
