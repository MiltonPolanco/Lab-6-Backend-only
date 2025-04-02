CREATE TABLE IF NOT EXISTS series (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    status VARCHAR(100) NOT NULL,
    episodes_watched INT DEFAULT 0,
    total_episodes INT DEFAULT 0,
    ranking INT DEFAULT 0
);
