SELECT r.name, rcp.name, m.name FROM resources r 
JOIN recipes_outputs ro ON r.id = ro.resources_id
JOIN recipes rcp ON ro.recipes_id = rcp.id
JOIN machines_recipes mr ON rcp.id = mr.recipes_id
JOIN machines m ON mr.machines_id = m.id
WHERE r.name = 'iron_plate';

SELECT rcp.name, m.name, r.name, ro.amount, ri.amount, rcp.production_time_s FROM resources r 
JOIN recipes_outputs ro ON r.id = ro.resources_id
JOIN recipes rcp ON ro.recipes_id = rcp.id
JOIN recipes_inputs ri ON rcp.id = ri.recipes_id
JOIN machines_recipes mr ON rcp.id = mr.recipes_id
JOIN machines m ON mr.machines_id = m.id
WHERE r.name = 'iron_plate';

--Select name and amount of resources that recipe uses as inputs
SELECT r.name AS resource_name, ri.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ri.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipes rcp
JOIN recipes_inputs ri ON rcp.id = ri.recipes_id
JOIN resources r ON ri.resources_id = r.id
JOIN machines_recipes mr ON rcp.id = mr.recipes_id
JOIN machines m ON mr.machines_id = m.id
WHERE rcp.id IN (
    --Select id of recipe(s), that have specified resource as an output
    SELECT rcp.id FROM recipes rcp
    JOIN recipes_outputs ro ON rcp.id = ro.recipes_id
    JOIN resources r ON ro.resources_id = r.id
    WHERE r.name = 'reinforced_iron_plate'
);

--Select name and amount of resources specified alongside recipe(s) that produce(s) them, machine(s) that use(s) that recipe and machine's speed
--Chooses the fastest production rate
SELECT rcp.id, r.name AS recipe_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipes rcp
JOIN recipes_outputs ro ON rcp.id = ro.recipes_id
JOIN resources r ON ro.resources_id = r.id
JOIN machines_recipes mr ON rcp.id = mr.recipes_id
JOIN machines m ON mr.machines_id = m.id
WHERE r.name = 'reinforced_iron_plate' AND rcp.default_choice = TRUE -- rcp.name IN (name_1, name_2)
ORDER BY rate
LIMIT 1;

SELECT rcp.id, rcp.name AS recipe_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate 
							FROM recipes rcp
							JOIN recipes_outputs ro ON rcp.id = ro.recipes_id
							JOIN resources r ON ro.resources_id = r.id
							JOIN machines_recipes mr ON rcp.id = mr.recipes_id
							JOIN machines m ON mr.machines_id = m.id
							WHERE r.name = 'iron_rod' AND (rcp.default_choice = TRUE OR rcp.name IN ())
							AND rcp.users_id = '1'
							AND ro.users_id = '1'
							AND r.users_id = '1'
							AND mr.users_id = '1'
							AND m.users_id = '1'
                            ORDER BY rcp.default_choice, rate
				            LIMIT 1;