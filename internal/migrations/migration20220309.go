package migrations

import (
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type Migration20220309 struct {
}

func (m Migration20220309) GetCreateAt() time.Time {
	return time.Date(2022, 3, 9, 0, 0, 0, 0, time.Local)
}

func (m Migration20220309) Upgrade(ormer orm.Ormer) {
	ormer.Raw("UPDATE `sys_integrate_setting` SET `type`='registry' WHERE `type`='harbor';").Exec()

	ormer.Raw(strings.ReplaceAll(`DROP PROCEDURE IF EXISTS <|SPIT|>ModifyHarborToRegistry<|SPIT|>;
delimiter $$
CREATE PROCEDURE <|SPIT|>ModifyHarborToRegistry<|SPIT|>()
BEGIN
    DECLARE HARBOREXISTS int DEFAULT 0;
    DECLARE REGISTRYEXISTS int DEFAULT 0;
    SELECT count(1) INTO @HARBOREXISTS FROM information_schema.COLUMNS WHERE TABLE_NAME='project_env' AND COLUMN_NAME='harbor';
    SELECT count(1) INTO @REGISTRYEXISTS FROM information_schema.COLUMNS WHERE TABLE_NAME='project_env' AND COLUMN_NAME='registry';
    IF @HARBOREXISTS>0 AND @REGISTRYEXISTS=0 #存在harbor列 不存在registry列时 直接修改列名
    THEN
        ALTER TABLE <|SPIT|>project_env<|SPIT|> CHANGE COLUMN <|SPIT|>harbor<|SPIT|> <|SPIT|>registry<|SPIT|> bigint(20) NOT NULL DEFAULT 0;
    ELSEIF  @HARBOREXISTS>0 AND @REGISTRYEXISTS>0 #harbor列和registry都存在时迁移数据并删除harbor列
    THEN
        UPDATE <|SPIT|>project_env<|SPIT|> SET <|SPIT|>registry<|SPIT|>=<|SPIT|>harbor<|SPIT|>;
        ALTER TABLE <|SPIT|>project_env<|SPIT|> DROP COLUMN <|SPIT|>harbor<|SPIT|>;
    END IF;
END;
$$
delimiter ;
CALL <|SPIT|>ModifyHarborToRegistry<|SPIT|>;
DROP PROCEDURE IF EXISTS <|SPIT|>ModifyHarborToRegistry<|SPIT|>;`, "<|SPIT|>", "`")).Exec()

}
