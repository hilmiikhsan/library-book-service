package BookStock

const (
	queryInsertNewBookStock = `
		INSERT INTO book_stocks
		(
			book_id,
			total_stock,
			available_stock
		) VALUES (?, ?, ?)
	`

	queryFindBookStockByID = `
		SELECT
			bs.id,
			bs.book_id,
			bs.total_stock,
			bs.available_stock,
			bs.created_at,
			bs.updated_at,
			b.id as book_id,
			b.title as book_title
		FROM book_stocks bs
		JOIN books b ON bs.book_id = b.id
		WHERE bs.id = ?
	`

	queryFindAllBookStock = `
		SELECT
			bs.id,
			bs.book_id,
			bs.total_stock,
			bs.available_stock,
			b.title AS book_title
		FROM book_stocks bs
		JOIN books b ON bs.book_id = b.id
		ORDER BY bs.updated_at DESC
		LIMIT ?
		OFFSET ?
	`

	queryUpdateBookStock = `
		UPDATE book_stocks
		SET
			book_id = ?,
			total_stock = ?,
			available_stock = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	queryDeleteBookStockByID = `
		DELETE FROM book_stocks WHERE id = ?
	`

	queryCountBookByBookID = `
		SELECT COUNT(id) FROM book_stocks WHERE book_id = ?
	`

	queryDecrementAvailableStock = `
		UPDATE book_stocks
		SET 
			available_stock = available_stock - ?
		WHERE book_id = ? 
		AND available_stock >= ?
	`

	queryLockBookStock = `
		SELECT 1
		FROM book_stocks
		WHERE book_id = ?
		FOR UPDATE
	`

	queryIncrementAvailableStock = `
		UPDATE book_stocks
		SET
			available_stock = available_stock + ?
		WHERE book_id = ?
	`

	queryLockBookStockReturned = `
		SELECT available_stock
		FROM book_stocks
		WHERE book_id = ?
		FOR UPDATE
	`
)
