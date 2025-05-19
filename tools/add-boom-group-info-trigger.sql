-- 为Postgres的 boom_group 表添加触发器,当有数据插入和更新时想指定的消息队列发送消息
DROP TRIGGER IF EXISTS boom_group_info_notify_trigger on boom_group_info;

-- 创建触发器函数
CREATE OR REPLACE FUNCTION boom_group_info_trigger() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        PERFORM pg_notify('boomGroupInfo',json_build_object('oprt', tg_op, 'data', row_to_json(NEW))::text);
    ELSEIF TG_OP = 'DELETE' THEN
        PERFORM pg_notify('boomGroupInfo',json_build_object('oprt', tg_op, 'id', NEW.token)::text);
    END IF;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- 创建触发器
CREATE TRIGGER boom_group_info_notify_trigger
    AFTER INSERT OR UPDATE OR DELETE ON boom_group_info
    FOR EACH ROW
    EXECUTE FUNCTION boom_group_info_trigger();


