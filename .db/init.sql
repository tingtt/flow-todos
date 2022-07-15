SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

--
-- Database: `flow-todos`
--

CREATE DATABASE IF NOT EXISTS `flow-todos` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `flow-todos`;

-- --------------------------------------------------------

--
-- Table structure for table `repeat_models`
--

CREATE TABLE `repeat_models` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `until` DATE DEFAULT NULL,
  `unit` VARCHAR(7) NOT NULL CHECK(`unit` IN('day','week','month')),
  `every_other` INT UNSIGNED DEFAULT NULL,
  `date` TINYINT(5) UNSIGNED DEFAULT NULL CHECK(`date` <= 31),
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

--
-- Table structure for table `repeat_days`
--

CREATE TABLE `repeat_days` (
  `repeat_model_id` BIGINT UNSIGNED,
  `day` TINYINT(3) UNSIGNED NOT NULL,
  `time` TIME DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`repeat_model_id`, `day`),
  FOREIGN KEY (`repeat_model_id`) REFERENCES `repeat_models` (`id`) ON DELETE CASCADE
);

--
-- Table structure for table `todos`
--

CREATE TABLE `todos` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `description` VARCHAR(255) DEFAULT NULL,
  `date` DATE DEFAULT NULL,
  `time` TIME DEFAULT NULL,
  `execution_time` INT DEFAULT NULL COMMENT 'minute',
  `sprint_id` BIGINT UNSIGNED DEFAULT NULL,
  `project_id` BIGINT UNSIGNED DEFAULT NULL,
  `completed` TINYINT(1) NOT NULL DEFAULT '0',
  `repeat_model_id` BIGINT UNSIGNED DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`repeat_model_id`) REFERENCES `repeat_models` (`id`) ON DELETE RESTRICT,
  PRIMARY KEY (id)
);