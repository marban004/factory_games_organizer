SELECT r.name, rcp.name, m.name FROM resources r 
JOIN recipies_outputs ro ON r.id = ro.resources_id
JOIN recipies rcp ON ro.recipies_id = rcp.id
JOIN machines_recipies mr ON rcp.id = mr.recipies_id
JOIN machines m ON mr.machines_id = m.id
WHERE r.name = 'iron_plate';