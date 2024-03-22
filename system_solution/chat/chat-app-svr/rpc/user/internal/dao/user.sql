CREATE schema "chat";

CREATE SEQUENCE "chat".user_id_seq MINVALUE 1 MAXVALUE 99999999999 INCREMENT BY 1 START WITH 1;
CREATE TABLE "chat"."user"
(
    "id"           int4                                       NOT NULL DEFAULT nextval('"chat".user_id_seq'::regclass),
    "nick_name"    varchar(255) COLLATE "pg_catalog"."default",
    "phone_number" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
    "profile"      varchar(255) COLLATE "pg_catalog"."default",
    "password"     varchar(255) COLLATE "pg_catalog"."default",
    "create_time"  timestamp(6)                                        DEFAULT CURRENT_TIMESTAMP,
    "update_time"  timestamp(6)                                        DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "offical_user_pkey" PRIMARY KEY ("id")
)
;
COMMENT
ON COLUMN "chat"."user"."nick_name" IS '昵称';

COMMENT
ON COLUMN "chat"."user"."phone_number" IS '手机号';

COMMENT
ON COLUMN "chat"."user"."profile" IS '头像';

COMMENT
ON COLUMN "chat"."user"."password" IS '密码';

