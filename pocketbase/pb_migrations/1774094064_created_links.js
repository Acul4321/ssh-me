/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = new Collection({
    "createRule": "profile.user.id = @request.auth.id",
    "deleteRule": "profile.user.id = @request.auth.id",
    "fields": [
      {
        "autogeneratePattern": "[a-z0-9]{15}",
        "hidden": false,
        "id": "text3208210256",
        "max": 15,
        "min": 15,
        "name": "id",
        "pattern": "^[a-z0-9]+$",
        "presentable": false,
        "primaryKey": true,
        "required": true,
        "system": true,
        "type": "text"
      },
      {
        "cascadeDelete": true,
        "collectionId": "pbc_3414089001",
        "hidden": false,
        "id": "relation2170006031",
        "maxSelect": 1,
        "minSelect": 0,
        "name": "profile",
        "presentable": false,
        "required": true,
        "system": false,
        "type": "relation"
      },
      {
        "autogeneratePattern": "",
        "hidden": false,
        "id": "text961728715",
        "max": 0,
        "min": 0,
        "name": "platform",
        "pattern": "",
        "presentable": false,
        "primaryKey": false,
        "required": true,
        "system": false,
        "type": "text"
      },
      {
        "autogeneratePattern": "",
        "hidden": false,
        "id": "text245846248",
        "max": 0,
        "min": 0,
        "name": "label",
        "pattern": "",
        "presentable": false,
        "primaryKey": false,
        "required": false,
        "system": false,
        "type": "text"
      },
      {
        "exceptDomains": null,
        "hidden": false,
        "id": "url4101391790",
        "name": "url",
        "onlyDomains": null,
        "presentable": false,
        "required": true,
        "system": false,
        "type": "url"
      },
      {
        "hidden": false,
        "id": "number1169138922",
        "max": null,
        "min": null,
        "name": "sort_order",
        "onlyInt": false,
        "presentable": false,
        "required": false,
        "system": false,
        "type": "number"
      },
      {
        "hidden": false,
        "id": "bool2058414169",
        "name": "visible",
        "presentable": false,
        "required": false,
        "system": false,
        "type": "bool"
      }
    ],
    "id": "pbc_449060851",
    "indexes": [],
    "listRule": "",
    "name": "links",
    "system": false,
    "type": "base",
    "updateRule": "profile.user.id = @request.auth.id",
    "viewRule": ""
  });

  return app.save(collection);
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_449060851");
  return app.delete(collection);
})
