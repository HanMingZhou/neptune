create table jwt_blacklists
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null,
    jwt        text        null comment 'jwt'
);

create index idx_jwt_blacklists_deleted_at
    on jwt_blacklists (deleted_at);

INSERT INTO aiInfra.jwt_blacklists (id, created_at, updated_at, deleted_at, jwt) VALUES (56, '2026-03-25 23:44:04.668', '2026-03-25 23:44:04.668', null, 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiYzNlOWYxMTMtZjJmOC00ODU2LThmYzEtMDgwODYwNTE0NTU5IiwiSUQiOjEsIlVzZXJuYW1lIjoiYWRtaW4iLCJOaWNrTmFtZSI6IiIsIkF1dGhvcml0eUlkIjo4ODgsIk5hbWVzcGFjZSI6ImhteiIsIkJ1ZmZlclRpbWUiOjg2NDAwLCJpc3MiOiJxbVBsdXMiLCJhdWQiOlsiR1ZBIl0sImV4cCI6MTc3NTA0OTM0MSwibmJmIjoxNzc0NDQ0NTQxfQ.dYsH2Knp01gg6gOqKzYexficv4xhJg0vKRCp1dPOy1A');
INSERT INTO aiInfra.jwt_blacklists (id, created_at, updated_at, deleted_at, jwt) VALUES (57, '2026-03-25 23:44:16.981', '2026-03-25 23:44:16.981', null, 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiYzNlOWYxMTMtZjJmOC00ODU2LThmYzEtMDgwODYwNTE0NTU5IiwiSUQiOjEsIlVzZXJuYW1lIjoiYWRtaW4iLCJOaWNrTmFtZSI6IiIsIkF1dGhvcml0eUlkIjo4ODgsIk5hbWVzcGFjZSI6ImhteiIsIkJ1ZmZlclRpbWUiOjg2NDAwLCJpc3MiOiJxbVBsdXMiLCJhdWQiOlsiR1ZBIl0sImV4cCI6MTc3NTA0OTM0MSwibmJmIjoxNzc0NDQ0NTQxfQ.dYsH2Knp01gg6gOqKzYexficv4xhJg0vKRCp1dPOy1A');
INSERT INTO aiInfra.jwt_blacklists (id, created_at, updated_at, deleted_at, jwt) VALUES (58, '2026-03-25 23:46:30.301', '2026-03-25 23:46:30.301', null, 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiYzNlOWYxMTMtZjJmOC00ODU2LThmYzEtMDgwODYwNTE0NTU5IiwiSUQiOjEsIlVzZXJuYW1lIjoiYWRtaW4iLCJOaWNrTmFtZSI6IiIsIkF1dGhvcml0eUlkIjo4ODgsIk5hbWVzcGFjZSI6ImhteiIsIkJ1ZmZlclRpbWUiOjg2NDAwLCJpc3MiOiJxbVBsdXMiLCJhdWQiOlsiR1ZBIl0sImV4cCI6MTc3NTA1ODI1NiwibmJmIjoxNzc0NDUzNDU2fQ.K-RtBCNg-qFjWd39HrjNa3HThYkjjqycSOKiBX7k1xU');
INSERT INTO aiInfra.jwt_blacklists (id, created_at, updated_at, deleted_at, jwt) VALUES (59, '2026-03-25 23:46:52.137', '2026-03-25 23:46:52.137', null, 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiYzNlOWYxMTMtZjJmOC00ODU2LThmYzEtMDgwODYwNTE0NTU5IiwiSUQiOjEsIlVzZXJuYW1lIjoiYWRtaW4iLCJOaWNrTmFtZSI6IiIsIkF1dGhvcml0eUlkIjo4ODgsIk5hbWVzcGFjZSI6ImhteiIsIkJ1ZmZlclRpbWUiOjg2NDAwLCJpc3MiOiJxbVBsdXMiLCJhdWQiOlsiR1ZBIl0sImV4cCI6MTc3NTA1ODI1NiwibmJmIjoxNzc0NDUzNDU2fQ.K-RtBCNg-qFjWd39HrjNa3HThYkjjqycSOKiBX7k1xU');
INSERT INTO aiInfra.jwt_blacklists (id, created_at, updated_at, deleted_at, jwt) VALUES (60, '2026-03-27 23:05:08.500', '2026-03-27 23:05:08.500', null, 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiYzNlOWYxMTMtZjJmOC00ODU2LThmYzEtMDgwODYwNTE0NTU5IiwiSUQiOjEsIlVzZXJuYW1lIjoiYWRtaW4iLCJOaWNrTmFtZSI6IiIsIkF1dGhvcml0eUlkIjo4ODgsIk5hbWVzcGFjZSI6ImhteiIsIkJ1ZmZlclRpbWUiOjg2NDAwLCJpc3MiOiJxbVBsdXMiLCJhdWQiOlsiR1ZBIl0sImV4cCI6MTc3NTIyNzk4NCwibmJmIjoxNzc0NjIzMTg0fQ.kMXOx3Ak5uwyuly36EWdJiMz-4QDTN9LB2JGEcS5OwA');
INSERT INTO aiInfra.jwt_blacklists (id, created_at, updated_at, deleted_at, jwt) VALUES (61, '2026-03-28 16:22:51.531', '2026-03-28 16:22:51.531', null, 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiYzNlOWYxMTMtZjJmOC00ODU2LThmYzEtMDgwODYwNTE0NTU5IiwiSUQiOjEsIlVzZXJuYW1lIjoiYWRtaW4iLCJOaWNrTmFtZSI6IiIsIkF1dGhvcml0eUlkIjo4ODgsIk5hbWVzcGFjZSI6ImhteiIsIkJ1ZmZlclRpbWUiOjg2NDAwLCJpc3MiOiJxbVBsdXMiLCJhdWQiOlsiR1ZBIl0sImV4cCI6MTc3NTIyODcwOCwibmJmIjoxNzc0NjIzOTA4fQ.LwmqDZo-IYYohS1-yfwaqH3N23PqNz31YIEh7d7CR1Y');
