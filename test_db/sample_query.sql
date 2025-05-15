SELECT r.name, rcp.name, m.name FROM resources r 
JOIN recipies_outputs ro ON r.id = ro.resources_id
JOIN recipies rcp ON ro.recipies_id = rcp.id
JOIN machines_recipies mr ON rcp.id = mr.recipies_id
JOIN machines m ON mr.machines_id = m.id
WHERE r.name = 'iron_plate';

SELECT rcp.name, m.name, r.name, ro.amount, ri.amount, rcp.production_time_s FROM resources r 
JOIN recipies_outputs ro ON r.id = ro.resources_id
JOIN recipies rcp ON ro.recipies_id = rcp.id
JOIN recipies_inputs ri ON rcp.id = ri.recipies_id
JOIN machines_recipies mr ON rcp.id = mr.recipies_id
JOIN machines m ON mr.machines_id = m.id
WHERE r.name = 'iron_plate';

--Select name and amount of resources that recipie uses as inputs
SELECT r.name AS resource_name, ri.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ri.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipies rcp
JOIN recipies_inputs ri ON rcp.id = ri.recipies_id
JOIN resources r ON ri.resources_id = r.id
JOIN machines_recipies mr ON rcp.id = mr.recipies_id
JOIN machines m ON mr.machines_id = m.id
WHERE rcp.id IN (
    --Select id of recipie(s), that have specified resource as an output
    SELECT rcp.id FROM recipies rcp
    JOIN recipies_outputs ro ON rcp.id = ro.recipies_id
    JOIN resources r ON ro.resources_id = r.id
    WHERE r.name = 'reinforced_iron_plate'
);

--Select name and amount of resources specified alongside recipie(s) that produce(s) them, machine(s) that use(s) that recipie and machine's speed
--Chooses the fastest production rate
SELECT rcp.id, r.name AS recipie_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipies rcp
JOIN recipies_outputs ro ON rcp.id = ro.recipies_id
JOIN resources r ON ro.resources_id = r.id
JOIN machines_recipies mr ON rcp.id = mr.recipies_id
JOIN machines m ON mr.machines_id = m.id
WHERE r.name = 'reinforced_iron_plate' AND rcp.default_choice = TRUE -- rcp.name IN (name_1, name_2)
ORDER BY rate
LIMIT 1;