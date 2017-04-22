-- CREATE TABLE Trivia (
--     Id              int         not null AUTO_INCREMENT PRIMARY KEY,
--     Channel         varchar(20) not null,
--     CurrentQuestion int         not null
-- ) ENGINE=INNODB;

-- CREATE TABLE QuestionLineup (
--     question_id          int(255)          not null,
--     trivia_id            int               not null
-- ) ENGINE=INNODB;

-- CREATE TABLE question (
--     id                int(255)          not null AUTO_INCREMENT PRIMARY KEY,
--     Category          varchar(255)      not null,
--     Difficulty        varchar(255)      not null,
--     Text              varchar(255)      not null
-- ) ENGINE=INNODB;

-- CREATE TABLE QuestionAnswers (
--     Id                int(255)          not null AUTO_INCREMENT PRIMARY KEY,
--     Question          int          not null,
--     Text              varchar(255) not null,
--     Correct           bool         not null
-- ) ENGINE=INNODB;

-- CREATE TABLE TriviaScore (
--     Trivia int not null,
--     UserId varchar(50) not null,
--     Score int not null
-- ) ENGINE=INNODB;