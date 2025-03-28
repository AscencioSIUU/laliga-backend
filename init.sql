-- Cambiar a mysql_native_password para root
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'example';

-- Cambiar a mysql_native_password para el usuario "laliga"
ALTER USER 'laliga'@'%' IDENTIFIED WITH mysql_native_password BY 'laliga';

FLUSH PRIVILEGES;


CREATE TABLE IF NOT EXISTS matches (
  id INT AUTO_INCREMENT PRIMARY KEY,
  team_a VARCHAR(100),
  team_b VARCHAR(100),
  score_a INT,
  score_b INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
