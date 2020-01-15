alter table "resource_votes" rename column "account" to "user_id";
alter table "resource_votes" rename column "resource" to "resource_id";
alter table resource_votes rename constraint resource_votes_account_fkey to resource_votes_user_fkey;
