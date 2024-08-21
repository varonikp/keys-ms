CREATE TABLE users (
    id 	  int PRIMARY KEY NOT NULL AUTO_INCREMENT,
    login MEDIUMTEXT NOT NULL,
    password MEDIUMTEXT NOT NUll,
    has_admin_role boolean not null DEFAULT false
);

CREATE TABLE softwares (
    id 	 int PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name MEDIUMTEXT NOT NULL
);

CREATE TABLE licenses (
    id 	 int PRIMARY KEY NOT NULL AUTO_INCREMENT,
    software_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expire_at TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id)
);

ALTER TABLE `licenses` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
