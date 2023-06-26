CREATE DATABASE posts;
CREATE ROLE postsadmin WITH LOGIN PASSWORD 'postkaadmin';
GRANT ALL PRIVILEGES ON DATABASE posts TO postsadmin;
\c posts postgres
GRANT ALL ON SCHEMA public TO postsadmin;

CREATE DATABASE users;
CREATE ROLE usersadmin WITH LOGIN PASSWORD 'userkaadmin';
GRANT ALL PRIVILEGES ON DATABASE users TO usersadmin;
\c users postgres
GRANT ALL ON SCHEMA public TO usersadmin;