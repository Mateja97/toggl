DROP TABLE IF EXISTS question;
CREATE TABLE question (
id INTEGER PRIMARY KEY AUTOINCREMENT,
body TEXT
);
DROP TABLE IF EXISTS option;
CREATE TABLE option (
id INTEGER PRIMARY KEY AUTOINCREMENT,
questionid INTEGER,
body TEXT,
correct INTEGER,
FOREIGN KEY (questionid) REFERENCES question(id)
);

INSERT INTO question (id,body)
VALUES(1,"Where does the sun set?");

INSERT INTO question (id,body)
VALUES(2,"But what is the ultimate question?");

INSERT INTO option (questionid,body,correct)
VALUES(1,"East",1);

INSERT INTO option (questionid,body,correct)
VALUES(1,"West",0);

INSERT INTO option (questionid,body,correct)
VALUES(2,"Question1",0);

INSERT INTO option (questionid,body,correct)
VALUES(2,"Question2",1);