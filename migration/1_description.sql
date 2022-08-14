CREATE TABLE Groups(id_group serial ,owner_id varchar(100) NOT NULL,name varchar(100),type varchar(15) NOT NULL,private boolean,thumbnail varchar(255),description varchar(255),created_at timestamp not null default now(),updated_at timestamp not null default now(),deleted_at timestamp,PRIMARY KEY (id_group));
CREATE TABLE Groups_Users(id_group serial,user_id varchar(100),join_at timestamp not null default now(),PRIMARY KEY (id_group,user_id),CONSTRAINT FK_Groups_Groups_Users FOREIGN KEY (id_group) REFERENCES Groups (id_group));
CREATE TABLE Messages(id_mess serial,user_sender varchar(100) ,content text,id_group serial NOT NULL ,created_at timestamp not null default now(),updated_at timestamp not null default now(),deleted_at timestamp,PRIMARY KEY (id_mess),CONSTRAINT FK_Groups_Messages FOREIGN KEY (id_group) REFERENCES Groups (id_group));
CREATE TABLE Messages_Delete(id_mess serial,user_deleted varchar(100),PRIMARY KEY (id_mess,user_deleted),CONSTRAINT FK_Messages_Messages_Delete FOREIGN KEY (id_mess) REFERENCES Messages (id_mess));
CREATE TABLE UserDetail(user_id varchar(100),role varchar(15),created_at timestamp not null default now(),updated_at timestamp not null default now(),deleted_at timestamp,PRIMARY KEY (user_id));
CREATE TABLE ONLINE(id_onl serial,hostname varchar(100),socket_id varchar(100),user_id varchar(100) ,log_at timestamp not null default now(),CONSTRAINT FK_ONLINE_USER FOREIGN KEY (user_id) REFERENCES UserDetail (user_id),PRIMARY KEY (id_onl));
CREATE TABLE Article (id_article serial,content text,title text,thumbnail text,version integer default 1,create_by varchar(100),update_by varchar(100),created_at timestamp not null default now(),updated_at timestamp not null default now(),primary key (id_article),constraint FK_Article_CreateBy FOREIGN KEY (create_by) references UserDetail(user_id),constraint FK_Article_UpdateBy FOREIGN KEY (update_by) references UserDetail(user_id));
create table Category (id_category serial,name text,parentId integer default -1,num integer default 0,version integer default 1,create_by varchar(100),update_by varchar(100),created_at timestamp not null default now(),updated_at timestamp not null default now(),constraint FK_Comment_CreateBy FOREIGN KEY (create_by) references UserDetail(user_id),constraint FK_Comment_UpdateBy FOREIGN KEY (update_by) references UserDetail(user_id),primary key (id_category));
CREATE TABLE Comment (id_cmt serial,id_article integer,content text,type text,parentId integer default -1,num integer default 0,version integer default 1,create_by varchar(100),update_by varchar(100),created_at timestamp not null default now(),updated_at timestamp not null default now(),primary key (id_cmt),constraint FK_Comment_Article FOREIGN KEY (id_article) REFERENCES Article(id_article),constraint FK_Comment_CreateBy FOREIGN KEY (create_by) references UserDetail(user_id),constraint FK_Comment_UpdateBy FOREIGN KEY (update_by) references UserDetail(user_id));
create table Reaction (id_reaction serial,id_article integer,type integer,create_by varchar(100),update_by varchar(100),created_at timestamp not null default now(),updated_at timestamp not null default now(),constraint FK_Reaction_CreateBy FOREIGN KEY (create_by) references UserDetail(user_id),constraint FK_Reaction_UpdateBy FOREIGN KEY (update_by) references UserDetail(user_id),constraint FK_Reaction_Article FOREIGN KEY (id_article) REFERENCES Article(id_article),primary key (id_reaction));

Alter table Messages add column parentID integer;
Alter table Messages add column numChild integer default 0;
ALTER table Messages add CONSTRAINT FK_Groups_Messages_Child FOREIGN KEY (parentID) REFERENCES Messages(id_mess);
Alter table Messages add column type text;
update messages set type = 'TEXT' WHERE type IS NULL ;
ALTER TABLE messages ALTER COLUMN type TYPE integer USING (type::integer);
ALTER TABLE messages ALTER COLUMN type SET default 1;

alter table Article add column slug text;
alter table article add column id_category serial;
alter table article add constraint FK_ARTICLE_Category FOREIGN KEY (id_category) references Category(id_category);
alter table public.category add column slug text;

alter table Article add column num_react integer default 0;
alter table Article add column num_cmt integer default 0;
alter table Article add column num_share integer default 0;

truncate table Messages cascade;

alter table public.userdetail add column avatar text default '';

ALTER TABLE UserDetail DROP COLUMN IF EXISTS avatar;
alter table UserDetail add column avatar text default '';

create unique index online_socket_id_uindex on online (socket_id);