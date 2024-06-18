-- commands to create grok-the-casbah database

create table article (
id int primary key,                                                      
title text not null,                                                      
body text not null,
timestamp int not null                                                        
);

create table comment (
id int primary key,                                                      
article_id int not null,                                                 
body text not null,                                                       
timestamp int not null,                                                   
foreign key (article_id) references article(id)                           
);

-- add some dummy data
INSERT INTO article (id, title, body, timestamp) VALUES (1, 'First Article', 'This is the first article.', strftime('%s', 'now'));
INSERT INTO article (id, title, body, timestamp) VALUES (2, 'Second Article', 'This is the second article.', strftime('%s', 'now'));
INSERT INTO comment (id, article_id, body, timestamp) VALUES (1, 1, 'This is a comment for article 1.', strftime('%s', 'now'));
INSERT INTO comment (id, article_id, body, timestamp) VALUES (2, 1, 'This is another comment for article 1', strftime('%s', 'now'));
INSERT INTO comment (id, article_id, body, timestamp) VALUES (3, 2, 'This is a comment for article 2.', strftime('%s', 'now'));
INSERT INTO comment (id, article_id, body, timestamp) VALUES (4, 2, 'This is another comment for article 2', strftime('%s', 'now'));
