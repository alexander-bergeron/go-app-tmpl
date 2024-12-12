-- Create the users table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    version integer NOT NULL DEFAULT 1
);

INSERT INTO users (username, email, first_name, last_name) VALUES
('Cow31337Killer', 'ckilla@hotmail.com', 'Cow', 'Killer'),
('Durial321', 'backslash@yahoo.com', 'Durial', '321'),
('BigRedJapan', 'brj@gmail.com', 'BigRed', 'Japan');
