-- ユーザー {"insertTiming":"新規ユーザー作成時"}
CREATE TABLE User (
  -- ユーザーID 
  UserID STRING(MAX) NOT NULL,
  -- サーバー内用ユーザーID 
  ServerUserID STRING(MAX),
  -- 公開ユーザーID 
  PublicUserID STRING(MAX),
  -- 作成日時 
  CreatedTime TIMESTAMP,
  -- 更新日時 
  UpdatedTime TIMESTAMP,
) PRIMARY KEY (UserID);

