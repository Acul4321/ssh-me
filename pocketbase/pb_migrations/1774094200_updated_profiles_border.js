/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_3414089001");

  const field = new Field({
    "autogeneratePattern": "",
    "hidden": false,
    "id": "text1234567890",
    "max": 0,
    "min": 0,
    "name": "border",
    "pattern": "",
    "presentable": false,
    "primaryKey": false,
    "required": false,
    "system": false,
    "type": "text"
  });

  collection.fields.add(field);
  return app.save(collection);
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_3414089001");
  collection.fields.removeById("text1234567890");
  return app.save(collection);
});
