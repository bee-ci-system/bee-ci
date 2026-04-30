CREATE OR REPLACE FUNCTION bee_schema.builds_trigger() RETURNS TRIGGER AS
$$
BEGIN
    PERFORM pg_notify('builds_channel', row_to_json(NEW)::TEXT);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER builds_notify_trigger
    AFTER INSERT OR UPDATE
    ON bee_schema.builds
    FOR EACH ROW
EXECUTE FUNCTION bee_schema.builds_trigger();
