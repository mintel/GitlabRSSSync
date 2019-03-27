# Gitlab RSS Sync
Create Gitlab issues from RSS Feeds

## Config file
```$yaml
interval: 300 // Interval in seconds to check the RSS feeds.
feeds:
  - id: test //Specify a feed ID that is used internally for duplicate detection.
    feed_url: http://example.com/rss.xml // The Feed URL.
    name: Test Feed // A User friendly display name.
    gitlab_project_id: 12345 // The Gitlab project ID to create issues under.
    labels: // A list of labels to add to created Issues.
      - TestLabel
   - id: ...
```

## Docker
A Docker image is made available on [DockerHub](https://hub.docker.com/r/adamhf/gitlabrsssync)

### Required Environment Variables
* GITLAB_API_TOKEN - Gitlab personal access token that will be used to create Issues NOTE: You must have access to create
issues in the projects you specify in the config file.
* CONFIG_DIR - The directory the application should look for config.yaml in.
* DATA_DIR - The directory the application should look for (or create) the state.db in.

### Volume mounts
Make sure the location of your DATA_DIR environment variable is set to a persistant volume / mount as the database
that is contained within it stores the state of which RSS items have already been synced.

### Run it
`docker run -e GITLAB_API_TOKEN=<INSERT_TOKEN> -e DATA_DIR=/data -e CONFIG_DIR=/app -v <PATH_TO_DATA_DIR>:/data -v <PATH_TO_CONFIG_DIR>/config adamhf/rss-sync:latest`