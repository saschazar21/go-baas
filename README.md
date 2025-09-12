<div align="center">
  <h1><a href="https://go-baas.netlify.app/" rel="noopener noreferrer">BaaS</a></h1>
  <strong>A RESTful web API for storing boolean values, written in Go. 🌐⚫⚪</strong>
  <br />
  <br />
  <a href="https://go-baas.netlify.app"><img src="https://img.shields.io/badge/netlify-deployed-green" alt="Netlify Status"></a> <a href="LICENSE"><img alt="License" src="https://img.shields.io/github/license/saschazar21/go-baas"></a>
  <br />
  <br />
  <a href="https://app.netlify.com/start/deploy?repository=https://github.com/saschazar21/go-baas"><img src="https://www.netlify.com/img/deploy/button.svg" alt="Deploy to Netlify"></a>
  <br />
  <br />
  <br />
</div>

## What is it?

`go-baas` is a simple RESTful web API for storing boolean values. It is written in Go and uses a [Redis database](https://redis.io) for persistence. The API is designed to be as simple as possible, with only two endpoints: one for creating a new boolean value, and one for retrieving/updating/toggling/deleting it.

The whole project is ready to be deployed on [Netlify](https://www.netlify.com/), but may be deployed on any other platform that supports Go and Redis.

### But why?

Mainly [for shits and giggles.](https://en.wiktionary.org/wiki/for_shits_and_giggles) 😄

## How to use it?

The API is designed to be as simple as possible. It only has two endpoints:

### `/api/v1/booleans`

- `POST /api/v1/booleans` to create a new boolean value:

  ```json
  {
    "label": "an optional label for the boolean value",
    "value": true
  }
  ```

  A `curl` request would look like the following:

  ```bash
  curl -X POST https://go-baas.netlify.app/api/v1/booleans -d '{"label": "optional label", "value": true}' -H "Content-Type: application/json"
  ```

  > ℹ️ By adding `expires_in` or `expires_at` as query parameter, the boolean value will be deleted after the specified time. The value for `expires_in` is in seconds from now, while `expires_at` is the desired expiration date as a Unix epoch in seconds. If both are provided, `expires_at` will be used.

  The response will be the newly created boolean value, including its unique ID:

  ```json
  {
    "id": "a unique ID",
    "label": "an optional label for the boolean value",
    "value": true
  }
  ```

### `/api/v1/booleans/:id`

- `GET /api/v1/booleans/:id` to retrieve a boolean value:

  ```bash
  curl -X GET https://go-baas.netlify.app/api/v1/booleans/:id
  ```

- `PUT /api/v1/booleans/:id` to update a boolean value:

  ```json
  {
    "label": "change to a new label",
    "value": false
  }
  ```

  A `curl` request would look like the following:

  ```bash
  curl -X PUT https://go-baas.netlify.app/api/v1/booleans/:id -d '{"label": "changed label", "value": false}' -H "Content-Type: application/json"
  ```

- `PATCH /api/v1/booleans/:id` to toggle a boolean value:

  ```bash
  curl -X PATCH https://go-baas.netlify.app/api/v1/booleans/:id
  ```

- `DELETE /api/v1/booleans/:id` to delete a boolean value:

  ```bash
  curl -X DELETE https://go-baas.netlify.app/api/v1/booleans/:id
  ```

## How to deploy it?

The project is ready to be deployed on Netlify. Just click the "Deploy to Netlify" button above, and follow the instructions. You will need to provide your Redis connection string as an environment variable.

If you want to deploy it on a different platform, you will need to set up its deploy environment and a Redis database. It's best to check out the deployment docs of the preferred platform.

## License

Licensed under the MIT license.

Copyright ©️ 2025 [Sascha Zarhuber](https://sascha.work)
