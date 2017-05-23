DROP TABLE IF EXISTS "admins" CASCADE;
DROP TABLE IF EXISTS "actions" CASCADE;
DROP TABLE IF EXISTS "videos" CASCADE;
DROP TABLE IF EXISTS "player_position_game_team" CASCADE;
DROP TABLE IF EXISTS "games" CASCADE;
DROP TABLE IF EXISTS "player_team" CASCADE;
DROP TABLE IF EXISTS "coach_team" CASCADE;
DROP TABLE IF EXISTS "metrics" CASCADE;
DROP TABLE IF EXISTS "teams" CASCADE;
DROP TABLE IF EXISTS "zones" CASCADE;
DROP TABLE IF EXISTS "sports" CASCADE;
DROP TABLE IF EXISTS "players" CASCADE;
DROP TABLE IF EXISTS "locations" CASCADE;
DROP TABLE IF EXISTS "field_types" CASCADE;
DROP TABLE IF EXISTS "categories" CASCADE;
DROP TABLE IF EXISTS "coaches" CASCADE;
DROP TABLE IF EXISTS "actions_type" CASCADE;
DROP TABLE IF EXISTS "seasons" CASCADE;
DROP TABLE IF EXISTS "positions" CASCADE;
DROP TABLE IF EXISTS "movements_type" CASCADE;



-- -----------------------------------------------------
-- Table "admins"
-- -----------------------------------------------------

CREATE TABLE "admins" (
  "id" SERIAL PRIMARY KEY,
  "email" VARCHAR(256) NOT NULL,
  "pass_hash" VARCHAR(256) NOT NULL,
  "token_reset" VARCHAR(256) NULL,
  "token_login" VARCHAR(256) NULL);



-- -----------------------------------------------------
-- Table "seasons"
-- -----------------------------------------------------

CREATE TABLE "seasons" (
  "id" SERIAL PRIMARY KEY,
  "years" VARCHAR(10) NOT NULL);


-- -----------------------------------------------------
-- Table "sports"
-- -----------------------------------------------------

CREATE TABLE "sports" (
  "id" VARCHAR(3) PRIMARY KEY,
  "name" VARCHAR(256) NOT NULL);


-- -----------------------------------------------------
-- Table "categories"
-- -----------------------------------------------------

CREATE TABLE "categories" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(256) NULL);


-- -----------------------------------------------------
-- Table "teams"
-- -----------------------------------------------------

CREATE TABLE "teams" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(256) NULL,
  "city" VARCHAR(256) NULL,
  "id_sport" VARCHAR(3) NULL,
  "id_category" INT NULL,
  CONSTRAINT "fk_sport_team"
    FOREIGN KEY ("id_sport")
    REFERENCES "sports" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_team_category"
    FOREIGN KEY ("id_category")
    REFERENCES "categories" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE);


-- -----------------------------------------------------
-- Table "players"
-- -----------------------------------------------------

CREATE TABLE "players" (
  "id" SERIAL PRIMARY KEY,
  "number" INT NOT NULL,
  "email" VARCHAR(256) NULL,
  "fname" VARCHAR(256) NULL,
  "lname" VARCHAR(256) NULL,
  "pass_hash" VARCHAR(256) NULL,
  "token_request" VARCHAR(256) NULL,
  "token_reset" VARCHAR(256) NULL,
  "token_login" VARCHAR(256) NULL);


-- -----------------------------------------------------
-- Table "locations"
-- -----------------------------------------------------

CREATE TABLE "locations" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(256) NOT NULL,
  "city" VARCHAR(256) NOT NULL,
  "address" VARCHAR(256) NULL,
  "inside_outside" VARCHAR(256) NOT NULL);



-- -----------------------------------------------------
-- Table "field_types"
-- -----------------------------------------------------

CREATE TABLE "field_types" (
  "id" SERIAL PRIMARY KEY,
  "type" VARCHAR(256) NOT NULL,
  "description" VARCHAR(256) NULL);



-- -----------------------------------------------------
-- Table "games"
-- -----------------------------------------------------

CREATE TABLE "games" (
  "id" SERIAL PRIMARY KEY,
  "id_team" INT NOT NULL,
  "status" VARCHAR(50) NOT NULL,
  "opposing_team" VARCHAR(100) NOT NULL,
  "id_season" INT NOT NULL,
  "id_location" INT NOT NULL,
  "field_condition" VARCHAR(45) NULL,
  "temperature" VARCHAR(45) NULL,
  "degree" VARCHAR(10) NULL,
  "date" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT "fk_game_team"
    FOREIGN KEY ("id_team")
    REFERENCES "teams" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_game_team_location"
    FOREIGN KEY ("id_location")
    REFERENCES "locations" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_game_team_season"
    FOREIGN KEY ("id_season")
    REFERENCES "seasons" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE);



-- -----------------------------------------------------
-- Table "video"
-- -----------------------------------------------------

CREATE TABLE "videos" (
  "id" SERIAL PRIMARY KEY,
  "path" TEXT NOT NULL,
  "part" INT NOT NULL DEFAULT 1,
  "completed" INT NOT NULL DEFAULT 0,
  "id_game" INT NOT NULL,
  CONSTRAINT "fk_game_id"
    FOREIGN KEY ("id_game")
    REFERENCES "games" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE);



