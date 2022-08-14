create table request (id_request serial,id_Group integer,id_userInvite varchar(100), create_by varchar(100),update_by varchar(100),status integer default 1,created_at timestamp not null default now(),updated_at timestamp not null default now(),constraint FK_Request_CreateBy FOREIGN KEY (create_by) references UserDetail(user_id),constraint FK_Request_UpdateBy FOREIGN KEY (update_by) references UserDetail(user_id),constraint FK_Request_Invite   foreign key (id_userInvite) references UserDetail(user_id),constraint FK_REQUEST_GROUP   FOREIGN KEY (id_Group) REFERENCES groups(id_group),primary key (id_request));
