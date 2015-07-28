# Notifications Release
This release deploys the [notification service](https://github.com/cloudfoundry-incubator/notifications) as an application onto the CloudFoundry platform.
A running CF is required to deploy.
The application will be deployed into the `system` org, and the `notifications-service` space.
The service registers itself at the address matching http://notifications.$CF_APP_DOMAIN.

# Prerequisites
1. Running UAA. This requirement is typically satisfied by having [CloudFoundry](https://github.com/cloudfoundry/cf-release) deployed.
1. Running MySQL instance. One option is to deploy the [CloudFoundry MySQL release](https://github.com/cloudfoundry/cf-mysql-release).

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

2. Add DB properties to `./bosh-lite/notifications-db-stub.yml` file for your running
   MySQL instance as follows:
  ```yaml
  properties:
    notifications:
      database:
        url: tcp://user:password@example.com:3306/dbname
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

# AWS
In addition to the stub file properties specified in the BOSH-Lite manifest, an AWS manifest stub file will require
some extra infrastructure specific fields. Included below is an example:
```yaml
infrastructure_properties:
  availability_zone: us-east-1a
```