-- -----------------------------------------------------
-- Table "positions"
-- -----------------------------------------------------

CREATE TABLE "positions" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(45) NOT NULL);



-- -----------------------------------------------------
-- Table "player_position_game_team"
-- -----------------------------------------------------

CREATE TABLE "player_position_game_team" (
  "id" SERIAL PRIMARY KEY,
  "id_player" INT NOT NULL,
  "id_game" INT NOT NULL,
  "id_position" INT NOT NULL,
  "id_team" INT NOT NULL,
  CONSTRAINT "fk_player_position"
    FOREIGN KEY ("id_player")
    REFERENCES "players" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_game"
    FOREIGN KEY ("id_game")
    REFERENCES "games" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_position"
    FOREIGN KEY ("id_position")
    REFERENCES "positions" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_team"
    FOREIGN KEY ("id_team")
    REFERENCES "teams" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE);



-- -----------------------------------------------------
-- Table "zones"
-- -----------------------------------------------------

CREATE TABLE "zones" (
  "id" VARCHAR(3) PRIMARY KEY,
  "name" VARCHAR(45) NOT NULL);



-- -----------------------------------------------------
-- Table "movements_type"
-- -----------------------------------------------------

CREATE TABLE "movements_type" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(10) NOT NULL);



-- -----------------------------------------------------
-- Table "actions_type"
-- -----------------------------------------------------

CREATE TABLE "actions_type" (
  "id" VARCHAR(5) PRIMARY KEY,
  "name" VARCHAR(256) NOT NULL,
  "description" TEXT NULL,
  "id_movement_type" INT NOT NULL,
  CONSTRAINT "fk_movement_type"
    FOREIGN KEY ("id_movement_type")
    REFERENCES "movements_type" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE);



-- -----------------------------------------------------
-- Table "actions"
-- -----------------------------------------------------

CREATE TABLE "actions" (
  "id" SERIAL PRIMARY KEY ,
  "id_action_type" VARCHAR(5) NOT NULL,
  "id_reception_type" VARCHAR(5) NOT NULL
  "id_player" INT  NOT NULL,
  "id_zone" VARCHAR(3) NOT NULL,
  "id_game" INT NOT NULL,
  "is_positive" INT NOT NULL,
  "x1" FLOAT NOT NULL,
  "y1" FLOAT NOT NULL,
  "x2" FLOAT NOT NULL,
  "y2" FLOAT NOT NULL,
  "x3" FLOAT NOT NULL,
  "y3" FLOAT NOT NULL,
  "time" FLOAT NOT NULL,
  "home_score" INT NOT NULL DEFAULT 0,
  "guest_score" INT NOT NULL DEFAULT 0,
  CONSTRAINT "fk_action_position"
    FOREIGN KEY ("id_zone")
    REFERENCES "zones" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_action_name"
    FOREIGN KEY ("id_action_type")
    REFERENCES "actions_type" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_action_game"
    FOREIGN KEY ("id_game")
    REFERENCES "games" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_reception_type"
    FOREIGN KEY ("id_reception_type")
    REFERENCES "reception_type" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_action_player"
    FOREIGN KEY ("id_player")
    REFERENCES "players" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE);



-- -----------------------------------------------------
-- Table "player_team"
-- -----------------------------------------------------

CREATE TABLE "player_team" (
  "id" SERIAL PRIMARY KEY ,
  "id_player" INT NOT NULL,
  "id_team" INT NOT NULL,
  CONSTRAINT "fk_player"
    FOREIGN KEY ("id_player")
    REFERENCES "players" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_team"
    FOREIGN KEY ("id_team")
    REFERENCES "teams" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE);

-- -----------------------------------------------------
-- Table "reception_type"
-- -----------------------------------------------------

CREATE TABLE "reception_type" (
  "id" VARCHAR(5) PRIMARY KEY,
  "name" VARCHAR(256) NOT NULL);

-- -----------------------------------------------------
-- Table "coaches"
-- -----------------------------------------------------

CREATE TABLE "coaches" (
  "id" SERIAL PRIMARY KEY ,
  "fname" VARCHAR(256) NOT NULL,
  "lname" VARCHAR(256) NOT NULL,
  "email" VARCHAR(256) NULL,
  "pass_hash" VARCHAR(256) NULL,
  "token_request" VARCHAR(256) NULL,
  "token_reset" VARCHAR(256) NULL,
  "token_login" VARCHAR(256) NULL);



-- -----------------------------------------------------
-- Table "coach_team"
-- -----------------------------------------------------

CREATE TABLE "coach_team" (
  "id" SERIAL PRIMARY KEY ,
  "id_coach" INT NOT NULL,
  "id_team" INT NOT NULL,
  CONSTRAINT "fk_mandat_coach"
    FOREIGN KEY ("id_coach")
    REFERENCES "coaches" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "fk_mandat_team"
    FOREIGN KEY ("id_team")
    REFERENCES "teams" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE);


  
-- -----------------------------------------------------
-- Table "metrics"
-- -----------------------------------------------------

CREATE TABLE "metrics" (
  "name" VARCHAR(256) NOT NULL,
  "equation" TEXT NOT NULL,
  "id_team" INT NOT NULL,
  PRIMARY KEY ("name", "equation"),
  CONSTRAINT "fk_metric_team"
    FOREIGN KEY ("id_team")
    REFERENCES "teams" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE);