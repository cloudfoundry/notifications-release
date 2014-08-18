# Notifications Release
This release deploys the [notification service](https://github.com/cloudfoundry-incubator/notifications) as an application onto the CloudFoundry platform.
A running CF is required to deploy.
The application will be deployed into the `notifications-service` org, and the `notifications-service` space.
The service registers itself at the address matching http://notifications.$CF_APP_DOMAIN.

# UAA Client
Notifications requires a UAA client to boot. The client can be created with the following properties:
```yaml
scope: uaa.none
client_id: notifications
authorized_grant_types: client_credentials
authorities: scim.read,cloud_controller.admin
```

# Bosh-lite

1. Add SMTP properties to `./bosh-lite/notifications-smtp-stub.yml` file as follows:
  ```yaml
  properties:
    notifications:
      smtp:
        host: stmp.example.com
        port: 587
        user: my-user-name
        pass: my-password
  ```

2. Add DB properties to `./bosh-lite/notifications-db-stub.yml` file as follows:
  ```yaml
  properties:
    notifications:
      database_url: tcp://user:password@example.com:3306/dbname
  ```


3. Generate manifest:
  ```bash
  ./bosh-lite/make_manifest
  ```

4. Create and upload release:
  ```bash
  bosh create release
  bosh upload release
  ```

5. Deploy
  ```bash
  bosh deploy
  bosh run errand deploy-notifications
  ```
