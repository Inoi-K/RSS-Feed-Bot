package database

import "context"

func SetUp(ctx context.Context) error {
	query := "create table if not exists source\n(\n    id    serial\n        constraint sources_pk\n            primary key,\n    title text,\n    url   text,\n    rss   text\n);"
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	query = "create table if not exists chat\n(\n    id   integer    not null\n        constraint table1_pk\n            primary key,\n    lang varchar(2) not null\n);"
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	query = "create table if not exists chat_source\n(\n    chatid   integer              not null\n        constraint usersource_users_id_fk\n            references chat,\n    sourceid integer              not null\n        constraint usersource_sources_id_fk\n            references source,\n    isactive boolean default true not null,\n    constraint usersource_pkey\n        primary key (chatid, sourceid)\n);"
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
