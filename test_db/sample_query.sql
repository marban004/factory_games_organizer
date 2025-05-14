SELECT r.name, rcp.name, m.name FROM resources r 
JOIN recipies_outputs ro ON r.id = ro.resources_id
JOIN recipies rcp ON ro.recipies_id = rcp.id
JOIN machines_recipies mr ON rcp.id = mr.recipies_id
JOIN machines m ON mr.machines_id = m.id
WHERE r.name = 'iron_plate';

SELECT rcp.name, m.name, r.name, ro.amount, ri.amount, rcp.production_time, rcp.production_time_unit FROM resources r 
JOIN recipies_outputs ro ON r.id = ro.resources_id
JOIN recipies rcp ON ro.recipies_id = rcp.id
JOIN recipies_inputs ri ON rcp.id = ri.recipies_id
JOIN machines_recipies mr ON rcp.id = mr.recipies_id
JOIN machines m ON mr.machines_id = m.id
WHERE r.name = 'iron_plate';

--Select name and amount of resources that recipie uses as inputs
SELECT r.name, ri.amount, rcp.production_time, rcp.production_time_unit, m.name, m.speed FROM recipies rcp 
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
SELECT r.name, ro.amount, rcp.production_time, rcp.production_time_unit, m.name, m.speed FROM recipies rcp
JOIN recipies_outputs ro ON rcp.id = ro.recipies_id
JOIN resources r ON ro.resources_id = r.id
JOIN machines_recipies mr ON rcp.id = mr.recipies_id
JOIN machines m ON mr.machines_id = m.id
WHERE r.name = 'reinforced_iron_plate';