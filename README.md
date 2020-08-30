# Refract

Refract hopes to expose [primitive](https://github.com/fogleman/primitive) as a web service and allow easy uploading and manipulation of images.

## API

---

### Image

The image api is used for retriveing the raw imags from the database.
Make a GET request against `/image/{uid}` to receive the database row.
`/image/{uid}/children`
gets all images, that have been tagged as the 'child' of that image.

Admins can also execute DELETE commands on the same URL

/upload accepts POST requests and uploads an image to the database. Provided correct image structure.

/upload/{url} in the works

---

#### Known Issues

Database accepts duplicate files.

#### TODO

Need random image endpoint
need search tags
add notion of tags
Need to see how many pages of images have been uploaded
Need to block duplicates...

swap extension to filename --->

learn react hooks? i guess

Profiles
CRUD profile,
pass in profile and image id -> generate new image marked as child of original image
