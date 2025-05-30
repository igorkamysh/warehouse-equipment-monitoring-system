CREATE TABLE IF NOT EXISTS parkings(
  id SERIAL,
  name text NOT NULL UNIQUE,
  mac_addr varchar(20) NOT NULL UNIQUE,
  machines integer DEFAULT 0,
  capacity integer DEFAULT 0,
  state integer DEFAULT 1,

  CHECK (state IN (0, 1)),
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS machines(
	id varchar(16) NOT NULL,
	state integer DEFAULT 0,
  parking_id integer DEFAULT 0,
  voltage integer DEFAULT 0,
  ip_addr varchar(22) NOT NULL,
		
	CHECK (state IN (0, 1, 2)),
  CHECK (voltage >= 0),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users(
  id SERIAL,
  name text NOT NULL,
	phone_number varchar(11) NOT NULL UNIQUE,
	job_position varchar(8) NOT NULL,
  password varchar (128) NOT NUll,

	CHECK (job_position IN ('worker', 'admin')),
  CHECK (LENGTH(password) >= 8),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sessions(
  id SERIAL,
  state integer DEFAULT 0,
  machine_id varchar(16) NOT NULL,
  worker_id integer NOT NULL,
  datetime_start bigint NOT NULL,
  datetime_finish bigint NOT NULL,

  CHECK (state IN (0, 1, 2)),

  PRIMARY KEY (id),

  FOREIGN KEY (machine_id) REFERENCES machines (id) ON DELETE CASCADE,
  FOREIGN KEY (worker_id) REFERENCES users (id) ON DELETE CASCADE
);
