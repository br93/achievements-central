-- FILE : 000002_soft_delete_rules.down.sql

DROP RULE "_soft_deletion" ON "tb_user_profiles";
DROP RULE "_soft_deletion" ON "tb_users";
DROP RULE "_delete_users" on "tb_users";