DROP USER IF EXISTS 'calculator_microservice';
CREATE USER 'calculator_microservice'@'%' IDENTIFIED BY 'yixnhg64G0.*hafc2^';
GRANT SELECT ON users_data.* TO 'calculator_microservice'@'%';
CREATE USER 'crud_microservice'@'%' IDENTIFIED BY 'juG56#ian>LK90';
GRANT INSERT, SELECT, UPDATE, DELETE ON users_data.* TO 'crud_microservice'@'%';