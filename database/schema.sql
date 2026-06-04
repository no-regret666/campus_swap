-- 校园闲置物品交换系统数据库设计，共 11 张表，满足课程设计不少于 8 张表要求。
-- 可用于 MySQL 8.x / PostgreSQL 参考建表，演示版后端使用 data/db.json 存储同构数据。

CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY,
  student_no VARCHAR(32) NOT NULL UNIQUE,
  password_hash VARCHAR(128) NOT NULL,
  name VARCHAR(50) NOT NULL,
  role VARCHAR(20) NOT NULL DEFAULT 'student',
  campus VARCHAR(50),
  major VARCHAR(80),
  credit INT NOT NULL DEFAULT 90,
  avatar VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  icon VARCHAR(20),
  sort_no INT DEFAULT 0
);

CREATE TABLE items (
  id VARCHAR(36) PRIMARY KEY,
  owner_id VARCHAR(36) NOT NULL,
  category_id VARCHAR(36) NOT NULL,
  title VARCHAR(120) NOT NULL,
  description TEXT NOT NULL,
  item_condition VARCHAR(50) NOT NULL,
  campus VARCHAR(50) NOT NULL,
  location VARCHAR(120) NOT NULL,
  want VARCHAR(160) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'available',
  views INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (owner_id) REFERENCES users(id),
  FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE item_tags (
  id VARCHAR(36) PRIMARY KEY,
  item_id VARCHAR(36) NOT NULL,
  tag_name VARCHAR(40) NOT NULL,
  FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE TABLE exchanges (
  id VARCHAR(36) PRIMARY KEY,
  item_id VARCHAR(36) NOT NULL,
  requester_id VARCHAR(36) NOT NULL,
  owner_id VARCHAR(36) NOT NULL,
  offered_item VARCHAR(160) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  meet_time VARCHAR(50),
  meet_place VARCHAR(120),
  message TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (item_id) REFERENCES items(id),
  FOREIGN KEY (requester_id) REFERENCES users(id),
  FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE messages (
  id VARCHAR(36) PRIMARY KEY,
  exchange_id VARCHAR(36) NOT NULL,
  sender_id VARCHAR(36) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (exchange_id) REFERENCES exchanges(id),
  FOREIGN KEY (sender_id) REFERENCES users(id)
);

CREATE TABLE favorites (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  item_id VARCHAR(36) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (user_id, item_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE TABLE reports (
  id VARCHAR(36) PRIMARY KEY,
  reporter_id VARCHAR(36) NOT NULL,
  item_id VARCHAR(36) NOT NULL,
  reason TEXT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'open',
  handled_by VARCHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (reporter_id) REFERENCES users(id),
  FOREIGN KEY (item_id) REFERENCES items(id),
  FOREIGN KEY (handled_by) REFERENCES users(id)
);

CREATE TABLE ratings (
  id VARCHAR(36) PRIMARY KEY,
  exchange_id VARCHAR(36) NOT NULL,
  rater_id VARCHAR(36) NOT NULL,
  rated_user_id VARCHAR(36) NOT NULL,
  score INT NOT NULL,
  content TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (exchange_id) REFERENCES exchanges(id),
  FOREIGN KEY (rater_id) REFERENCES users(id),
  FOREIGN KEY (rated_user_id) REFERENCES users(id)
);

CREATE TABLE notifications (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  title VARCHAR(100) NOT NULL,
  content TEXT,
  read_flag BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE audit_logs (
  id VARCHAR(36) PRIMARY KEY,
  operator_id VARCHAR(36),
  action VARCHAR(80) NOT NULL,
  detail TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (operator_id) REFERENCES users(id)
);
