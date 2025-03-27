/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_955655590")

  // update collection data
  unmarshal({
    "listRule": "@request.auth.id != \"\" && assignees.id ?= @request.auth.id && id ?= @collection.lessons.course.id",
    "viewRule": "@request.auth.id != \"\" && assignees.id ?= @request.auth.id && id ?= @collection.lessons.course.id"
  }, collection)

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_955655590")

  // update collection data
  unmarshal({
    "listRule": null,
    "viewRule": null
  }, collection)

  return app.save(collection)
})
