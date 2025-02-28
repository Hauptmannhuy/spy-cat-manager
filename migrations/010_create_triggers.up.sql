
CREATE OR REPLACE FUNCTION check_target_count()
RETURNS TRIGGER AS $$
DECLARE 
    target_count INT;
BEGIN
    SELECT COUNT(*) INTO target_count
    FROM targets 
    WHERE mission_id = NEW.mission_id;
    
    IF target_count >= 3 THEN
        RAISE EXCEPTION 'Mission reached maximum amount of targets';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;



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
        RAISE EXCEPTION 'Cannot add note to mission, mission is completed';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


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
  WHERE id = OLD.id;
  IF target_completed THEN
    RAISE EXCEPTION 'Cannot delete, target is completed';
  END IF;
  RETURN OLD;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER check_target_before_delete
BEFORE DELETE on targets
FOR EACH ROW
EXECUTE FUNCTION check_target_before_delete();





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
$$ LANGUAGE plpgsql;


CREATE TRIGGER check_mission_before_add_target
BEFORE INSERT on targets
FOR EACH ROW
EXECUTE FUNCTION check_mission_before_add_target();


CREATE OR REPLACE FUNCTION check_mission_before_assign_spy()
RETURNS TRIGGER AS $$
DECLARE 
  mission_completed BOOLEAN;
BEGIN 
  SELECT completed INTO mission_completed 
  FROM missions
  WHERE id = NEW.id; 

  IF mission_completed THEN
    RAISE EXCEPTION 'Cannot assign spy, mission is completed';
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER mission_spy_check
BEFORE INSERT OR UPDATE ON missions
FOR EACH ROW
WHEN (NEW.spy_id IS NOT NULL) 
EXECUTE FUNCTION check_mission_before_assign_spy();

CREATE OR REPLACE FUNCTION check_mission_before_delete()
RETURNS TRIGGER AS $$
DECLARE
  assigned_spy_id INTEGER;
BEGIN

  SELECT spy_id INTO assigned_spy_id
  FROM missions
  WHERE OLD.id = id;

  IF assigned_spy_id IS NOT NULL THEN
    RAISE EXCEPTION 'Cannot delete mission, already assigned to a cat';
  END IF;

    RETURN OLD;

END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER check_mission_before_delete
BEFORE DELETE ON missions
FOR EACH ROW
EXECUTE FUNCTION check_mission_before_delete();



CREATE OR REPLACE FUNCTION mark_complete_mission()
RETURNS TRIGGER AS $$
DECLARE
  completed_targets INTEGER;
BEGIN
  SELECT COUNT(*) INTO completed_targets
  FROM targets
  WHERE mission_id = OLD.mission_id AND completed = TRUE;

  if completed_targets = 3 THEN
    UPDATE missions SET completed = TRUE WHERE id = OLD.mission_id;
  END IF;
  RETURN OLD;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER mark_complete_mission
AFTER UPDATE ON targets
FOR EACH ROW
EXECUTE FUNCTION mark_complete_mission();