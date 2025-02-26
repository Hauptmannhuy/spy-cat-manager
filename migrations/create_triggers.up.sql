

CREATE OR REPLACE FUNCTION check_target_count()
RETURNS TRIGGER AS $$
BEGIN
    DECLARE target_count INT;
    SELECT COUNT(*) INTO target_count
    FROM targets 
    WHERE mission_id = NEW.mission_id;
    
    IF target_count >= 3 THEN
        RAISE EXCEPTION = 'Не может быть больше трех целей';
    END IF;
END;

CREATE TRIGGER check_target_count
BEFORE INSERT ON targets
FOR EACH ROW
EXECUTE FUNCTION check_target_count();


CREATE OR REPLACE FUNCTION check_completion_status() 
RETURNS TRIGGER AS $$
DECLARE
    mission_completed BOOLEAN;
    target_completed BOOLEAN;
BEGIN
    SELECT completed INTO target_completed 
    FROM targets 
    WHERE id = NEW.target_id;

    SELECT completed INTO mission_completed 
    FROM missions 
    WHERE id = (SELECT mission_id FROM targets WHERE id = NEW.target_id);

    IF target_completed OR mission_completed THEN
        RAISE EXCEPTION 'Нельзя обновлять заметки, если цель или миссия уже завершены';
    END IF;

    RETURN NEW;
END;
-- $$ LANGUAGE plpgsql;

CREATE TRIGGER check_complition_status
BEFORE UPDATE ON notes
FOR EACH ROW
EXECUTE FUNCTION check_completion_status();



CREATE OR REPLACE FUNCTION check_target_before_delete()
RETURNS TRIGGER AS $$
DECLARE
  target_completed BOOLEAN;
BEGIN
  SELECT completed INTO target_completed
  FROM targets
  WHERE id = NEW.id;
  IF target_completed THEN
    RAISE EXCEPTION 'Cannot delete, target is completed';
  END IF;
  RETURN DELETE
END;

CREATE TRIGGER check_target_before_delete
BEFORE DELETE on targets
FOR EACH ROW
EXECUTE FUNCTION check_target_before_delete();






-- check_mission_before_add_target
CREATE OR REPLACE FUNCTION check_mission_before_add_target()
RETURNS TRIGGER AS $$
DECLARE
  mission_completed BOOLEAN;
BEGIN
  SELECT completed INTO mission_completed
  FROM missions
  WHERE id = NEW.mission_id;

  if mission_completed THEN
    RAISE EXCEPTION 'Cannot add new target, mission is completed';
  END IF;
  RETURN NEW;
END;

CREATE TRIGGER check_mission_before_add_target
BEFORE INSERT on targets
FOR EACH ROW
EXECUTE FUNCTION check_mission_before_add_target();