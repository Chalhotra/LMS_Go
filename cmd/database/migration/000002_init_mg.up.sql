ALTER TABLE books ADD current_borrowed_count INT DEFAULT 0;
UPDATE books b
SET current_borrowed_count = (
    SELECT COUNT(*)
    FROM checkouts c
    WHERE c.book_id = b.id AND c.return_date IS NULL
);

ALTER TABLE books ADD available_quantity INT AS (quantity - current_borrowed_count) STORED;