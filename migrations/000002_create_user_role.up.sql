CREATE ROLE rw_user WITH LOGIN PASSWORD 'rwuser';
GRANT SELECT, INSERT ON ALL TABLES IN SCHEMA public TO rw_user;