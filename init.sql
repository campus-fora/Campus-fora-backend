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

CREATE DATABASE likes;
CREATE ROLE likesadmin WITH LOGIN PASSWORD 'likekaadmin';
GRANT ALL PRIVILEGES ON DATABASE likes TO likesadmin;
\c likes postgres
GRANT ALL ON SCHEMA public TO likesadmin;

CREATE DATABASE stats;
CREATE ROLE statsadmin WITH LOGIN PASSWORD 'statskaadmin';
GRANT ALL PRIVILEGES ON DATABASE stats TO statsadmin;
\c stats postgres
GRANT ALL ON SCHEMA public TO statsadmin;