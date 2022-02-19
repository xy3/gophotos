## Photo gallery API in Go

I tried to keep the amount of time working on this to below 3 hours but it was taking too long to finish.

### Paths:

Base Path: `/api/v1`
Host: `http://localhost:8090`


- /photos
  - GET
    If `photo_id` is set, this endpoint returns a single photo. Otherwise, it returns a paginated list of results
  - PUT
   Upload a photo
  - DELETE
   Delete a photo
- /photos/download
  - GET
   Download a photo
- /users
  - GET
  Get the current user information
  - DELETE
  Delete the current user
- /users/signup
  - POST
   Create an account

This is not complete documentation, it is work in progress and will be updated shortly :) 
