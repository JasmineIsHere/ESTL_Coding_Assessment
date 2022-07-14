DROP TABLE  `employees` IF EXISTS;

CREATE TABLE `employees` (
                             `id` varchar(16) NOT NULL,
                             `login` varchar(128) NOT NULL,
                             `name` varchar(128) NOT NULL,
                             `salary` double NOT NULL,
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `login` (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;