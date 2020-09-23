# Prism

Prism hopes to expose [primitive](https://github.com/fogleman/primitive) as a web service and allow easy uploading and manipulation of images.

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

### Profiles

Profiles make the life blood of the app. They are a set of options that can be called on primitive. This endpoint will be
basic functionanlity
/profile/create
/profile/read

### Refract

Refract will take in an image and config, and spit out a image transformed by primitive
/refract?profile={name}&image={uid}
/refract and then the request just has those in it?

when someone calls refract for the first time i should check if that file has already been refracted by this profile.
theres a few ways to do it.
i think take profile
profile should have a list of

#### Known Issues

#### TODO

Profiles
CRUD profile, UD!
pass in profile and image id -> generate new image marked as child of original image
