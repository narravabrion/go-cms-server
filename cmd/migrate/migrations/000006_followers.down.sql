DROP TABLE IF EXISTS followers;


SELECT 
    p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username 
FROM posts p 
LEFT JOIN users u on p.user_id = u.id 
JOIN followers f ON f.follower_id = p.user_id OR p.user_id = 9
WHERE f.user_id = 9 OR P.user_id = 9
GROUP BY p.id, u.username
ORDER BY p.created_at DESC LIMIT 20;


below are my tables for posts, users and followers

user are relates to posts in that one user can create many posts and one post can only belong to one user 
users can follow each other once. The follower_id represents the user that has been followed


write a query to populate a users feed such that it returns all the post content plust the author of the post but the posts should only be from the authors that the user follows


go_cms-# \d  posts
 id         | bigint                      |           | not null | nextval('posts_id_seq'::regclass)
 title      | text                        |           | not null | 
 user_id    | bigint                      |           | not null | 
 content    | text                        |           | not null | 
 created_at | timestamp(0) with time zone |           | not null | now()
 tags       | character varying(128)[]    |           |          | 
 updated_at | timestamp(0) with time zone |           | not null | now()
 version    | integer                     |           |          | 0

go_cms-# \d users
 id         | bigint                      |           | not null | nextval('users_id_seq'::regclass)
 email      | citext                      |           | not null | 
 username   | character varying(255)      |           |          | 
 password   | bytea                       |           | not null | 
 created_at | timestamp(0) with time zone |           | not null | now()

go_cms-# \d followers
 user_id     | bigint                      |           | not null | 
 follower_id | bigint                      |           | not null | 
 created_at  | timestamp(0) with time zone |           | not null | now()



SELECT p.id AS post_id, p.title, p.content, p.created_at, p.updated_at, p.tags, p.version, u.id AS author_id, u.username AS author_username, u.email AS author_email FROM posts p
JOIN users u ON p.user_id = u.id
JOIN followers f ON f.follower_id = p.user_id  
WHERE f.user_id = $1
ORDER BY p.created_at DESC;


SELECT 
    u.id AS author_id, 
    u.username, 
    COUNT(p.id) AS total_posts, 
    ARRAY_AGG(p.title ORDER BY p.created_at DESC) AS post_titles
FROM posts p 
LEFT JOIN users u ON p.user_id = u.id 
JOIN followers f ON f.follower_id = p.user_id  
WHERE f.user_id = 9 OR p.user_id = 9  
GROUP BY u.id, u.username
ORDER BY total_posts DESC;
