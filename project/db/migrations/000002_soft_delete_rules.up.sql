-- FILE : 000002_soft_delete_rules.up.sql

CREATE RULE "_soft_deletion" 
    AS ON DELETE TO "tb_user_profiles"
    DO INSTEAD UPDATE tb_user_profiles SET deleted_at = NOW() WHERE id = old.id AND deleted_at IS NULL;

CREATE RULE "_soft_deletion" 
    AS ON DELETE TO "tb_users"
    DO INSTEAD UPDATE tb_users SET deleted_at = NOW(), is_active = FALSE WHERE id = old.id AND deleted_at IS NULL;

CREATE RULE "_delete_users" 
    AS ON UPDATE TO "tb_users"
    WHERE old.deleted_at IS NULL AND new.deleted_at IS NOT NULL
    DO ALSO UPDATE tb_user_profiles SET deleted_at = NOW(), steam_user = NULL, gog_user = NULL WHERE user_id = old.id;