CREATE TABLE `user`(
  id INT NOT NULL,
  email VARCHAR,
  phone_number VARCHAR,
  `password` VARCHAR,
  first_name VARCHAR,
  middle_name VARCHAR,
  last_name VARCHAR,
  date_of_birth DATE,
  hashtag_name VARCHAR,
  gender VARCHAR,
  `language` VARCHAR,
  country VARCHAR,
  profile_description VARCHAR,
  avatar_image VARCHAR,
  cover_image VARCHAR,
  validate_code VARCHAR,
  account_status VARCHAR,
  verification_status VARCHAR,
  create_at DATETIME,
  last_update DATETIME,
  last_active DATETIME,
  device_token VARCHAR,
  job_name VARCHAR,
  time_zone VARCHAR,
  PRIMARY KEY(id)
);

CREATE TABLE message(
  id INT NOT NULL,
  sender_id INT,
  content VARCHAR,
  create_at DATETIME,
  last_update DATETIME,
  channel_id INT,
  is_edited BOOLEAN,
  `status` VARCHAR,
  PRIMARY KEY(id)
);

CREATE TABLE media(
  id INT NOT NULL,
  sender_id INT,
  message_id INT,
  channel_id INT,
  `type` VARCHAR,
  `url` VARCHAR,
  create_at INT,
  PRIMARY KEY(id)
);

CREATE TABLE `channel`(
  id INT NOT NULL,
  `name` VARCHAR,
  creator_id INT,
  create_at DATETIME,
  invite_code VARCHAR,
  last_message_id INT,
  channel_image VARCHAR,
  `type` VARCHAR,
  is_allow_invvite_code BOOLEAN,
  PRIMARY KEY(id)
);

CREATE TABLE channel_member(
  id INT NOT NULL,
  user_id INT,
  channel_id INT,
  join_at DATETIME,
  `role` VARCHAR,
  inchannel_name VARCHAR,
  PRIMARY KEY(id)
);

CREATE TABLE friend(
  id INT NOT NULL,
  first_user_id INT,
  second_user_id INT,
  create_at DATETIME,
  PRIMARY KEY(id)
);

CREATE TABLE friend_request(
  id INT NOT NULL,
  sender_id INT,
  receiver_id INT,
  create_at DATETIME,
  PRIMARY KEY(id)
);

ALTER TABLE message
  ADD CONSTRAINT id_sender_id FOREIGN KEY (sender_id) REFERENCES `user` (id);

ALTER TABLE media
  ADD CONSTRAINT id_sender_id FOREIGN KEY (sender_id) REFERENCES `user` (id);

ALTER TABLE media
  ADD CONSTRAINT id_message_id FOREIGN KEY (message_id) REFERENCES message (id);

ALTER TABLE media
  ADD CONSTRAINT id_channel_id FOREIGN KEY (channel_id) REFERENCES `channel` (id)
  ;

ALTER TABLE `channel`
  ADD CONSTRAINT id_creator_id FOREIGN KEY (creator_id) REFERENCES `user` (id);

ALTER TABLE `channel`
  ADD CONSTRAINT id_last_message_id
    FOREIGN KEY (last_message_id) REFERENCES message (id);

ALTER TABLE message
  ADD CONSTRAINT id_channel_id FOREIGN KEY (channel_id) REFERENCES `channel` (id)
  ;

ALTER TABLE channel_member
  ADD CONSTRAINT id_channel_id FOREIGN KEY (channel_id) REFERENCES `channel` (id)
  ;

ALTER TABLE channel_member
  ADD CONSTRAINT id_user_id FOREIGN KEY (user_id) REFERENCES `user` (id);

ALTER TABLE friend
  ADD CONSTRAINT id_first_user_id
    FOREIGN KEY (first_user_id) REFERENCES `user` (id);

ALTER TABLE friend
  ADD CONSTRAINT id_second_user_id
    FOREIGN KEY (second_user_id) REFERENCES `user` (id);

ALTER TABLE friend_request
  ADD CONSTRAINT id_sender_id FOREIGN KEY (sender_id) REFERENCES `user` (id);

ALTER TABLE friend_request
  ADD CONSTRAINT id_receiver_id FOREIGN KEY (receiver_id) REFERENCES `user` (id)
  ;
