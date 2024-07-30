CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    isAdmin TINYINT(1) DEFAULT 0,
    admin_request_status ENUM('pending', 'approved', 'denied') NULL
);

CREATE TABLE books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    isbn VARCHAR(25) UNIQUE NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    quantity INT DEFAULT 1,
    available TINYINT(1) GENERATED ALWAYS AS (quantity > 0) STORED
);

CREATE TABLE checkouts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    book_id INT,
    checkout_date DATETIME NOT NULL,
    due_date DATETIME NOT NULL,
    return_date DATETIME NULL,
    fine DECIMAL(5,2) DEFAULT 0.00,
    checkout_status ENUM('pending', 'approved', 'denied') DEFAULT 'pending',
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO books (isbn, title, author, quantity)
VALUES 
('978-3-16-148410-0', 'Harry Potter and the Philosophers Stone', 'J.K. Rowling', 20),
('978-3-16-148411-7', '1984', 'George Orwell', 15),
('978-3-16-148412-4', 'To Kill a Mockingbird', 'Harper Lee', 25),
('978-3-16-148413-1', 'The Great Gatsby', 'F. Scott Fitzgerald', 18),
('978-3-16-148414-8', 'Moby Dick', 'Herman Melville', 22),
('978-3-16-148415-5', 'Pride and Prejudice', 'Jane Austen', 30),
('978-3-16-148416-2', 'The Catcher in the Rye', 'J.D. Salinger', 28),
('978-3-16-148417-9', 'The Hobbit', 'J.R.R. Tolkien', 24),
('978-3-16-148418-6', 'The Lord of the Rings', 'J.R.R. Tolkien', 26),
('978-3-16-148419-3', 'The Hunger Games', 'Suzanne Collins', 35);


INSERT INTO users(username, password, isAdmin) VALUES ("admin", "$2a$10$LlWH7dGow.Xf5YoJHVE7M.K5yKgaMtF3HSrdnzD3iCXwJpnwHWWoe", 2);