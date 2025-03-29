-- Cambiar a mysql_native_password para root
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'example';

-- Cambiar a mysql_native_password para el usuario "laliga"
ALTER USER 'laliga'@'%' IDENTIFIED WITH mysql_native_password BY 'laliga';

FLUSH PRIVILEGES;

-- Crear la tabla matches si no existe
CREATE TABLE IF NOT EXISTS matches (
  id INT AUTO_INCREMENT PRIMARY KEY,
  team_a VARCHAR(100),  -- Aquí se almacenará el valor del equipo local (HomeTeam)
  team_b VARCHAR(100),  -- Aquí se almacenará el valor del equipo visitante (AwayTeam)
  match_date DATE,
  score_a INT,
  score_b INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Agregar columnas adicionales para soporte de endpoints PATCH
ALTER TABLE matches
  ADD COLUMN IF NOT EXISTS goals INT DEFAULT 0,
  ADD COLUMN IF NOT EXISTS yellowcards INT DEFAULT 0,
  ADD COLUMN IF NOT EXISTS redcards INT DEFAULT 0,
  ADD COLUMN IF NOT EXISTS extratime INT DEFAULT 0;
