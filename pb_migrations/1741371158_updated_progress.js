/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_1649388127")

  // update collection data
  unmarshal({
    "listRule": "@request.auth.id != \"\" && assignee = @request.auth.id",
    "updateRule": "@request.auth.id != \"\" && assignee = @request.auth.id && @request.body.course:isset = false && @request.body.assignee:isset = false",
    "viewRule": "@request.auth.id != \"\" && assignee = @request.auth.id"
  }, collection)

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_1649388127")

  // update collection data
  unmarshal({
    "listRule": null,
    "updateRule": null,
    "viewRule": null
  }, collection)

  return app.save(collection)
})
