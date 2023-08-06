CREATE SCHEMA MTGREPORTS;
USE MTGREPORTS;

DROP TABLE IF EXISTS cards_details;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS prices;

CREATE TABLE `cards` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `set_name` varchar(255) NOT NULL,
    `collector_number` varchar(255) NOT NULL,
    `foil` tinyint NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_cards_name` (`name`),
    UNIQUE INDEX `unique_idx` (`set_name`, `collector_number`, `foil`)
) AUTO_INCREMENT = 1 DEFAULT CHARSET = latin1;

CREATE TABLE `cards_details` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `card_id` int unsigned NOT NULL,
    `last_price` decimal(10,2) NOT NULL DEFAULT 0,
    `old_price` decimal(10,2) NOT NULL DEFAULT 0,
    `price_change` decimal(10,2) NOT NULL DEFAULT 0,
    `last_update` datetime NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_card_details_card_id_last_update` (`card_id`, `last_update`),
    CONSTRAINT `fk_card_id`
        FOREIGN KEY (`card_id`)
        REFERENCES `cards` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) AUTO_INCREMENT = 1 DEFAULT CHARSET = latin1;

CREATE TABLE `prices` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `old_price` decimal(10,2) NOT NULL DEFAULT 0,
    `new_price` decimal(10,2) NOT NULL DEFAULT 0,
    `price_change` decimal(10,2) NOT NULL DEFAULT 0,
    `last_update` datetime NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_last_update` (`last_update`)
) AUTO_INCREMENT = 1 DEFAULT CHARSET = latin1;



