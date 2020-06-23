-- Calculate PR Open metric for repository
DROP PROCEDURE IF EXISTS calculate_metric_open_pr;
DELIMITER $$  
CREATE PROCEDURE calculate_metric_open_pr(
   IN repositoryId INT
)
	main: BEGIN
	   	DELETE FROM metric WHERE repository_id = repositoryId AND type = 5;
	   	
		SELECT @minPrOpen := MIN(created_hour) FROM pull_request WHERE repository_id = repositoryId;
		
		IF @minPrOpen IS NULL THEN
			LEAVE main;
		END IF;
		
		SET @hour := STR_TO_DATE(@minPrOpen, '%Y%m%d%H');
		SET @end := NOW();
		
		metric_loop: LOOP
		IF @hour < @end THEN
				INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
				SELECT repository_id, 'master', 5, COUNT(*), DATE_FORMAT(@hour, '%Y'), DATE_FORMAT(@hour, '%Y%m'), DATE_FORMAT(@hour, '%Y%m%d'), DATE_FORMAT(@hour, '%Y%m%d%H')
				FROM pull_request
				WHERE repository_id = repositoryId AND created_hour <= DATE_FORMAT(@hour, '%Y%m%d%H')
					AND (closed_hour = 0 OR closed_hour IS NULL OR closed_hour > DATE_FORMAT(@hour, '%Y%m%d%H'));
				
				SET @hour := DATE_ADD(@hour, INTERVAL 1 HOUR);
		ELSE
			LEAVE metric_loop;
		END IF;
   END LOOP metric_loop;
END $$

-- Calculate for all repos
DROP PROCEDURE IF EXISTS calculate_metric_open_pr_all_repos;
DELIMITER $$
CREATE PROCEDURE calculate_metric_open_pr_all_repos()
BEGIN
		DECLARE repositoryId INT DEFAULT NULL;
		DECLARE done TINYINT DEFAULT FALSE;
		
		# cursor over repository id
		DECLARE cur
		CURSOR FOR
		SELECT id
		FROM repository;
		
		DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
		
		DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
	   	FETCH NEXT FROM cur INTO repositoryId;
		
		OPEN cur;
		
		main_loop: LOOP
			FETCH NEXT FROM cur INTO repositoryId;
			
			IF done THEN
				LEAVE main_loop;
			ELSE 
				CALL calculate_metric_open_pr(repositoryId);
			END IF;
		END LOOP;
		
		CLOSE cur;
END $$

--Run
CALL calculate_metric_open_pr_all_repos()