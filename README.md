## Photo gallery API in Go

### Approach

I used a DDD (Domain-driven design) approach when developing this API. I focused solely on doing exactly what was required for this project - an API for storing and retrieving photos.

### Usage

To run the server run:

```
cd cmd
go run main.go
```

### API Specification

You can view the OpenAPI 3.0 specification here: https://app.swaggerhub.com/apis/xy3/photo-gallery/1.0.0

If you use Insomnia you can import the spec using this button:

[![Run in Insomnia}](https://insomnia.rest/images/run.svg)](https://insomnia.rest/run/?label=Photo%20gallery&uri=https%3A%2F%2Fraw.githubusercontent.com%2Fxy3%2Fgophotos%2Fmain%2Fphotos-v1.yaml)


### Functionality

A general flow would look like this:

A user:
1. Signs up on `POST /user/signup` with `email` and `password`
2. Checks their authentication status using `POST /user/signin`
3. Adds BasicAuth to their next requests as a header
4. Uploads a photo to `PUT /photo` with formdata containing a file with the key `photo`
5. Downloads or deletes a photo at `GET /photo` and `DELETE /photo` providing the `photo_id`
6. Gets or updates information of a photo at `GET /photo/info` or `PATCH /photo/info` providing the `photo_id`
7. Lists photo information at `GET /photo/list` with `page` and `pageSize` for pagination
8. Gets account information at `GET /user`
9. Deletes their account at `DELETE /user`

When a user signs up, a directory is created for them under the `photos_storage` directory, and subsequent photo uploads are stored in that directory using their file hash as the file name - to avoid collisions.

### Thanks!

Thanks, this was fun!

Theodore Coyne Morgan

hi@theodore.ie
