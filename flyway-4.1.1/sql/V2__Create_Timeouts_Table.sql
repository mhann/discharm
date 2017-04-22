CREATE TABLE "timeouts" (
	"id" serial NOT NULL,
	"guild_id" integer NOT NULL,
	"target_user_id" integer NOT NULL,
	"creator_user_id" integer NOT NULL,
	CONSTRAINT timeouts_pk PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);