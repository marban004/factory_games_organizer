DROP USER IF EXISTS 'users_microservice';
FLUSH PRIVILEGES;
CREATE USER 'users_microservice'@'%' IDENTIFIED BY 'bxu7%^yhag##KKL';
GRANT INSERT, SELECT, UPDATE, DELETE ON users.users TO 'users_microservice'@'%';