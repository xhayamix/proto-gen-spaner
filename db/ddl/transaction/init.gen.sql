CREATE TABLE User (
  UserID STRING(MAX) NOT NULL,
  PublicUserID STRING(MAX),
  CreatedTime TIMESTAMP,
  UpdatedTime TIMESTAMP,
) PRIMARY KEY (UserID);

